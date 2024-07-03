package context_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	context "github.com/zeta-chain/zetacore/zetaclient/context"
)

func Test_NewAppContext(t *testing.T) {
	t.Run("should create new app context with empty config", func(t *testing.T) {
		testCfg := config.NewConfig()

		appContext := context.New(testCfg)
		require.NotNil(t, appContext)

		// assert config
		require.Equal(t, testCfg, appContext.Config())

		// assert keygen
		keyGen := appContext.GetKeygen()
		require.Equal(t, observertypes.Keygen{}, keyGen)

		// assert enabled external chains
		require.Empty(t, appContext.GetEnabledExternalChains())

		// assert external chain params
		require.Empty(t, appContext.GetEnabledExternalChainParams())

		// assert current tss pubkey
		require.Equal(t, "", appContext.GetCurrentTssPubkey())

		// assert crosschain flags
		require.Equal(t, observertypes.CrosschainFlags{}, appContext.GetCrossChainFlags())

		// assert additional chains
		require.Empty(t, appContext.GetAdditionalChains())
	})
}

func Test_SetGetConfig(t *testing.T) {
	t.Run("should create new app context with empty config", func(t *testing.T) {
		oldCfg := config.NewConfig()
		appContext := context.New(oldCfg)
		require.NotNil(t, appContext)
		require.Equal(t, oldCfg, appContext.Config())

		// set new config
		evmChain := chains.Ethereum
		newCfg := config.NewConfig()
		newCfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
			Chain: evmChain,
		}
		appContext.SetConfig(newCfg)
		require.Equal(t, newCfg, appContext.Config())
	})
}

func Test_UpdateAndGetters(t *testing.T) {
	// use evm and btc chains for testing
	evmChain := chains.Ethereum
	btcChain := chains.BitcoinMainnet

	// create sample parameters
	keyGen := sample.Keygen(t)
	chainsEnabled := []chains.Chain{evmChain, btcChain}
	chainParamMap := map[int64]*observertypes.ChainParams{
		evmChain.ChainId: sample.ChainParams(evmChain.ChainId),
		btcChain.ChainId: sample.ChainParams(btcChain.ChainId),
	}
	btcNetParams := &chaincfg.MainNetParams
	tssPubKey := "tsspubkeytest"
	ccFlags := *sample.CrosschainFlags()
	additionalChains := []chains.Chain{
		sample.Chain(1),
		sample.Chain(2),
		sample.Chain(3),
	}
	headerSupportedChains := sample.HeaderSupportedChains()

	// feed app context fields
	appContext := context.New(config.NewConfig())
	appContext.Update(
		*keyGen,
		tssPubKey,
		chainsEnabled,
		chainParamMap,
		btcNetParams,
		ccFlags,
		additionalChains,
		headerSupportedChains,
		log.Logger,
	)

	t.Run("should get keygen", func(t *testing.T) {
		result := appContext.GetKeygen()
		require.Equal(t, *keyGen, result)
	})
	t.Run("should get current tss pubkey", func(t *testing.T) {
		result := appContext.GetCurrentTssPubkey()
		require.Equal(t, tssPubKey, result)
	})
	t.Run("should get external enabled chains", func(t *testing.T) {
		result := appContext.GetEnabledExternalChains()
		require.Equal(t, chainsEnabled, result)
	})
	t.Run("should get enabled BTC chains", func(t *testing.T) {
		result := appContext.GetEnabledBTCChains()
		require.Equal(t, []chains.Chain{btcChain}, result)
	})
	t.Run("should get enabled external chain params", func(t *testing.T) {
		result := appContext.GetEnabledExternalChainParams()
		require.Equal(t, chainParamMap, result)
	})
	t.Run("should get external chain params by chain id", func(t *testing.T) {
		for _, chain := range chainsEnabled {
			result, found := appContext.GetExternalChainParams(chain.ChainId)
			require.True(t, found)
			require.Equal(t, chainParamMap[chain.ChainId], result)
		}
	})
	t.Run("should get btc network params", func(t *testing.T) {
		result := appContext.GetBTCNetParams()
		require.Equal(t, btcNetParams, result)
	})
	t.Run("should get crosschain flags", func(t *testing.T) {
		result := appContext.GetCrossChainFlags()
		require.Equal(t, ccFlags, result)
	})
	t.Run("should get additional chains", func(t *testing.T) {
		result := appContext.GetAdditionalChains()
		require.Equal(t, additionalChains, result)
	})
	t.Run("should get block header enabled chains", func(t *testing.T) {
		for _, chain := range headerSupportedChains {
			result, found := appContext.GetBlockHeaderEnabledChains(chain.ChainId)
			require.True(t, found)
			require.Equal(t, chain, result)
		}
	})
}

func TestIsOutboundObservationEnabled(t *testing.T) {
	// create test chain params and flags
	evmChain := chains.Ethereum
	ccFlags := sample.CrosschainFlags()
	verificationFlags := sample.HeaderSupportedChains()
	chainParams := &observertypes.ChainParams{
		ChainId:     evmChain.ChainId,
		IsSupported: true,
	}

	t.Run("should return true if chain is supported and outbound flag is enabled", func(t *testing.T) {
		appContext := makeAppContext(evmChain, chainParams, *ccFlags, verificationFlags)

		require.True(t, appContext.IsOutboundObservationEnabled(*chainParams))
	})
	t.Run("should return false if chain is not supported yet", func(t *testing.T) {
		paramsUnsupported := &observertypes.ChainParams{ChainId: evmChain.ChainId, IsSupported: false}
		appContextUnsupported := makeAppContext(evmChain, paramsUnsupported, *ccFlags, verificationFlags)

		require.False(t, appContextUnsupported.IsOutboundObservationEnabled(*paramsUnsupported))
	})
	t.Run("should return false if outbound flag is disabled", func(t *testing.T) {
		flagsDisabled := ccFlags
		flagsDisabled.IsOutboundEnabled = false
		coreContextDisabled := makeAppContext(evmChain, chainParams, *flagsDisabled, verificationFlags)

		require.False(t, coreContextDisabled.IsOutboundObservationEnabled(*chainParams))
	})
}

func TestIsInboundObservationEnabled(t *testing.T) {
	// create test chain params and flags
	evmChain := chains.Ethereum
	ccFlags := sample.CrosschainFlags()
	verificationFlags := sample.HeaderSupportedChains()
	chainParams := &observertypes.ChainParams{
		ChainId:     evmChain.ChainId,
		IsSupported: true,
	}

	t.Run("should return true if chain is supported and inbound flag is enabled", func(t *testing.T) {
		appContext := makeAppContext(evmChain, chainParams, *ccFlags, verificationFlags)

		require.True(t, appContext.IsInboundObservationEnabled(*chainParams))
	})

	t.Run("should return false if chain is not supported yet", func(t *testing.T) {
		paramsUnsupported := &observertypes.ChainParams{ChainId: evmChain.ChainId, IsSupported: false}
		appContextUnsupported := makeAppContext(evmChain, paramsUnsupported, *ccFlags, verificationFlags)

		require.False(t, appContextUnsupported.IsInboundObservationEnabled(*paramsUnsupported))
	})

	t.Run("should return false if inbound flag is disabled", func(t *testing.T) {
		flagsDisabled := ccFlags
		flagsDisabled.IsInboundEnabled = false
		appContextDisabled := makeAppContext(evmChain, chainParams, *flagsDisabled, verificationFlags)

		require.False(t, appContextDisabled.IsInboundObservationEnabled(*chainParams))
	})
}

// makeAppContext makes a test app context with provided chain params and flags
func makeAppContext(
	evmChain chains.Chain,
	evmChainParams *observertypes.ChainParams,
	ccFlags observertypes.CrosschainFlags,
	headerSupportedChains []lightclienttypes.HeaderSupportedChain,
) *context.AppContext {
	// create config
	cfg := config.NewConfig()
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Chain: evmChain,
	}

	// create app context
	appContext := context.New(cfg)
	newChainParams := make(map[int64]*observertypes.ChainParams)
	newChainParams[evmChain.ChainId] = evmChainParams

	// feed app context fields
	appContext.Update(
		observertypes.Keygen{},
		"",
		[]chains.Chain{evmChain},
		newChainParams,
		&chaincfg.RegressionNetParams,
		ccFlags,
		[]chains.Chain{},
		headerSupportedChains,
		zerolog.Logger{},
	)
	return appContext
}
