package utils

import "os"

func generate(characters []string, length int, current string, result *[]string) {
	// Base case: if the current word has reached the desired length
	if len(current) == length {
		*result = append(*result, current)
		return
	}

	// Recursively append characters to form the wordlist
	for _, char := range characters {
		generate(characters, length, current+char, result)
	}
}

func GenerateWordlist(filename string) error {
	characters := "abcdef0123456789"

	// Convert the string to []string
	var strArray []string
	for _, char := range characters {
		strArray = append(strArray, string(char))
	}
	var wordlist []string

	generate(strArray, 5, "", &wordlist)

	// Create or open the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write each word to the file
	for _, word := range wordlist {
		_, err := file.WriteString(word + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
