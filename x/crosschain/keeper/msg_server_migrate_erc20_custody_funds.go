package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// MigrateERC20CustodyFunds migrates the funds from the current TSS to the new TSS
func (k msgServer) MigrateERC20CustodyFunds(
	goCtx context.Context,
	msg *types.MsgMigrateERC20CustodyFunds,
) (*types.MsgMigrateERC20CustodyFundsResponse, error) {
	_ = sdk.UnwrapSDKContext(goCtx)

	return &types.MsgMigrateERC20CustodyFundsResponse{}, nil
}
