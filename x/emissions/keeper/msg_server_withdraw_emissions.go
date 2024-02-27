package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k msgServer) WithdrawEmission(goCtx context.Context, msg *types.MsgWithdrawEmission) (*types.MsgWithdrawEmissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.CreateWithdrawEmissions(ctx, msg.Creator, msg.Amount)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUnableToCreateWithdrawEmissions, err.Error())
	}
	return &types.MsgWithdrawEmissionResponse{}, nil
}
