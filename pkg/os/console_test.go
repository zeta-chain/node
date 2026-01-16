package os_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	zetaos "github.com/zeta-chain/node/pkg/os"
)

func Test_PromptPassword(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{
			name:   "Valid password",
			input:  " pass123\n",
			output: "pass123",
		},
		{
			name:   "Empty password",
			input:  "\n",
			output: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe to simulate stdin
			r, w, err := os.Pipe()
			require.NoError(t, err)

			// Write the test input to the pipe
			_, err = w.Write([]byte(tt.input))
			require.NoError(t, err)
			w.Close() // Close the write end of the pipe

			// Backup the original stdin and restore it after the test
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			// Redirect stdin to the read end of the pipe
			os.Stdin = r

			// Call the function with the test case data
			password, err := zetaos.PromptPassword("anyTitle")

			// Check the returned passwords
			require.NoError(t, err)
			require.Equal(t, tt.output, password)
		})
	}
}

// Test function for PromptPasswords
func Test_PromptPasswords(t *testing.T) {
	tests := []struct {
		name           string
		passwordTitles []string
		input          string
		expected       []string
	}{
		{
			name:           "Single password prompt",
			passwordTitles: []string{"HotKey"},
			input:          " pass123\n",
			expected:       []string{"pass123"},
		},
		{
			name:           "Multiple password prompts",
			passwordTitles: []string{"HotKey", "TSS", "RelayerKey"},
			input:          "pass_hotkey\npass_tss\npass_relayer\n",
			expected:       []string{"pass_hotkey", "pass_tss", "pass_relayer"},
		},
		{
			name:           "Empty input for passwords is allowed",
			passwordTitles: []string{"HotKey", "TSS", "RelayerKey"},
			input:          "\n\n\n",
			expected:       []string{"", "", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe to simulate stdin
			r, w, err := os.Pipe()
			require.NoError(t, err)

			// Write the test input to the pipe
			_, err = w.Write([]byte(tt.input))
			require.NoError(t, err)
			w.Close() // Close the write end of the pipe

			// Backup the original stdin and restore it after the test
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()

			// Redirect stdin to the read end of the pipe
			os.Stdin = r

			// Call the function with the test case data
			passwords, err := zetaos.PromptPasswords(tt.passwordTitles)

			// Check the returned passwords
			require.NoError(t, err)
			require.Equal(t, tt.expected, passwords)
		})
	}
}
