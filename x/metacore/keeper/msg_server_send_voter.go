package keeper

import (
	"context"
	"fmt"
	"github.com/Meta-Protocol/metacore/common"
	"strconv"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SendVoter(goCtx context.Context, msg *types.MsgSendVoter) (*types.MsgSendVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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
		bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
		send.Broadcaster = uint64(bftTime.Nanosecond() % len(send.Signers))

		// validate; abort if  failed
		abort := false
		recvChain, err := common.ParseChain(send.ReceiverChain)
		if err != nil {
			abort = true
		}
		_, err = common.NewAddress(send.Receiver, recvChain)
		if err != nil {
			abort = true
		}

		var chain common.Chain // the chain for outbound
		if abort {
			chain, err = common.ParseChain(send.SenderChain)
			if err != nil {
				send.StatusMessage = fmt.Sprintf("cannot parse sender chain: %s", send.SenderChain)
				send.Status = types.SendStatus_Aborted
				goto END
			}
			send.Status = types.SendStatus_Abort
			send.StatusMessage = fmt.Sprintf("cannot parse recever address/chain")
		} else {
			chain = recvChain
		}
		gasPrice, isFound := k.GetGasPrice(ctx, chain.String())
		if !isFound {
			send.StatusMessage = fmt.Sprintf("no gas price found: chain %s", send.ReceiverChain)
			send.Status = types.SendStatus_Aborted
			goto END
		}
		mi := gasPrice.MedianIndex
		medianPrice := gasPrice.Prices[mi]
		price := float64(medianPrice) * 1.5 // 1.5x Median gas; in wei
		send.GasPrice = fmt.Sprintf("%.0f", price)
		gasLimit := float64(90_000) //TODO: let user supply this
		exchangeRate := 1.0         // Zeta/ETH ratio; TODO: this information should come from oracle or onchain pool.
		gasFeeInZeta := price * gasLimit * exchangeRate
		mBurnt, err := strconv.ParseFloat(send.MBurnt, 64)
		if err != nil {
			send.StatusMessage = fmt.Sprintf("MBurnt cannot parse")
			send.Status = types.SendStatus_Aborted
			goto END
		}
		if gasFeeInZeta > mBurnt {
			send.StatusMessage = fmt.Sprintf("feeInZeta(%f) more than mBurnt (%f)", gasFeeInZeta, mBurnt)
			send.Status = types.SendStatus_Aborted
			goto END
		}
		send.MMint = fmt.Sprintf("%.0f", mBurnt-gasFeeInZeta)

		nonce, found := k.GetChainNonces(ctx, chain.String())
		if !found {
			send.StatusMessage = fmt.Sprintf("cannot find receiver chain nonce: %s", chain)
			send.Status = types.SendStatus_Aborted
			goto END
		}

		send.Nonce = nonce.Nonce
		nonce.Nonce++
		k.SetChainNonces(ctx, nonce)

	}

END:
	k.SetSend(ctx, send)
	return &types.MsgSendVoterResponse{}, nil
}
