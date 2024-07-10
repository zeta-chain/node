package context_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
)

func TestNew(t *testing.T) {
	var (
		testCfg = config.New(false)
		logger  = zerolog.Nop()
	)

	t.Run("should create new zetacore context with empty config", func(t *testing.T) {
		appContext := context.New(testCfg, logger)
		require.NotNil(t, appContext)

		// assert keygen
		keyGen := appContext.GetKeygen()
		require.Equal(t, observertypes.Keygen{}, keyGen)

		// assert enabled chains
		require.Empty(t, len(appContext.GetEnabledChains()))

		// assert external chains
		require.Empty(t, len(appContext.GetEnabledExternalChains()))

		// assert current tss pubkey
		require.Equal(t, "", appContext.GetCurrentTssPubKey())

		// assert btc chain params
		chain, btcChainParams, btcChainParamsFound := appContext.GetBTCChainParams()
		require.Equal(t, chains.Chain{}, chain)
		require.False(t, btcChainParamsFound)
		require.Nil(t, btcChainParams)

		// assert evm chain params
		allEVMChainParams := appContext.GetAllEVMChainParams()
		require.Empty(t, allEVMChainParams)
	})

	t.Run("should return nil chain params if chain id is not found", func(t *testing.T) {
		// create config with btc config
		testCfg := config.New(false)
		testCfg.BitcoinConfig = config.BTCConfig{
			RPCUsername: "test_user",
			RPCPassword: "test_password",
		}

		// create zetacore context with 0 chain id
		appContext := context.New(testCfg, logger)
		require.NotNil(t, appContext)

		// assert btc chain params
		chain, btcChainParams, btcChainParamsFound := appContext.GetBTCChainParams()
		require.Equal(t, chains.Chain{}, chain)
		require.False(t, btcChainParamsFound)
		require.Nil(t, btcChainParams)
	})

	t.Run("should create new zetacore context with config containing evm chain params", func(t *testing.T) {
		testCfg := config.New(false)
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
		appContext := context.New(testCfg, logger)
		require.NotNil(t, appContext)

		// assert evm chain params
		allEVMChainParams := appContext.GetAllEVMChainParams()
		require.Equal(t, 2, len(allEVMChainParams))
		require.Equal(t, &observertypes.ChainParams{}, allEVMChainParams[1])
		require.Equal(t, &observertypes.ChainParams{}, allEVMChainParams[2])

		evmChainParams1, found := appContext.GetEVMChainParams(1)
		require.True(t, found)
		require.Equal(t, &observertypes.ChainParams{}, evmChainParams1)

		evmChainParams2, found := appContext.GetEVMChainParams(2)
		require.True(t, found)
		require.Equal(t, &observertypes.ChainParams{}, evmChainParams2)
	})

	t.Run("should create new zetacore context with config containing btc config", func(t *testing.T) {
		testCfg := config.New(false)
		testCfg.BitcoinConfig = config.BTCConfig{
			RPCUsername: "test username",
			RPCPassword: "test password",
			RPCHost:     "test host",
			RPCParams:   "test params",
		}
		appContext := context.New(testCfg, logger)
		require.NotNil(t, appContext)
	})
}

func TestAppContextUpdate(t *testing.T) {
	var (
		testCfg = config.New(false)
		logger  = zerolog.Nop()
	)

	t.Run("should update zetacore context after being created from empty config", func(t *testing.T) {
		appContext := context.New(testCfg, logger)
		require.NotNil(t, appContext)

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
		evmChainParamsToUpdate := map[int64]*observertypes.ChainParams{
			1: {
				ChainId: 1,
			},
			2: {
				ChainId: 2,
			},
		}
		btcChainParamsToUpdate := &observertypes.ChainParams{
			ChainId: 3,
		}
		tssPubKeyToUpdate := "tsspubkeytest"
		crosschainFlags := sample.CrosschainFlags()
		verificationFlags := sample.HeaderSupportedChains()

		require.NotNil(t, crosschainFlags)
		appContext.Update(
			&keyGenToUpdate,
			enabledChainsToUpdate,
			evmChainParamsToUpdate,
			btcChainParamsToUpdate,
			tssPubKeyToUpdate,
			*crosschainFlags,
			[]chains.Chain{},
			verificationFlags,
			false,
		)

		// assert keygen updated
		keyGen := appContext.GetKeygen()
		require.Equal(t, keyGenToUpdate, keyGen)

		// assert enabled chains updated
		require.Equal(t, enabledChainsToUpdate, appContext.GetEnabledChains())

		// assert enabled external chains
		require.Equal(t, enabledChainsToUpdate[0:2], appContext.GetEnabledExternalChains())

		// assert current tss pubkey updated
		require.Equal(t, tssPubKeyToUpdate, appContext.GetCurrentTssPubKey())

		// assert btc chain params still empty because they were not specified in config
		chain, btcChainParams, btcChainParamsFound := appContext.GetBTCChainParams()
		require.Equal(t, chains.Chain{}, chain)
		require.False(t, btcChainParamsFound)
		require.Nil(t, btcChainParams)

		// assert evm chain params still empty because they were not specified in config
		allEVMChainParams := appContext.GetAllEVMChainParams()
		require.Empty(t, allEVMChainParams)

		ccFlags := appContext.GetCrossChainFlags()
		require.Equal(t, *crosschainFlags, ccFlags)

		verFlags := appContext.GetAllHeaderEnabledChains()
		require.Equal(t, verificationFlags, verFlags)
	})

	t.Run(
		"should update zetacore context after being created from config with evm and btc chain params",
		func(t *testing.T) {
			testCfg := config.New(false)
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

			appContext := context.New(testCfg, logger)
			require.NotNil(t, appContext)

			keyGenToUpdate := observertypes.Keygen{
				Status:         observertypes.KeygenStatus_KeyGenSuccess,
				GranteePubkeys: []string{"testpubkey1"},
			}
			enabledChainsToUpdate := []chains.Chain{
				{
					ChainName: 1,
					ChainId:   1,
				},
				{
					ChainName: 2,
					ChainId:   2,
				},
			}
			evmChainParamsToUpdate := map[int64]*observertypes.ChainParams{
				1: {
					ChainId: 1,
				},
				2: {
					ChainId: 2,
				},
			}

			testBtcChain := chains.BitcoinTestnet
			btcChainParamsToUpdate := &observertypes.ChainParams{
				ChainId: testBtcChain.ChainId,
			}
			tssPubKeyToUpdate := "tsspubkeytest"
			crosschainFlags := sample.CrosschainFlags()
			verificationFlags := sample.HeaderSupportedChains()
			require.NotNil(t, crosschainFlags)
			appContext.Update(
				&keyGenToUpdate,
				enabledChainsToUpdate,
				evmChainParamsToUpdate,
				btcChainParamsToUpdate,
				tssPubKeyToUpdate,
				*crosschainFlags,
				[]chains.Chain{},
				verificationFlags,
				false,
			)

			// assert keygen updated
			keyGen := appContext.GetKeygen()
			require.Equal(t, keyGenToUpdate, keyGen)

			// assert enabled chains updated
			require.Equal(t, enabledChainsToUpdate, appContext.GetEnabledChains())

			// assert current tss pubkey updated
			require.Equal(t, tssPubKeyToUpdate, appContext.GetCurrentTssPubKey())

			// assert btc chain params
			chain, btcChainParams, btcChainParamsFound := appContext.GetBTCChainParams()
			require.Equal(t, testBtcChain, chain)
			require.True(t, btcChainParamsFound)
			require.Equal(t, btcChainParamsToUpdate, btcChainParams)

			// assert evm chain params
			allEVMChainParams := appContext.GetAllEVMChainParams()
			require.Equal(t, evmChainParamsToUpdate, allEVMChainParams)

			evmChainParams1, found := appContext.GetEVMChainParams(1)
			require.True(t, found)
			require.Equal(t, evmChainParamsToUpdate[1], evmChainParams1)

			evmChainParams2, found := appContext.GetEVMChainParams(2)
			require.True(t, found)
			require.Equal(t, evmChainParamsToUpdate[2], evmChainParams2)

			ccFlags := appContext.GetCrossChainFlags()
			require.Equal(t, ccFlags, *crosschainFlags)

			verFlags := appContext.GetAllHeaderEnabledChains()
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
		appContext := makeAppContext(evmChain, chainParams, ccFlags, verificationFlags)

		require.True(t, appContext.IsOutboundObservationEnabled(*chainParams))
	})
	t.Run("should return false if chain is not supported yet", func(t *testing.T) {
		paramsUnsupported := &observertypes.ChainParams{ChainId: evmChain.ChainId, IsSupported: false}
		appContextUnsupported := makeAppContext(evmChain, paramsUnsupported, ccFlags, verificationFlags)

		require.False(t, appContextUnsupported.IsOutboundObservationEnabled(*paramsUnsupported))
	})
	t.Run("should return false if outbound flag is disabled", func(t *testing.T) {
		flagsDisabled := ccFlags
		flagsDisabled.IsOutboundEnabled = false
		coreContextDisabled := makeAppContext(evmChain, chainParams, flagsDisabled, verificationFlags)

		require.False(t, coreContextDisabled.IsOutboundObservationEnabled(*chainParams))
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
		appContext := makeAppContext(evmChain, chainParams, ccFlags, verificationFlags)

		require.True(t, appContext.IsInboundObservationEnabled(*chainParams))
	})

	t.Run("should return false if chain is not supported yet", func(t *testing.T) {
		paramsUnsupported := &observertypes.ChainParams{ChainId: evmChain.ChainId, IsSupported: false}
		appContextUnsupported := makeAppContext(evmChain, paramsUnsupported, ccFlags, verificationFlags)

		require.False(t, appContextUnsupported.IsInboundObservationEnabled(*paramsUnsupported))
	})

	t.Run("should return false if inbound flag is disabled", func(t *testing.T) {
		flagsDisabled := ccFlags
		flagsDisabled.IsInboundEnabled = false
		appContextDisabled := makeAppContext(evmChain, chainParams, flagsDisabled, verificationFlags)

		require.False(t, appContextDisabled.IsInboundObservationEnabled(*chainParams))
	})
}

func TestGetBTCChainAndConfig(t *testing.T) {
	logger := zerolog.Nop()

	emptyConfig := config.New(false)
	nonEmptyConfig := config.New(true)

	assertEmpty := func(t *testing.T, chain chains.Chain, btcConfig config.BTCConfig, enabled bool) {
		assert.Empty(t, chain)
		assert.Empty(t, btcConfig)
		assert.False(t, enabled)
	}

	for _, tt := range []struct {
		name   string
		cfg    config.Config
		setup  func(app *context.AppContext)
		assert func(t *testing.T, chain chains.Chain, btcConfig config.BTCConfig, enabled bool)
	}{
		{
			name:   "no btc config",
			cfg:    emptyConfig,
			setup:  nil,
			assert: assertEmpty,
		},
		{
			name:   "btc config exists, but not chain params are set",
			cfg:    nonEmptyConfig,
			setup:  nil,
			assert: assertEmpty,
		},
		{
			name: "btc config exists but chain is invalid",
			cfg:  nonEmptyConfig,
			setup: func(app *context.AppContext) {
				app.Update(
					&observertypes.Keygen{},
					[]chains.Chain{},
					nil,
					&observertypes.ChainParams{ChainId: 123},
					"",
					observertypes.CrosschainFlags{},
					[]chains.Chain{},
					nil,
					true,
				)
			},
			assert: assertEmpty,
		},
		{
			name: "btc config exists and chain params are set",
			cfg:  nonEmptyConfig,
			setup: func(app *context.AppContext) {
				app.Update(
					&observertypes.Keygen{},
					[]chains.Chain{},
					nil,
					&observertypes.ChainParams{ChainId: chains.BitcoinMainnet.ChainId},
					"",
					observertypes.CrosschainFlags{},
					[]chains.Chain{},
					nil,
					true,
				)
			},
			assert: func(t *testing.T, chain chains.Chain, btcConfig config.BTCConfig, enabled bool) {
				assert.Equal(t, chains.BitcoinMainnet.ChainId, chain.ChainId)
				assert.Equal(t, "smoketest", btcConfig.RPCUsername)
				assert.True(t, enabled)
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// Given app context
			appContext := context.New(tt.cfg, logger)

			// And optional setup
			if tt.setup != nil {
				tt.setup(appContext)
			}

			// ACT
			chain, btcConfig, enabled := appContext.GetBTCChainAndConfig()

			// ASSERT
			tt.assert(t, chain, btcConfig, enabled)
		})
	}
}

func TestGetBlockHeaderEnabledChains(t *testing.T) {
	// ARRANGE
	// Given app config
	appContext := context.New(config.New(false), zerolog.Nop())

	// That was eventually updated
	appContext.Update(
		&observertypes.Keygen{},
		[]chains.Chain{},
		nil,
		&observertypes.ChainParams{ChainId: chains.BitcoinMainnet.ChainId},
		"",
		observertypes.CrosschainFlags{},
		[]chains.Chain{},
		[]lightclienttypes.HeaderSupportedChain{
			{ChainId: 1, Enabled: true},
		},
		true,
	)

	// ACT #1 (found)
	chain, found := appContext.GetBlockHeaderEnabledChains(1)

	// ASSERT #1
	assert.True(t, found)
	assert.Equal(t, int64(1), chain.ChainId)
	assert.True(t, chain.Enabled)

	// ACT #2 (not found)
	chain, found = appContext.GetBlockHeaderEnabledChains(2)

	// ASSERT #2
	assert.False(t, found)
	assert.Empty(t, chain)
}

func TestGetAdditionalChains(t *testing.T) {
	// ARRANGE
	// Given app config
	appContext := context.New(config.New(false), zerolog.Nop())

	additionalChains := []chains.Chain{
		sample.Chain(1),
		sample.Chain(2),
		sample.Chain(3),
	}

	// That was eventually updated
	appContext.Update(
		&observertypes.Keygen{},
		[]chains.Chain{},
		nil,
		&observertypes.ChainParams{},
		"",
		observertypes.CrosschainFlags{},
		additionalChains,
		[]lightclienttypes.HeaderSupportedChain{
			{ChainId: 1, Enabled: true},
		},
		true,
	)

	// ACT
	found := appContext.GetAdditionalChains()

	// ASSERT
	assert.EqualValues(t, additionalChains, found)
}

func makeAppContext(
	evmChain chains.Chain,
	evmChainParams *observertypes.ChainParams,
	ccFlags observertypes.CrosschainFlags,
	headerSupportedChains []lightclienttypes.HeaderSupportedChain,
) *context.AppContext {
	// create config
	cfg := config.New(false)
	logger := zerolog.Nop()
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Chain: evmChain,
	}

	// create zetacore context
	coreContext := context.New(cfg, logger)
	evmChainParamsMap := make(map[int64]*observertypes.ChainParams)
	evmChainParamsMap[evmChain.ChainId] = evmChainParams

	// feed chain params
	coreContext.Update(
		&observertypes.Keygen{},
		[]chains.Chain{evmChain},
		evmChainParamsMap,
		nil,
		"",
		ccFlags,
		[]chains.Chain{},
		headerSupportedChains,
		true,
	)

	return coreContext
}
