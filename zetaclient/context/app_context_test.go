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

// getTestAppContext creates a test app context with provided chain params and flags
func getTestAppContext(
	evmChain chains.Chain,
	evmChainParams *observertypes.ChainParams,
	btcChainParams *observertypes.ChainParams,
	ccFlags *observertypes.CrosschainFlags,
	headerSupportedChains []lightclienttypes.HeaderSupportedChain,
) *context.AppContext {
	// create config
	cfg := config.NewConfig()
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Chain: evmChain,
	}
	if btcChainParams != nil {
		cfg.BitcoinConfig = config.BTCConfig{
			RPCUsername: "test",
		}
	}

	// create app context
	appContext := context.NewAppContext(cfg)
	newChainParams := make(map[int64]*observertypes.ChainParams)
	newChainParams[evmChain.ChainId] = evmChainParams
	newChainParams[btcChainParams.ChainId] = btcChainParams

	// create crosschain flags if not provided
	if ccFlags == nil {
		ccFlags = sample.CrosschainFlags()
	}

	// feed chain params
	appContext.Update(
		cfg,
		observertypes.Keygen{},
		[]chains.Chain{evmChain},
		newChainParams,
		&chaincfg.RegressionNetParams,
		"",
		*ccFlags,
		headerSupportedChains,
		true,
		zerolog.Logger{},
	)
	return appContext
}

func TestNewAppContext(t *testing.T) {
	t.Run("should create new app context with empty config", func(t *testing.T) {
		testCfg := config.NewConfig()

		zetaContext := context.NewAppContext(testCfg)
		require.NotNil(t, zetaContext)

		// assert keygen
		keyGen := zetaContext.GetKeygen()
		require.Equal(t, observertypes.Keygen{}, keyGen)

		// assert external chains
		require.Empty(t, len(zetaContext.GetEnabledExternalChains()))

		// assert current tss pubkey
		require.Equal(t, "", zetaContext.GetCurrentTssPubkey())

		// assert external chain params
		externalChainParams := zetaContext.GetEnabledExternalChainParams()
		require.Empty(t, externalChainParams)
	})

	t.Run("should return nil chain params if chain id is not found", func(t *testing.T) {
		// create config with btc config
		testCfg := config.NewConfig()
		testCfg.BitcoinConfig = config.BTCConfig{
			RPCUsername: "test_user",
			RPCPassword: "test_password",
		}

		// create app context with 0 chain id
		zetaContext := context.NewAppContext(testCfg)
		require.NotNil(t, zetaContext)
	})

	t.Run("should create new app context with config containing evm chain params", func(t *testing.T) {
		testCfg := config.NewConfig()
		testCfg.EVMChainConfigs = map[int64]config.EVMConfig{
			1: {
				Chain: chains.Chain{
					ChainName: 1,
					ChainId:   1,
				},
			},
			2: {
				Chain: chains.Chain{
					ChainName: 2,
					ChainId:   2,
				},
			},
		}
		zetaContext := context.NewAppContext(testCfg)
		require.NotNil(t, zetaContext)

		// assert external chain params
		externalChainParams := zetaContext.GetEnabledExternalChainParams()
		require.Equal(t, 2, len(externalChainParams))
		require.Equal(t, &observertypes.ChainParams{}, externalChainParams[1])
		require.Equal(t, &observertypes.ChainParams{}, externalChainParams[2])

		chainParams1, found := zetaContext.GetExternalChainParams(1)
		require.True(t, found)
		require.Equal(t, &observertypes.ChainParams{}, chainParams1)

		chainParams2, found := zetaContext.GetExternalChainParams(2)
		require.True(t, found)
		require.Equal(t, &observertypes.ChainParams{}, chainParams2)
	})

	t.Run("should create new app context with config containing btc config", func(t *testing.T) {
		testCfg := config.NewConfig()
		testCfg.BitcoinConfig = config.BTCConfig{
			RPCUsername: "test username",
			RPCPassword: "test password",
			RPCHost:     "test host",
			RPCParams:   "test params",
		}
		zetaContext := context.NewAppContext(testCfg)
		require.NotNil(t, zetaContext)
	})
}

func TestUpdateAppContext(t *testing.T) {
	t.Run("should update app context after being created from empty config", func(t *testing.T) {
		testCfg := config.NewConfig()

		zetaContext := context.NewAppContext(testCfg)
		require.NotNil(t, zetaContext)

		keyGenToUpdate := observertypes.Keygen{
			Status:         observertypes.KeygenStatus_KeyGenSuccess,
			GranteePubkeys: []string{"testpubkey1"},
		}
		enabledChainsToUpdate := []chains.Chain{
			{
				ChainName:  1,
				ChainId:    1,
				IsExternal: true,
			},
			{
				ChainName:  2,
				ChainId:    2,
				IsExternal: true,
			},
			chains.ZetaChainTestnet,
		}
		newChainParamsToUpdate := map[int64]*observertypes.ChainParams{
			1: {
				ChainId: 1,
			},
			2: {
				ChainId: 2,
			},
			3: {
				ChainId: 3,
			},
		}
		tssPubKeyToUpdate := "tsspubkeytest"
		crosschainFlags := sample.CrosschainFlags()
		verificationFlags := sample.HeaderSupportedChains()

		require.NotNil(t, crosschainFlags)
		zetaContext.Update(
			testCfg,
			keyGenToUpdate,
			enabledChainsToUpdate,
			newChainParamsToUpdate,
			&chaincfg.RegressionNetParams,
			tssPubKeyToUpdate,
			*crosschainFlags,
			verificationFlags,
			false,
			log.Logger,
		)

		// assert keygen updated
		keyGen := zetaContext.GetKeygen()
		require.Equal(t, keyGenToUpdate, keyGen)

		// assert enabled external chains
		require.Equal(t, enabledChainsToUpdate[0:2], zetaContext.GetEnabledExternalChains())

		// assert current tss pubkey updated
		require.Equal(t, tssPubKeyToUpdate, zetaContext.GetCurrentTssPubkey())

		// assert evm chain params still empty because they were not specified in config
		externalChainParams := zetaContext.GetEnabledExternalChainParams()
		require.Empty(t, externalChainParams)

		ccFlags := zetaContext.GetCrossChainFlags()
		require.Equal(t, *crosschainFlags, ccFlags)

		verFlags := zetaContext.GetAllHeaderEnabledChains()
		require.Equal(t, verificationFlags, verFlags)
	})

	t.Run(
		"should update app context after being created from config with evm and btc chain params",
		func(t *testing.T) {
			testCfg := config.NewConfig()
			testCfg.EVMChainConfigs = map[int64]config.EVMConfig{
				1: {
					Chain: chains.Chain{
						ChainName: 1,
						ChainId:   1,
					},
				},
				2: {
					Chain: chains.Chain{
						ChainName: 2,
						ChainId:   2,
					},
				},
			}
			testCfg.BitcoinConfig = config.BTCConfig{
				RPCUsername: "test username",
				RPCPassword: "test password",
				RPCHost:     "test host",
				RPCParams:   "test params",
			}

			zetaContext := context.NewAppContext(testCfg)
			require.NotNil(t, zetaContext)

			keyGenToUpdate := observertypes.Keygen{
				Status:         observertypes.KeygenStatus_KeyGenSuccess,
				GranteePubkeys: []string{"testpubkey1"},
			}
			enabledChainsToUpdate := []chains.Chain{
				{
					ChainName:  1,
					ChainId:    1,
					IsExternal: true,
				},
				{
					ChainName:  2,
					ChainId:    2,
					IsExternal: true,
				},
			}
			newChainParamsToUpdate := map[int64]*observertypes.ChainParams{
				1: {
					ChainId: 1,
				},
				2: {
					ChainId: 2,
				},
				chains.BitcoinTestnet.ChainId: {
					ChainId: chains.BitcoinTestnet.ChainId,
				},
			}

			tssPubKeyToUpdate := "tsspubkeytest"
			crosschainFlags := sample.CrosschainFlags()
			verificationFlags := sample.HeaderSupportedChains()
			require.NotNil(t, crosschainFlags)
			zetaContext.Update(
				testCfg,
				keyGenToUpdate,
				enabledChainsToUpdate,
				newChainParamsToUpdate,
				&chaincfg.RegressionNetParams,
				tssPubKeyToUpdate,
				*crosschainFlags,
				verificationFlags,
				false,
				log.Logger,
			)

			// assert keygen updated
			keyGen := zetaContext.GetKeygen()
			require.Equal(t, keyGenToUpdate, keyGen)

			// assert enabled chains updated
			require.Equal(t, enabledChainsToUpdate, zetaContext.GetEnabledExternalChains())

			// assert current tss pubkey updated
			require.Equal(t, tssPubKeyToUpdate, zetaContext.GetCurrentTssPubkey())

			// assert external chain params
			externalChainParams := zetaContext.GetEnabledExternalChainParams()
			require.Equal(t, newChainParamsToUpdate, externalChainParams)

			chainParams1, found := zetaContext.GetExternalChainParams(1)
			require.True(t, found)
			require.Equal(t, newChainParamsToUpdate[1], chainParams1)

			chainParams2, found := zetaContext.GetExternalChainParams(2)
			require.True(t, found)
			require.Equal(t, newChainParamsToUpdate[2], chainParams2)

			ccFlags := zetaContext.GetCrossChainFlags()
			require.Equal(t, ccFlags, *crosschainFlags)

			verFlags := zetaContext.GetAllHeaderEnabledChains()
			require.Equal(t, verFlags, verificationFlags)
		},
	)
}

func TestIsOutboundObservationEnabled(t *testing.T) {
	// create test chain params and flags
	evmChain := chains.Ethereum
	ccFlags := *sample.CrosschainFlags()
	verificationFlags := sample.HeaderSupportedChains()
	chainParams := &observertypes.ChainParams{
		ChainId:     evmChain.ChainId,
		IsSupported: true,
	}

	t.Run("should return true if chain is supported and outbound flag is enabled", func(t *testing.T) {
		appCTX := getTestAppContext(evmChain, chainParams, nil, &ccFlags, verificationFlags)
		require.True(t, context.IsOutboundObservationEnabled(appCTX, *chainParams))
	})
	t.Run("should return false if chain is not supported yet", func(t *testing.T) {
		paramsUnsupported := &observertypes.ChainParams{
			ChainId:     evmChain.ChainId,
			IsSupported: false,
		}
		appCTXUnsupported := getTestAppContext(evmChain, paramsUnsupported, nil, &ccFlags, verificationFlags)
		require.False(t, context.IsOutboundObservationEnabled(appCTXUnsupported, *paramsUnsupported))
	})
	t.Run("should return false if outbound flag is disabled", func(t *testing.T) {
		flagsDisabled := ccFlags
		flagsDisabled.IsOutboundEnabled = false
		appCTXDisabled := getTestAppContext(evmChain, chainParams, nil, &flagsDisabled, verificationFlags)
		require.False(t, context.IsOutboundObservationEnabled(appCTXDisabled, *chainParams))
	})
}

func TestIsInboundObservationEnabled(t *testing.T) {
	// create test chain params and flags
	evmChain := chains.Ethereum
	ccFlags := *sample.CrosschainFlags()
	verificationFlags := sample.HeaderSupportedChains()
	chainParams := &observertypes.ChainParams{
		ChainId:     evmChain.ChainId,
		IsSupported: true,
	}

	t.Run("should return true if chain is supported and inbound flag is enabled", func(t *testing.T) {
		appCTX := getTestAppContext(evmChain, chainParams, nil, &ccFlags, verificationFlags)
		require.True(t, context.IsInboundObservationEnabled(appCTX, *chainParams))
	})
	t.Run("should return false if chain is not supported yet", func(t *testing.T) {
		paramsUnsupported := &observertypes.ChainParams{
			ChainId:     evmChain.ChainId,
			IsSupported: false,
		}
		appCTXUnsupported := getTestAppContext(evmChain, paramsUnsupported, nil, &ccFlags, verificationFlags)
		require.False(t, context.IsInboundObservationEnabled(appCTXUnsupported, *paramsUnsupported))
	})
	t.Run("should return false if inbound flag is disabled", func(t *testing.T) {
		flagsDisabled := ccFlags
		flagsDisabled.IsInboundEnabled = false
		appCTXDisabled := getTestAppContext(evmChain, chainParams, nil, &flagsDisabled, verificationFlags)
		require.False(t, context.IsInboundObservationEnabled(appCTXDisabled, *chainParams))
	})
}
