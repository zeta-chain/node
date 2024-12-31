package os

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptPassword prompts the user for a password with the given title
func PromptPassword(passwordTitle string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	return readPassword(reader, passwordTitle)
}

// PromptPasswords is a convenience function that prompts the user for multiple passwords
func PromptPasswords(passwordTitles []string) ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	passwords := make([]string, len(passwordTitles))

	// iterate over password titles and prompt for each
	for i, title := range passwordTitles {
		password, err := readPassword(reader, title)
		if err != nil {
			return nil, err
		}
		passwords[i] = password
	}

	return passwords, nil
}

// readPassword is a helper function that reads a password from bufio.Reader
func readPassword(reader *bufio.Reader, passwordTitle string) (string, error) {
	const delimiter = '\n'

	// prompt for password
	fmt.Printf("%s Password: ", passwordTitle)
	password, err := reader.ReadString(delimiter)
	if err != nil {
		return "", err
	}

	// trim leading and trailing spaces
	return strings.TrimSpace(password), nil
}
