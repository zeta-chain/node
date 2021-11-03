package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) NonceVoter(goCtx context.Context, msg *types.MsgNonceVoter) (*types.MsgNonceVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	chain := msg.Chain
	chainNonce, isFound := k.GetChainNonces(ctx, chain)
	if isFound && chainNonce.Nonce == msg.Nonce {
		chainNonce.Signers = append(chainNonce.Signers, msg.Creator)
	} else if !isFound {
		chainNonce = types.ChainNonces{
			Creator: msg.Creator,
			Index:   msg.Chain,
			Chain:   msg.Chain,
			Nonce:   msg.Nonce,
			Signers: []string{msg.Creator},
		}
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("chainNonce vote msg does not match state: %v vs %v", msg, chainNonce))
	}

	// TODO:
	if len(chainNonce.Signers) == 2 {
		chainNonce.FinalizedHeight = uint64(ctx.BlockHeader().Height)
	}

	k.SetChainNonces(ctx, chainNonce)
	return &types.MsgNonceVoterResponse{}, nil
}
