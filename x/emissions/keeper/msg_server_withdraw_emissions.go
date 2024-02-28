package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/cmd/zetacored/config"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

// WithdrawEmission create a withdraw emission object , which is then process at endblock
// The withdraw emission object is created and stored
// using the address of the creator as the index key ,therefore, if more that one withdraw requests are created in a block on thr last one would be processed.
// Creating a withdraw does not guarantee that the emission will be processed
// All withdraws for a block are deleted at the end of the block irrespective of whether they were processed or not.
func (k msgServer) WithdrawEmission(goCtx context.Context, msg *types.MsgWithdrawEmission) (*types.MsgWithdrawEmissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if the creator address is valid
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidAddress, err.Error())
	}

	// check if undistributed rewards pool has enough balance to process this request.
	// This is just a basic check , the actual processing at endblock might still fail if the pool balance gets affected .
	undistributedRewardsBalanced := k.GetBankKeeper().GetBalance(ctx, types.UndistributedObserverRewardsPoolAddress, config.BaseDenom)
	if undistributedRewardsBalanced.Amount.LT(msg.Amount) {
		return nil, errorsmod.Wrap(types.ErrRewardsPoolDoesNotHaveEnoughBalance, " rewards pool does not have enough balance to process this request")
	}

	// create a withdraw emission object
	// CreateWithdrawEmissions makes sure that enough withdrawable emissions are available before creating the withdraw object
	err = k.CreateWithdrawEmissions(ctx, msg.Creator, msg.Amount)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrUnableToCreateWithdrawEmissions, err.Error())
	}
	return &types.MsgWithdrawEmissionResponse{}, nil
}
