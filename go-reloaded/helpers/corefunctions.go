package helpers

import (
	"strconv"
	"strings"
)

func FormatInput(text string) []string {
	apostropheCount, isEven := IsApostropheCounter(text)
	word := ""
	result := []string{}
	withinBrackets := false
	for _, char := range text {
		if !isEven && apostropheCount == 1 && IsApostrophe(string(char)) {
			word += string(char)
			continue
		} else if IsApostrophe(string(char)) &&  !withinBrackets {
			apostropheCount--
			if len(word) != 0 {
				result = append(result, word)
				word = ""
			}
			result = append(result, string(char))
		} else if IsPunctuation(string(char)) && !withinBrackets && !IsApostrophe(string(char)) {
			if len(word) != 0 {
				result = append(result, word)
				word = ""
			}
			result = append(result, string(char))
		} else if char == ' ' && !withinBrackets {
			if len(word) != 0 {
				result = append(result, word)
				word = ""
			}
		} else if char == '(' {
			if len(word) != 0 {
				result = append(result, word)
				word = ""
			}
			word += string(char)
			withinBrackets = true
		} else if char == ')' {
			word += string(char)
			withinBrackets = false
			if len(word) != 0 {
				result = append(result, word)
				word = ""
			}
		} else if char == '\n' {
			if len(word) != 0 {
				result = append(result, word)
				word = ""
			}
			result = append(result, "\n")
		} else {
			word += string(char)
		}
	}
	if len(word) != 0 {
		result = append(result, word)
	}

	return result
}

func ProcessWords(words []string) []string {
	result := []string{}
	for i := 0; i < len(words); i++ {
		word := words[i]
		if IsFlag(word) {
			word = RemoveWhitespace(word)
			switch {
			case strings.HasPrefix(strings.ToLower(word), "(hex)"):
				if i > 0 && IsEditable(result[len(result)-1]) {
					result[len(result)-1] = ConvertHex(result[len(result)-1])
				}
			case strings.HasPrefix(strings.ToLower(word), "(bin)"):
				if i > 0 && IsEditable(result[len(result)-1]) {
					result[len(result)-1] = ConvertBin(result[len(result)-1])
				}
			case strings.HasPrefix(strings.ToLower(word), "(up"):
				count := 1
				if strings.Contains(word, ",") {
					parts := strings.Split(word, ",")
					if len(parts) > 1 {
						count, _ = strconv.Atoi(strings.TrimSuffix(parts[1], ")"))
					}
				}
				for j := len(result) - 1; j >= 0 && count > 0; j-- {
					if IsEditable(result[j]) && result[j] != "\n" && !IsPunctuation(result[j]) && !HasNumbers(result[j]) {
						result[j] = Upper(result[j])
						count--
					} 
				}
			case strings.HasPrefix(strings.ToLower(word), "(low"):
				count := 1
				if strings.Contains(word, ",") {
					parts := strings.Split(word, ",")
					if len(parts) > 1 {
						count, _ = strconv.Atoi(strings.TrimSuffix(parts[1], ")"))
					}
				}
				for j := len(result) - 1; j >= 0 && count > 0; j-- {
					if IsEditable(result[j]) && result[j] != "\n" && !IsPunctuation(result[j])  && !HasNumbers(result[j]) {
						result[j] = Lower(result[j])
						count--
					} else if IsWordPart(result[j]) {
						result[j] = Lower(result[j])
					}
				}
			case strings.HasPrefix(strings.ToLower(word), "(cap"):
				count := 1
				if strings.Contains(word, ",") {
					parts := strings.Split(word, ",")
					if len(parts) > 1 {
						count, _ = strconv.Atoi(strings.TrimSuffix(parts[1], ")"))
					}
				}
				for j := len(result) - 1; j >= 0 && count > 0; j-- {
					if  IsEditable(result[j]) && result[j] != "\n" && !IsPunctuation(result[j])  && !HasNumbers(result[j]) {
						result[j] = Capitalize(result[j])
						count--
					} else if j>0 && IsWordPart(result[j]) && IsApostrophe(result[j-1])  && !HasNumbers(result[j])  {
						result[j] = Lower(result[j])
					}
				}
			}
		} else {
			result = append(result, word)
		}
	}

	return result
}

func AtoAN(words []string) []string {
	vowels := []rune{'a', 'e', 'i', 'o', 'u', 'h'}
	for i := 0; i < len(words); i++ {
		if strings.ToLower(words[i]) == "a" {
			for _, vowel := range vowels {
				if i < len(words)-1 && strings.ToLower(string(words[i+1][0])) == string(vowel) {
					words[i] = "an"
				}
			}
		}
	}
	return words
}

func FormatOutput(words []string) string {
	apostropheCount, isEven := IsApostropheCounter(words)
	res := ""
	isClosingQuote := false
	for i := 0; i < len(words); i++ {
		word := words[i]

		if word == "\n" {
			res += "\n"
			continue
		}

		if IsPunctuation(word) {
			if apostropheCount == 1 && IsApostrophe(word) && !isEven {
				res += word
				continue
			} else if IsApostrophe(word) && i < len(words)-1 && apostropheCount > 1 {
				apostropheCount--

				nextWord := words[i+1]
				if IsWordPart(nextWord) {
					res += word
					isClosingQuote = !isClosingQuote
					continue
				}
				if isClosingQuote {
					res += word + " "
					isClosingQuote = !isClosingQuote
				} else {
					res += " " + word
					isClosingQuote = !isClosingQuote
				}

				continue
			} else if i < len(words)-1 && !IsPunctuation(words[i+1]) {
				res += word + " "
			} else {
				res += word
			}
			continue
		}

		res += word
		if i < len(words)-1 && words[i+1] != "\n" && !IsPunctuation(words[i+1]) {
			res += " "
		}
	}

	return res
}
