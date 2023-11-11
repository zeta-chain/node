package v4_test

import (
	"fmt"
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	v4 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v4"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMigrateStore(t *testing.T) {
	k, ctx, _, _ := keepertest.CrosschainKeeper(t)
	amount := SetCrossRandomTx(10, ctx, *k)
	err := v4.MigrateStore(ctx, k.GetStoreKey(), k.GetCodec())
	assert.NoError(t, err)
	abortedZetaAmount, found := k.GetAbortedZetaAmount(ctx)
	assert.True(t, found)
	assert.Equal(t, amount, abortedZetaAmount.Amount)

}

func SetCrossRandomTx(maxlen int, ctx sdk.Context, k keeper.Keeper) sdkmath.Uint {
	total := sdkmath.ZeroUint()

	r := rand.New(rand.NewSource(9))
	for i := 0; i < r.Intn(maxlen); i++ {
		amount := sdkmath.NewUint(uint64(r.Uint32()))
		k.SetCrossChainTx(ctx, types.CrossChainTx{
			Index:            fmt.Sprintf("%d", i),
			CctxStatus:       &types.Status{Status: types.CctxStatus_Aborted},
			OutboundTxParams: []*types.OutboundTxParams{{Amount: amount}},
		})
		total = total.Add(amount)
	}
	return total
}
