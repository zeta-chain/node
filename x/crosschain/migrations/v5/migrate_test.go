package v5_test

import (
	"fmt"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	v5 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v5"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("TestMigrateStore", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctxList := CrossChainTxList(100)
		v5ZetaAccountingAmount := math.ZeroUint()
		v4ZetaAccountingAmount := math.ZeroUint()
		for _, cctx := range cctxList {
			if cctx.CctxStatus.Status != crosschaintypes.CctxStatus_Aborted || cctx.GetCurrentOutTxParam().CoinType != common.CoinType_Zeta {
				continue
			}
			v5ZetaAccountingAmount = v5ZetaAccountingAmount.Add(crosschainkeeper.GetAbortedAmount(cctx))
			v4ZetaAccountingAmount = v4ZetaAccountingAmount.Add(cctx.GetCurrentOutTxParam().Amount)
			k.SetCrossChainTx(ctx, cctx)
		}
		assert.True(t, v5ZetaAccountingAmount.GT(v4ZetaAccountingAmount))
		// Previously set the zeta accounting
		k.SetZetaAccounting(ctx, crosschaintypes.ZetaAccounting{
			AbortedZetaAmount: v4ZetaAccountingAmount,
		})
		err := v5.MigrateStore(ctx, k)
		require.NoError(t, err)
		zetaAccounting, found := k.GetZetaAccounting(ctx)
		require.True(t, found)
		require.True(t, v5ZetaAccountingAmount.Equal(zetaAccounting.AbortedZetaAmount))
	})

}

func CrossChainTxList(count int) []crosschaintypes.CrossChainTx {
	cctxList := make([]crosschaintypes.CrossChainTx, count)
	i := 0
	r := rand.New(rand.NewSource(9))
	for ; i < count/2; i++ {
		amount := math.NewUint(uint64(r.Uint32()))
		cctxList[i] = crosschaintypes.CrossChainTx{
			Index:      fmt.Sprintf("%d", i),
			CctxStatus: &crosschaintypes.Status{Status: crosschaintypes.CctxStatus_Aborted},
			InboundTxParams: &crosschaintypes.InboundTxParams{
				Amount:   amount.Add(math.NewUint(uint64(r.Uint32()))),
				CoinType: common.CoinType_Zeta,
			},
			OutboundTxParams: []*crosschaintypes.OutboundTxParams{{
				Amount:   amount,
				CoinType: common.CoinType_Zeta,
			}},
		}
		for ; i < count; i++ {
			amount := math.NewUint(uint64(r.Uint32()))
			cctxList[i] = crosschaintypes.CrossChainTx{
				Index:      fmt.Sprintf("%d", i),
				CctxStatus: &crosschaintypes.Status{Status: crosschaintypes.CctxStatus_Aborted},
				InboundTxParams: &crosschaintypes.InboundTxParams{
					Amount:   amount,
					CoinType: common.CoinType_Zeta,
				},
				OutboundTxParams: []*crosschaintypes.OutboundTxParams{{
					Amount:   math.ZeroUint(),
					CoinType: common.CoinType_Zeta,
				}},
			}
		}
	}
	return cctxList
}
