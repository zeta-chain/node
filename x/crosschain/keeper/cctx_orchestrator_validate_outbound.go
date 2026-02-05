package keeper

import (
	"encoding/base64"
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	tmbytes "github.com/cometbft/cometbft/libs/bytes"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/cosmos/evm/x/vm/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	ccctxerror "github.com/zeta-chain/node/pkg/errors"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungiblekeeper "github.com/zeta-chain/node/x/fungible/keeper"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// ValidateOutboundZEVM processes the finalization of an outbound transaction if receiver is ZEVM.
// It takes deposit error and information if contract revert happened during deposit, to make a decision:
// - If the deposit was successful, the CCTX status is changed to OutboundMined.
// - If the deposit returned an internal error, but isContractReverted is false, the CCTX status is changed to Aborted.
// - If the deposit is reverted, the function tries to create a revert cctx with status PendingRevert.
// - If the creation of revert tx also fails it changes the status to Aborted.
// Note : Aborted CCTXs are not refunded in this function. The refund is done using a separate refunding mechanism.
// We do not return an error from this function, as all changes need to be persisted to the state.
// Instead we use a temporary context to make changes and then commit the context on for the happy path, i.e cctx is set to OutboundMined.
// New CCTX status after preprocessing is returned.
func (k Keeper) ValidateOutboundZEVM(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	depositErr error,
	shouldRevert bool,
) (newCCTXStatus types.CctxStatus) {
	if depositErr != nil && shouldRevert {
		tmpCtxRevert, commitRevert := ctx.CacheContext()
		// contract call reverted; should refund via a revert tx
		revertErr := k.processFailedOutboundOnExternalChain(
			tmpCtxRevert,
			cctx,
			types.CctxStatus_PendingOutbound,
			depositErr,
			cctx.InboundParams.Amount,
		)
		if revertErr != nil {
			// Error here would mean the outbound tx failed and we also failed to create a revert tx.
			// This is the only case where we set outbound and revert messages,
			// as both the outbound and the revert failed in the same block
			k.ProcessAbort(ctx, cctx, types.StatusMessages{
				StatusMessage:        "revert failed to be processed",
				ErrorMessageOutbound: ccctxerror.NewCCTXErrorJSONMessage("", depositErr),
				ErrorMessageRevert:   ccctxerror.NewCCTXErrorJSONMessage("", revertErr),
			})
			return types.CctxStatus_Aborted
		}

		commitRevert()
		return types.CctxStatus_PendingRevert
	}
	k.processSuccessfulOutbound(ctx, cctx, "", false)
	return types.CctxStatus_OutboundMined
}

// processSuccessfulOutbound processes a successful outbound transaction. It does the following things in one function:
//
//  1. Change the status of the CCTX from
//     - PendingRevert to Reverted
//     - PendingOutbound to OutboundMined
//  2. Set the finalization status of the current outbound tx to executed
//  3. Emit an event for the successful outbound transaction if flag is provided
//
// This function sets CCTX status, in cases where the outbound tx is successful, but tx itself fails
// This is done because HandleValidOutbound does not set the cctx status
// For cases where the outbound tx is unsuccessful, the cctx status is automatically set to Aborted in the processFailedOutboundObservers function, so we can just return and error to trigger that
func (k Keeper) processSuccessfulOutbound(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	valueReceived string,
	emitEvent bool,
) {
	oldStatus := cctx.CctxStatus.Status
	switch oldStatus {
	case types.CctxStatus_PendingRevert:
		cctx.SetReverted()
	case types.CctxStatus_PendingOutbound:
		cctx.SetOutboundMined()
	default:
		return
	}
	cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	newStatus := cctx.CctxStatus.Status.String()
	if emitEvent {
		EmitOutboundSuccess(ctx, valueReceived, oldStatus.String(), newStatus, cctx.Index)
	}
}

// processFailedOutboundOnExternalChain processes the failed outbound transaction where the receiver is an external chain.
func (k Keeper) processFailedOutboundOnExternalChain(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	oldStatus types.CctxStatus,
	depositErr error,
	inputAmount math.Uint,
) error {
	switch oldStatus {
	case types.CctxStatus_PendingOutbound:
		if _, found := k.zetaObserverKeeper.GetSupportedChainFromChainID(
			ctx,
			cctx.InboundParams.SenderChainId,
		); !found {
			return observertypes.ErrSupportedChains
		}

		gasLimit, err := k.GetRevertGasLimit(ctx, *cctx)
		if err != nil {
			return cosmoserrors.Wrap(err, "GetRevertGasLimit")
		}
		if gasLimit == 0 {
			// use same gas limit of outbound as a fallback -- should not happen
			gasLimit = cctx.OutboundParams[0].CallOptions.GasLimit
		}
		// create new OutboundParams for the revert
		err = cctx.AddRevertOutbound(gasLimit)
		if err != nil {
			return cosmoserrors.Wrap(err, "AddRevertOutbound")
		}

		// pay revert outbound gas fee
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

		receiverBytes, err := chains.DecodeAddressFromChainID(
			cctx.GetCurrentOutboundParam().ReceiverChainId,
			cctx.GetCurrentOutboundParam().Receiver,
			k.GetAuthorityKeeper().GetAdditionalChainList(ctx),
		)
		if err != nil {
			return errors.Wrap(err, "failed to decode receiver address")
		}

		// validate data of the revert outbound
		err = k.validateOutbound(
			ctx,
			cctx.GetCurrentOutboundParam().ReceiverChainId,
			cctx.InboundParams.CoinType,
			cctx.GetCurrentOutboundParam().Amount.BigInt(),
			receiverBytes,
		)
		if err != nil {
			return errors.Wrap(err, "failed to validate ZRC20 withdrawal")
		}

		err = k.SetObserverOutboundInfo(ctx, cctx.InboundParams.SenderChainId, cctx)
		if err != nil {
			return err
		}
		// Not setting the finalization status here, the required changes have been made while creating the revert tx
		cctx.SetPendingRevert(types.StatusMessages{
			StatusMessage:        "outbound failed",
			ErrorMessageOutbound: ccctxerror.NewCCTXErrorJSONMessage("", depositErr),
		})
	case types.CctxStatus_PendingRevert:
		cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
		return errors.Wrap(depositErr, "revert failed to be processed on connected chain")
	default:
		return fmt.Errorf("unexpected cctx status %s", cctx.CctxStatus.Status)
	}
	return nil
}

// ValidateOutboundObservers processes the finalization of an outbound transaction based on the ballot status.
// The state is committed only if the individual steps are successful.
func (k Keeper) ValidateOutboundObservers(
	ctx sdk.Context,
	cctx *types.CrossChainTx,
	ballotStatus observertypes.BallotStatus,
	valueReceived string,
) error {
	// temporary context ensure we don't end up with inconsistent state
	tmpCtx, commit := ctx.CacheContext()

	switch ballotStatus {
	case observertypes.BallotStatus_BallotFinalized_SuccessObservation:
		k.processSuccessfulOutbound(tmpCtx, cctx, valueReceived, true)
	case observertypes.BallotStatus_BallotFinalized_FailureObservation:
		k.processFailedOutboundObservers(tmpCtx, cctx, valueReceived)
	}

	err := cctx.Validate()
	if err != nil {
		return err
	}
	commit()
	return nil
}

// processFailedOutboundObservers processes a failed outbound transaction for observers. It does the following things in one function:
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
// This is done because HandleValidOutbound does not set the cctx status
// For cases where the outbound tx is unsuccessful, the cctx status is automatically set to Aborted in the processFailedOutboundObservers function, so we can just return and error to trigger that
func (k Keeper) processFailedOutboundObservers(ctx sdk.Context, cctx *types.CrossChainTx, valueReceived string) {
	oldStatus := cctx.CctxStatus.Status
	// The following logic is used to handler the mentioned conditions separately. The reason being
	// All admin tx is created using a policy message, there is no associated inbound tx, therefore we do not need any revert logic
	// For transactions which originated from ZEVM, we can process the outbound in the same block as there is no TSS signing required for the revert
	// For all other transactions we need to create a revert tx and set the status to pending revert

	if cctx.ProtocolContractVersion == types.ProtocolContractVersion_V2 {
		err := k.processFailedOutboundV2(ctx, cctx)

		// if the revert failed to be processed, we process the abort of the cctx
		if err != nil {
			k.ProcessAbort(ctx, cctx, types.StatusMessages{
				StatusMessage: "outbound failed and revert failed",
				ErrorMessageRevert: ccctxerror.NewCCTXErrorJSONMessage(
					"revert tx failed to be executed",
					err,
				),
			})
		}
		return
	}

	if cctx.InboundParams.CoinType == coin.CoinType_Cmd {
		// if the cctx is of coin type cmd or the sender chain is zeta chain, then we do not revert, the cctx is aborted
		cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
		cctx.SetAbort(types.StatusMessages{
			StatusMessage: "outbound failed for admin cmd",
			ErrorMessageOutbound: ccctxerror.NewCCTXErrorJSONMessage(
				"admin cmd outbound failed to be executed on connected chain",
				nil,
			),
			ErrorMessageRevert: ccctxerror.NewCCTXErrorJSONMessage("admin cmd outbound can't be reverted", nil),
		})
	} else if chains.IsZetaChain(cctx.InboundParams.SenderChainId, k.GetAuthorityKeeper().GetAdditionalChainList(ctx)) {
		switch cctx.InboundParams.CoinType {
		// Try revert if the coin-type is ZETA
		case coin.CoinType_Zeta:
			{
				err := k.processFailedZETAOutboundOnZEVM(ctx, cctx)
				if err != nil {
					cctx.SetAbort(types.StatusMessages{
						StatusMessage: "zeta outbound failed and revert failed on ZetaChain",
						ErrorMessageOutbound: ccctxerror.NewCCTXErrorJSONMessage(
							"outbound failed to be executed on connected chain",
							nil,
						),
						ErrorMessageRevert: ccctxerror.NewCCTXErrorJSONMessage("revert failed", err),
					})
				}
			}
		// For all other coin-types, we do not revert, the cctx is aborted
		default:
			{
				cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
				cctx.SetAbort(types.StatusMessages{
					StatusMessage: "outbound failed and revert can't be processed",
					ErrorMessageOutbound: ccctxerror.NewCCTXErrorJSONMessage(
						"outbound failed to be executed on connected chain",
						nil,
					),
					ErrorMessageRevert: ccctxerror.NewCCTXErrorJSONMessage(
						fmt.Sprintf(
							"revert on ZetaChain is not supported for cctx v1 with coin type %s",
							cctx.InboundParams.CoinType,
						),
						nil,
					),
				})
			}
		}
	} else {
		// We add a hardcoded message here as the error from the connected chain is not available,
		err := k.processFailedOutboundOnExternalChain(
			ctx,
			cctx,
			oldStatus,
			errors.New("outbound failed to be executed on connected chain"),
			cctx.GetCurrentOutboundParam().Amount,
		)
		if err != nil {
			cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
			cctx.SetAbort(types.StatusMessages{
				StatusMessage: "outbound failed and revert failed on connected chain",
				ErrorMessageOutbound: ccctxerror.NewCCTXErrorJSONMessage(
					"outbound failed to be executed on connected chain",
					nil,
				),
				ErrorMessageRevert: ccctxerror.NewCCTXErrorJSONMessage("revert failed", err),
			})
		}
	}
	newStatus := cctx.CctxStatus.Status.String()
	EmitOutboundFailure(ctx, valueReceived, oldStatus.String(), newStatus, cctx.Index)
}

// processFailedZETAOutboundOnZEVM processes a failed ZETA outbound on ZEVM
func (k Keeper) processFailedZETAOutboundOnZEVM(ctx sdk.Context, cctx *types.CrossChainTx) error {
	indexBytes, err := cctx.GetCCTXIndexBytes()
	if err != nil {
		// Return err to save the failed outbound and set to aborted
		return fmt.Errorf("failed reverting GetCCTXIndexBytes: %s", err.Error())
	}
	// Finalize the older outbound tx
	cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed

	// create new OutboundParams for the revert. We use the fixed gas limit for revert when calling zEVM

	if err := cctx.AddRevertOutbound(fungiblekeeper.ZEVMGasLimitDepositAndCall.Uint64()); err != nil {
		// Return err to save the failed outbound ad set to aborted
		return fmt.Errorf("failed AddRevertOutbound: %s", err.Error())
	}

	// Trying to revert the transaction, this would get set to a finalized status in the same block as this does not need a TSS singing
	// The outbound failed due to majority of observer voting failure.We do not have the exact reason for the failure available on chain.
	cctx.SetPendingRevert(types.StatusMessages{
		StatusMessage: "outbound failed",
		ErrorMessageOutbound: ccctxerror.NewCCTXErrorJSONMessage(
			"outbound failed to be executed on connected chain",
			nil,
		),
	})
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

	cctx.SetReverted()

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

// processFailedOutboundV2 processes a failed outbound transaction for protocol version 2
// for revert, in V2 we have some assumption simplifying the logic
// - sender chain is ZetaChain for regular outbound (not revert outbound)
// - all coin type use the same workflow
// TODO: consolidate logic with above function
// https://github.com/zeta-chain/node/issues/2627
func (k Keeper) processFailedOutboundV2(ctx sdk.Context, cctx *types.CrossChainTx) error {
	switch cctx.CctxStatus.Status {
	case types.CctxStatus_PendingOutbound:
		// check the sender is ZetaChain
		zetaChain, err := chains.ZetaChainFromCosmosChainID(ctx.ChainID())
		if err != nil {
			return errors.Wrap(err, "failed to get ZetaChain chainID")
		}
		if cctx.InboundParams.SenderChainId != zetaChain.ChainId {
			return fmt.Errorf(
				"sender chain for withdraw cctx is not ZetaChain expected %d got %d",
				zetaChain.ChainId,
				cctx.InboundParams.SenderChainId,
			)
		}

		//  get the chain ID of the connected chain
		chainID := cctx.GetCurrentOutboundParam().ReceiverChainId

		// add revert outbound
		if err := cctx.AddRevertOutbound(fungiblekeeper.ZEVMGasLimitDepositAndCall.Uint64()); err != nil {
			// Return err to save the failed outbound ad set to aborted
			return errors.Wrap(err, "failed AddRevertOutbound")
		}

		// update status
		cctx.SetPendingRevert(types.StatusMessages{
			StatusMessage: "outbound failed",
			ErrorMessageOutbound: ccctxerror.NewCCTXErrorJSONMessage(
				"outbound tx failed to be executed on connected chain",
				nil,
			),
		})

		// use a temporary context to not commit any state change in case of error
		tmpCtx, commit := ctx.CacheContext()

		// process the revert on ZEVM
		to := ethcommon.HexToAddress(cctx.GetCurrentOutboundParam().Receiver)
		evmTxResponse, err := k.fungibleKeeper.ProcessRevert(
			tmpCtx,
			cctx.InboundParams.Sender,
			cctx.GetCurrentOutboundParam().Amount.BigInt(),
			chainID,
			cctx.InboundParams.CoinType,
			cctx.InboundParams.Asset,
			to,
			cctx.RevertOptions.CallOnRevert,
			cctx.RevertOptions.RevertMessage,
		)
		if fungibletypes.IsContractReverted(evmTxResponse, err) {
			// this is a contract revert, we commit the state to save the emitted logs related to revert
			commit()
			return errors.Wrap(err, "revert transaction reverted")
		} else if err != nil {
			// this should not happen and we don't commit the state to avoid inconsistent state
			return errors.Wrap(err, "revert transaction could not be processed")
		}

		// a withdrawal event in the logs could generate cctxs for outbound transactions.
		if evmTxResponse != nil {
			logs := evmtypes.LogsToEthereum(evmTxResponse.Logs)
			if len(logs) > 0 {
				tmpCtx = tmpCtx.WithValue(InCCTXIndexKey, cctx.Index)
				txOrigin := cctx.InboundParams.TxOrigin
				if txOrigin == "" {
					txOrigin = cctx.InboundParams.Sender
				}

				// process logs to process cctx events initiated during the contract call
				if err = k.ProcessLogs(tmpCtx, logs, to, txOrigin); err != nil {
					// this happens if the cctx events are not processed correctly with invalid withdrawals
					// in this situation we want the CCTX to be reverted, we don't commit the state so the contract call is not persisted
					// the contract call is considered as reverted
					return errors.Wrap(types.ErrCannotProcessWithdrawal, err.Error())
				}
			}
		}

		// commit state change from the deposit and eventual cctx events
		commit()

		// tx is reverted
		cctx.SetReverted()

		// add event for tendermint transaction hash format
		if len(ctx.TxBytes()) > 0 {
			hash := tmbytes.HexBytes(tmtypes.Tx(ctx.TxBytes()).Hash())
			ethTxHash := ethcommon.BytesToHash(hash)
			cctx.GetCurrentOutboundParam().Hash = ethTxHash.String()
			// #nosec G115 always positive
			cctx.GetCurrentOutboundParam().ObservedExternalHeight = uint64(ctx.BlockHeight())
		}
		cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
	case types.CctxStatus_PendingRevert:
		cctx.GetCurrentOutboundParam().TxFinalizationStatus = types.TxFinalizationStatus_Executed
		return errors.New("revert failed to be processed on connected chain")
	default:
		return fmt.Errorf("unexpected cctx status %s", cctx.CctxStatus.Status)
	}

	return nil
}
