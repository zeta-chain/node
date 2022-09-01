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
	var cctx types.CrossChainTx
	cctx, isFound := k.GetCrossChainTx(ctx, index)

	// Validate
	recvChain, _ := parseChainAndAddress(cctx.OutBoundTxParams.ReceiverChain, cctx.OutBoundTxParams.Receiver)
	//if err != nil {
	//	send.StatusMessage = err.Error()
	//	send.Status = types.SendStatus_PendingRevert
	//	abort = true
	//}
	//
	//var chain common.Chain // the chain for outbound
	//if abort {
	//	chain, err = common.ParseChain(send.SenderChain)
	//	if err != nil {
	//		send.StatusMessage = fmt.Sprintf("cannot parse sender chain: %s", send.SenderChain)
	//		send.Status = types.SendStatus_Aborted
	//		goto EPILOGUE
	//	}
	//	send.Status = types.SendStatus_PendingRevert
	//} else {
	//	chain = recvChain
	//}

	if isFound {
		if isDuplicateSigner(msg.Creator, cctx.Signers) {
			return nil, sdkerrors.Wrap(types.ErrDuplicateMsg, fmt.Sprintf("signer %s double signing!!", msg.Creator))
		}
		cctx.Signers = append(cctx.Signers, msg.Creator)
	} else {
		cctx = k.createNewCCTX(ctx, msg, index)
	}

	hasEnoughVotes := k.hasSuperMajorityValidators(ctx, cctx.Signers)
	if hasEnoughVotes {
		err := k.finalizeInbound(ctx, cctx, recvChain.String())
		if err != nil {
			cctx.CctxStatus.Status = types.CctxStatus_Aborted
			cctx.CctxStatus.StatusMessage = err.Error()
			ctx.Logger().Error(err.Error())
		}
	}
	k.SetCrossChainTx(ctx, cctx)
	return &types.MsgVoteOnObservedInboundTxResponse{}, nil
}

// updates gas price, gas fee, zeta to mint, and nonce
// returns ok?

func (k msgServer) finalizeInbound(ctx sdk.Context, cctx types.CrossChainTx, receiveChain string) error {
	cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()
	cctx.InBoundTxParams.InBoundTxFinalizedHeight = uint64(ctx.BlockHeader().Height)

	cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
	k.UpdateLastBlockHeight(ctx, &cctx)
	bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
	cctx.OutBoundTxParams.Broadcaster = uint64(bftTime.Nanosecond() % len(cctx.Signers))

	err := k.updateCctx(ctx, receiveChain, &cctx)
	if err != nil {
		return err
	}
	k.EmitEventSendFinalized(ctx, &send)
	return err
}
func (k msgServer) createNewCCTX(ctx sdk.Context, msg *types.MsgVoteOnObservedInboundTx, index string) types.CrossChainTx {
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
func (k msgServer) updateCctx(ctx sdk.Context, receiveChain string, cctx *types.CrossChainTx) error {
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, receiveChain)
	if !isFound {
		cctx.CctxStatus.Status = types.CctxStatus_Aborted
		return sdkerrors.Wrap(types.ErrUnableToGetGasPrice, fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutBoundTxParams.ReceiverChain, cctx.LogIdentifierForCCTX()))
	}
	gasLimit := sdk.NewUint(cctx.OutBoundTxParams.OutBoundTxGasLimit)

	gasFeeInZeta, err := k.computeFeeInZeta(ctx, medianGasPrice, gasLimit, receiveChain, cctx)
	if err != nil {
		return err
	}
	// Check parse for Uint zetaBurntInt in validate Basic
	cctx.OutBoundTxParams.OutBoundTxGasPrice = medianGasPrice.String()
	zetaBurnt := sdk.NewUintFromString(cctx.ZetaBurnt)
	//zetaBurntInt, ok := big.NewInt(0).SetString(send.ZetaBurnt, 0)
	//if !ok {
	//	send.StatusMessage = fmt.Sprintf("ZetaBurnt cannot parse")
	//	send.Status = types.SendStatus_Aborted
	//	return false
	//}
	if gasFeeInZeta.GT(zetaBurnt) {
		return sdkerrors.Wrap(types.ErrNotEnoughZetaBurnt, fmt.Sprintf("feeInZeta(%s) more than mBurnt (%s) | Identifiers : %s ", gasFeeInZeta, zetaBurnt, cctx.LogIdentifierForCCTX()))
	}
	//if gasFeeInZeta.Cmp(zetaBurntInt) > 0 {
	//
	//}

	cctx.ZetaMint = zetaBurnt.Sub(gasFeeInZeta).String()
	//send.ZetaMint = fmt.Sprintf("%d", big.NewInt(0).Sub(zetaBurntInt, gasFeeInZeta))

	nonce, found := k.GetChainNonces(ctx, receiveChain)
	if !found {
		return sdkerrors.Wrap(types.ErrCannotFindReceiverNonce, fmt.Sprintf("Chain(%s) | Identifiers : %s ", receiveChain, cctx.LogIdentifierForCCTX()))
	}
	cctx.OutBoundTxParams.OutBoundTxTSSNonce = nonce.Nonce
	nonce.Nonce++
	k.SetChainNonces(ctx, nonce)
	return nil
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

func (k msgServer) EmitEventCCTXCreated(ctx sdk.Context, cctx *types.CrossChainTx) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.InboundCreated,
			sdk.NewAttribute(types.CctxIndex, cctx.Index),
			sdk.NewAttribute(types.Sender, cctx.InBoundTxParams.Sender),
			sdk.NewAttribute(types.SenderChain, cctx.InBoundTxParams.SenderChain),
			sdk.NewAttribute(types.InTxHash, cctx.InBoundTxParams.InBoundTxObservedHash),
			sdk.NewAttribute(types.Receiver, cctx.OutBoundTxParams.Receiver),
			sdk.NewAttribute(types.ReceiverChain, cctx.OutBoundTxParams.ReceiverChain),
			sdk.NewAttribute(types.ZetaBurnt, cctx.ZetaBurnt),
			sdk.NewAttribute(types.NewStatus, cctx.CctxStatus.String()),
			sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
		),
	)
}

func (k msgServer) EmitEventSendFinalized(ctx sdk.Context, cctx *types.CrossChainTx) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.InboundFinalized,
			sdk.NewAttribute(types.CctxIndex, cctx.Index),
			sdk.NewAttribute(types.Sender, cctx.InBoundTxParams.Sender),
			sdk.NewAttribute(types.SenderChain, cctx.InBoundTxParams.SenderChain),
			sdk.NewAttribute(types.InTxHash, cctx.InBoundTxParams.InBoundTxObservedHash),
			sdk.NewAttribute(types.InBlockHeight, fmt.Sprintf("%d", cctx.InBoundTxParams.InBoundTxObservedHeight)),
			sdk.NewAttribute(types.Receiver, cctx.OutBoundTxParams.Receiver),
			sdk.NewAttribute(types.ReceiverChain, cctx.OutBoundTxParams.ReceiverChain),
			sdk.NewAttribute(types.ZetaBurnt, cctx.ZetaBurnt),
			sdk.NewAttribute(types.ZetaMint, cctx.ZetaMint),
			sdk.NewAttribute(types.RelayedMessage, cctx.RelayedMessage),
			sdk.NewAttribute(types.NewStatus, cctx.CctxStatus.Status.String()),
			sdk.NewAttribute(types.StatusMessage, cctx.CctxStatus.StatusMessage),
			sdk.NewAttribute(types.Identifiers, cctx.LogIdentifierForCCTX()),
		),
	)
}

// returns feeInZeta (uint uuzeta), and whether to abort zeta-tx
// TODO : Add unit test
func (k msgServer) computeFeeInZeta(ctx sdk.Context, price sdk.Uint, gasLimit sdk.Uint, receiveChain string, cctx *types.CrossChainTx) (sdk.Uint, error) {

	rate, isFound := k.GetZetaConversionRate(ctx, receiveChain)
	if !isFound {
		return sdk.ZeroUint(), sdkerrors.Wrap(types.ErrUnableToGetConversionRate, fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutBoundTxParams.ReceiverChain, cctx.LogIdentifierForCCTX()))
	}
	//exchangeRateInt, ok := big.NewInt(0).SetString(rate.ZetaConversionRates[rate.MedianIndex], 0)
	//if !ok {
	//	return sdk.ZeroUint(),sdkerrors.Wrap(types.ErrFloatParseError,fmt.Sprintf("median exchange rate %s |Identifiers : %s ", rate.ZetaConversionRates[rate.MedianIndex],cctx.LogIdentifierForCCTX()))
	//}
	medianRate := rate.ZetaConversionRates[rate.MedianIndex]
	uintMedianRate := sdk.NewUintFromString(medianRate)
	staticValue := sdk.NewUintFromString("1000000000000000000")
	gasFeeInZeta := price.Mul(gasLimit).Mul(uintMedianRate).Quo(staticValue)
	gasFeeInZetaIncludingProtocolFee := gasFeeInZeta.Add(staticValue)
	// price*gasLimit*exchangeRate/1e18
	//gasFeeInZeta := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(0).Mul(price, gasLimit), exchangeRateInt), OneEighteen)
	// add protocol flat fee: 1 ZETA
	//gasFeeInZeta.Add(gasFeeInZeta, OneEighteen)
	return gasFeeInZetaIncludingProtocolFee, nil
}
