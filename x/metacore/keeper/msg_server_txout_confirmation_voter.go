package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) TxoutConfirmationVoter(goCtx context.Context, msg *types.MsgTxoutConfirmationVoter) (*types.MsgTxoutConfirmationVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	found := k.HasTxout(ctx, msg.TxoutId)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("txoutId %d does not exist", msg.TxoutId))
	}

	txoutConf := types.TxoutConfirmation{
		Creator:           "",
		Index:             "",
		TxoutId:           msg.TxoutId,
		TxHash:            msg.TxHash,
		MMint:             msg.MMint,
		DestinationAsset:  msg.DestinationAsset,
		DestinationAmount: msg.DestinationAmount,
		ToAddress:         msg.ToAddress,
		BlockHeight:       msg.BlockHeight,
		Signers:           nil,
		FinalizedHeight:   0,
	}

	hashTxoutConf := crypto.Keccak256Hash([]byte(txoutConf.String()))
	txoutConf2, found := k.GetTxoutConfirmation(ctx, hashTxoutConf.Hex())
	if !found {
		txoutConf.Index = hashTxoutConf.Hex()
		txoutConf.Signers = append(txoutConf.Signers, msg.Creator)
	} else {
		txoutConf2.Signers = append(txoutConf2.Signers, msg.Creator)
		txoutConf = txoutConf2
	}

	k.SetTxoutConfirmation(ctx, txoutConf)


	if len(txoutConf.Signers) == 2 { // TODO: fix threshold
		txoutConf.FinalizedHeight = uint64(ctx.BlockHeader().Height)
		k.RemoveTxout(ctx, txoutConf.TxoutId)
	}

	return &types.MsgTxoutConfirmationVoterResponse{}, nil
}
