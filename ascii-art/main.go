package main

import (
	"fmt"
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

func main() {
	var input, bannerName, outputFile string
	args := os.Args[1:]
	if len(args) < 1 || len(args) > 3 {
		fmt.Println("Usage: go run main.go [OPTION] [STRING] [BANNER]")
		return
	}
	if strings.HasPrefix(args[0], "--output=") {
		parts := strings.SplitN(args[0], "=", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid flag format. Use: --output=<filename>")
			return
		}
		outputFile = parts[1]
		args = args[1:]
	}
	if len(args) == 1 {
		input = args[0]
		bannerName = "standard"
	} else if len(args) == 2 {
		input = args[0]
		bannerName = args[1]
	}
	banner, exists := banners[bannerName]
	if !exists {
		fmt.Printf("Error: Banner '%s' not found.\n", bannerName)
		return
	}
	processedLines := handleNewlines(input)
	asciiArt := generateAsciiArt(processedLines, banner)

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(asciiArt), 0664)
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	} else {
		fmt.Print(asciiArt)
	}
}

func handleNewlines(input string) []string {
	return strings.Split(input, "\\n")
}

func generateAsciiArt(lines []string, banner Banner) string {
	result := ""
	for _, line := range lines {
		if line == "" {
			result += "\n"
			continue
		}
		result += processLine(line, banner)
	}
	return result
}

func processLine(line string, banner Banner) string {
	result := ""
	for i := 1; i <= banner.lineHeight; i++ {
		res := ""
		for _, letter := range line {
			res += getLine(1+int(letter-32)*(banner.lineHeight+1)+i, banner.filePath)
		}
		result += res + "\n"
	}
	return result
}

func getLine(num int, filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading banner file:", err)
		os.Exit(1)
	}
	lines := strings.Split(string(content), "\n")
	if num-1 < len(lines) {
		return strings.ReplaceAll(lines[num-1], "\r", "")
	}
	return ""
}
