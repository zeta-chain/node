package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultLoggers(t *testing.T) {
	defaultLoggers := DefaultLoggers()
	require.NotNil(t, defaultLoggers)
	require.NotNil(t, defaultLoggers.Std)
	require.NotNil(t, defaultLoggers.Compliance)
}
