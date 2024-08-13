package os

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptPasswords prompts the user for passwords with the given titles
func PromptPasswords(passwordTitles []string) ([]string, error) {
	reader := bufio.NewReader(os.Stdin)
	passwords := make([]string, len(passwordTitles))

	// iterate over password titles and prompt for each
	for i, title := range passwordTitles {
		fmt.Printf("%s Password: ", title)
		password, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		// trim delimiters
		password = strings.TrimSpace(password)
		passwords[i] = password
	}

	return passwords, nil
}
