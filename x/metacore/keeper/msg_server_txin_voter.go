package keeper

import (
	"context"
	"fmt"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateTxinVoter(goCtx context.Context, msg *types.MsgCreateTxinVoter) (*types.MsgCreateTxinVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetTxinVoter(ctx, msg.Index)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("index %v already set", msg.Index))
	}

	var txinVoter = types.TxinVoter{
		Index:            msg.Index,
		Creator:          msg.Creator,
		TxHash:           msg.TxHash,
		SourceAsset:      msg.SourceAsset,
		SourceAmount:     msg.SourceAmount,
		MBurnt:           msg.MBurnt,
		DestinationAsset: msg.DestinationAsset,
		FromAddress:      msg.FromAddress,
		ToAddress:        msg.ToAddress,
		BlockHeight:      msg.BlockHeight,
	}

	k.SetTxinVoter(ctx, txinVoter)

	// Create Txin, add signers to it
	txin, isFound := k.GetTxin(ctx, msg.TxHash)
	if isFound { // txin already created; add signer to it
		for _, s := range txin.Signers {
			if s == msg.Creator {
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("txin index %s already set from signer %s", msg.TxHash, msg.Creator))
			}
		}
		txin.Signers = append(txin.Signers, msg.Creator)
	} else { // first signer for TxHash
		txin = types.Txin{
			Creator:          msg.Creator,
			Index:            msg.TxHash,
			TxHash:           msg.TxHash,
			SourceAsset:      msg.SourceAsset,
			SourceAmount:     msg.SourceAmount,
			MBurnt:           msg.MBurnt,
			DestinationAsset: msg.DestinationAsset,
			FromAddress:      msg.FromAddress,
			ToAddress:        msg.ToAddress,
			Signers:          []string{msg.Creator},
		}
	}
	k.SetTxin(ctx, txin)

	// see if the txin reached consensus. If so, create the corresponding txout.
	if len(txin.Signers) == 2 { // the first time that txin reaches consensus

	}

	return &types.MsgCreateTxinVoterResponse{}, nil
}
