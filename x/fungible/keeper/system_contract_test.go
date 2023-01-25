package keeper_test

import (
	"fmt"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"testing"
)

func TestKeeper_GetSystemContract(t *testing.T) {
	keeper, ctx := keepertest.FungibleKeeper(t)
	keeper.SetSystemContract(ctx, types.SystemContract{SystemContract: "test"})
	val, b := keeper.GetSystemContract(ctx)
	fmt.Println(val, b)
}
