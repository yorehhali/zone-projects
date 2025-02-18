package ascii

import (
	"fmt"
	"os"
	"strings"
)

type Banner struct {
	filePath   string
	lineHeight int
}

var Banners = map[string]Banner{
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

func GenAscii(input, bannerName string) string{
	banner, exists := Banners[bannerName]
	if !exists {
		fmt.Printf("Error: Banner '%s' not found.\n", bannerName)
		return "Error: Banner not found: " + bannerName
	}
	processedLines := HandleNewlines(input)
	asciiArt := GenerateAsciiArt(processedLines, banner)
	return asciiArt
}
func HandleNewlines(input string) []string {
	return strings.Split(input, "\n")
}

func GenerateAsciiArt(lines []string, banner Banner) string {
	result := ""
	for _, line := range lines {
		if line == "" {
			result += "\n"
			continue
		}
		result += ProcessLine(line, banner)
	}
	return result
}

func ProcessLine(line string, banner Banner) string {
	result := ""
	for i := 1; i <= banner.lineHeight; i++ {
		res := ""
		for _, letter := range line {
			if letter<32 || letter>127 || letter == '\n'{
				continue
			}
			res += GetLine(1+int(letter-32)*(banner.lineHeight+1)+i, banner.filePath)
		}
		result += res + "\n"
	}
	return result
}

func GetLine(num int, filePath string) string {
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
