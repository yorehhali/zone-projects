package main

import (
	"fmt"
	"os"
	"strings"
	"goreloaded/core" 
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . input.txt result.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	if !strings.HasSuffix(inputFile, ".txt") {
		fmt.Println("Invalid file extension:", inputFile," Please provide a .txt input file.")
		return
	}
	if !strings.HasSuffix(outputFile, ".txt") {
		fmt.Println("Invalid file extension:", outputFile," Please provide a .txt output file.")
		return
	}

	content, err := core.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		return
	}

	words := core.SplitContent(content)
	words = core.AtoAN(words)
	fixedWords := core.FixFlags(words)
	processedWords := core.ApplyFlags(fixedWords)
	fixedWords = core.FixFlags(processedWords)
	processedWords = core.ApplyFlags(fixedWords)
	
	outputContent := core.FormatOutput(content,processedWords)

	err = core.WriteFile(outputFile, outputContent)
	if err != nil {
		fmt.Println("Error writing to the file:", err)
		return
	}

	fmt.Println("File processing complete. Output written to", outputFile)
}
