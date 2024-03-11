package keeper

import (
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) GetInbound(ctx sdk.Context, msg *types.MsgVoteOnObservedInboundTx) types.CrossChainTx {

	// get the latest TSS to set the TSS public key in the CCTX
	tssPub := ""
	tss, tssFound := k.zetaObserverKeeper.GetTSS(ctx)
	if tssFound {
		tssPub = tss.TssPubkey
	}
	return CreateNewCCTX(ctx, msg, msg.Digest(), tssPub, types.CctxStatus_PendingInbound, msg.SenderChainId, msg.ReceiverChain)
}

func CreateNewCCTX(
	ctx sdk.Context,
	msg *types.MsgVoteOnObservedInboundTx,
	index string,
	tssPubkey string,
	s types.CctxStatus,
	senderChainID,
	receiverChainID int64,
) types.CrossChainTx {
	if msg.TxOrigin == "" {
		msg.TxOrigin = msg.Sender
	}
	inboundParams := &types.InboundTxParams{
		Sender:                          msg.Sender,
		SenderChainId:                   senderChainID,
		TxOrigin:                        msg.TxOrigin,
		Asset:                           msg.Asset,
		Amount:                          msg.Amount,
		InboundTxObservedHash:           msg.InTxHash,
		InboundTxObservedExternalHeight: msg.InBlockHeight,
		InboundTxFinalizedZetaHeight:    0,
		InboundTxBallotIndex:            index,
	}

	outBoundParams := &types.OutboundTxParams{
		Receiver:                         msg.Receiver,
		ReceiverChainId:                  receiverChainID,
		OutboundTxHash:                   "",
		OutboundTxTssNonce:               0,
		OutboundTxGasLimit:               msg.GasLimit,
		OutboundTxGasPrice:               "",
		OutboundTxBallotIndex:            "",
		OutboundTxObservedExternalHeight: 0,
		Amount:                           sdkmath.ZeroUint(),
		TssPubkey:                        tssPubkey,
	}
	status := &types.Status{
		Status:              s,
		StatusMessage:       "",
		LastUpdateTimestamp: ctx.BlockHeader().Time.Unix(),
		IsAbortRefunded:     false,
	}
	newCctx := types.CrossChainTx{
		Creator:          msg.Creator,
		Index:            index,
		ZetaFees:         sdkmath.ZeroUint(),
		RelayedMessage:   msg.Message,
		CctxStatus:       status,
		InboundTxParams:  inboundParams,
		OutboundTxParams: []*types.OutboundTxParams{outBoundParams},
		EventIndex:       msg.EventIndex,
		CoinType:         msg.CoinType,
	}
	return newCctx
}

// UpdateNonce sets the CCTX outbound nonce to the next nonce, and updates the nonce of blockchain state.
// It also updates the PendingNonces that is used to track the unfulfilled outbound txs.
func (k Keeper) UpdateNonce(ctx sdk.Context, receiveChainID int64, cctx *types.CrossChainTx) error {
	chain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, receiveChainID)
	if chain == nil {
		return zetaObserverTypes.ErrSupportedChains
	}

	nonce, found := k.GetObserverKeeper().GetChainNonces(ctx, chain.ChainName.String())
	if !found {
		return cosmoserrors.Wrap(types.ErrCannotFindReceiverNonce, fmt.Sprintf("Chain(%s) | Identifiers : %s ", chain.ChainName.String(), cctx.LogIdentifierForCCTX()))
	}

	// SET nonce
	cctx.GetCurrentOutTxParam().OutboundTxTssNonce = nonce.Nonce
	tss, found := k.zetaObserverKeeper.GetTSS(ctx)
	if !found {
		return cosmoserrors.Wrap(types.ErrCannotFindTSSKeys, fmt.Sprintf("Chain(%s) | Identifiers : %s ", chain.ChainName.String(), cctx.LogIdentifierForCCTX()))
	}

	p, found := k.GetObserverKeeper().GetPendingNonces(ctx, tss.TssPubkey, receiveChainID)
	if !found {
		return cosmoserrors.Wrap(types.ErrCannotFindPendingNonces, fmt.Sprintf("chain_id %d, nonce %d", receiveChainID, nonce.Nonce))
	}

	// #nosec G701 always in range
	if p.NonceHigh != int64(nonce.Nonce) {
		return cosmoserrors.Wrap(types.ErrNonceMismatch, fmt.Sprintf("chain_id %d, high nonce %d, current nonce %d", receiveChainID, p.NonceHigh, nonce.Nonce))
	}

	nonce.Nonce++
	p.NonceHigh++
	k.GetObserverKeeper().SetChainNonces(ctx, nonce)
	k.GetObserverKeeper().SetPendingNonces(ctx, p)
	return nil
}

// GetRevertGasLimit returns the gas limit for the revert transaction in a CCTX
// It returns 0 if there is no error but the gas limit can't be determined from the CCTX data
func (k Keeper) GetRevertGasLimit(ctx sdk.Context, cctx *types.CrossChainTx) (uint64, error) {
	if cctx.InboundTxParams == nil {
		return 0, nil
	}

	if cctx.CoinType == common.CoinType_Gas {
		// get the gas limit of the gas token
		fc, found := k.fungibleKeeper.GetGasCoinForForeignCoin(ctx, cctx.InboundTxParams.SenderChainId)
		if !found {
			return 0, types.ErrForeignCoinNotFound
		}
		gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(fc.Zrc20ContractAddress))
		if err != nil {
			return 0, errors.Wrap(fungibletypes.ErrContractCall, err.Error())
		}
		return gasLimit.Uint64(), nil
	} else if cctx.CoinType == common.CoinType_ERC20 {
		// get the gas limit of the associated asset
		fc, found := k.fungibleKeeper.GetForeignCoinFromAsset(ctx, cctx.InboundTxParams.Asset, cctx.InboundTxParams.SenderChainId)
		if !found {
			return 0, types.ErrForeignCoinNotFound
		}
		gasLimit, err := k.fungibleKeeper.QueryGasLimit(ctx, ethcommon.HexToAddress(fc.Zrc20ContractAddress))
		if err != nil {
			return 0, errors.Wrap(fungibletypes.ErrContractCall, err.Error())
		}
		return gasLimit.Uint64(), nil
	}

	return 0, nil
}

func IsPending(cctx types.CrossChainTx) bool {
	// pending inbound is not considered a "pending" state because it has not reached consensus yet
	return cctx.CctxStatus.Status == types.CctxStatus_PendingOutbound || cctx.CctxStatus.Status == types.CctxStatus_PendingRevert
}

// GetAbortedAmount returns the amount to refund for a given CCTX .
// If the CCTX has an outbound transaction, it returns the amount of the outbound transaction.
// If OutTxParams is nil or the amount is zero, it returns the amount of the inbound transaction.
// This is because there might be a case where the transaction is set to be aborted before paying gas or creating an outbound transaction.In such a situation we can refund the entire amount that has been locked in connector or TSS
func GetAbortedAmount(cctx types.CrossChainTx) sdkmath.Uint {
	if cctx.OutboundTxParams != nil && !cctx.GetCurrentOutTxParam().Amount.IsZero() {
		return cctx.GetCurrentOutTxParam().Amount
	}
	if cctx.InboundTxParams != nil {
		return cctx.InboundTxParams.Amount
	}

	return sdkmath.ZeroUint()
}

// ProcessInbound processes the inbound CCTX.
// It does a conditional dispatch to ProcessZEVMDeposit or ProcessCrosschainMsgPassing based on the receiver chain.
func (k Keeper) ProcessInbound(ctx sdk.Context, cctx *types.CrossChainTx) {
	if common.IsZetaChain(cctx.GetCurrentOutTxParam().ReceiverChainId) {
		k.ProcessZEVMDeposit(ctx, cctx)
	} else {
		k.ProcessCrosschainMsgPassing(ctx, cctx)
	}
}

// ProcessZEVMDeposit processes the EVM deposit CCTX. A deposit is a cctx which has Zetachain as the receiver chain.
// If the deposit is successful, the CCTX status is changed to OutboundMined.
// If the deposit returns an internal error i.e if HandleEVMDeposit() returns an error, but isContractReverted is false, the CCTX status is changed to Aborted.
// If the deposit is reverted, the function tries to create a revert cctx with status PendingRevert.
// If the creation of revert tx also fails it changes the status to Aborted.
// Note : Aborted CCTXs are not refunded in this function. The refund is done using a separate refunding mechanism.
// We do not return an error from this function , as all changes need to be persisted to the state.
// Instead we use a temporary context to make changes and then commit the context on for the happy path ,i.e cctx is set to OutboundMined.
func (k Keeper) ProcessZEVMDeposit(ctx sdk.Context, cctx *types.CrossChainTx) {
	tmpCtx, commit := ctx.CacheContext()
	isContractReverted, err := k.HandleEVMDeposit(tmpCtx, cctx)

	if err != nil && !isContractReverted { // exceptional case; internal error; should abort CCTX
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
		return
	} else if err != nil && isContractReverted { // contract call reverted; should refund
		revertMessage := err.Error()
		senderChain := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, cctx.InboundTxParams.SenderChainId)
		if senderChain == nil {
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "invalid sender chain")
			return
		}

		gasLimit, err := k.GetRevertGasLimit(ctx, cctx)
		if err != nil {
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, fmt.Sprintf("can't get revert tx gas limit,%s", err.Error()))
			return
		}
		if gasLimit == 0 {
			// use same gas limit of outbound as a fallback -- should not happen
			gasLimit = cctx.GetCurrentOutTxParam().OutboundTxGasLimit
		}

		// create new OutboundTxParams for the revert
		revertTxParams := &types.OutboundTxParams{
			Receiver:           cctx.InboundTxParams.Sender,
			ReceiverChainId:    cctx.InboundTxParams.SenderChainId,
			Amount:             cctx.InboundTxParams.Amount,
			OutboundTxGasLimit: gasLimit,
		}
		cctx.OutboundTxParams = append(cctx.OutboundTxParams, revertTxParams)

		// we create a new cached context, and we don't commit the previous one with EVM deposit
		tmpCtxRevert, commitRevert := ctx.CacheContext()
		err = func() error {
			err := k.PayGasAndUpdateCctx(
				tmpCtxRevert,
				senderChain.ChainId,
				cctx,
				cctx.InboundTxParams.Amount,
				false,
			)
			if err != nil {
				return err
			}
			// Update nonce using senderchain id as this is a revert tx and would go back to the original sender
			return k.UpdateNonce(tmpCtxRevert, senderChain.ChainId, cctx)
		}()
		if err != nil {
			cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, fmt.Sprintf("deposit revert message: %s err : %s", revertMessage, err.Error()))
			return
		}
		commitRevert()
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingRevert, revertMessage)
		return
	}
	// successful HandleEVMDeposit;
	commit()
	cctx.CctxStatus.ChangeStatus(types.CctxStatus_OutboundMined, "Remote omnichain contract call completed")
	return
}

// ProcessCrosschainMsgPassing processes the CCTX for crosschain message passing. A crosschain message passing is a cctx which has a non-Zetachain as the receiver chain.
// If the crosschain message passing is successful, the CCTX status is changed to PendingOutbound.
// If the crosschain message passing returns an error, the CCTX status is changed to Aborted.
// We do not return an error from this function , as all changes need to be persisted to the state.
// Instead we use a temporary context to make changes and then commit the context on for the happy path ,i.e cctx is set to PendingOutbound.
func (k Keeper) ProcessCrosschainMsgPassing(ctx sdk.Context, cctx *types.CrossChainTx) {
	tmpCtx, commit := ctx.CacheContext()
	outboundReceiverChainID := cctx.GetCurrentOutTxParam().ReceiverChainId
	err := func() error {
		err := k.PayGasAndUpdateCctx(
			tmpCtx,
			outboundReceiverChainID,
			cctx,
			cctx.InboundTxParams.Amount,
			false,
		)
		if err != nil {
			return err
		}
		return k.UpdateNonce(tmpCtx, outboundReceiverChainID, cctx)
	}()
	if err != nil {
		// do not commit anything here as the CCTX should be aborted
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
		return
	}
	commit()
	cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingOutbound, "")
	return
}

func (k Keeper) SaveInbound(ctx sdk.Context, cctx *types.CrossChainTx) {
	EmitEventInboundFinalized(ctx, cctx)
	k.AddFinalizedInbound(ctx,
		cctx.GetInboundTxParams().InboundTxObservedHash,
		cctx.GetInboundTxParams().SenderChainId,
		cctx.EventIndex)
	// #nosec G701 always positive
	cctx.InboundTxParams.InboundTxFinalizedZetaHeight = uint64(ctx.BlockHeight())
	cctx.InboundTxParams.TxFinalizationStatus = types.TxFinalizationStatus_Executed
	k.RemoveInTxTrackerIfExists(ctx, cctx.InboundTxParams.SenderChainId, cctx.InboundTxParams.InboundTxObservedHash)
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, *cctx)
}
