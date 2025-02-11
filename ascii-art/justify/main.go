package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetTerminalWidth() int {
	if os.Getenv("TERM") == "" {
		fmt.Println("Not running in a terminal")
		return 80
	}
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running stty command:", err)
		return 80
	}
	fields := strings.Fields(string(out))
	if len(fields) == 2 {
		width, err := strconv.Atoi(fields[1])
		if err == nil {
			return width
		}
	}
	return 80
}




func JustifyText(text, mode string, width int) string {
	lines := strings.Split(text, "\n")
	var result []string

	for _, line := range lines {
		lineLen := len(line)

		switch mode {
		case "right":
			padding := width - lineLen
			if padding > 0 {
				line = strings.Repeat(" ", padding) + line
			}
		case "center":
			padding := (width - lineLen) / 2
			if padding > 0 {
				line = strings.Repeat(" ", padding) + line
			}
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func main() {
	justify := flag.String("justify", "left", "Justify text: left, right, or center")
	flag.Parse()

	text := "This is a sample text that will be justified in the terminal output."
	width := GetTerminalWidth()

	fmt.Println(JustifyText(text, *justify, width))
}
