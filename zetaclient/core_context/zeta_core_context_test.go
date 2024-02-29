package corecontext

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

func TestNewZetaCoreContext(t *testing.T) {
	t.Run("should create new zeta core context with empty config", func(t *testing.T) {
		testCfg := config.NewConfig()

		zetaContext := NewZetaCoreContext(testCfg)
		require.NotNil(t, zetaContext)

		// assert keygen
		keyGen, keyGenFound := zetaContext.GetKeygen()
		require.False(t, keyGenFound)
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
		testCfg.EVMChainConfigs = map[int64]*config.EVMConfig{
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
		zetaContext := NewZetaCoreContext(testCfg)
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
		testCfg.BitcoinConfig = &config.BTCConfig{
			RPCUsername: "test username",
			RPCPassword: "test password",
			RPCHost:     "test host",
			RPCParams:   "test params",
		}
		zetaContext := NewZetaCoreContext(testCfg)
		require.NotNil(t, zetaContext)

		// assert btc chain params panic because chain params are not yet updated
		assertPanic(t, func() {
			zetaContext.GetBTCChainParams()
		}, "BTCChain is missing for chainID 0")
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
