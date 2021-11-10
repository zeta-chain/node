package keeper

import (
	"context"
	"fmt"
	"github.com/Meta-Protocol/metacore/common"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/rand"
	"strconv"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SendVoter(goCtx context.Context, msg *types.MsgSendVoter) (*types.MsgSendVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	chain, err := common.ParseChain(msg.ReceiverChain)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("cannot parse chain %s", msg.ReceiverChain))
	}
	nonce, isFound := k.GetChainNonces(ctx, chain.String())
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("no chain nonce"))
	}
	index := msg.Digest()
	send, isFound := k.GetSend(ctx, index)
	if isFound { // send exists; add creator to signers
		send.Signers = append(send.Signers, msg.Creator)
	} else {
		send = types.Send{
			Creator:             msg.Creator,
			Index:               index,
			Sender:              msg.Sender,
			SenderChain:         msg.SenderChain,
			Receiver:            msg.Receiver,
			ReceiverChain:       msg.ReceiverChain,
			MBurnt:              msg.MBurnt,
			MMint:               msg.MMint,
			Message:             msg.Message,
			InTxHash:            msg.InTxHash,
			InBlockHeight:       msg.InBlockHeight,
			FinalizedMetaHeight: 0,
			Signers:             []string{msg.Creator},
			Status:              types.SendStatus_Created,
			Nonce:               0,
		}
	}

	//TODO: proper super majority needed
	if len(send.Signers) == 2 {
		send.FinalizedMetaHeight = uint64(ctx.BlockHeader().Height)
		send.Status = types.SendStatus_Finalized
		lastblock, isFound := k.GetLastBlockHeight(ctx, msg.SenderChain)
		if !isFound {
			lastblock = types.LastBlockHeight{
				Creator:           msg.Creator,
				Index:             msg.SenderChain,
				Chain:             msg.SenderChain,
				LastSendHeight:    msg.InBlockHeight,
				LastReceiveHeight: 0,
			}
		} else {
			lastblock.LastSendHeight = msg.InBlockHeight
		}
		k.SetLastBlockHeight(ctx, lastblock)

		send.Broadcaster = uint64(rand.Intn(len(send.Signers)))
		gasPrice, isFound := k.GetGasPrice(ctx, send.ReceiverChain)
		if !isFound {
			return  nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("no gas price found: chain %s", send.ReceiverChain))
		}
		price := float64(gasPrice.Median) * 1.5 // 1.5x Median gas; in wei
		send.GasPrice = fmt.Sprintf("%.0f", price)
		gasLimit := float64(90_000) //TODO: let user supply this
		exchangeRate := 1.0 // Zeta/ETH ratio; TODO: this information should come from oracle or onchain pool.
		gasFeeInZeta := price*gasLimit * exchangeRate
		mBurnt, err  := strconv.ParseFloat(send.MBurnt, 64)
		if err != nil {
			return  nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("MBurnt parse error %s", send.MBurnt))
		}
		if gasFeeInZeta > mBurnt {
			//TODO: this send should be garbage collected
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("burnt ZETA not enough to pay gas on %s: fee %.0f, burnt %.0f", send.ReceiverChain, gasFeeInZeta, mBurnt))
		}
		send.MMint = fmt.Sprintf("%.0f", mBurnt - gasFeeInZeta)

		send.Nonce = nonce.Nonce
		nonce.Nonce++
		k.SetChainNonces(ctx, nonce)


	}

	k.SetSend(ctx, send)

	return &types.MsgSendVoterResponse{}, nil
}
