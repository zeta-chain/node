package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k msgServer) ZetaConversionRateVoter(goCtx context.Context, msg *types.MsgZetaConversionRateVoter) (*types.MsgZetaConversionRateVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgZetaConversionRateVoterResponse{}, nil
}
