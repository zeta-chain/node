package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/cmd/zetacored/config"
	"github.com/zeta-chain/node/x/emissions/types"
)

// WithdrawEmission allows the user to withdraw from their withdrawable emissions.
// on a successful withdrawal, the amount is transferred from the undistributed rewards pool to the user's account.
// if the amount to be withdrawn is greater than the available withdrawable emission, the max available amount is withdrawn.
// if the pool does not have enough balance to process this request, an error is returned.
func (k msgServer) WithdrawEmission(
	goCtx context.Context,
	msg *types.MsgWithdrawEmission,
) (*types.MsgWithdrawEmissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if the creator address is valid
	address, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, err.Error())
	}

	// check if the undistributed rewards pool has enough balance to process this request.
	// This is just a preliminary check, the actual processing at endblock might still fail if the pool balance gets affected.
	undistributedRewardsBalance := k.GetBankKeeper().
		GetBalance(ctx, types.UndistributedObserverRewardsPoolAddress, config.BaseDenom)
	if undistributedRewardsBalance.Amount.LT(msg.Amount) {
		return nil, errorsmod.Wrap(
			types.ErrRewardsPoolDoesNotHaveEnoughBalance,
			" rewards pool does not have enough balance to process this request",
		)
	}

	err = k.RemoveWithdrawableEmission(ctx, msg.Creator, msg.Amount)
	if err != nil {
		return nil, errorsmod.Wrap(
			types.ErrUnableToWithdrawEmissions,
			fmt.Sprintf("error while removing withdrawable emission for address %s : %s", msg.Creator, err),
		)
	}

	err = k.GetBankKeeper().
		SendCoinsFromModuleToAccount(ctx, types.UndistributedObserverRewardsPool, address, sdk.NewCoins(sdk.NewCoin(config.BaseDenom, msg.Amount)))
	if err != nil {
		ctx.Logger().
			Error(fmt.Sprintf("Error while processing withdraw of emission to address %s for amount %s : err %s", address, msg.Amount, err))
		return nil, errorsmod.Wrap(types.ErrUnableToWithdrawEmissions, err.Error())
	}

	return &types.MsgWithdrawEmissionResponse{}, nil
}
