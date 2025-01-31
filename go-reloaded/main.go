package main

import (
	"fmt"
	"os"
	"goreloaded/helpers"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run . input_file.txt output_file.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]
	if !strings.HasSuffix(inputFile, ".txt") || !strings.HasSuffix(outputFile, ".txt") {
		fmt.Println("Usage: go run . input_file.txt output_file.txt")
		return
	}

	content, err := helpers.ReadFromFile(inputFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	words := helpers.FormatInput(content)
	words = helpers.AtoAN(words)
	words = helpers.ProcessWords(words)
	
	formattedText := helpers.RemoveDuplicateSpaces(helpers.FormatOutput(words))

	err = helpers.WriteToFile(outputFile, formattedText)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
	fmt.Println("Output written to", outputFile)
}
