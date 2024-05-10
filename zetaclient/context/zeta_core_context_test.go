package context_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	context "github.com/zeta-chain/zetacore/zetaclient/context"
)

func assertPanic(t *testing.T, f func(), errorLog string) {
	defer func() {
		r := recover()
		if r != nil {
			require.Contains(t, r, errorLog)
		}
	}()
	f()
}

func getTestCoreContext(
	evmChain chains.Chain,
	evmChainParams *observertypes.ChainParams,
	ccFlags observertypes.CrosschainFlags,
	headerSupportedChains []lightclienttypes.HeaderSupportedChain,
) *context.ZetaCoreContext {
	// create config
	cfg := config.NewConfig()
	cfg.EVMChainConfigs[evmChain.ChainId] = config.EVMConfig{
		Chain: evmChain,
	}
	// create core context
	coreContext := context.NewZetaCoreContext(cfg)
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
		headerSupportedChains,
		true,
		zerolog.Logger{},
	)
	return coreContext
}

func TestNewZetaCoreContext(t *testing.T) {
	t.Run("should create new zetacore context with empty config", func(t *testing.T) {
		testCfg := config.NewConfig()

		zetaContext := context.NewZetaCoreContext(testCfg)
		require.NotNil(t, zetaContext)

		// assert keygen
		keyGen := zetaContext.GetKeygen()
		require.Equal(t, observertypes.Keygen{}, keyGen)

		// assert enabled chains
		require.Empty(t, len(zetaContext.GetEnabledChains()))

		// assert external chains
		require.Empty(t, len(zetaContext.GetEnabledExternalChains()))

		// assert current tss pubkey
		require.Equal(t, "", zetaContext.GetCurrentTssPubkey())

		// assert btc chain params
		chain, btcChainParams, btcChainParamsFound := zetaContext.GetBTCChainParams()
		require.Equal(t, chains.Chain{}, chain)
		require.False(t, btcChainParamsFound)
		require.Equal(t, &observertypes.ChainParams{}, btcChainParams)

		// assert evm chain params
		allEVMChainParams := zetaContext.GetAllEVMChainParams()
		require.Empty(t, allEVMChainParams)
	})

	t.Run("should create new zetacore context with config containing evm chain params", func(t *testing.T) {
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
		zetaContext := context.NewZetaCoreContext(testCfg)
		require.NotNil(t, zetaContext)

		// assert evm chain params
		allEVMChainParams := zetaContext.GetAllEVMChainParams()
		require.Equal(t, 2, len(allEVMChainParams))
		require.Equal(t, &observertypes.ChainParams{}, allEVMChainParams[1])
		require.Equal(t, &observertypes.ChainParams{}, allEVMChainParams[2])

		evmChainParams1, found := zetaContext.GetEVMChainParams(1)
		require.True(t, found)
		require.Equal(t, &observertypes.ChainParams{}, evmChainParams1)

		evmChainParams2, found := zetaContext.GetEVMChainParams(2)
		require.True(t, found)
		require.Equal(t, &observertypes.ChainParams{}, evmChainParams2)
	})

	t.Run("should create new zetacore context with config containing btc config", func(t *testing.T) {
		testCfg := config.NewConfig()
		testCfg.BitcoinConfig = config.BTCConfig{
			RPCUsername: "test username",
			RPCPassword: "test password",
			RPCHost:     "test host",
			RPCParams:   "test params",
		}
		zetaContext := context.NewZetaCoreContext(testCfg)
		require.NotNil(t, zetaContext)

		// assert btc chain params panic because chain params are not yet updated
		assertPanic(t, func() {
			zetaContext.GetBTCChainParams()
		}, "BTCChain is missing for chainID 0")
	})
}

func TestUpdateZetaCoreContext(t *testing.T) {
	t.Run("should update core context after being created from empty config", func(t *testing.T) {
		testCfg := config.NewConfig()

		zetaContext := context.NewZetaCoreContext(testCfg)
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
			chains.ZetaTestnetChain,
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
		loggers := clientcommon.DefaultLoggers()
		crosschainFlags := sample.CrosschainFlags()
		verificationFlags := sample.HeaderSupportedChains()

		require.NotNil(t, crosschainFlags)
		zetaContext.Update(
			&keyGenToUpdate,
			enabledChainsToUpdate,
			evmChainParamsToUpdate,
			btcChainParamsToUpdate,
			tssPubKeyToUpdate,
			*crosschainFlags,
			verificationFlags,
			false,
			loggers.Std,
		)

		// assert keygen updated
		keyGen := zetaContext.GetKeygen()
		require.Equal(t, keyGenToUpdate, keyGen)

		// assert enabled chains updated
		require.Equal(t, enabledChainsToUpdate, zetaContext.GetEnabledChains())

		// assert enabled external chains
		require.Equal(t, enabledChainsToUpdate[0:2], zetaContext.GetEnabledExternalChains())

		// assert current tss pubkey updated
		require.Equal(t, tssPubKeyToUpdate, zetaContext.GetCurrentTssPubkey())

		// assert btc chain params still empty because they were not specified in config
		chain, btcChainParams, btcChainParamsFound := zetaContext.GetBTCChainParams()
		require.Equal(t, chains.Chain{}, chain)
		require.False(t, btcChainParamsFound)
		require.Equal(t, &observertypes.ChainParams{}, btcChainParams)

		// assert evm chain params still empty because they were not specified in config
		allEVMChainParams := zetaContext.GetAllEVMChainParams()
		require.Empty(t, allEVMChainParams)

		ccFlags := zetaContext.GetCrossChainFlags()
		require.Equal(t, *crosschainFlags, ccFlags)

		verFlags := zetaContext.GetAllHeaderEnabledChains()
		require.Equal(t, verificationFlags, verFlags)
	})

	t.Run("should update core context after being created from config with evm and btc chain params", func(t *testing.T) {
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

		zetaContext := context.NewZetaCoreContext(testCfg)
		require.NotNil(t, zetaContext)

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

		testBtcChain := chains.BtcTestNetChain
		btcChainParamsToUpdate := &observertypes.ChainParams{
			ChainId: testBtcChain.ChainId,
		}
		tssPubKeyToUpdate := "tsspubkeytest"
		crosschainFlags := sample.CrosschainFlags()
		verificationFlags := sample.HeaderSupportedChains()
		require.NotNil(t, crosschainFlags)
		loggers := clientcommon.DefaultLoggers()
		zetaContext.Update(
			&keyGenToUpdate,
			enabledChainsToUpdate,
			evmChainParamsToUpdate,
			btcChainParamsToUpdate,
			tssPubKeyToUpdate,
			*crosschainFlags,
			verificationFlags,
			false,
			loggers.Std,
		)

		// assert keygen updated
		keyGen := zetaContext.GetKeygen()
		require.Equal(t, keyGenToUpdate, keyGen)

		// assert enabled chains updated
		require.Equal(t, enabledChainsToUpdate, zetaContext.GetEnabledChains())

		// assert current tss pubkey updated
		require.Equal(t, tssPubKeyToUpdate, zetaContext.GetCurrentTssPubkey())

		// assert btc chain params
		chain, btcChainParams, btcChainParamsFound := zetaContext.GetBTCChainParams()
		require.Equal(t, testBtcChain, chain)
		require.True(t, btcChainParamsFound)
		require.Equal(t, btcChainParamsToUpdate, btcChainParams)

		// assert evm chain params
		allEVMChainParams := zetaContext.GetAllEVMChainParams()
		require.Equal(t, evmChainParamsToUpdate, allEVMChainParams)

		evmChainParams1, found := zetaContext.GetEVMChainParams(1)
		require.True(t, found)
		require.Equal(t, evmChainParamsToUpdate[1], evmChainParams1)

		evmChainParams2, found := zetaContext.GetEVMChainParams(2)
		require.True(t, found)
		require.Equal(t, evmChainParamsToUpdate[2], evmChainParams2)

		ccFlags := zetaContext.GetCrossChainFlags()
		require.Equal(t, ccFlags, *crosschainFlags)

		verFlags := zetaContext.GetAllHeaderEnabledChains()
		require.Equal(t, verFlags, verificationFlags)
	})
}

func TestIsOutboundObservationEnabled(t *testing.T) {
	// create test chain params and flags
	evmChain := chains.EthChain
	ccFlags := *sample.CrosschainFlags()
	verificationFlags := sample.HeaderSupportedChains()
	chainParams := &observertypes.ChainParams{
		ChainId:     evmChain.ChainId,
		IsSupported: true,
	}

	t.Run("should return true if chain is supported and outbound flag is enabled", func(t *testing.T) {
		coreCTX := getTestCoreContext(evmChain, chainParams, ccFlags, verificationFlags)
		require.True(t, context.IsOutboundObservationEnabled(coreCTX, *chainParams))
	})
	t.Run("should return false if chain is not supported yet", func(t *testing.T) {
		paramsUnsupported := &observertypes.ChainParams{
			ChainId:     evmChain.ChainId,
			IsSupported: false,
		}
		coreCTXUnsupported := getTestCoreContext(evmChain, paramsUnsupported, ccFlags, verificationFlags)
		require.False(t, context.IsOutboundObservationEnabled(coreCTXUnsupported, *paramsUnsupported))
	})
	t.Run("should return false if outbound flag is disabled", func(t *testing.T) {
		flagsDisabled := ccFlags
		flagsDisabled.IsOutboundEnabled = false
		coreCTXDisabled := getTestCoreContext(evmChain, chainParams, flagsDisabled, verificationFlags)
		require.False(t, context.IsOutboundObservationEnabled(coreCTXDisabled, *chainParams))
	})
}

func TestIsInboundObservationEnabled(t *testing.T) {
	// create test chain params and flags
	evmChain := chains.EthChain
	ccFlags := *sample.CrosschainFlags()
	verificationFlags := sample.HeaderSupportedChains()
	chainParams := &observertypes.ChainParams{
		ChainId:     evmChain.ChainId,
		IsSupported: true,
	}

	t.Run("should return true if chain is supported and inbound flag is enabled", func(t *testing.T) {
		coreCTX := getTestCoreContext(evmChain, chainParams, ccFlags, verificationFlags)
		require.True(t, context.IsInboundObservationEnabled(coreCTX, *chainParams))
	})
	t.Run("should return false if chain is not supported yet", func(t *testing.T) {
		paramsUnsupported := &observertypes.ChainParams{
			ChainId:     evmChain.ChainId,
			IsSupported: false,
		}
		coreCTXUnsupported := getTestCoreContext(evmChain, paramsUnsupported, ccFlags, verificationFlags)
		require.False(t, context.IsInboundObservationEnabled(coreCTXUnsupported, *paramsUnsupported))
	})
	t.Run("should return false if inbound flag is disabled", func(t *testing.T) {
		flagsDisabled := ccFlags
		flagsDisabled.IsInboundEnabled = false
		coreCTXDisabled := getTestCoreContext(evmChain, chainParams, flagsDisabled, verificationFlags)
		require.False(t, context.IsInboundObservationEnabled(coreCTXDisabled, *chainParams))
	})
}
