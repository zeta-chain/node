package observer

import (
	"context"
	"encoding/hex"
	"strconv"

	cosmosmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/metrics"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// maxSignaturesPerTicker is the maximum number of signatures to process on a ticker.
const maxSignaturesPerTicker = 100

// ObserveInbound observes the Solana chain for inbounds and posts votes to zetacore.
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	logger := ob.Logger().Inbound

	// Scan from gateway's 1st signature if last scanned transaction is absent in the database.
	// The 1st gateway signature is typically the program initialization.
	if ob.LastTxScanned() == "" {
		sig, err := ob.solanaRepo.GetFirstSignature(ctx)
		if err != nil {
			return err
		}
		ob.WithLastTxScanned(sig.String())
	}

	// Get last finalized slot.
	lastSlot, errSlot := ob.solanaRepo.GetSlot(ctx, rpc.CommitmentFinalized)
	if errSlot != nil {
		logger.Err(errSlot).Msg("failed to get last finalized slot")
	}

	// Get all signatures for the gateway address since the last scanned signature.
	lastSig := solana.MustSignatureFromBase58(ob.LastTxScanned())
	sigs, err := ob.solanaRepo.GetSignaturesSince(ctx, lastSig)
	if err != nil {
		logger.Err(err).Send()
		return err
	}

	if len(sigs) == 0 {
		// Update metrics if there are no new signatures.
		if errSlot == nil {
			ob.WithLastBlockScanned(lastSlot)
		}
	} else {
		logger.Info().Int("signatures", len(sigs)).Msg("got inbound signatures")
	}

	// Iterate over the signatures from oldest to latest to filter inbound events.
	for i := len(sigs) - 1; i >= 0; i-- {
		sig := sigs[i]

		// Process only successfull transactions.
		if sig.Err == nil {
			txResult, err := ob.solanaRepo.GetTransaction(ctx, sig.Signature, rpc.CommitmentFinalized)
			if errors.Is(err, repo.ErrUnsupportedTxVersion) {
				logger.Warn().
					Stringer("tx_signature", sig.Signature).
					Msg("skipping unsupported transaction")
			} else if err != nil {
				return err
			} else {
				events, err := FilterInboundEvents(txResult, ob.gatewayID, ob.Chain().ChainId, logger)
				if err != nil {
					// Log the error but continue processing other transactions
					logger.Error().
						Err(err).
						Stringer("tx_signature", sig.Signature).
						Msg("observe inbound: error filtering events, skipping")
					continue
				}

				err = ob.VoteInboundEvents(ctx, events, false, false)
				if err != nil {
					return errors.Wrapf(err,
						"error voting on events for transaction %s, will retry",
						sig.Signature.String(),
					)
				}
			}
		}

		// signature scanned; save last scanned signature to both memory and db, ignore db error
		err = ob.SaveLastTxScanned(sig.Signature.String(), sig.Slot)
		if err != nil {
			logger.Error().
				Err(err).
				Stringer("tx_signature", sig.Signature).
				Msg("error saving last signature")
		}

		logger.Info().
			Stringer("tx_signature", sig.Signature).
			Uint64("tx_slot", sig.Slot).
			Msg("last scanned signature")

		// Take a rest if the maximum number of signatures per ticker has been reached.
		if len(sigs)-i >= maxSignaturesPerTicker {
			break
		}
	}

	return nil
}

// VoteInboundEvents posts votes for inbound events to zetacore.
func (ob *Observer) VoteInboundEvents(
	ctx context.Context,
	events []*clienttypes.InboundEvent,
	fromTracker bool,
	isInternalTracker bool,
) error {
	for _, event := range events {
		msg := ob.BuildInboundVoteMsgFromEvent(event)
		if msg != nil {
			if fromTracker {
				metrics.InboundObservationsTrackerTotal.
					WithLabelValues(ob.Chain().Name, strconv.FormatBool(isInternalTracker)).
					Inc()
			} else {
				metrics.InboundObservationsBlockScanTotal.WithLabelValues(ob.Chain().Name).Inc()
			}
			_, err := ob.ZetaRepo().VoteInbound(ctx,
				ob.Logger().Inbound,
				msg,
				zetacore.PostVoteInboundExecutionGasLimit,
				ob.WatchMonitoringError,
			)
			if err != nil {
				return err
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
		ob.ZetaRepo().GetOperatorAddress(),
		event.Sender,
		event.SenderChainID,
		event.Sender,
		event.Receiver,
		ob.ZetaRepo().ZetaChain().ChainId,
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
			false, ob.Chain().ChainId, event.TxHash, event.Sender, event.Receiver, &event.CoinType)
		return false
	default:
		ob.Logger().Inbound.Error().Interface("category", category).Msg("unreachable code, got InboundCategory")
		return false
	}
}
