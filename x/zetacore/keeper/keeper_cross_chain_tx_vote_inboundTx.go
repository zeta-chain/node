package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k msgServer) VoteOnObservedInboundTx(goCtx context.Context, msg *types.MsgVoteOnObservedInboundTx) (*types.MsgVoteOnObservedInboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if !k.isAuthorized(ctx, msg.Creator) {
		return nil, sdkerrors.Wrap(types.ErrNotBondedValidator, fmt.Sprintf("signer %s is not a bonded validator", msg.Creator))
	}

	index := msg.Digest()
	var cctx types.CrossChainTx
	cctx, isFound := k.GetCrossChainTx(ctx, index)
	if isFound {
		if cctx.CctxStatus.Status == types.CctxStatus_Aborted {
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}
		if isDuplicateSigner(msg.Creator, cctx.Signers) {
			return nil, sdkerrors.Wrap(types.ErrDuplicateMsg, fmt.Sprintf("signer %s double signing!!", msg.Creator))
		}
		cctx.Signers = append(cctx.Signers, msg.Creator)
	} else {
		// We can return directlu from here as new CCTX has not been created yet
		if !k.IsChainSupported(ctx, msg.ReceiverChain) || !k.IsChainSupported(ctx, msg.ReceiverChain) {
			return nil, sdkerrors.Wrap(types.ErrUnsupportedChain, "Receiving chain is not supported")
		}
		cctx = k.createNewCCTX(ctx, msg, index)
	}

	hasEnoughVotes := k.hasSuperMajorityValidators(ctx, cctx.Signers)
	if hasEnoughVotes {
		err := k.FinalizeInbound(ctx, cctx, msg.ReceiverChain)
		if err != nil {
			cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_Aborted, err.Error(), cctx.LogIdentifierForCCTX())
			ctx.Logger().Error(err.Error())
			k.SetCrossChainTx(ctx, cctx)
			return &types.MsgVoteOnObservedInboundTxResponse{}, nil
		}
		cctx.CctxStatus.ChangeStatus(&ctx, types.CctxStatus_PendingOutbound, "Status Changed to Pending Outbound", cctx.LogIdentifierForCCTX())
	}
	k.SetCrossChainTx(ctx, cctx)
	return &types.MsgVoteOnObservedInboundTxResponse{}, nil
}

func (k msgServer) FinalizeInbound(ctx sdk.Context, cctx types.CrossChainTx, receiveChain string) error {
	cctx.InBoundTxParams.InBoundTxFinalizedZetaHeight = uint64(ctx.BlockHeader().Height)
	k.UpdateLastBlockHeight(ctx, &cctx)
	bftTime := ctx.BlockHeader().Time // we use BFTTime of the current block as random number
	cctx.OutBoundTxParams.Broadcaster = uint64(bftTime.Nanosecond() % len(cctx.Signers))

	err := k.updatePrices(ctx, receiveChain, &cctx)
	if err != nil {
		return err
	}
	err = k.updateNonce(ctx, receiveChain, &cctx)
	if err != nil {
		return err
	}
	EmitEventSendFinalized(ctx, &cctx)
	return nil
}

func (k Keeper) updatePrices(ctx sdk.Context, receiveChain string, cctx *types.CrossChainTx) error {
	medianGasPrice, isFound := k.GetMedianGasPriceInUint(ctx, receiveChain)
	if !isFound {
		return sdkerrors.Wrap(types.ErrUnableToGetGasPrice, fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutBoundTxParams.ReceiverChain, cctx.LogIdentifierForCCTX()))
	}
	gasLimit := sdk.NewUint(cctx.OutBoundTxParams.OutBoundTxGasLimit)
	rate, isFound := k.GetZetaConversionRate(ctx, receiveChain)
	if !isFound {
		return sdkerrors.Wrap(types.ErrUnableToGetConversionRate, fmt.Sprintf(" chain %s | Identifiers : %s ", cctx.OutBoundTxParams.ReceiverChain, cctx.LogIdentifierForCCTX()))
	}
	medianRate := rate.ZetaConversionRates[rate.MedianIndex]
	uintmedianRate := sdk.NewUintFromString(medianRate)
	// Calculate Gas FEE
	gasFeeInZeta := CalculateFee(medianGasPrice, gasLimit, uintmedianRate)

	cctx.OutBoundTxParams.OutBoundTxGasPrice = medianGasPrice.String()

	// Set ZetaBurnt and ZetaMint
	zetaBurnt := cctx.ZetaBurnt
	if gasFeeInZeta.GT(zetaBurnt) {
		return sdkerrors.Wrap(types.ErrNotEnoughZetaBurnt, fmt.Sprintf("feeInZeta(%s) more than mBurnt (%s) | Identifiers : %s ", gasFeeInZeta, zetaBurnt, cctx.LogIdentifierForCCTX()))
	}
	cctx.ZetaMint = zetaBurnt.Sub(gasFeeInZeta)
	return nil
}
func (k Keeper) updateNonce(ctx sdk.Context, receiveChain string, cctx *types.CrossChainTx) error {
	nonce, found := k.GetChainNonces(ctx, receiveChain)
	if !found {
		return sdkerrors.Wrap(types.ErrCannotFindReceiverNonce, fmt.Sprintf("Chain(%s) | Identifiers : %s ", receiveChain, cctx.LogIdentifierForCCTX()))
	}

	// SET nonce
	cctx.OutBoundTxParams.OutBoundTxTSSNonce = nonce.Nonce
	nonce.Nonce++
	k.SetChainNonces(ctx, nonce)
	return nil
}

func (k msgServer) UpdateLastBlockHeight(ctx sdk.Context, msg *types.CrossChainTx) {
	lastblock, isFound := k.GetLastBlockHeight(ctx, msg.InBoundTxParams.SenderChain)
	if !isFound {
		lastblock = types.LastBlockHeight{
			Creator:           msg.Creator,
			Index:             msg.InBoundTxParams.SenderChain, // ?
			Chain:             msg.InBoundTxParams.SenderChain,
			LastSendHeight:    msg.InBoundTxParams.InBoundTxObservedExternalHeight,
			LastReceiveHeight: 0,
		}
	} else {
		lastblock.LastSendHeight = msg.InBoundTxParams.InBoundTxObservedExternalHeight
	}
	k.SetLastBlockHeight(ctx, lastblock)
}

func CalculateFee(price, gasLimit, rate sdk.Uint) sdk.Uint {
	gasFee := price.Mul(gasLimit).Mul(rate)
	gasFee = reducePrecision(gasFee)
	return gasFee.Add(types.GetProtocolFee())
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Utils
// These functions should always remain under private scope

func (k Keeper) createNewCCTX(ctx sdk.Context, msg *types.MsgVoteOnObservedInboundTx, index string) types.CrossChainTx {
	inboundParams := &types.InBoundTxParams{
		Sender:                          msg.Sender,
		SenderChain:                     msg.SenderChain,
		InBoundTxObservedHash:           msg.InTxHash,
		InBoundTxObservedExternalHeight: msg.InBlockHeight,
		InBoundTxFinalizedZetaHeight:    0,
	}

	outBoundParams := &types.OutBoundTxParams{
		Receiver:                         msg.Receiver,
		ReceiverChain:                    msg.ReceiverChain,
		Broadcaster:                      0,
		OutBoundTxHash:                   "",
		OutBoundTxTSSNonce:               0,
		OutBoundTxGasLimit:               msg.GasLimit,
		OutBoundTxGasPrice:               "",
		OutBoundTXReceiveIndex:           "",
		OutBoundTxFinalizedZetaHeight:    0,
		OutBoundTxObservedExternalHeight: 0,
	}
	status := &types.Status{
		Status:              types.CctxStatus_PendingInbound,
		StatusMessage:       "",
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
	}
	newCctx := types.CrossChainTx{
		Creator:          msg.Creator,
		Index:            index,
		ZetaBurnt:        sdk.NewUintFromString(msg.ZetaBurnt),
		ZetaMint:         sdk.Uint{},
		RelayedMessage:   msg.Message,
		Signers:          []string{},
		CctxStatus:       status,
		InBoundTxParams:  inboundParams,
		OutBoundTxParams: outBoundParams,
	}
	EmitEventCCTXCreated(ctx, &newCctx)
	return newCctx
}
