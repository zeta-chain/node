package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultLoggers(t *testing.T) {
	defaultLoggers := DefaultLoggers()
	require.NotNil(t, defaultLoggers)
	require.NotNil(t, defaultLoggers.Std)
	require.NotNil(t, defaultLoggers.Compliance)
}
