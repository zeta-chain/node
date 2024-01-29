package emissions_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestDistributeObserverRewards(t *testing.T) {
	k, ctx, _, zk := keepertest.EmisionKeeper(t)
	k.AddObserverEmission(ctx, "zetavaloper1", sdk.NewIntFromUint64(1000000000000000000))
	k.AddObserverEmission(ctx, "zetavaloper2", sdk.NewIntFromUint64(1000000000000000000))
	val, found := k.GetWithdrawableEmission(ctx, "zetavaloper1")
	assert.True(t, found)
	assert.Equal(t, sdk.NewIntFromUint64(1000000000000000000), val.Amount)
	zk.ObserverKeeper.SetObserverSet(ctx, types.ObserverSet{
		ObserverList: []string{"zetavaloper1", "zetavaloper2"},
	})
	v, found := zk.ObserverKeeper.GetObserverSet(ctx)
	assert.True(t, found)
	fmt.Println(v)
}
