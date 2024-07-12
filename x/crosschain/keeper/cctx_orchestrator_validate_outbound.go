package keeper

import (
	"encoding/base64"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

/*
ValidateOutboundZEVM processes the finalization of an outbound transaction if receiver is ZEVM.
It takes deposit error and information if contract revert happened during deposit, to make a decision:

- If the deposit was successful, the CCTX status is changed to OutboundMined.

- If the deposit returned an internal error, but isContractReverted is false, the CCTX status is changed to Aborted.

- If the deposit is reverted, the function tries to create a revert cctx with status PendingRevert.

- If the creation of revert tx also fails it changes the status to Aborted.

Note : Aborted CCTXs are not refunded in this function. The refund is done using a separate refunding mechanism.
We do not return an error from this function, as all changes need to be persisted to the state.
Instead we use a temporary context to make changes and then commit the context on for the happy path, i.e cctx is set to OutboundMined.
New CCTX status after preprocessing is returned.
*/
func (k Keeper) ValidateOutboundZEVM(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	depositErr error,
	isContractReverted bool,
) (newCCTXStatus types.CctxStatus) {
	if depositErr != nil && isContractReverted {
		tmpCtxRevert, commitRevert := ctx.CacheContext()
		// contract call reverted; should refund via a revert tx
		err := k.validateFailedOutbound(
			tmpCtxRevert,
			cctx,
			types.CctxStatus_PendingOutbound,
			depositErr.Error(),
			cctx.InboundParams.Amount,
		)
		if err != nil {
			cctx.SetAbort(err.Error())
			return types.CctxStatus_Aborted
		}

		commitRevert()
		return types.CctxStatus_PendingRevert
	}
	k.validateSuccessfulOutbound(ctx, cctx, "", false)
	return types.CctxStatus_OutboundMined
}

// ValidateOutboundObservers processes the finalization of an outbound transaction based on the ballot status.
// The state is committed only if the individual steps are successful.
func (k Keeper) ValidateOutboundObservers(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	ballotStatus observertypes.BallotStatus,
	valueReceived string,
) error {
	tmpCtx, commit := ctx.CacheContext()
	err := func() error {
		switch ballotStatus {
		case observertypes.BallotStatus_BallotFinalized_SuccessObservation:
			k.validateSuccessfulOutbound(tmpCtx, cctx, valueReceived, true)
		case observertypes.BallotStatus_BallotFinalized_FailureObservation:
			err := k.validateFailedOutboundObservers(tmpCtx, cctx, valueReceived)
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if err != nil {
		return err
	}
	err = cctx.Validate()
	if err != nil {
		return err
	}
	commit()
	return nil
}

// validateFailedOutboundObservers processes a failed outbound transaction for observers. It does the following things in one function:
//
// 1. For Admin Tx or a withdrawal from Zeta chain, it aborts the CCTX
//
// 2. For other CCTX
//   - If the CCTX is in PendingOutbound, it creates a revert tx and sets the finalization status of the current outbound tx to executed
//   - If the CCTX is in PendingRevert, it sets the Status to Aborted
//
// 3. Emit an event for the failed outbound transaction
//
// 4. Set the finalization status of the current outbound tx to executed. If a revert tx is is created, the finalization status is not set, it would get set when the revert is processed via a subsequent transaction
//
// This function sets CCTX status , in cases where the outbound tx is successful, but tx itself fails
// This is done because SaveSuccessfulOutbound does not set the cctx status
// For cases where the outbound tx is unsuccessful, the cctx status is automatically set to Aborted in the validateFailedOutboundObservers function, so we can just return and error to trigger that
func (k Keeper) validateFailedOutboundObservers(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) error {
	oldStatus := cctx.CctxStatus.Status
	// The following logic is used to handler the mentioned conditions separately. The reason being
	// All admin tx is created using a policy message, there is no associated inbound tx, therefore we do not need any revert logic
	// For transactions which originated from ZEVM, we can process the outbound in the same block as there is no TSS signing required for the revert
	// For all other transactions we need to create a revert tx and set the status to pending revert

	if cctx.InboundParams.CoinType == coin.CoinType_Cmd {
		// if the cctx is of coin type cmd or the sender chain is zeta chain, then we do not revert, the cctx is aborted
		cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
		cctx.SetAbort("Outbound failed")
	} else if chains.IsZetaChain(cctx.InboundParams.SenderChainId, k.GetAuthorityKeeper().GetAdditionalChainList(ctx)) {
		switch cctx.InboundParams.CoinType {
		// Try revert if the coin-type is ZETA
		case coin.CoinType_Zeta:
			{
				err := k.validateFailedOutboundObserversForZEVM(ctx, cctx)
				if err != nil {
					return cosmoserrors.Wrap(err, "validateFailedOutboundObserversForZEVMTx")
				}
			}
		// For all other coin-types, we do not revert, the cctx is aborted
		default:
			{
				cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
				cctx.SetAbort("Outbound failed")
			}
		}
	} else {
		err := k.validateFailedOutbound(ctx, cctx, oldStatus, "Outbound failed, start revert", cctx.GetCurrentOutboundParam().Amount)
		if err != nil {
			return cosmoserrors.Wrap(err, "validateFailedOutbound")
		}
	}
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundFailure(ctx, valueReceived, oldStatus.String(), newStatus, cctx.Index)
	return nil
}

// validateFailedOutbound processes the failed outbound transaction
func (k Keeper) validateFailedOutbound(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	oldStatus types.CctxStatus,
	revertMsg string,
	inputAmount math.Uint,
) error {
	switch oldStatus {
	case types.CctxStatus_PendingOutbound:
		if _, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(ctx, cctx.InboundParams.SenderChainId); !found {
			return observertypes.ErrSupportedChains
		}

		gasLimit, err := k.GetRevertGasLimit(ctx, *cctx)
		if err != nil {
			return cosmoserrors.Wrap(err, "GetRevertGasLimit")
		}
		if gasLimit == 0 {
			// use same gas limit of outbound as a fallback -- should not happen
			gasLimit = cctx.OutboundParams[0].GasLimit
		}

		// create new OutboundParams for the revert
		err = cctx.AddRevertOutbound(gasLimit)
		if err != nil {
			return cosmoserrors.Wrap(err, "AddRevertOutbound")
		}

		err = k.PayGasAndUpdateCctx(
			ctx,
			cctx.InboundParams.SenderChainId,
			cctx,
			inputAmount,
			false,
		)
		if err != nil {
			return err
		}
		err = k.SetObserverOutboundInfo(ctx, cctx.InboundParams.SenderChainId, cctx)
		if err != nil {
			return err
		}
		// Not setting the finalization status here, the required changes have been made while creating the revert tx
		cctx.SetPendingRevert(revertMsg)
	case types.CctxStatus_PendingRevert:
		cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
		cctx.SetAbort("Outbound failed: revert failed; abort TX")
	}
	return nil
}

// validateSuccessfulOutbound processes a successful outbound transaction. It does the following things in one function:
//
//  1. Change the status of the CCTX from
//     - PendingRevert to Reverted
//     - PendingOutbound to OutboundMined
//
//  2. Set the finalization status of the current outbound tx to executed
//
//  3. Emit an event for the successful outbound transaction if flag is provided
//
// This function sets CCTX status, in cases where the outbound tx is successful, but tx itself fails
// This is done because SaveSuccessfulOutbound does not set the cctx status
// For cases where the outbound tx is unsuccessful, the cctx status is automatically set to Aborted in the validateFailedOutboundObservers function, so we can just return and error to trigger that
func (k Keeper) validateSuccessfulOutbound(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	valueReceived string,
	emitEvent bool,
) {
	oldStatus := cctx.CctxStatus.Status
	switch oldStatus {
	case types.CctxStatus_PendingRevert:
		cctx.SetReverted("Outbound succeeded, revert executed")
	case types.CctxStatus_PendingOutbound:
		cctx.SetOutboundMined("Outbound succeeded, mined")
	default:
		return
	}
	cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	newStatus := cctx.CctxStatus.Status.String()
	if emitEvent {
		EmitOutboundSuccess(ctx, valueReceived, oldStatus.String(), newStatus, cctx.Index)
	}
}

// validateFailedOutboundObserversForZEVM processes the failed outbound transaction for ZEVM
func (k Keeper) validateFailedOutboundObserversForZEVM(ctx sdk.Context, cctx *types.CrossChainTx) error {
	indexBytes, err := cctx.GetCCTXIndexBytes()
	if err != nil {
		// Return err to save the failed outbound and set to aborted
		return fmt.Errorf("failed reverting GetCCTXIndexBytes: %s", err.Error())
	}
	// Finalize the older outbound tx
	cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed

	// create new OutboundParams for the revert. We use the fixed gas limit for revert when calling zEVM
	err = cctx.AddRevertOutbound(fungiblekeeper.ZEVMGasLimitDepositAndCall.Uint64())
	if err != nil {
		// Return err to save the failed outbound ad set to aborted
		return fmt.Errorf("failed AddRevertOutbound: %s", err.Error())
	}

	// Trying to revert the transaction this would get set to a finalized status in the same block as this does not need a TSS singing
	cctx.SetPendingRevert("Outbound failed, trying revert")
	data, err := base64.StdEncoding.DecodeString(cctx.RelayedMessage)
	if err != nil {
		return fmt.Errorf("failed decoding relayed message: %s", err.Error())
	}

	// Fetch the original sender and receiver from the CCTX , since this is a revert the sender with be the receiver in the new tx
	originalSender := ethcommon.HexToAddress(cctx.InboundParams.Sender)
	// This transaction will always have two outbounds, the following logic is just an added precaution.
	// The contract call or token deposit would go the original sender.
	originalReceiver := ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
	orginalReceiverChainID := cctx.GetCurrentOutboundParam().ReceiverChainId
	if len(cctx.OutboundParams) == 2 {
		// If there are 2 outbound tx, then the original receiver is the receiver in the first outbound tx
		originalReceiver = ethcommon.HexToAddress(cctx.OutboundParams[0].Receiver)
		orginalReceiverChainID = cctx.OutboundParams[0].ReceiverChainId
	}

	// Call evm to revert the transaction
	// If revert fails, we set it to abort directly there is no way to refund here as the revert failed
	_, err = k.fungibleKeeper.ZETARevertAndCallContract(
		ctx,
		originalSender,
		originalReceiver,
		cctx.InboundParams.SenderChainId,
		orginalReceiverChainID,
		cctx.GetCurrentOutboundParam().Amount.BigInt(),
		data,
		indexBytes,
	)
	if err != nil {
		return fmt.Errorf("failed ZETARevertAndCallContract: %s", err.Error())
	}

	cctx.SetReverted("Outbound failed, revert executed")
	if len(ctx.TxBytes()) > 0 {
		// add event for tendermint transaction hash format
		hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
		ethTxHash := ethcommon.BytesToHash(hash)
		cctx.GetCurrentOutboundParam().Hash = ethTxHash.String()
		// #nosec G115 always positive
		cctx.GetCurrentOutboundParam().ObservedExternalHeight = uint64(ctx.BlockHeight())
	}
	cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	return nil
}
