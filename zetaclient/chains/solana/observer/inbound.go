package observer

import (
	"context"
	"encoding/hex"

	cosmosmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

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
		ob.Logger().Inbound.Err(err).Msg("error calling GetSignaturesForAddressUntil")
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
					Stringer("tx_signature", sig.Signature).
					Msg("observe inbound: skip unsupported transaction")
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
						Str("tx_signature", sigString).
						Msg("observe inbound: error filtering events, skipping")
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
			ob.Logger().Inbound.Error().
				Err(err).
				Str("tx_signature", sigString).
				Msg("observe inbound: error saving last sig")
		}

		ob.Logger().Inbound.Info().
			Str("tx_signature", sigString).
			Uint64("tx_slot", sig.Slot).
			Msg("observe inbound: last scanned sig")

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
//   - takes at least one event (the first) per token (SOL or SPL or call) per transaction.
//   - takes at most 3 events (one SOL + one SPL + one call) per transaction.
//   - ignores exceeding events.
//   - assigns indices based on instruction position in the transaction
func FilterInboundEvents(
	txResult *rpc.GetTransactionResult,
	gatewayID solana.PublicKey,
	senderChainID int64,
	logger zerolog.Logger,
) ([]*clienttypes.InboundEvent, error) {
	if txResult.Meta.Err != nil {
		return nil, errors.Errorf("transaction failed with error: %v", txResult.Meta.Err)
	}

	parser, err := NewInboundEventParser(txResult, gatewayID, senderChainID, logger)
	if err != nil {
		return nil, err
	}

	if err := parser.Parse(); err != nil {
		return nil, err
	}

	return parser.GetEvents(), nil
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
		uint64(event.Index),
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
