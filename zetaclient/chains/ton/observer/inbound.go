package observer

import (
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// maximum number of transactions to process on a ticker
	// TODO: move to config
	// https://github.com/zeta-chain/node/issues/3086
	maxTransactionsPerTick = 100
)

// ObserveGateway observes Gateway's account for new transactions.
// Due to TON's architecture we have to scan for all net-new transactions.
// The main purpose is to observe inbounds from TON.
// Note that we might also have *outbounds* here (if a signer broadcasts a tx, it will be observed here).
func (ob *Observer) ObserveGateway(ctx context.Context) error {
	if err := ob.ensureLastScannedTX(ctx); err != nil {
		return errors.Wrap(err, "unable to ensure last scanned tx")
	}

	// extract logicalTime and tx hash from last scanned tx
	lt, hashBits, err := liteapi.TransactionHashFromString(ob.LastTxScanned())
	if err != nil {
		return errors.Wrapf(err, "unable to parse last scanned tx %q", ob.LastTxScanned())
	}

	txs, err := ob.client.GetTransactionsSince(ctx, ob.gateway.AccountID(), lt, hashBits)
	if err != nil {
		return errors.Wrap(err, "unable to get transactions")
	}

	switch {
	case len(txs) == 0:
		// noop
		return nil
	case len(txs) > maxTransactionsPerTick:
		ob.Logger().Inbound.Info().
			Msgf("observeGateway: got %d transactions. Taking first %d", len(txs), maxTransactionsPerTick)

		txs = txs[:maxTransactionsPerTick]
	default:
		ob.Logger().Inbound.Info().Msgf("observeGateway: got %d transactions", len(txs))
	}

	for i := range txs {
		var skip bool

		tx, err := ob.gateway.ParseTransaction(txs[i])
		switch {
		case errors.Is(err, toncontracts.ErrParse) || errors.Is(err, toncontracts.ErrUnknownOp):
			skip = true
		case err != nil:
			// should not happen
			return errors.Wrap(err, "unexpected error")
		case tx.ExitCode != 0:
			skip = true
			ob.Logger().Inbound.Warn().Fields(txLogFields(tx)).Msg("observeGateway: observed a failed tx")
		}

		if skip {
			tx = &toncontracts.Transaction{Transaction: txs[i]}
			txHash := liteapi.TransactionToHashString(tx.Transaction)
			ob.Logger().Inbound.Warn().Str("transaction.hash", txHash).Msg("observeGateway: skipping tx")
			ob.setLastScannedTX(tx)
			continue
		}

		// Should not happen
		//goland:noinspection GoDfaConstantCondition
		if tx == nil {
			return errors.New("tx is nil")
		}

		// As we might have outbounds here, let's ensure outbound tracker.
		// TON signer broadcasts ExtInMsgInfo with `src=null, dest=gateway`, so it will be observed here
		if tx.IsOutbound() {
			if err = ob.addOutboundTracker(ctx, tx); err != nil {
				ob.Logger().Inbound.
					Error().Err(err).
					Fields(txLogFields(tx)).
					Msg("observeGateway: unable to add outbound tracker")

				return errors.Wrap(err, "unable to add outbound tracker")
			}

			ob.setLastScannedTX(tx)
			continue
		}

		// Ok, let's process a new inbound tx
		if _, err := ob.voteInbound(ctx, tx); err != nil {
			ob.Logger().Inbound.
				Error().Err(err).
				Fields(txLogFields(tx)).
				Msg("observeGateway: unable to vote for inbound tx")

			return errors.Wrapf(err, "unable to vote for inbound tx %s", tx.Hash().Hex())
		}

		ob.setLastScannedTX(tx)
	}

	return nil
}

// ObserveInboundTrackers handles adhoc trackers that were somehow missed by
func (ob *Observer) ObserveInboundTrackers(ctx context.Context) error {
	trackers, err := ob.ZetacoreClient().GetInboundTrackersForChain(ctx, ob.Chain().ChainId)
	if err != nil {
		return errors.Wrap(err, "unable to get inbound trackers")
	}

	// noop
	if len(trackers) == 0 {
		return nil
	}

	gatewayAccountID := ob.gateway.AccountID()

	// a single error should not block other trackers
	for _, tracker := range trackers {
		txHash := tracker.TxHash

		lt, hash, err := liteapi.TransactionHashFromString(txHash)
		if err != nil {
			ob.logSkippedTracker(txHash, "unable_to_parse_hash", err)
			continue
		}

		raw, err := ob.client.GetTransaction(ctx, gatewayAccountID, lt, hash)
		if err != nil {
			ob.logSkippedTracker(txHash, "unable_to_get_tx", err)
			continue
		}

		tx, err := ob.gateway.ParseTransaction(raw)

		switch {
		case errors.Is(err, toncontracts.ErrParse) || errors.Is(err, toncontracts.ErrUnknownOp):
			ob.logSkippedTracker(txHash, "unrelated_tx", err)
			continue
		case err != nil:
			// should not happen
			ob.logSkippedTracker(txHash, "unexpected_error", err)
			continue
		case tx.ExitCode != 0:
			ob.logSkippedTracker(txHash, "failed_tx", nil)
			continue
		case tx.IsOutbound():
			ob.logSkippedTracker(txHash, "outbound_tx", nil)
			continue
		}

		if _, err := ob.voteInbound(ctx, tx); err != nil {
			ob.logSkippedTracker(txHash, "vote_failed", err)
			continue
		}
	}

	return nil
}

// Sends PostVoteInbound to zetacore
func (ob *Observer) voteInbound(ctx context.Context, tx *toncontracts.Transaction) (string, error) {
	// noop
	if tx.Operation == toncontracts.OpDonate {
		ob.Logger().Inbound.Info().Fields(txLogFields(tx)).Msg("Thank you rich folk for your donation!")
		return "", nil
	}

	// TODO: Add compliance check
	// https://github.com/zeta-chain/node/issues/2916

	blockHeader, err := ob.client.GetBlockHeader(ctx, tx.BlockID, 0)
	if err != nil {
		return "", errors.Wrapf(err, "unable to get block header %s", tx.BlockID.String())
	}

	sender, amount, memo, err := extractInboundData(tx)
	if err != nil {
		return "", err
	}

	seqno := blockHeader.MinRefMcSeqno

	return ob.voteDeposit(ctx, tx, sender, amount, memo, seqno)
}

// extractInboundData parses Gateway tx into deposit (TON sender, amount, memo)
func extractInboundData(tx *toncontracts.Transaction) (string, math.Uint, []byte, error) {
	switch tx.Operation {
	case toncontracts.OpDeposit:
		d, err := tx.Deposit()
		if err != nil {
			return "", math.NewUint(0), nil, err
		}

		return d.Sender.ToRaw(), d.Amount, d.Memo(), nil
	case toncontracts.OpDepositAndCall:
		d, err := tx.DepositAndCall()
		if err != nil {
			return "", math.NewUint(0), nil, err
		}

		return d.Sender.ToRaw(), d.Amount, d.Memo(), nil
	default:
		return "", math.NewUint(0), nil, fmt.Errorf("unknown operation %d", tx.Operation)
	}
}

func (ob *Observer) voteDeposit(
	ctx context.Context,
	tx *toncontracts.Transaction,
	sender string,
	amount math.Uint,
	memo []byte,
	seqno uint32,
) (string, error) {
	const (
		eventIndex    = 0 // not a smart contract call
		coinType      = coin.CoinType_Gas
		asset         = "" // empty for gas coin
		gasLimit      = 0
		retryGasLimit = zetacore.PostVoteInboundExecutionGasLimit
	)

	var (
		operatorAddress = ob.ZetacoreClient().GetKeys().GetOperatorAddress()
		inboundHash     = liteapi.TransactionHashToString(tx.Lt, ton.Bits256(tx.Hash()))
	)

	// TODO: use protocol contract v2 for deposit
	// https://github.com/zeta-chain/node/issues/2967

	msg := zetacore.GetInboundVoteMessage(
		sender,
		ob.Chain().ChainId,
		sender,
		sender,
		ob.ZetacoreClient().Chain().ChainId,
		amount,
		hex.EncodeToString(memo),
		inboundHash,
		uint64(seqno),
		gasLimit,
		coinType,
		asset,
		operatorAddress.String(),
		eventIndex,
	)

	return ob.PostVoteInbound(ctx, msg, retryGasLimit)
}

func (ob *Observer) ensureLastScannedTX(ctx context.Context) error {
	// noop
	if ob.LastTxScanned() != "" {
		return nil
	}

	rawTX, _, err := ob.client.GetFirstTransaction(ctx, ob.gateway.AccountID())
	if err != nil {
		return err
	}

	ob.setLastScannedTX(&toncontracts.Transaction{Transaction: *rawTX})

	return nil
}

func (ob *Observer) setLastScannedTX(tx *toncontracts.Transaction) {
	txHash := liteapi.TransactionToHashString(tx.Transaction)

	ob.WithLastTxScanned(txHash)

	if err := ob.WriteLastTxScannedToDB(txHash); err != nil {
		ob.Logger().Inbound.Error().
			Err(err).
			Fields(txLogFields(tx)).
			Msgf("setLastScannedTX: unable to WriteLastTxScannedToDB")

		return
	}

	ob.Logger().Inbound.Info().
		Fields(txLogFields(tx)).
		Msg("setLastScannedTX: WriteLastTxScannedToDB")
}

func (ob *Observer) logSkippedTracker(hash string, reason string, err error) {
	ob.Logger().Inbound.Warn().
		Str("transaction.hash", hash).
		Str("skip_reason", reason).
		Err(err).
		Msg("Skipping tracker")
}

func txLogFields(tx *toncontracts.Transaction) map[string]any {
	return map[string]any{
		"transaction.hash":           liteapi.TransactionToHashString(tx.Transaction),
		"transaction.ton.lt":         tx.Lt,
		"transaction.ton.hash":       tx.Hash().Hex(),
		"transaction.ton.block_id":   tx.BlockID.BlockID.String(),
		"transaction.ton.is_inbound": tx.IsInbound(),
		"transaction.ton.op_code":    tx.Operation,
		"transaction.ton.exit_code":  tx.ExitCode,
	}
}
