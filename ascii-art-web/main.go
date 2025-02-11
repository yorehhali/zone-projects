package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Banner struct {
	filePath   string
	lineHeight int
}

var banners = map[string]Banner{
	"thinkertoy": {"banners/thinkertoy.txt", 8},
	"standard":   {"banners/standard.txt", 8},
	"shadow":     {"banners/shadow.txt", 8},
	"phoenix":    {"banners/phoenix.txt", 7},
	"blocks":     {"banners/blocks.txt", 11},
	"arob":       {"banners/arob.txt", 8},
	"coins":      {"banners/coins.txt", 8},
	"fire":       {"banners/fire.txt", 9},
	"jacky":      {"banners/jacky.txt", 8},
	"small":      {"banners/small.txt", 5},
}

type PageData struct {
	Art string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir("frontend"))))
	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        tmpl := template.Must(template.ParseFiles("frontend/index.html"))
        tmpl.Execute(w, nil)
        return
    }

    if r.Method == http.MethodPost {
        text := r.FormValue("text")
        bannerName := r.FormValue("banner")

        banner, exists := banners[bannerName]
        if !exists {
            http.Error(w, "Invalid banner", http.StatusBadRequest)
            return
        }

        processedLines := handleNewlines(text)
        asciiArt := generateAsciiArt(processedLines, banner)

        data := PageData{Art: asciiArt}
        tmpl := template.Must(template.ParseFiles("frontend/index.html"))
        tmpl.Execute(w, data)
        return
    }

    http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}


func handleNewlines(input string) []string {
	return strings.Split(input, "\\n")
}

func generateAsciiArt(lines []string, banner Banner) string {
	result := ""
	for _, line := range lines {
		if line != "" {  // Avoid adding empty lines
			result += processLine(line, banner)
		}
	}
	return result
}


func processLine(line string, banner Banner) string {
	result := "\n"
	for i := 1; i <= banner.lineHeight; i++ {
		res := ""
		for _, letter := range line {
			res += getLine(1+int(letter-32)*(banner.lineHeight+1)+i, banner.filePath)
		}
		result += res+ "\n" 
	}
	result += "\n" // Only add a newline at the end of the processed line
	return result
}


func getLine(num int, filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "Error reading banner file."
	}

	lines := strings.Split(string(content), "\n")
	if num-1 < len(lines) {
		return strings.ReplaceAll(lines[num-1], "\r", "")
	}
	return ""
}
