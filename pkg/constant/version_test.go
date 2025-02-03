package constant

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/mod/semver"
)

func TestVersionWithoutMake(t *testing.T) {
	require.True(t, semver.IsValid(versionWhenBuiltWithoutMake))
}
