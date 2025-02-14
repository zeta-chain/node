package constant

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"
)

func TestVersionWithoutMake(t *testing.T) {
	require.True(t, semver.IsValid(versionWhenBuiltWithoutMake))
}

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "version without v prefix",
			input:    "5.0.0",
			expected: "v5.0.0",
		},
		{
			name:     "version with v prefix",
			input:    "v5.0.0",
			expected: "v5.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, NormalizeVersion(tt.input))
		})
	}
}
