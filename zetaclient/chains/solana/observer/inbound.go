package observer

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"

	cosmosmath "cosmossdk.io/math"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	solanacontracts "github.com/zeta-chain/node/pkg/contracts/solana"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	solanarpc "github.com/zeta-chain/node/zetaclient/chains/solana/rpc"
	"github.com/zeta-chain/node/zetaclient/compliance"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// MaxSignaturesPerTicker is the maximum number of signatures to process on a ticker
	MaxSignaturesPerTicker = 100
)

// WatchInbound watches Solana chain for inbounds on a ticker.
func (ob *Observer) WatchInbound(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	ticker, err := clienttypes.NewDynamicTicker(
		fmt.Sprintf("Solana_WatchInbound_%d", ob.Chain().ChainId),
		ob.ChainParams().InboundTicker,
	)
	if err != nil {
		ob.Logger().Inbound.Error().Err(err).Msg("error creating ticker")
		return err
	}
	defer ticker.Stop()

	ob.Logger().Inbound.Info().Msgf("WatchInbound started for chain %d", ob.Chain().ChainId)
	sampledLogger := ob.Logger().Inbound.Sample(&zerolog.BasicSampler{N: 10})

	for {
		select {
		case <-ticker.C():
			if !app.IsInboundObservationEnabled() {
				sampledLogger.Info().
					Msgf("WatchInbound: inbound observation is disabled for chain %d", ob.Chain().ChainId)
				continue
			}
			err := ob.ObserveInbound(ctx)
			if err != nil {
				ob.Logger().Inbound.Err(err).Msg("WatchInbound: observeInbound error")
			}
		case <-ob.StopChannel():
			ob.Logger().Inbound.Info().Msgf("WatchInbound stopped for chain %d", ob.Chain().ChainId)
			return nil
		}
	}
}

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

	// get all signatures for the gateway address since last scanned signature
	lastSig := solana.MustSignatureFromBase58(ob.LastTxScanned())
	signatures, err := solanarpc.GetSignaturesForAddressUntil(ctx, ob.solClient, ob.gatewayID, lastSig, pageLimit)
	if err != nil {
		ob.Logger().Inbound.Err(err).Msg("error GetSignaturesForAddressUntil")
		return err
	}
	if len(signatures) > 0 {
		ob.Logger().Inbound.Info().Msgf("ObserveInbound: got %d signatures for chain %d", len(signatures), chainID)
	}

	// loop signature from oldest to latest to filter inbound events
	for i := len(signatures) - 1; i >= 0; i-- {
		sig := signatures[i]
		sigString := sig.Signature.String()

		// process successfully signature only
		if sig.Err == nil {
			maxTxVersion := rpc.MaxSupportedTransactionVersion0
			txResult, skip, err := solanarpc.GetTransactionWithMaxVersion(
				ctx,
				ob.solClient,
				sig.Signature,
				&maxTxVersion,
			)
			if err != nil {
				// we have to re-scan this signature on next ticker
				return errors.Wrapf(err, "error GetTransaction for chain %d sig %s", chainID, sigString)
			}

			switch {
			case skip:
				ob.Logger().Inbound.Warn().Msgf("ObserveInbound: skip unsupported transaction sig %s", sigString)
			default:
				// filter inbound events and vote
				err = ob.FilterInboundEventsAndVote(ctx, txResult)
				if err != nil {
					// we have to re-scan this signature on next ticker
					return errors.Wrapf(err, "error FilterInboundEventAndVote for chain %d sig %s", chainID, sigString)
				}
			}
		}

		// signature scanned; save last scanned signature to both memory and db, ignore db error
		if err := ob.SaveLastTxScanned(sigString, sig.Slot); err != nil {
			ob.Logger().
				Inbound.Error().
				Err(err).
				Msgf("ObserveInbound: error saving last sig %s for chain %d", sigString, chainID)
		}
		ob.Logger().
			Inbound.Info().
			Msgf("ObserveInbound: last scanned sig is %s for chain %d in slot %d", sigString, chainID, sig.Slot)

		// take a rest if max signatures per ticker is reached
		if len(signatures)-i >= MaxSignaturesPerTicker {
			break
		}
	}

	return nil
}

// FilterInboundEventsAndVote filters inbound events from a txResult and post a vote.
func (ob *Observer) FilterInboundEventsAndVote(ctx context.Context, txResult *rpc.GetTransactionResult) error {
	// filter inbound events from txResult
	events, err := ob.FilterInboundEvents(txResult)
	if err != nil {
		return errors.Wrapf(err, "error FilterInboundEvent")
	}

	// build inbound vote message from events and post to zetacore
	for _, event := range events {
		msg := ob.BuildInboundVoteMsgFromEvent(event)
		if msg != nil {
			_, err = ob.PostVoteInbound(ctx, msg, zetacore.PostVoteInboundExecutionGasLimit)
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
func (ob *Observer) FilterInboundEvents(txResult *rpc.GetTransactionResult) ([]*clienttypes.InboundEvent, error) {
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
	events := make([]*clienttypes.InboundEvent, 0)

	// loop through instruction list to filter the 1st valid event
	for i, instruction := range tx.Message.Instructions {
		// get the program ID
		programPk, err := tx.Message.Program(instruction.ProgramIDIndex)
		if err != nil {
			ob.Logger().
				Inbound.Err(err).
				Msgf("no program found at index %d for sig %s", instruction.ProgramIDIndex, tx.Signatures[0])
			continue
		}

		// skip instructions that are irrelevant to the gateway program invocation
		if !programPk.Equals(ob.gatewayID) {
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
					SenderChainID: ob.Chain().ChainId,
					Sender:        deposit.Sender,
					Receiver:      deposit.Sender, // receiver is pulled out from memo
					TxOrigin:      deposit.Sender,
					Amount:        deposit.Amount,
					Memo:          deposit.Memo,
					BlockNumber:   deposit.Slot, // instead of using block, Solana explorer uses slot for indexing
					TxHash:        tx.Signatures[0].String(),
					Index:         0, // hardcode to 0 for Solana, not a EVM smart contract call
					CoinType:      coin.CoinType_Gas,
					Asset:         deposit.Asset,
				})
				ob.Logger().Inbound.Info().
					Msgf("FilterInboundEvents: deposit detected in sig %s instruction %d", tx.Signatures[0], i)
			}
		} else {
			ob.Logger().Inbound.Warn().
				Msgf("FilterInboundEvents: multiple deposits detected in sig %s instruction %d", tx.Signatures[0], i)
		}

		// try parsing the instruction as a 'deposit_spl_token' if not seen yet
		if !seenDepositSPL {
			deposit, err := solanacontracts.ParseInboundAsDepositSPL(tx, i, txResult.Slot)
			if err != nil {
				return nil, errors.Wrap(err, "error ParseInboundAsDepositSPL")
			} else if deposit != nil {
				seenDepositSPL = true
				events = append(events, &clienttypes.InboundEvent{
					SenderChainID: ob.Chain().ChainId,
					Sender:        deposit.Sender,
					Receiver:      deposit.Sender, // receiver is pulled out from memo
					TxOrigin:      deposit.Sender,
					Amount:        deposit.Amount,
					Memo:          deposit.Memo,
					BlockNumber:   deposit.Slot, // instead of using block, Solana explorer uses slot for indexing
					TxHash:        tx.Signatures[0].String(),
					Index:         0, // hardcode to 0 for Solana, not a EVM smart contract call
					CoinType:      coin.CoinType_ERC20,
					Asset:         deposit.Asset,
				})
				ob.Logger().Inbound.Info().
					Msgf("FilterInboundEvents: SPL deposit detected in sig %s instruction %d", tx.Signatures[0], i)
			}
		} else {
			ob.Logger().Inbound.Warn().
				Msgf("FilterInboundEvents: multiple SPL deposits detected in sig %s instruction %d", tx.Signatures[0], i)
		}
	}

	return events, nil
}

// BuildInboundVoteMsgFromEvent builds a MsgVoteInbound from an inbound event
func (ob *Observer) BuildInboundVoteMsgFromEvent(event *clienttypes.InboundEvent) *crosschaintypes.MsgVoteInbound {
	// compliance check. Return nil if the inbound contains restricted addresses
	if compliance.DoesInboundContainsRestrictedAddress(event, ob.Logger()) {
		return nil
	}

	// donation check
	if bytes.Equal(event.Memo, []byte(constant.DonationMessage)) {
		ob.Logger().Inbound.Info().
			Msgf("thank you rich folk for your donation! tx %s chain %d", event.TxHash, event.SenderChainID)
		return nil
	}

	return zetacore.GetInboundVoteMessage(
		event.Sender,
		event.SenderChainID,
		event.Sender,
		event.Sender,
		ob.ZetacoreClient().Chain().ChainId,
		cosmosmath.NewUint(event.Amount),
		hex.EncodeToString(event.Memo),
		event.TxHash,
		event.BlockNumber,
		0,
		event.CoinType,
		event.Asset,
		ob.ZetacoreClient().GetKeys().GetOperatorAddress().String(),
		0, // not a smart contract call
	)
}
