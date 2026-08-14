package cli

import (
	"bytes"
	"math/big"
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

	t.Run("a value above 256 bits is rejected, not panicked on", func(t *testing.T) {
		// sdkmath.Uint's constructor panics past 256 bits, so an extra-zeros typo would take
		// down zetatool mid-drain instead of reporting a bad flag
		maxUint256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

		v, err := parseEVMMaxAmount(maxUint256.String())
		require.NoError(t, err)
		require.Equal(t, maxUint256.String(), v.String())

		require.NotPanics(t, func() {
			_, err = parseEVMMaxAmount(new(big.Int).Add(maxUint256, big.NewInt(1)).String())
		})
		require.ErrorContains(t, err, "exceeds the maximum 256-bit")
	})

	t.Run("rejects non-decimal and negative values", func(t *testing.T) {
		for _, s := range []string{"0x64", "1e18", "abc", "-1", "1.5"} {
			_, err := parseEVMMaxAmount(s)
			require.ErrorContains(t, err, "invalid --evm-max-amount", "input %q", s)
		}
	})
}

func TestWarnIfCapped(t *testing.T) {
	t.Run("silent when nothing is capped", func(t *testing.T) {
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.ZeroUint(), 0)
		require.Empty(t, buf.String())
	})

	t.Run("an unset cap reads as uncapped, not as zero", func(t *testing.T) {
		// "btc 0 sats total" would read as "sends no BTC", which is the opposite of the truth:
		// BTC is swept in full when only the EVM cap is set
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.NewUint(1_000), 0)

		out := buf.String()
		require.Contains(t, out, "REHEARSAL PAYLOAD")
		require.Contains(t, out, "evm: 1000 wei per chain")
		require.Contains(t, out, "btc: uncapped")
	})

	t.Run("btc-only cap", func(t *testing.T) {
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.ZeroUint(), 100_000)

		out := buf.String()
		require.Contains(t, out, "evm: uncapped")
		require.Contains(t, out, "btc: 100000 sats total")
	})

	t.Run("does not claim nothing is drained", func(t *testing.T) {
		// a cap is an upper bound: a chain holding less than it is swept in full
		var buf bytes.Buffer
		warnIfCapped(&buf, sdkmath.NewUint(1_000), 100_000)
		require.Contains(t, buf.String(), "except on chains already below the cap")
	})
}
