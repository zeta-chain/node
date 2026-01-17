package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
)

func TestGetConfig(t *testing.T) {
	t.Run("Get default config if not specified", func(t *testing.T) {
		cfg, err := config.GetConfig(chains.Ethereum, "")
		require.NoError(t, err)
		require.Equal(t, "https://zetachain-mainnet.g.allthatnode.com:443/archive/tendermint", cfg.ZetaChainRPC)

		cfg, err = config.GetConfig(chains.Sepolia, "")
		require.NoError(t, err)
		require.Equal(t, "https://zetachain-athens.g.allthatnode.com/archive/tendermint", cfg.ZetaChainRPC)

		cfg, err = config.GetConfig(chains.GoerliLocalnet, "")
		require.NoError(t, err)
		require.Equal(t, "http://127.0.0.1:26657", cfg.ZetaChainRPC)
	})
}
