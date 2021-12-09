package keeper

import (
	"context"
	"fmt"
	"github.com/Meta-Protocol/metacore/common"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"math/big"
	"strconv"

	"github.com/Meta-Protocol/metacore/x/metacore/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SendVoter(goCtx context.Context, msg *types.MsgSendVoter) (*types.MsgSendVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !isBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	index := msg.Digest()
	send, isFound := k.GetSend(ctx, index)

	if isDuplicateSigner(msg.Creator, send.Signers) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrorInvalidSigner, fmt.Sprintf("signer %s double signing!!", msg.Creator))
	}

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
			RecvHash:            "",
			IndexTxList:         -1,
		}
	}

	if hasSuperMajorityValidators(len(send.Signers), validators) {
		inTx, _ := k.GetInTx(ctx, send.InTxHash)
		inTx.Index = send.InTxHash
		inTx.SendHash = send.Index
		k.SetInTx(ctx, inTx)

		tx := &types.Tx{
			SendHash:   send.Index,
			RecvHash:   "",
			InTxHash:   send.InTxHash,
			InTxChain:  send.SenderChain,
			OutTxHash:  "",
			OutTxChain: "",
		}
		txList, found := k.GetTxList(ctx)
		if !found {
			txList = types.TxList{
				Creator: "",
				Tx:      []*types.Tx{tx},
			}
		} else {
			txList.Tx = append(txList.Tx, tx)
			send.IndexTxList = int64(len(txList.Tx) - 1)
		}
		k.SetTxList(ctx, txList)

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
		{
			gasPrice, isFound := k.GetGasPrice(ctx, recvChain.String())
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
			mBurnt, ok := big.NewInt(0).SetString(send.MBurnt, 10)
			if !ok {
				send.StatusMessage = fmt.Sprintf("MBurnt cannot parse")
				send.Status = types.SendStatus_Aborted
				goto END
			}
			mMint, ok := big.NewInt(0).SetString(send.MMint, 10)
			if !ok {
				send.StatusMessage = fmt.Sprintf("MMint cannot parse")
				send.Status = types.SendStatus_Aborted
				goto END
			}
			gasFee := big.NewInt(int64(gasFeeInZeta))
			toMint := big.NewInt(0).Sub(mBurnt, gasFee)
			if toMint.Cmp(mMint) < 0 { // not enough burnt
				abort = true
			}
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
		if abort {
			send.MMint = fmt.Sprintf("%.0f", mBurnt-gasFeeInZeta)
		} // if not abort, then MMint is small enough that we can mint.

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
