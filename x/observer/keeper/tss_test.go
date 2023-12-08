package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func createTSS(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.TSS {
	tssList := make([]types.TSS, n)
	for i := 0; i < n; i++ {
		tss := types.TSS{
			TssPubkey:           "tssPubkey",
			TssParticipantList:  []string{"tssParticipantList"},
			OperatorAddressList: []string{"operatorAddressList"},
			KeyGenZetaHeight:    int64(100 + i),
			FinalizedZetaHeight: int64(110 + i),
		}
		keeper.SetTSS(ctx, tss)
		keeper.SetTSSHistory(ctx, tss)
		tssList[i] = tss
	}
	return tssList
}

func TestTSSGet(t *testing.T) {
	keeper, ctx := keepertest.ObserverKeeper(t)
	tssSaved := createTSS(keeper, ctx, 1)
	tss, found := keeper.GetTSS(ctx)
	assert.True(t, found)
	assert.Equal(t, tssSaved[len(tssSaved)-1], tss)

}
func TestTSSRemove(t *testing.T) {
	keeper, ctx := keepertest.ObserverKeeper(t)
	_ = createTSS(keeper, ctx, 1)
	keeper.RemoveTSS(ctx)
	_, found := keeper.GetTSS(ctx)
	assert.False(t, found)
}
