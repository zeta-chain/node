package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// Casts a vote on an outbound transaction observed on a connected chain (after
// it has been broadcasted to and finalized on a connected chain). If this is
// the first vote, a new ballot is created. When a threshold of votes is
// reached, the ballot is finalized. When a ballot is finalized, the outbound
// transaction is processed.
//
// If the observation is successful, the difference between zeta burned
// and minted is minted by the bank module and deposited into the module
// account.
//
// If the observation is unsuccessful, the logic depends on the previous
// status.
//
// If the previous status was `PendingOutbound`, a new revert transaction is
// created. To cover the revert transaction fee, the required amount of tokens
// submitted with the CCTX are swapped using a Uniswap V2 contract instance on
// ZetaChain for the ZRC20 of the gas token of the receiver chain. The ZRC20
// tokens are then
// burned. The nonce is updated. If everything is successful, the CCTX status is
// changed to `PendingRevert`.
//
// If the previous status was `PendingRevert`, the CCTX is aborted.
//
// ```mermaid
// stateDiagram-v2
//
//	state observation <<choice>>
//	state success_old_status <<choice>>
//	state fail_old_status <<choice>>
//	PendingOutbound --> observation: Finalize outbound
//	observation --> success_old_status: Observation succeeded
//	success_old_status --> Reverted: Old status is PendingRevert
//	success_old_status --> OutboundMined: Old status is PendingOutbound
//	observation --> fail_old_status: Observation failed
//	fail_old_status --> PendingRevert: Old status is PendingOutbound
//	fail_old_status --> Aborted: Old status is PendingRevert
//	PendingOutbound --> Aborted: Finalize outbound error
//
// ```
//
// Only observer validators are authorized to broadcast this message.
func (k msgServer) ProveOutboundTx(goCtx context.Context, msg *types.MsgProveOutboundTx) (*types.MsgProveOutboundTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// verify proof
	blockHash := eth.HexToHash(msg.BlockHash)
	res, found := k.zetaObserverKeeper.GetBlockHeader(ctx, blockHash.Bytes())
	if !found {
		return nil, sdkerrors.Wrapf(observerTypes.ErrBlockHeaderNotFound, "Block header not found: %s", msg.BlockHash)
	}
	chainID := res.ChainId

	var header ethtypes.Header
	err := rlp.DecodeBytes(res.Header, &header)
	if err != nil {
		return nil, sdkerrors.Wrapf(observerTypes.ErrBlockHeaderNotFound, "cannot decode block header: %s", msg.BlockHash)
	}
	var txx ethtypes.Transaction
	var receipt ethtypes.Receipt
	txxVal, err := msg.TxProof.Verify(header.TxHash, int(msg.TxIndex))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrProofVerificationFail, "cannot verify tx proof: %s", err)
	}
	receiptVal, err := msg.ReceiptProof.Verify(header.ReceiptHash, int(msg.TxIndex))
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrProofVerificationFail, "cannot verify receipt proof: %s", err)
	}
	err = txx.UnmarshalBinary(txxVal)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrProofVerificationFail, "cannot unmarshal tx: %s", err)
	}
	err = receipt.UnmarshalBinary(receiptVal)
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrProofVerificationFail, "cannot unmarshal receipt: %s", err)
	}
	// end verify proof

	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCannotFindTSSKeys, "Tss not found")
	}
	nonce := txx.Nonce()
	cctxToNonce, found := k.GetNonceToCctx(ctx, tss.TssPubkey, chainID, int64(nonce))
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCannotFindCctx, "Cctx not found")
	}
	index := cctxToNonce.CctxIndex
	cctx, found := k.GetCrossChainTx(ctx, index)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCannotFindCctx, "Cctx not found")
	}
	status := cctx.CctxStatus.Status
	if status != types.CctxStatus_PendingOutbound && status != types.CctxStatus_PendingRevert {
		return nil, sdkerrors.Wrapf(types.ErrStatusNotPending, "Cctx status is not pending")
	}

	// this is important; otherwise smoketest "Goerli->Goerli Message Passing (revert fail)"
	// will appear to be successful but it's not correct, because the same outtx are used
	// accepted here. The scond one should be rejected because the nonce is not correct
	if cctx.GetCurrentOutTxParam().OutboundTxTssNonce != nonce {
		return nil, sdkerrors.Wrapf(types.ErrNonceMismatch, "Nonce mismatch")
	}

	if receipt.Status == ethtypes.ReceiptStatusFailed {

	}

	cctx.GetCurrentOutTxParam().OutboundTxHash = txx.Hash().Hex()
	cctx.CctxStatus.LastUpdateTimestamp = ctx.BlockHeader().Time.Unix()

	// FinalizeOutbound sets final status for a successful vote
	// FinalizeOutbound updates CCTX Prices and Nonce for a revert

	tmpCtx, commit := ctx.CacheContext()
	err = func() error { //err = FinalizeOutbound(k, ctx, &cctx, msg, ballot.BallotStatus)
		cctx.GetCurrentOutTxParam().OutboundTxObservedExternalHeight = uint64(res.Height)
		oldStatus := cctx.CctxStatus.Status
		switch receipt.Status {
		case ethtypes.ReceiptStatusSuccessful:
			switch oldStatus {
			case types.CctxStatus_PendingRevert:
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_Reverted, "")
			case types.CctxStatus_PendingOutbound:
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_OutboundMined, "")
			}
			newStatus := cctx.CctxStatus.Status.String()
			err := ctx.EventManager().EmitTypedEvents(&types.EventOutboundSuccess{
				MsgTypeUrl: sdk.MsgTypeURL(&types.MsgVoteOnObservedOutboundTx{}),
				CctxIndex:  cctx.Index,
				//ZetaMinted: msg.ZetaMinted.String(), // FIXME: this is not the correct amount
				OldStatus: oldStatus.String(),
				NewStatus: newStatus,
			})
			if err != nil {
				ctx.Logger().Error("Error emitting MsgVoteOnObservedOutboundTx :", err)
			}
		case ethtypes.ReceiptStatusFailed:
			switch oldStatus {
			case types.CctxStatus_PendingOutbound:
				// create new OutboundTxParams for the revert
				cctx.OutboundTxParams = append(cctx.OutboundTxParams, &types.OutboundTxParams{
					Receiver:           cctx.InboundTxParams.Sender,
					ReceiverChainId:    cctx.InboundTxParams.SenderChainId,
					Amount:             cctx.InboundTxParams.Amount,
					CoinType:           cctx.InboundTxParams.CoinType,
					OutboundTxGasLimit: cctx.OutboundTxParams[0].OutboundTxGasLimit, // NOTE(pwu): revert gas limit = initial outbound gas limit set by user;
				})
				err := k.PayGasInZetaAndUpdateCctx(tmpCtx, cctx.InboundTxParams.SenderChainId, &cctx)
				if err != nil {
					return err
				}
				err = k.UpdateNonce(tmpCtx, cctx.InboundTxParams.SenderChainId, &cctx)
				if err != nil {
					return err
				}
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_PendingRevert, "Outbound failed, start revert")
			case types.CctxStatus_PendingRevert:
				cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, "Outbound failed: revert failed; abort TX")
			}
			//newStatus := cctx.CctxStatus.Status.String()
			//EmitOutboundFailure(ctx, msg, oldStatus.String(), newStatus, cctx) // FIXME: enable this
		}
		return nil
	}()
	if err != nil {
		// do not commit tmpCtx
		cctx.CctxStatus.ChangeStatus(types.CctxStatus_Aborted, err.Error())
		ctx.Logger().Error(err.Error())
		k.RemoveOutTxTracker(ctx, chainID, nonce)
		k.RemoveFromPendingNonces(ctx, tss.TssPubkey, chainID, int64(nonce))
		k.RemoveOutTxTracker(ctx, chainID, nonce)
		k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
		return &types.MsgProveOutboundTxResponse{}, nil
	}
	commit()
	// Set the ballot index to the finalized ballot
	//cctx.GetCurrentOutTxParam().OutboundTxBallotIndex = ballotIndex // no ballot
	k.RemoveFromPendingNonces(ctx, tss.TssPubkey, chainID, int64(nonce))
	k.RemoveOutTxTracker(ctx, chainID, nonce)
	k.SetCctxAndNonceToCctxAndInTxHashToCctx(ctx, cctx)
	return &types.MsgProveOutboundTxResponse{}, nil
}
