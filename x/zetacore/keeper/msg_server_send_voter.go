package keeper

import (
	"context"
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
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
		return nil, sdkerrors.Wrap(types.ErrDuplicateMsg, fmt.Sprintf("signer %s double signing!!", msg.Creator))
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
			LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
		}
		k.EmitEventSendCreated(ctx, &send)
	}

	if hasSuperMajorityValidators(len(send.Signers), validators) {
		send.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
		k.UpdateTxList(ctx, &send)

		send.FinalizedMetaHeight = uint64(ctx.BlockHeader().Height)
		send.Status = types.SendStatus_Finalized
		k.UpdateLastBlockHeigh(ctx, msg)

		bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
		send.Broadcaster = uint64(bftTime.Nanosecond() % len(send.Signers))

		// validate receiver address & chain; abort if failed
		recvChain, abort := k.validateReceiver(&send)
		// quick check: is enough zeta burnt to cover destination gas fee?
		abort = k.validate(ctx, recvChain, send, abort)

		k.processSend(ctx, abort, &send, recvChain)
	}

END:
	k.SetSend(ctx, send)
	return &types.MsgSendVoterResponse{}, nil
}

func (k msgServer) validate(ctx sdk.Context, recvChain common.Chain, send *types.Send, abort bool) bool {
	{
		price, ok := k.findGasPrice(ctx, recvChain, &send)
		if !ok {

		}
		send.GasPrice = fmt.Sprintf("%.0f", price)
		gasLimit := float64(250_000) //TODO: let user supply this
		var gasFeeInZeta float64
		gasFeeInZeta, abort = k.computeFeeInZeta(ctx, price, gasLimit, send.ReceiverChain, &send)

		mBurnt, ok := big.NewInt(0).SetString(send.MBurnt, 10)
		if !ok {
			send.StatusMessage = fmt.Sprintf("MBurnt cannot parse")
			send.Status = types.SendStatus_Aborted
			//goto END
		}
		mMint, ok := big.NewInt(0).SetString(send.MMint, 10)
		if !ok {
			send.StatusMessage = fmt.Sprintf("MMint cannot parse")
			send.Status = types.SendStatus_Aborted
			//goto END
		}
		gasFee := big.NewInt(int64(gasFeeInZeta))
		toMint := big.NewInt(0).Sub(mBurnt, gasFee)
		if toMint.Cmp(mMint) < 0 { // not enough burnt
			abort = true
			send.StatusMessage = fmt.Sprintf("wanted %d, but can only mint %d", mMint, toMint)
		}
	}
	return abort
}

func (k msgServer) processSend(ctx sdk.Context, abort bool, send *types.Send, recvChain common.Chain) {
	var chain common.Chain // the chain for outbound
	var err error
	if abort {
		chain, err = common.ParseChain(send.SenderChain)
		if err != nil {
			send.StatusMessage = fmt.Sprintf("cannot parse sender chain: %s", send.SenderChain)
			send.Status = types.SendStatus_Aborted
			return
		}
		send.Status = types.SendStatus_Revert
	} else {
		chain = recvChain
	}
	gasPrice, isFound := k.GetGasPrice(ctx, chain.String())
	if !isFound {
		send.StatusMessage = fmt.Sprintf("no gas price found: chain %s", send.ReceiverChain)
		send.Status = types.SendStatus_Aborted
		return
	}
	mi := gasPrice.MedianIndex
	medianPrice := gasPrice.Prices[mi]
	price := float64(medianPrice)
	send.GasPrice = fmt.Sprintf("%.0f", price)
	gasLimit := float64(250_000) //TODO: let user supply this
	exchangeRate := 1.0          // Zeta/ETH ratio; TODO: this information should come from oracle or onchain pool.
	gasFeeInZeta := price * gasLimit * exchangeRate
	mBurnt, err := strconv.ParseFloat(send.MBurnt, 64)
	if err != nil {
		send.StatusMessage = fmt.Sprintf("MBurnt cannot parse")
		send.Status = types.SendStatus_Aborted
		return
	}
	if gasFeeInZeta > mBurnt {
		send.StatusMessage = fmt.Sprintf("feeInZeta(%f) more than mBurnt (%f)", gasFeeInZeta, mBurnt)
		send.Status = types.SendStatus_Aborted
		return
	}
	if abort {
		send.MMint = fmt.Sprintf("%.0f", mBurnt-gasFeeInZeta)
	} // if not abort, then MMint is small enough that we can mint.

	nonce, found := k.GetChainNonces(ctx, chain.String())
	if !found {
		send.StatusMessage = fmt.Sprintf("cannot find receiver chain nonce: %s", chain)
		send.Status = types.SendStatus_Aborted
		return
	}

	send.Nonce = nonce.Nonce
	nonce.Nonce++
	k.SetChainNonces(ctx, nonce)
}

// returns (gas price in wei per unit gas, and ok?
func (k msgServer) findGasPrice(ctx sdk.Context, recvChain common.Chain, send *types.Send) (float64, bool) {
	gasPrice, isFound := k.GetGasPrice(ctx, recvChain.String())
	if !isFound {
		send.StatusMessage = fmt.Sprintf("no gas price found: chain %s", send.ReceiverChain)
		send.Status = types.SendStatus_Aborted
		return 0, false
	}
	mi := gasPrice.MedianIndex
	medianPrice := gasPrice.Prices[mi]
	price := float64(medianPrice)
	return price, true
}

// returns abort?
func (k msgServer) validateReceiver(send *types.Send) (common.Chain, bool) {
	recvChain, err := common.ParseChain(send.ReceiverChain)
	if err != nil {
		send.StatusMessage = fmt.Sprintf("cannot parse receiver chain")
		return recvChain, true
	}
	_, err = common.NewAddress(send.Receiver, recvChain)
	if err != nil {
		send.StatusMessage = fmt.Sprintf("cannot parse receiver address")
		return recvChain, true
	}
	return recvChain, false
}

func (k msgServer) UpdateLastBlockHeigh(ctx sdk.Context, msg *types.MsgSendVoter) {
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
}

func (k msgServer) EmitEventSendCreated(ctx sdk.Context, send *types.Send) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
			sdk.NewAttribute(sdk.AttributeKeyAction, "SendCreated"),
			sdk.NewAttribute("Index", send.Index),
			sdk.NewAttribute("Sender", send.Sender),
			sdk.NewAttribute("SenderChain", send.SenderChain),
			sdk.NewAttribute("Receiver", send.Receiver),
			sdk.NewAttribute("ReceiverChain", send.ReceiverChain),
			sdk.NewAttribute("MBurnt", send.MBurnt),
			sdk.NewAttribute("MMint", send.MMint),
			sdk.NewAttribute("Message", send.Message),
			sdk.NewAttribute("InTxHash", send.InTxHash),
			sdk.NewAttribute("InBlockHeight", fmt.Sprintf("%d", send.InBlockHeight)),
		),
	)
}

func (k msgServer) UpdateTxList(ctx sdk.Context, send *types.Send) {
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
}

// returns feeInZeta, and whether to abort zeta-tx
func (k msgServer) computeFeeInZeta(ctx sdk.Context, price float64, gasLimit float64, chain string, send *types.Send) (float64, bool) {
	abort := false
	rate, isFound := k.GetZetaConversionRate(ctx, chain)
	if !isFound {
		send.StatusMessage = fmt.Sprintf("Zeta conversion rate not found")

	}
	exchangeRate, err := strconv.ParseFloat(rate.ZetaConversionRates[rate.MedianIndex], 64)
	if err != nil {
		send.StatusMessage = fmt.Sprintf("median exchange rate %s cannot parse into float", rate.ZetaConversionRates[rate.MedianIndex])
		abort = true
	}
	gasFeeInZeta := price * gasLimit * exchangeRate

	return gasFeeInZeta, abort
}
