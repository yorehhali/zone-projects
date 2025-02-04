package main

import (
	"fmt"
	"os"
)

// Function to create a test case file
func createTestCase(fileName, content string) error {
	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", fileName, err)
	}
	defer file.Close()

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("could not write to file %s: %w", fileName, err)
	}

	return nil
}

func main() {
	// Define the test cases with their respective file names and content

	testCases := []struct {
		fileName string
		content  string
	}{
		{
			fileName: "sample1.txt",
			content:  "If I make you BREAKFAST IN BED (low, 3) just say thank you instead of: how (cap) did you get in my house (up, 2) ?",
		},
		{
			fileName: "sample2.txt",
			content:  "I have to pack 101 (bin) outfits. Packed 1a (hex) just to be sure.",
		},
		{
			fileName: "sample3.txt",
			content:  "Don not be sad ,because sad backwards is das . And das not good.",
		},
		{
			fileName: "sample4.txt",
			content:  "harold wilson (cap, 2) : ' I am a optimist ,but a optimist who carries a raincoat . '",
		},
	}

	// Create the test case files
	for _, testCase := range testCases {
		err := createTestCase(testCase.fileName, testCase.content)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", testCase.fileName, err)
		} else {
			fmt.Printf("Test case file %s created successfully\n", testCase.fileName)
		}
	}
}
