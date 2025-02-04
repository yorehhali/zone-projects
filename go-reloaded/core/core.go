package core

import (
	"os"
	"strconv"
	"strings"
	"unicode"
)

func RemoveDuplicateSpaces(input string) string {
	var result string
	lastWasSpace := false
	for _, ch := range input {
		if ch == ' ' {
			if !lastWasSpace {
				result += string(ch)
				lastWasSpace = true
			}
		} else {
			result += string(ch)
			lastWasSpace = false
		}
	}
	return strings.TrimSpace(result)
}
func IsPunctuation(char string) bool {
	puncts := []string{".", ",", ":", ";", "!", "?"}
	for _, punct := range puncts {
		if punct == char {
			return true
		}
	}
	return false
}
func ApostropheCounter(text string) (int, bool) {
	counter := 0
	for _, char := range text {
		if char == '\'' {
			counter++
		}
	}
	return counter, counter%2 == 0
}
func WordPartsEditable(word string) bool {
	for _, char := range word {
		if IsEditable(char) {
			return true
		}
	}
	return false
}
func IsEmoji(r rune) bool {
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map Symbols
		(r >= 0x1F700 && r <= 0x1F77F) || // Alchemical Symbols
		(r >= 0x1F780 && r <= 0x1F7FF) || // Geometric Shapes Extended
		(r >= 0x1F800 && r <= 0x1F8FF) || // Supplemental Arrows-C
		(r >= 0x1F900 && r <= 0x1F9FF) || // Supplemental Symbols and Pictographs
		(r >= 0x1FA00 && r <= 0x1FA6F) || // Chess Symbols
		(r >= 0x1FA70 && r <= 0x1FAFF) // Symbols and Pictographs Extended-A
}
func IsEditable(char rune) bool {
	return !IsEmoji(char) && unicode.IsLetter(char) && (unicode.ToUpper(char) != char || unicode.ToLower(char) != char) || unicode.IsNumber(char)
}

var validFlags = []string{"hex", "bin", "up", "cap", "low"}

func SplitContent(text string) []string {
	word := ""
	res := []string{}
	for i, char := range text {
		if char == '\'' {
			if (i > 0 && text[i-1] == ' ') || (i < len(text)-1 && text[i+1] == ' ') || (i < len(text)-1 && text[i+1] == '\'') || (i > 0 && text[i-1] == '\'') {
				if word != "" {
					res = append(res, word)
					word = ""
				}
				res = append(res, "'") // Add the quote as a separate token
			} else {
				word += string(char)
			}
		} else if IsPunctuation(string(char)) {
			if word != "" {
				res = append(res, word)
				word = ""
			}
			res = append(res, string(char))
		} else if char == ' ' {
			if word != "" {
				res = append(res, word)
				word = ""
			}
		} else if char == '(' || char == ')' || char == '\n' {
			if word != "" {
				res = append(res, word)
				word = ""
			}
			if char == '(' {
				res = append(res, "(")
			} else if char == ')' {
				res = append(res, ")")
			} else {
				res = append(res, "\n")
			}
		} else {
			word += string(char)
		}
	}
	if word != "" {
		res = append(res, word)
	}
	return res
}

func AtoAN(words []string) []string {
	vowels := []rune{'a', 'e', 'i', 'o', 'u', 'h'}
	for i := 0; i < len(words); i++ {
		if words[i] == "a" {
			for _, vowel := range vowels {
				if i < len(words)-1 && strings.ToLower(string(words[i+1][0])) == string(vowel) {
					words[i] = "an"
				}
			}
		} else if words[i] == "A" {
			for _, vowel := range vowels {
				if i < len(words)-1 && strings.ToLower(string(words[i+1][0])) == string(vowel) {
					words[i] = "An"
				}
			}
		}
	}
	return words
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
	res := ""
	for _, char := range word {
		if IsEditable(char) {
			res += string(unicode.ToUpper(char))
		} else {
			res += string(char)
		}
	}
	return res
}

func Lower(word string) string {
	res := ""
	for _, char := range word {
		if IsEditable(char) {
			res += string(unicode.ToLower(char))
		} else {
			res += string(char)
		}
	}
	return res
}

func Capitalize(word string) string {
	if len(word) == 0 {
		return word
	}
	res := ""
	for index, char := range word {
		if IsEditable(char) && index == 0 {
			res += string(unicode.ToUpper(char))
		} else if IsEditable(char) && index > 0 {
			res += string(unicode.ToLower(char))
		} else {
			res += string(char)
		}
	}
	return res
}

func GetFlagFunction(FlagName string) func(string) string {
	switch FlagName {
	case "hex":
		return ConvertHex
	case "bin":
		return ConvertBin
	case "up":
		return Upper
	case "low":
		return Lower
	case "cap":
		return Capitalize
	default:
		return func(s string) string { return s }
	}
}

func IsShortFlag(word string) bool {
	return strings.HasPrefix(word, "(") && strings.HasSuffix(word, ")") && len(word) <= 5
}

func IsLongFlag(word string) bool {
	return strings.HasPrefix(word, "(") && strings.HasSuffix(word, ")") && len(word) > 5
}

func FixFlags(words []string) []string {
	NewResult := []string{}
	i := 0
	for i < len(words) {
		if i+2 < len(words) && words[i] == "(" && ValidateFlag(words[i+1]) && words[i+2] == ")" {
			NewResult = append(NewResult, words[i]+words[i+1]+words[i+2])
			i += 3
			continue
		} else if i+4 < len(words) && words[i] == "(" && ValidateFlag(words[i+1]) && words[i+2] == "," && IsNumber(words[i+3]) && words[i+4] == ")" {
			NewResult = append(NewResult, words[i]+words[i+1]+words[i+2]+words[i+3]+words[i+4])
			i += 5
			continue
		}
		NewResult = append(NewResult, words[i])
		i++
	}
	return NewResult
}

func ValidateFlag(FlagName string) bool {
	for _, validate := range validFlags {
		if FlagName == validate {
			return true
		}
	}
	return false
}

func IsNumber(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}

func ApplyFlags(words []string) []string {
	result := []string{}

	for i := 0; i < len(words); i++ {
		word := words[i]

		if IsShortFlag(word) {
			if len(result) == 0 {
				continue
			}
			newFlag := strings.Trim(word, "()")
			flagName := newFlag
			flagRepeat := 1
			flagRepeat = min(flagRepeat, len(result))
			count := flagRepeat
			for j := len(result) - 1; j >= 0 && count > 0; j-- {
				if WordPartsEditable(result[j]) {
					result[j] = GetFlagFunction(flagName)(result[j])
					count--
				}
			}
		} else if IsLongFlag(word) {
			newFlag := strings.Trim(word, "()")
			parts := strings.Split(newFlag, ",")
			if len(parts) != 2 {
				continue
			}
			flagName := parts[0]
			flagRepeat, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}
			flagRepeat = min(flagRepeat, len(result))
			count := flagRepeat
			for j := len(result) - 1; j >= 0 && count > 0; j-- {
				if WordPartsEditable(result[j]) {
					result[j] = GetFlagFunction(flagName)(result[j])
					count--
				}
			}
		} else {
			result = append(result, word)
		}
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func FormatOutput(text string, words []string) string {
	res := ""
	openingQuote := true
	counter, isEven := ApostropheCounter(text)
	for i, word := range words {
		nextWord := ""
		if i < len(words)-1 {
			nextWord = words[i+1]
		}
		if word == "'" {
			if !isEven && counter == 1 {
				res += " " + "'" + " "
				continue
			}
			counter--
			if openingQuote {
				if i > 0 && !strings.HasSuffix(res, " ") {
					res += " "

				}

			} else {
				if nextWord != "" && !IsPunctuation(nextWord) {
					res += "'" + " "
					openingQuote = !openingQuote
					continue
				}
			}
			res += "'"
			openingQuote = !openingQuote
		} else if i < len(words)-1 && (word == "\n" || nextWord == "\n" ||
			IsPunctuation(nextWord) || nextWord == ")" || nextWord == "'") {
			res += word
		} else {
			if word == "(" {
				res += word
				continue
			}
			res += word + " "
		}
	}
	return RemoveDuplicateSpaces(res)
}

func ReadFile(name string) (string, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func WriteFile(name, content string) error {
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		return err
	}
	return nil
}
