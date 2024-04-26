package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

func TestDefaultVerificationFlags(t *testing.T) {
	t.Run("default verification flags is all disabled", func(t *testing.T) {
		flags := DefaultVerificationFlags()
		for _, f := range flags {
			switch f.ChainId {
			case chains.EthChain.ChainId:
				require.False(t, f.Enabled)
			case chains.BscMainnetChain.ChainId:
				require.False(t, f.Enabled)
			case chains.SepoliaChain.ChainId:
				require.False(t, f.Enabled)
			case chains.BscTestnetChain.ChainId:
				require.False(t, f.Enabled)
			case chains.GoerliLocalnetChain.ChainId:
				require.False(t, f.Enabled)
			case chains.GoerliChain.ChainId:
				require.False(t, f.Enabled)
			default:
				require.False(t, f.Enabled, "unexpected chain id")
			}
		}
	})
}
