package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestGetParams(t *testing.T) {
	k, ctx, _, _ := testkeeper.FungibleKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	assert.EqualValues(t, params, k.GetParams(ctx))
}
