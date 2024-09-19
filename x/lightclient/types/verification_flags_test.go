package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
)

func TestDefaultVerificationFlags(t *testing.T) {
	t.Run("default verification flags is all disabled", func(t *testing.T) {
		flags := DefaultHeaderSupportedChains()
		for _, f := range flags {
			switch f.ChainId {
			case chains.Ethereum.ChainId:
				require.False(t, f.Enabled)
			case chains.BscMainnet.ChainId:
				require.False(t, f.Enabled)
			case chains.Sepolia.ChainId:
				require.False(t, f.Enabled)
			case chains.BscTestnet.ChainId:
				require.False(t, f.Enabled)
			case chains.GoerliLocalnet.ChainId:
				require.False(t, f.Enabled)
			case chains.Goerli.ChainId:
				require.False(t, f.Enabled)
			default:
				require.False(t, f.Enabled, "unexpected chain id")
			}
		}
	})
}
