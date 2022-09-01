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

var (
	OneEighteen, _ = big.NewInt(0).SetString("1000000000000000000", 10)
)

func (k msgServer) SendVoter(goCtx context.Context, msg *types.MsgVoteOnObservedInboundTx) (*types.MsgVoteOnObservedInboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.isAuthorized(ctx, msg.Creator) {
		return nil, sdkerrors.Wrap(types.ErrNotBondedValidator, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}
	index := msg.Digest()
	var cctx  types.CrossChainTx
	cctx, isFound := k.GetCrossChainTx(ctx, index)


	// Validate
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

	if isFound {
		if isDuplicateSigner(msg.Creator, cctx.Signers) {
			return nil, sdkerrors.Wrap(types.ErrDuplicateMsg, fmt.Sprintf("signer %s double signing!!", msg.Creator))
		}
		cctx.Signers = append(cctx.Signers, msg.Creator)
	} else {
		cctx = k.createNewCCTX(ctx,msg,index)
	}

	hasEnoughVotes := k. hasSuperMajorityValidators(ctx,cctx.Signers)
	if hasEnoughVotes {

	}




	}

EPILOGUE:
	k.SetCrossChainTx(ctx, send)
	return &types.MsgSendVoterResponse{}, nil
}

// updates gas price, gas fee, zeta to mint, and nonce
// returns ok?

func (k msgServer)finalizeInbound(ctx sdk.Context,cctx types.CrossChainTx,receiveChain string) {
	cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
	cctx.InBoundTxParams.InBoundTxFinalizedHeight = uint64(ctx.BlockHeader().Height)

	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	k.UpdateLastBlockHeight(ctx, &cctx)
	bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
	cctx.OutBoundTxParams.Broadcaster = uint64(bftTime.Nanosecond() % len(cctx.Signers))

	k.updateSend(ctx,receiveChain, &cctx)
	k.EmitEventSendFinalized(ctx, &send)

}
func (k msgServer) createNewCCTX (ctx sdk.Context,msg *types.MsgVoteOnObservedInboundTx , index string) types.CrossChainTx {
	inboundParams := &types.InBoundTxParams{
		Sender:                   msg.Sender,
		SenderChain:              msg.SenderChain,
		InBoundTxObservedHash:    msg.InTxHash,
		InBoundTxObservedHeight:  msg.InBlockHeight,
		InBoundTxFinalizedHeight: 0,
	}

	outBoundParams := &types.OutBoundTxParams{
		Receiver:               msg.Receiver,
		ReceiverChain:          msg.ReceiverChain,
		Broadcaster:            0,
		OutBoundTxHash:         "",
		OutBoundTxTSSNonce:     0,
		OutBoundTxGasLimit:     msg.GasLimit,
		OutBoundTxGasPrice:     "",
		OutBoundTXReceiveIndex: "",
	}
	status := &types.Status{
		Status:              types.CctxStatus_PendingInbound,
		StatusMessage:       "",
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
	}
	newCctx := types.CrossChainTx{
		Creator:          msg.Creator,
		Index:            index,
		ZetaBurnt:        msg.ZetaBurnt,
		ZetaMint:         "",
		RelayedMessage:   msg.Message,
		Signers:          []string{msg.Creator},
		CctxStatus:       status,
		InBoundTxParams:  inboundParams,
		OutBoundTxParams: outBoundParams,
	}
	k.EmitEventCCTXCreated(ctx, &newCctx)
	return newCctx
}
func (k msgServer) updateSend(ctx sdk.Context, receiveChain string, cctx *types.CrossChainTx) error {
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, receiveChain)
	if !isFound {
		cctx.CctxStatus.Status = types.CctxStatus_Aborted
		return sdkerrors.Wrap(types.ErrUnableToGetGasPrice,fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutBoundTxParams.ReceiverChain ,cctx.LogIdentifierForCCTX()))
	}
	gasLimit := sdk.NewUint(cctx.OutBoundTxParams.OutBoundTxGasLimit)

	gasFeeInZeta, err := k.computeFeeInZeta(ctx, price, gasLimit, chain, send)
	if err!=nil{
		cctx.CctxStatus.Status = types.CctxStatus_Aborted
		return err
	}

	cctx.OutBoundTxParams.OutBoundTxGasPrice = medianGasPrice.String()
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

func (k msgServer) UpdateLastBlockHeight(ctx sdk.Context, msg *types.CrossChainTx) {
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

func (k msgServer) EmitEventCCTXCreated(ctx sdk.Context, send *types.CrossChainTx) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.SubTypeKey, types.InboundCreated),
			sdk.NewAttribute(types.SendHash, send.Index),
			sdk.NewAttribute(types.NewStatus, send.CctxStatus.String()),
		),
	)
}

func (k msgServer) EmitEventSendFinalized(ctx sdk.Context, send *types.Send) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "zetacore"),
			sdk.NewAttribute(types.SubTypeKey, types.InboundFinalized),
			sdk.NewAttribute(types.SendHash, send.Index),
			sdk.NewAttribute(types.Sender, send.Sender),
			sdk.NewAttribute(types.SenderChain, send.SenderChain),
			sdk.NewAttribute(types.Receiver, send.Receiver),
			sdk.NewAttribute(types.ReceiverChain, send.ReceiverChain),
			sdk.NewAttribute(types.ZetaBurnt, send.ZetaBurnt),
			sdk.NewAttribute(types.ZetaMint, send.ZetaMint),
			sdk.NewAttribute(types.Message, send.Message),
			sdk.NewAttribute(types.InTxHash, send.InTxHash),
			sdk.NewAttribute(types.InBlockHeight, fmt.Sprintf("%d", send.InBlockHeight)),
			sdk.NewAttribute(types.NewStatus, send.Status.String()),
			sdk.NewAttribute(types.StatusMessage, send.StatusMessage),
		),
	)
}

// returns feeInZeta (uint uuzeta), and whether to abort zeta-tx
func (k msgServer) computeFeeInZeta(ctx sdk.Context, price sdk.Uint, gasLimit sdk.Uint, chain string, cctx *types.CrossChainTx) (sdk.Uint, error) {

	rate, isFound := k.GetZetaConversionRate(ctx, chain)
	if !isFound {
		return sdk.ZeroUint(),sdkerrors.Wrap(types.ErrUnableToGetConversionRate,fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutBoundTxParams.ReceiverChain ,cctx.LogIdentifierForCCTX()))
	}
	exchangeRateInt, ok := big.NewInt(0).SetString(rate.ZetaConversionRates[rate.MedianIndex], 0)
	if !ok {
		return sdk.ZeroUint(),sdkerrors.Wrap(types.ErrFloatParseError,fmt.Sprintf("median exchange rate %s |Identifiers : %s ", rate.ZetaConversionRates[rate.MedianIndex],cctx.LogIdentifierForCCTX()))
	}
	medianRate := rate.ZetaConversionRates[rate.MedianIndex]
	uintMedianRate := rate


	// price*gasLimit*exchangeRate/1e18
	gasFeeInZeta := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(0).Mul(price, gasLimit), exchangeRateInt), OneEighteen)
	// add protocol flat fee: 1 ZETA
	gasFeeInZeta.Add(gasFeeInZeta, OneEighteen)
	return gasFeeInZeta, abort
}
