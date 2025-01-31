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
	"standard":   {"banners/standard.txt", 9},
	"thinkertoy": {"banners/thinkertoy.txt", 9},
	"blocks":     {"banners/blocks.txt", 12},
	"shadow":     {"banners/shadow.txt", 9},
}

var colors = map[string]string{
	"red":    "\033[31m",
	"green":  "\033[32m",
	"yellow": "\033[33m",
	"blue":   "\033[34m",
	"reset":  "\033[0m",
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <banner> <text> <color>")
		fmt.Println("Or: go run main.go <text> (defaults to 'standard' banner and 'reset' color)")
		return
	}

	// Assume the first argument is the input text by default
	input := concatenateArgs(os.Args[1:])
	bannerName := "standard" // Default banner
	colorName := "reset"     // Default color

	// If the last argument is a valid banner name, use it
	if len(os.Args) >= 3 {
		lastArg := os.Args[len(os.Args)-2]
		if _, exists := banners[lastArg]; exists {
			bannerName = lastArg
			input = concatenateArgs(os.Args[1 : len(os.Args)-2]) // Exclude the banner name from input
		} else {
			colorName = lastArg // If not a banner, treat as color
		}
	}

	// If the last argument is a valid color, use it
	if len(os.Args) >= 3 && isValidColor(os.Args[len(os.Args)-1]) {
		colorName = os.Args[len(os.Args)-1]
	}

	banner, exists := banners[bannerName]
	if !exists {
		fmt.Printf("Error: Banner '%s' not found.\n", bannerName)
		return
	}

	color, exists := colors[colorName]
	if !exists {
		color = colors["reset"]
	}

	processedLines := handleNewlines(input)
	generateAsciiArt(processedLines, banner, color)
}

func concatenateArgs(args []string) string {
	return strings.Join(args, " ")
}

func handleNewlines(input string) []string {
	return strings.Split(input, "\\n")
}

func generateAsciiArt(lines []string, banner Banner, color string) {
	for _, line := range lines {
		if line == "" {
			fmt.Println()
			continue
		}
		processLine(line, banner, color)
	}
}

func processLine(line string, banner Banner, color string) {
	for i := 0; i < banner.lineHeight; i++ {
		res := ""
		for _, letter := range line {
			res += getLine(1+int(letter-' ')*banner.lineHeight+i, banner.filePath)
		}
		fmt.Print(color) // Set color
		fmt.Println(res)
		fmt.Print(colors["reset"]) // Reset color after each line
	}
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

func isValidColor(colorName string) bool {
	_, exists := colors[colorName]
	return exists
}
