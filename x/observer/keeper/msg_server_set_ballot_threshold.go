package keeper

import (
	"context"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k msgServer) SetBallotThreshold(goCtx context.Context, msg *types.MsgSetBallotThreshold) (*types.MsgSetBallotThresholdResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Creator != fungibletypes.AdminAddress {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "invalid creator address (%s)", msg.Creator)
	}

	params := k.GetParams(ctx)
	thresholds := params.BallotThresholds

	chain := msg.Chain
	obChain := types.ParseCommonChaintoObservationChain(chain)
	if obChain == types.ObserverChain_Empty {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid chain: %s", chain)
	}
	threshold, err := sdk.NewDecFromStr(msg.Threshold)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid threshold (%s): %s", err, msg.Threshold)
	}

	found := false
	for _, v := range thresholds {
		if v.Chain == obChain {
			v.Threshold = threshold
			found = true
		}
	}
	if !found {
		thresholds = append(thresholds, &types.BallotThreshold{Chain: obChain, Threshold: threshold, Observation: types.ObservationType_InBoundTx})
		thresholds = append(thresholds, &types.BallotThreshold{Chain: obChain, Threshold: threshold, Observation: types.ObservationType_OutBoundTx})
	}

	k.SetParams(ctx, types.Params{BallotThresholds: thresholds})

	return &types.MsgSetBallotThresholdResponse{}, nil
}
