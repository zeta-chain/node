package corecontext_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	corecontext "github.com/zeta-chain/zetacore/zetaclient/core_context"
)

func TestNewZetaCoreContext(t *testing.T) {
	t.Run("should create new zeta core context with empty config", func(t *testing.T) {
		testCfg := config.NewConfig()

		zetaContext := corecontext.NewZetaCoreContext(testCfg)
		require.NotNil(t, zetaContext)

		// assert keygen
		keyGen := zetaContext.GetKeygen()
		require.Equal(t, observertypes.Keygen{}, keyGen)

		// assert enabled chains
		require.Empty(t, len(zetaContext.GetEnabledChains()))

		// assert current tss pubkey
		require.Equal(t, "", zetaContext.GetCurrentTssPubkey())

		// assert btc chain params
		chain, btcChainParams, btcChainParamsFound := zetaContext.GetBTCChainParams()
		require.Equal(t, common.Chain{}, chain)
		require.False(t, btcChainParamsFound)
		require.Equal(t, &observertypes.ChainParams{}, btcChainParams)

		// assert evm chain params
		allEVMChainParams := zetaContext.GetAllEVMChainParams()
		require.Empty(t, allEVMChainParams)
	})

	t.Run("should create new zeta core context with config containing evm chain params", func(t *testing.T) {
		testCfg := config.NewConfig()
		testCfg.EVMChainConfigs = map[int64]config.EVMConfig{
			1: {
				Chain: common.Chain{
					ChainName: 1,
					ChainId:   1,
				},
			},
			2: {
				Chain: common.Chain{
					ChainName: 2,
					ChainId:   2,
				},
			},
		}
		zetaContext := corecontext.NewZetaCoreContext(testCfg)
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

	t.Run("should create new zeta core context with config containing btc config", func(t *testing.T) {
		testCfg := config.NewConfig()
		testCfg.BitcoinConfig = config.BTCConfig{
			RPCUsername: "test username",
			RPCPassword: "test password",
			RPCHost:     "test host",
			RPCParams:   "test params",
		}
		zetaContext := corecontext.NewZetaCoreContext(testCfg)
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

		zetaContext := corecontext.NewZetaCoreContext(testCfg)
		require.NotNil(t, zetaContext)

		keyGenToUpdate := observertypes.Keygen{
			Status:         observertypes.KeygenStatus_KeyGenSuccess,
			GranteePubkeys: []string{"testpubkey1"},
		}
		enabledChainsToUpdate := []common.Chain{
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
		btcChainParamsToUpdate := &observertypes.ChainParams{
			ChainId: 3,
		}
		tssPubKeyToUpdate := "tsspubkeytest"
		loggers := clientcommon.DefaultLoggers()
		zetaContext.Update(
			&keyGenToUpdate,
			enabledChainsToUpdate,
			evmChainParamsToUpdate,
			btcChainParamsToUpdate,
			tssPubKeyToUpdate,
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

		// assert btc chain params still empty because they were not specified in config
		chain, btcChainParams, btcChainParamsFound := zetaContext.GetBTCChainParams()
		require.Equal(t, common.Chain{}, chain)
		require.False(t, btcChainParamsFound)
		require.Equal(t, &observertypes.ChainParams{}, btcChainParams)

		// assert evm chain params still empty because they were not specified in config
		allEVMChainParams := zetaContext.GetAllEVMChainParams()
		require.Empty(t, allEVMChainParams)
	})

	t.Run("should update core context after being created from config with evm and btc chain params", func(t *testing.T) {
		testCfg := config.NewConfig()
		testCfg.EVMChainConfigs = map[int64]config.EVMConfig{
			1: {
				Chain: common.Chain{
					ChainName: 1,
					ChainId:   1,
				},
			},
			2: {
				Chain: common.Chain{
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

		zetaContext := corecontext.NewZetaCoreContext(testCfg)
		require.NotNil(t, zetaContext)

		keyGenToUpdate := observertypes.Keygen{
			Status:         observertypes.KeygenStatus_KeyGenSuccess,
			GranteePubkeys: []string{"testpubkey1"},
		}
		enabledChainsToUpdate := []common.Chain{
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

		testBtcChain := common.BtcTestNetChain()
		btcChainParamsToUpdate := &observertypes.ChainParams{
			ChainId: testBtcChain.ChainId,
		}
		tssPubKeyToUpdate := "tsspubkeytest"
		loggers := clientcommon.DefaultLoggers()
		zetaContext.Update(
			&keyGenToUpdate,
			enabledChainsToUpdate,
			evmChainParamsToUpdate,
			btcChainParamsToUpdate,
			tssPubKeyToUpdate,
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
	})
}

func assertPanic(t *testing.T, f func(), errorLog string) {
	defer func() {
		r := recover()
		if r != nil {
			require.Contains(t, r, errorLog)
		}
	}()
	f()
}
