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
		keyGen, keyGenFound := zetaContext.GetKeygen()
		// assert keygen
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
}
