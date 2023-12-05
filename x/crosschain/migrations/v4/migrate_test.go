package v4_test

import (
	"fmt"
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	v4 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v4"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMigrateStore(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	amountZeta := SetRandomCctx(ctx, *k)
	err := v4.MigrateStore(ctx, k.GetStoreKey(), k.GetCodec())
	assert.NoError(t, err)
	zetaAccounting, found := k.GetZetaAccounting(ctx)
	assert.True(t, found)
	assert.Equal(t, amountZeta, zetaAccounting.AbortedZetaAmount)

}

func SetRandomCctx(ctx sdk.Context, k keeper.Keeper) sdkmath.Uint {
	totalZeta := sdkmath.ZeroUint()

	i := 0
	r := rand.New(rand.NewSource(9))
	for ; i < 10; i++ {
		amount := sdkmath.NewUint(uint64(r.Uint32()))
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index:      fmt.Sprintf("%d", i),
			CctxStatus: &types.Status{Status: types.CctxStatus_Aborted},
			OutboundTxParams: []*types.OutboundTxParams{{
				Amount:   amount,
				CoinType: common.CoinType_Zeta,
			}},
		})
		totalZeta = totalZeta.Add(amount)
	}
	return totalZeta
}
