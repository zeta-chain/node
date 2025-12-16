package observer

import (
	"context"
	"encoding/hex"
	"strconv"

	cosmosmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	addresslookuptable "github.com/gagliardetto/solana-go/programs/address-lookup-table"
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

// MaxSignaturesPerTicker is the maximum number of signatures to process on a ticker
const MaxSignaturesPerTicker = 100

// getRPCClient attempts to extract *rpc.Client from the SolanaClient interface.
// Returns nil if the client cannot be unwrapped to *rpc.Client.
func getRPCClient(client SolanaClient) *rpc.Client {
	if rpcClient, ok := client.(*rpc.Client); ok {
		return rpcClient
	}
	if wrapper, ok := client.(interface{ UnwrapClient() any }); ok {
		if rpcClient, ok := wrapper.UnwrapClient().(*rpc.Client); ok {
			return rpcClient
		}
	}
	return nil
}

// ProcessTransactionWithAddressLookups resolves address lookup tables in a versioned transaction.
// This must be called before filtering inbound events to ensure accounts are properly resolved.
func ProcessTransactionWithAddressLookups(ctx context.Context, tx *solana.Transaction, rpcClient *rpc.Client) error {
	lookups := tx.Message.GetAddressTableLookups()
	if lookups == nil || len(lookups) == 0 {
		// No address lookup tables in this transaction
		return nil
	}

	resolutions := make(map[solana.PublicKey]solana.PublicKeySlice)
	for _, lookup := range lookups {
		tableKey := lookup.AccountKey
		altState, err := addresslookuptable.GetAddressLookupTableStateWithOpts(
			ctx,
			rpcClient,
			tableKey,
			&rpc.GetAccountInfoOpts{Commitment: rpc.CommitmentConfirmed},
		)
		if err != nil {
			return errors.Wrapf(err, "error getting address lookup table state for %s", tableKey)
		}

		if altState == nil {
			return errors.Errorf("address lookup table %s not found", tableKey)
		}

		resolutions[tableKey] = altState.Addresses
	}

	if err := tx.Message.SetAddressTables(resolutions); err != nil {
		return errors.Wrap(err, "error setting address tables")
	}

	if err := tx.Message.ResolveLookups(); err != nil {
		return errors.Wrap(err, "error resolving lookups")
	}

	return nil
}

// ObserveInbound observes the Solana chain for inbounds and post votes to zetacore.
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	chainID := ob.Chain().ChainId
	pageLimit := repo.DefaultPageLimit

	// scan from gateway 1st signature if last scanned tx is absent in the database
	// the 1st gateway signature is typically the program initialization
	if ob.LastTxScanned() == "" {
		lastSig, err := ob.solanaRepo.GetFirstSignatureForAddress(ctx, ob.gatewayID, pageLimit)
		if err != nil {
			format := "error GetFirstSignatureForAddress for chain %d address %s"
			return errors.Wrapf(err, format, chainID, ob.gatewayID)
		}
		ob.WithLastTxScanned(lastSig.String())
	}

	// query last finalized slot
	lastSlot, errSlot := ob.solanaClient.GetSlot(ctx, rpc.CommitmentFinalized)
	if errSlot != nil {
		ob.Logger().Inbound.Err(errSlot).Msg("unable to get last slot")
	}

	// get all signatures for the gateway address since last scanned signature
	lastSig := solana.MustSignatureFromBase58(ob.LastTxScanned())
	signatures, err := ob.solanaRepo.GetSignaturesForAddressUntil(ctx, ob.gatewayID, lastSig, pageLimit)
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
		ob.Logger().Inbound.Info().Int("signatures", len(signatures)).Msg("got inbound signatures")
	}

	// loop signature from oldest to latest to filter inbound events
	for i := len(signatures) - 1; i >= 0; i-- {
		sig := signatures[i]
		sigString := sig.Signature.String()

		// process successfully signature only
		if sig.Err == nil {
			txResult, err := ob.solanaRepo.GetTransaction(ctx, sig.Signature)
			switch {
			case errors.Is(err, repo.ErrUnsupportedTxVersion):
				ob.Logger().Inbound.Warn().
					Stringer("tx_signature", sig.Signature).
					Msg("observe inbound: skip unsupported transaction")
			// just save the sig to last scanned txs
			case err != nil:
				// we have to re-scan this signature on next ticker
				return errors.Wrapf(err, "error GetTransaction for sig %s", sigString)
			default:
				// Process address lookup tables before filtering events
				tx, err := txResult.Transaction.GetTransaction()
				if err == nil {
					if rpcClient := getRPCClient(ob.solanaClient); rpcClient != nil {
						if err := ProcessTransactionWithAddressLookups(ctx, tx, rpcClient); err != nil {
							ob.Logger().Inbound.Warn().
								Err(err).
								Str("tx_signature", sigString).
								Msg("observe inbound: error processing address lookup tables, continuing anyway")
						}
					}
				}

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
				if err := ob.VoteInboundEvents(ctx, events, false, false); err != nil {
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
				metrics.InboundObservationsTrackerTotal.WithLabelValues(ob.Chain().Name, strconv.FormatBool(isInternalTracker)).
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
