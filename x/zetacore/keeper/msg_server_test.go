package keeper_test

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/zetacore/keeper"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.ZetacoreKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
