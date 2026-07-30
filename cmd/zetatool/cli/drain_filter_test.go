package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewChainFilter(t *testing.T) {
	t.Run("only-chains allows just the listed chains", func(t *testing.T) {
		f, err := newChainFilter("1,56", "")
		require.NoError(t, err)
		require.True(t, f.allow(1))
		require.True(t, f.allow(56))
		require.False(t, f.allow(137))
	})

	t.Run("exclude-chains drops the listed chains", func(t *testing.T) {
		f, err := newChainFilter("", "56, 137")
		require.NoError(t, err)
		require.True(t, f.allow(1))
		require.False(t, f.allow(56))
		require.False(t, f.allow(137))
	})

	t.Run("empty allows all", func(t *testing.T) {
		f, err := newChainFilter("", "")
		require.NoError(t, err)
		require.True(t, f.allow(1))
		require.True(t, f.allow(8332))
	})

	t.Run("both set is an error", func(t *testing.T) {
		_, err := newChainFilter("1", "56")
		require.Error(t, err)
	})

	t.Run("invalid chain id is an error", func(t *testing.T) {
		_, err := newChainFilter("1,abc", "")
		require.Error(t, err)
	})
}
