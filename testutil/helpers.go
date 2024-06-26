package testutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/stretchr/testify/assert"
)

const helpersFile = "testutil/helpers.go"

// Must a helper that terminates the program if the error is not nil.
func Must[T any](v T, err error) T {
	NoError(err)
	return v
}

// NoError terminates the program if the error is not nil.
func NoError(err error) {
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Printf("Unable to continue execution: %s.\nStacktrace:\n", err)

	for _, line := range assert.CallerInfo() {
		if strings.Contains(line, helpersFile) {
			continue
		}

		fmt.Println("  ", line)
	}

	os.Exit(1)
}
