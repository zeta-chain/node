package cli

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestParseEVMMaxAmount(t *testing.T) {
	t.Run("empty means no cap", func(t *testing.T) {
		v, err := parseEVMMaxAmount("")
		require.NoError(t, err)
		require.True(t, v.IsZero())
	})

	t.Run("whitespace is trimmed", func(t *testing.T) {
		v, err := parseEVMMaxAmount("  1000  ")
		require.NoError(t, err)
		require.Equal(t, sdkmath.NewUint(1000), v)
	})

	t.Run("parses a wei amount beyond int64", func(t *testing.T) {
		v, err := parseEVMMaxAmount("100000000000000000000")
		require.NoError(t, err)
		require.Equal(t, "100000000000000000000", v.String())
	})

	t.Run("an explicit zero is rejected", func(t *testing.T) {
		// a zero cap would pin a transfer that pays gas to move nothing; the operator almost
		// certainly meant to omit the flag
		_, err := parseEVMMaxAmount("0")
		require.ErrorContains(t, err, "must be positive")
	})

	t.Run("rejects non-decimal and negative values", func(t *testing.T) {
		for _, s := range []string{"0x64", "1e18", "abc", "-1", "1.5"} {
			_, err := parseEVMMaxAmount(s)
			require.ErrorContains(t, err, "invalid --evm-max-amount", "input %q", s)
		}
	})
}
