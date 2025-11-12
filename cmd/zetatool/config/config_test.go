package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/cmd/zetatool/config"
	"github.com/zeta-chain/node/pkg/chains"
)

func TestRead(t *testing.T) {
	t.Run("TestRead", func(t *testing.T) {
		c := config.Config{}
		err := c.Read("sample_config.json")
		require.NoError(t, err)

		require.Equal(t, "https://zetachain-testnet-grpc.itrocket.net:443", c.ZetaChainRPC)
		require.Equal(t, "https://ethereum-sepolia-rpc.publicnode.com", c.EthereumRPC)
		require.Equal(t, int64(101), c.ZetaChainID)
		require.Equal(t, "", c.BtcUser)
		require.Equal(t, "", c.BtcPassword)
		require.Equal(t, "", c.BtcHost)
		require.Equal(t, "", c.BtcParams)
		require.Equal(t, "", c.SolanaRPC)
		require.Equal(t, "https://bsc-testnet-rpc.publicnode.com", c.BscRPC)
		require.Equal(t, "https://polygon-amoy.gateway.tenderly.com", c.PolygonRPC)
		require.Equal(t, "https://base-sepolia-rpc.publicnode.com", c.BaseRPC)
		require.Equal(t, "https://sepolia-rollup.arbitrum.io/rpc", c.ArbitrumRPC)
		require.Equal(t, "https://sepolia.optimism.io", c.OptimismRPC)
		require.Equal(t, "https://avalanche-fuji-c-chain-rpc.publicnode.com", c.AvalancheRPC)
		require.Equal(t, "", c.WorldRPC)
	})
}

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

	t.Run("Get config from file if specified", func(t *testing.T) {
		cfg, err := config.GetConfig(chains.Ethereum, "sample_config.json")
		require.NoError(t, err)
		require.Equal(t, "https://zetachain-testnet-grpc.itrocket.net:443", cfg.ZetaChainRPC)

		cfg, err = config.GetConfig(chains.Sepolia, "sample_config.json")
		require.NoError(t, err)
		require.Equal(t, "https://zetachain-testnet-grpc.itrocket.net:443", cfg.ZetaChainRPC)
	})
}

func TestWorldRPCConfig(t *testing.T) {
	t.Run("TestnetConfig has World RPC", func(t *testing.T) {
		cfg := config.TestnetConfig()
		require.Equal(t, "https://worldchain-sepolia.g.alchemy.com/public", cfg.WorldRPC)
	})

	t.Run("MainnetConfig has World RPC", func(t *testing.T) {
		cfg := config.MainnetConfig()
		require.Equal(t, "https://worldchain-mainnet.g.alchemy.com/public", cfg.WorldRPC)
	})

	t.Run("DevnetConfig has empty World RPC", func(t *testing.T) {
		cfg := config.DevnetConfig()
		require.Equal(t, "", cfg.WorldRPC)
	})

	t.Run("PrivateNetConfig has empty World RPC", func(t *testing.T) {
		cfg := config.PrivateNetConfig()
		require.Equal(t, "", cfg.WorldRPC)
	})
}
