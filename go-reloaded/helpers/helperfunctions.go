package helpers

import (
	"os"
	"strconv"
	"strings"
	"unicode"
)

func IsEditable(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsPunct(char) && !unicode.IsNumber(char){
			return false
		}
	}
	return true
}
func HasNumbers(s string) bool {
	for _, char := range s {
		if unicode.IsNumber(char){
			return true
		}
	}
	return false
}
func ConvertHex(base string) string {
	val, err := strconv.ParseInt(base, 16, 64)
	if err != nil {
		return base
	}
	return strconv.Itoa(int(val))
}

func ConvertBin(base string) string {
	val, err := strconv.ParseInt(base, 2, 64)
	if err != nil {
		return base
	}
	return strconv.Itoa(int(val))
}

func Upper(word string) string {
	return strings.ToUpper(word)
}

func Lower(word string) string {
	return strings.ToLower(word)
}

func Capitalize(word string) string {
	res := ""
	for index, char := range word {
		if index == 0 {
			res += Upper(string(char))
		} else {
			res += Lower(string(char))
		}
	}
	return res
}

func RemoveWhitespace(input string) string {
	return strings.ReplaceAll(input, " ", "")
}

func ReadFromFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func WriteToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func IsPunctuation(char string) bool {
	punctuations := []string{".", ",", ":", ";", "?", "!", "'"}
	for _, punctuation := range punctuations {
		if punctuation == char {
			return true
		}
	}
	return false
}
func IsApostrophe(char string) bool {
	return char == "'"
}
func IsWordPart(word string) bool {
	return Lower(word) == "s" || Lower(word) == "t" || Lower(word) == "m" || Lower(word) == "d" || Lower(word) == "re" || Lower(word) == "ve" || Lower(word) == "ll"
}
func IsFlag(word string) bool {
	prefixes := []string{"(cap", "(up", "(low", "(hex", "(bin"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(word), prefix) {
			return true
		}
	}
	return false
}
func IsApostropheCounter(text interface{}) (int, bool) {
	counter := 0

	switch v := text.(type) {
	case string:
		for _, char := range v {
			if IsApostrophe(string(char)) {
				counter++
			}
		}
	case []string:
		for _, str := range v {
			for _, char := range str {
				if IsApostrophe(string(char)) {
					counter++
				}
			}
		}
	default:
		panic("unsupported input type")
	}

	return counter, counter%2 == 0
}
func RemoveDuplicateSpaces(input string) string {
	var result string
	lastCharWasSpace := false

	for _, char := range input {
		if char == ' ' {
			if !lastCharWasSpace {
				result += string(char)
				lastCharWasSpace = true
			}
		} else {
			result += string(char)
			lastCharWasSpace = false
		}
	}

	return strings.Trim(result, " ")
}
