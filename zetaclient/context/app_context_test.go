package context_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	context "github.com/zeta-chain/zetacore/zetaclient/context"
)

func Test_GetBTCChainID(t *testing.T) {
	// test chains and params
	evmChain := chains.Ethereum
	btcChain := chains.BitcoinMainnet
	evmParams := &observertypes.ChainParams{
		ChainId: evmChain.ChainId,
	}
	btcParams := &observertypes.ChainParams{
		ChainId: btcChain.ChainId,
	}

	t.Run("GetBTCChainID returns regnet chain id if btc chain is not enabled in config", func(t *testing.T) {
		// config without btc chain
		cfg := config.NewConfig()

		// create app context without BTC chain
		coreContext := getTestCoreContext(evmChain, evmParams, nil, nil, nil)
		appContext := context.NewAppContext(coreContext, cfg)

		btcChainID := appContext.GetBTCChainID()
		require.Equal(t, chains.BitcoinRegtest.ChainId, btcChainID)
	})
	t.Run("GetBTCChainID returns regnet chain id if btc chain params are not enabled", func(t *testing.T) {
		// config with btc chain
		cfg := config.NewConfig()
		cfg.BitcoinConfig = config.BTCConfig{
			RPCUsername: "user",
		}

		// create app context without BTC chain params
		coreContext := getTestCoreContext(evmChain, evmParams, nil, nil, nil)
		appContext := context.NewAppContext(coreContext, cfg)

		btcChainID := appContext.GetBTCChainID()
		require.Equal(t, chains.BitcoinRegtest.ChainId, btcChainID)
	})
	t.Run("GetBTCChainID returns btc chain id if enabled", func(t *testing.T) {
		// config with btc chain
		cfg := config.NewConfig()
		cfg.BitcoinConfig = config.BTCConfig{
			RPCUsername: "user",
		}

		// create app context with BTC chain
		coreContext := getTestCoreContext(evmChain, evmParams, btcParams, nil, nil)
		appContext := context.NewAppContext(coreContext, cfg)

		btcChainID := appContext.GetBTCChainID()
		require.Equal(t, btcChain.ChainId, btcChainID)
	})
}
