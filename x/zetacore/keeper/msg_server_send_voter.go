package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	"math/big"
)

func (k msgServer) SendVoter(goCtx context.Context, msg *types.MsgSendVoter) (*types.MsgSendVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	validators := k.StakingKeeper.GetAllValidators(ctx)
	if !isBondedValidator(msg.Creator, validators) {
		return nil, sdkerrors.Wrap(types.ErrNotBondedValidator, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
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
			ZetaBurnt:           msg.ZetaBurnt,
			Message:             msg.Message,
			InTxHash:            msg.InTxHash,
			InBlockHeight:       msg.InBlockHeight,
			FinalizedMetaHeight: 0,
			Signers:             []string{msg.Creator},
			Status:              types.SendStatus_PendingInbound,
			Nonce:               0,
			RecvHash:            "",
			IndexTxList:         -1,
			LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
			ZetaMint:            "",
		}
		k.EmitEventSendCreated(ctx, &send)
	}

	if hasSuperMajorityValidators(len(send.Signers), validators) {
		send.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
		k.UpdateTxList(ctx, &send)

		send.FinalizedMetaHeight = uint64(ctx.BlockHeader().Height)
		send.Status = types.SendStatus_PendingOutbound
		k.UpdateLastBlockHeigh(ctx, msg)

		bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
		send.Broadcaster = uint64(bftTime.Nanosecond() % len(send.Signers))

		abort := false
		// validate receiver address & chain; abort if failed
		recvChain, err := parseChainAndAddress(send.ReceiverChain, send.Receiver)
		if err != nil {
			send.StatusMessage = err.Error()
			send.Status = types.SendStatus_PendingRevert
			abort = true
		}

		var chain common.Chain // the chain for outbound
		if abort {
			chain, err = common.ParseChain(send.SenderChain)
			if err != nil {
				send.StatusMessage = fmt.Sprintf("cannot parse sender chain: %s", send.SenderChain)
				send.Status = types.SendStatus_Aborted
				goto EPILOGUE
			}
			send.Status = types.SendStatus_PendingRevert
		} else {
			chain = recvChain
		}

		k.updateSend(ctx, chain.String(), &send)
		k.EmitEventSendFinalized(ctx, &send)
	}

EPILOGUE:
	k.SetSend(ctx, send)
	return &types.MsgSendVoterResponse{}, nil
}

// updates gas price, gas fee, zeta to mint, and nonce
// returns ok?
func (k msgServer) updateSend(ctx sdk.Context, chain string, send *types.Send) bool {
	gasPrice, isFound := k.GetGasPrice(ctx, chain)
	if !isFound {
		send.StatusMessage = fmt.Sprintf("no gas price found: chain %s", send.ReceiverChain)
		send.Status = types.SendStatus_Aborted
		return false
	}
	mi := gasPrice.MedianIndex
	medianPrice := gasPrice.Prices[mi]
	send.GasPrice = fmt.Sprintf("%d", medianPrice)
	price := float64(medianPrice)
	gasLimit := float64(250_000) //TODO: let user supply this
	gasFeeInZeta, abort := k.computeFeeInZeta(ctx, price, gasLimit, chain, send)
	if abort {
		send.Status = types.SendStatus_Aborted
		return false
	}
	zetaBurntInt, ok := big.NewInt(0).SetString(send.ZetaBurnt, 0)
	if !ok {
		send.StatusMessage = fmt.Sprintf("ZetaBurnt cannot parse")
		send.Status = types.SendStatus_Aborted
		return false
	}
	if gasFeeInZeta.Cmp(zetaBurntInt) > 0 {
		send.StatusMessage = fmt.Sprintf("feeInZeta(%d) more than mBurnt (%d)", gasFeeInZeta, zetaBurntInt)
		send.Status = types.SendStatus_Aborted
		return false
	}
	send.ZetaMint = fmt.Sprintf("%d", big.NewInt(0).Sub(zetaBurntInt, gasFeeInZeta))

	nonce, found := k.GetChainNonces(ctx, chain)
	if !found {
		send.StatusMessage = fmt.Sprintf("cannot find receiver chain nonce: %s", chain)
		send.Status = types.SendStatus_Aborted
		return false
	}

	send.Nonce = nonce.Nonce
	nonce.Nonce++
	k.SetChainNonces(ctx, nonce)
	return true
}

// returns (chain,error)
// chain: the receiverChain if ok
func parseChainAndAddress(chain string, addr string) (common.Chain, error) {
	recvChain, err := common.ParseChain(chain)
	if err != nil {
		return recvChain, fmt.Errorf("cannot parse receiver chain %s", chain)
	}
	_, err = common.NewAddress(addr, recvChain)
	if err != nil {
		return recvChain, fmt.Errorf("cannot parse receiver addr %s", addr)
	}
	return recvChain, nil
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
			sdk.NewAttribute(types.SubType, types.InboundCreated),
			sdk.NewAttribute(types.SendHash, send.Index),
			//sdk.NewAttribute(types.Sender, send.Sender),
			//sdk.NewAttribute(types.SenderChain, send.SenderChain),
			//sdk.NewAttribute(types.Receiver, send.Receiver),
			//sdk.NewAttribute(types.ReceiverChain, send.ReceiverChain),
			//sdk.NewAttribute(types.ZetaBurnt, send.ZetaBurnt),
			//sdk.NewAttribute(types.Message, send.Message),
			//sdk.NewAttribute(types.InTxHash, send.InTxHash),
			//sdk.NewAttribute(types.InBlockHeight, fmt.Sprintf("%d", send.InBlockHeight)),
			sdk.NewAttribute(types.NewStatus, send.Status.String()),
		),
	)
}

func (k msgServer) EmitEventSendFinalized(ctx sdk.Context, send *types.Send) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
			sdk.NewAttribute(types.SubType, types.InboundFinalized),
			sdk.NewAttribute(types.SendHash, send.Index),
			sdk.NewAttribute(types.Sender, send.Sender),
			sdk.NewAttribute(types.SenderChain, send.SenderChain),
			sdk.NewAttribute(types.Receiver, send.Receiver),
			sdk.NewAttribute(types.ReceiverChain, send.ReceiverChain),
			sdk.NewAttribute(types.ZetaBurnt, send.ZetaBurnt),
			sdk.NewAttribute(types.Message, send.Message),
			sdk.NewAttribute(types.InTxHash, send.InTxHash),
			sdk.NewAttribute(types.InBlockHeight, fmt.Sprintf("%d", send.InBlockHeight)),
			sdk.NewAttribute(types.NewStatus, send.Status.String()),
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

// returns feeInZeta (uint uuzeta), and whether to abort zeta-tx
func (k msgServer) computeFeeInZeta(ctx sdk.Context, price float64, gasLimit float64, chain string, send *types.Send) (*big.Int, bool) {
	abort := false
	rate, isFound := k.GetZetaConversionRate(ctx, chain)
	if !isFound {
		send.StatusMessage = fmt.Sprintf("Zeta conversion rate not found")
		abort = true
	}
	exchangeRateInt, ok := big.NewInt(0).SetString(rate.ZetaConversionRates[rate.MedianIndex], 0)
	if !ok {
		send.StatusMessage = fmt.Sprintf("median exchange rate %s cannot parse into float", rate.ZetaConversionRates[rate.MedianIndex])
		abort = true
	}
	exchangeRateFloat, _ := big.NewFloat(0).SetInt(exchangeRateInt).Float64()
	exchangeRateFloat = exchangeRateFloat / 1.0e18 // 18 decimals

	gasFeeInZeta := price * gasLimit * exchangeRateFloat
	gasFeeInZetaInt, _ := big.NewFloat(0).SetFloat64(gasFeeInZeta).Int(nil)
	return gasFeeInZetaInt, abort
}
