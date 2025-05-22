package observer

import (
	"context"
	"encoding/hex"

	cosmosmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	solanarpc "github.com/zeta-chain/node/zetaclient/chains/solana/rpc"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/logs"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// MaxSignaturesPerTicker is the maximum number of signatures to process on a ticker
	MaxSignaturesPerTicker = 100
)

// ObserveInbound observes the Solana chain for inbounds and post votes to zetacore.
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	chainID := ob.Chain().ChainId
	pageLimit := solanarpc.DefaultPageLimit

	// scan from gateway 1st signature if last scanned tx is absent in the database
	// the 1st gateway signature is typically the program initialization
	if ob.LastTxScanned() == "" {
		lastSig, err := solanarpc.GetFirstSignatureForAddress(ctx, ob.solClient, ob.gatewayID, pageLimit)
		if err != nil {
			return errors.Wrapf(err, "error GetFirstSignatureForAddress for chain %d address %s", chainID, ob.gatewayID)
		}
		ob.WithLastTxScanned(lastSig.String())
	}

	// query last finalized slot
	lastSlot, errSlot := ob.solClient.GetSlot(ctx, rpc.CommitmentFinalized)
	if errSlot != nil {
		ob.Logger().Inbound.Err(errSlot).Msg("unable to get last slot")
	}

	// get all signatures for the gateway address since last scanned signature
	lastSig := solana.MustSignatureFromBase58(ob.LastTxScanned())
	signatures, err := solanarpc.GetSignaturesForAddressUntil(ctx, ob.solClient, ob.gatewayID, lastSig, pageLimit)
	if err != nil {
		ob.Logger().Inbound.Err(err).Msg("error GetSignaturesForAddressUntil")
		return err
	}

	// update metrics if no new signatures found
	if len(signatures) == 0 {
		if errSlot == nil {
			ob.WithLastBlockScanned(lastSlot)
		}
	} else {
		ob.Logger().Inbound.Info().
			Str(logs.FieldMethod, "ObserveInbound").
			Int("signatures", len(signatures)).
			Int64(logs.FieldChain, chainID).
			Msg("got wrong amount of signatures")
	}

	// loop signature from oldest to latest to filter inbound events
	for i := len(signatures) - 1; i >= 0; i-- {
		sig := signatures[i]
		sigString := sig.Signature.String()

		// process successfully signature only
		if sig.Err == nil {
			txResult, err := solanarpc.GetTransaction(ctx, ob.solClient, sig.Signature)
			switch {
			case errors.Is(err, solanarpc.ErrUnsupportedTxVersion):
				ob.Logger().Inbound.Warn().
					Stringer("tx.signature", sig.Signature).
					Msg("ObserveInbound: skip unsupported transaction")
			// just save the sig to last scanned txs
			case err != nil:
				// we have to re-scan this signature on next ticker
				return errors.Wrapf(err, "error GetTransaction for sig %s", sigString)
			default:
				// filter the events
				events, err := FilterInboundEvents(txResult, ob.gatewayID, ob.Chain().ChainId, ob.Logger().Inbound)
				if err != nil {
					// Log the error but continue processing other transactions
					ob.Logger().Inbound.Error().
						Err(err).
						Str("tx.signature", sigString).
						Msg("ObserveInbound: error filtering events, skipping")
					continue
				}

				// vote on the events
				if err := ob.VoteInboundEvents(ctx, events); err != nil {
					// return error to retry this transaction
					return errors.Wrapf(err, "error voting on events for transaction %s, will retry", sigString)
				}
			}
		}

		// signature scanned; save last scanned signature to both memory and db, ignore db error
		if err = ob.SaveLastTxScanned(sigString, sig.Slot); err != nil {
			ob.Logger().
				Inbound.Error().
				Err(err).
				Str("tx.signature", sigString).
				Msg("ObserveInbound: error saving last sig")
		}

		ob.Logger().
			Inbound.Info().
			Str("tx.signature", sigString).
			Uint64("tx.slot", sig.Slot).
			Msg("ObserveInbound: last scanned sig")

		// take a rest if max signatures per ticker is reached
		if len(signatures)-i >= MaxSignaturesPerTicker {
			break
		}
	}

	return nil
}

// VoteInboundEvents posts votes for inbound events to zetacore.
func (ob *Observer) VoteInboundEvents(ctx context.Context, events []*clienttypes.InboundEvent) error {
	for _, event := range events {
		msg := ob.BuildInboundVoteMsgFromEvent(event)
		if msg != nil {
			_, err := ob.PostVoteInbound(ctx, msg, zetacore.PostVoteInboundExecutionGasLimit)
			if err != nil {
				return errors.Wrapf(err, "error PostVoteInbound")
			}
		}
	}
	return nil
}

// FilterInboundEvents filters inbound events from a tx result.
// Note: for consistency with EVM chains, this method
//   - takes at one event (the first) per token (SOL or SPL) per transaction.
//   - takes at most two events (one SOL + one SPL) per transaction.
//   - ignores exceeding events.
func FilterInboundEvents(
	txResult *rpc.GetTransactionResult,
	gatewayID solana.PublicKey,
	senderChainID int64,
	logger zerolog.Logger,
) ([]*clienttypes.InboundEvent, error) {
	// unmarshal transaction
	tx, err := txResult.Transaction.GetTransaction()
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling transaction")
	}

	// there should be at least one instruction and one account, otherwise skip
	if len(tx.Message.Instructions) <= 0 {
		return nil, nil
	}

	// create event array to collect all events in the transaction
	seenDeposit := false
	seenDepositSPL := false
	seenCall := false
	events := make([]*clienttypes.InboundEvent, 0)

	// loop through instruction list to filter the 1st valid event
	for i, instruction := range tx.Message.Instructions {
		// get the program ID
		programPk, err := tx.Message.Program(instruction.ProgramIDIndex)
		if err != nil {
			logger.Err(err).
				Str(logs.FieldMethod, "FilterInboundEvents").
				Str("signature", tx.Signatures[0].String()).
				Uint16("index", instruction.ProgramIDIndex).
				Msg("no program found")
			continue
		}

		// skip instructions that are irrelevant to the gateway program invocation
		if !programPk.Equals(gatewayID) {
			continue
		}

		// try parsing the instruction as a 'deposit' if not seen yet
		if !seenDeposit {
			deposit, err := solanacontracts.ParseInboundAsDeposit(tx, i, txResult.Slot)
			if err != nil {
				return nil, errors.Wrap(err, "error ParseInboundAsDeposit")
			} else if deposit != nil {
				seenDeposit = true
				events = append(events, &clienttypes.InboundEvent{
					SenderChainID:    senderChainID,
					Sender:           deposit.Sender,
					Receiver:         deposit.Receiver,
					TxOrigin:         deposit.Sender,
					Amount:           deposit.Amount,
					Memo:             deposit.Memo,
					BlockNumber:      deposit.Slot, // instead of using block, Solana explorer uses slot for indexing
					TxHash:           tx.Signatures[0].String(),
					Index:            0, // hardcode to 0 for Solana, not a EVM smart contract call
					CoinType:         coin.CoinType_Gas,
					Asset:            deposit.Asset,
					IsCrossChainCall: deposit.IsCrossChainCall,
					RevertOptions:    deposit.RevertOptions,
				})
				logger.Info().
					Str(logs.FieldMethod, "FilterInboundEvents").
					Str("signature", tx.Signatures[0].String()).
					Int("instruction", i).
					Msg("deposit detected")
			}
		} else {
			logger.Warn().
				Str(logs.FieldMethod, "FilterInboundEvents").
				Str("signature", tx.Signatures[0].String()).
				Int("instruction", i).
				Msg("multiple deposits detected")
		}

		// try parsing the instruction as a 'deposit_spl_token' if not seen yet
		if !seenDepositSPL {
			deposit, err := solanacontracts.ParseInboundAsDepositSPL(tx, i, txResult.Slot)
			if err != nil {
				return nil, errors.Wrap(err, "error ParseInboundAsDepositSPL")
			} else if deposit != nil {
				seenDepositSPL = true
				events = append(events, &clienttypes.InboundEvent{
					SenderChainID:    senderChainID,
					Sender:           deposit.Sender,
					Receiver:         deposit.Receiver,
					TxOrigin:         deposit.Sender,
					Amount:           deposit.Amount,
					Memo:             deposit.Memo,
					BlockNumber:      deposit.Slot, // instead of using block, Solana explorer uses slot for indexing
					TxHash:           tx.Signatures[0].String(),
					Index:            0, // hardcode to 0 for Solana, not a EVM smart contract call
					CoinType:         coin.CoinType_ERC20,
					Asset:            deposit.Asset,
					IsCrossChainCall: deposit.IsCrossChainCall,
					RevertOptions:    deposit.RevertOptions,
				})
				logger.Info().
					Str(logs.FieldMethod, "FilterInboundEvents").
					Str("signature", tx.Signatures[0].String()).
					Int("instruction", i).
					Msg("SPL deposit detected")
			}
		} else {
			logger.Warn().
				Str("signature", tx.Signatures[0].String()).
				Int("instruction", i).
				Msg("multiple SPL deposits detected")
		}

		// try parsing the instruction as a 'call' if not seen yet
		if !seenCall {
			call, err := solanacontracts.ParseInboundAsCall(tx, i, txResult.Slot)
			if err != nil {
				return nil, errors.Wrap(err, "error ParseInboundAsCall")
			} else if call != nil {
				seenCall = true
				events = append(events, &clienttypes.InboundEvent{
					SenderChainID:    senderChainID,
					Sender:           call.Sender,
					Receiver:         call.Receiver,
					TxOrigin:         call.Sender,
					Amount:           call.Amount,
					Memo:             call.Memo,
					BlockNumber:      call.Slot, // instead of using block, Solana explorer uses slot for indexing
					TxHash:           tx.Signatures[0].String(),
					Index:            0, // hardcode to 0 for Solana, not a EVM smart contract call
					CoinType:         coin.CoinType_NoAssetCall,
					Asset:            call.Asset,
					IsCrossChainCall: call.IsCrossChainCall,
					RevertOptions:    call.RevertOptions,
				})
				logger.Info().
					Str(logs.FieldMethod, "FilterInboundEvents").
					Str("signature", tx.Signatures[0].String()).
					Int("instruction", i).
					Msg("call detected")
			}
		} else {
			logger.Warn().
				Str(logs.FieldMethod, "FilterInboundEvents").
				Str("signature", tx.Signatures[0].String()).
				Int("instruction", i).
				Msg("multiple calls detected")
		}
	}

	return events, nil
}

// BuildInboundVoteMsgFromEvent builds a MsgVoteInbound from an inbound event
func (ob *Observer) BuildInboundVoteMsgFromEvent(event *clienttypes.InboundEvent) *crosschaintypes.MsgVoteInbound {
	// check if the event is processable
	if !ob.IsEventProcessable(*event) {
		return nil
	}

	options := []crosschaintypes.InboundVoteOption{crosschaintypes.WithCrossChainCall(event.IsCrossChainCall)}
	if event.RevertOptions != nil {
		options = append(options, crosschaintypes.WithSOLRevertOptions(*event.RevertOptions))
	}

	// create inbound vote message
	return crosschaintypes.NewMsgVoteInbound(
		ob.ZetacoreClient().GetKeys().GetOperatorAddress().String(),
		event.Sender,
		event.SenderChainID,
		event.Sender,
		event.Receiver,
		ob.ZetacoreClient().Chain().ChainId,
		cosmosmath.NewUint(event.Amount),
		hex.EncodeToString(event.Memo),
		event.TxHash,
		event.BlockNumber,
		0,
		event.CoinType,
		event.Asset,
		0, // not a smart contract call
		crosschaintypes.ProtocolContractVersion_V2,
		false, // not used
		crosschaintypes.InboundStatus_SUCCESS,
		crosschaintypes.ConfirmationMode_SAFE,
		options...,
	)
}

// IsEventProcessable checks if the inbound event is processable
func (ob *Observer) IsEventProcessable(event clienttypes.InboundEvent) bool {
	logFields := map[string]any{logs.FieldTx: event.TxHash}

	switch category := event.Category(); category {
	case clienttypes.InboundCategoryProcessable:
		return true
	case clienttypes.InboundCategoryDonation:
		ob.Logger().Inbound.Info().Fields(logFields).Msg("thank you rich folk for your donation!")
		return false
	case clienttypes.InboundCategoryRestricted:
		compliance.PrintComplianceLog(ob.Logger().Inbound, ob.Logger().Compliance,
			false, ob.Chain().ChainId, event.TxHash, event.Sender, event.Receiver, event.CoinType.String())
		return false
	default:
		ob.Logger().Inbound.Error().Interface("category", category).Msg("unreachable code, got InboundCategory")
		return false
	}
}
