package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultVerificationFlags(t *testing.T) {
	t.Run("default verification flags is all disabled", func(t *testing.T) {
		flags := DefaultVerificationFlags()
		require.False(t, flags.EthTypeChainEnabled)
		require.False(t, flags.BtcTypeChainEnabled)
	})
}
