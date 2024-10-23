package observer

import (
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/math"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/pkg/ticker"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const (
	// MaxTransactionsPerTick is the maximum number of transactions to process on a ticker
	MaxTransactionsPerTick = 100
)

// watchInbound watches for new txs to Gateway's account.
func (ob *Observer) watchInbound(ctx context.Context) error {
	app, err := zctx.FromContext(ctx)
	if err != nil {
		return err
	}

	var (
		initialInterval = ticker.SecondsFromUint64(ob.ChainParams().InboundTicker)
		sampledLogger   = ob.Logger().Inbound.Sample(&zerolog.BasicSampler{N: 10})
	)

	ob.Logger().Inbound.Info().Msgf("WatchInbound started")

	task := func(ctx context.Context, t *ticker.Ticker) error {
		if !app.IsInboundObservationEnabled() {
			sampledLogger.Info().Msg("WatchInbound: inbound observation is disabled")
			return nil
		}

		if err := ob.observeGateway(ctx); err != nil {
			ob.Logger().Inbound.Err(err).Msg("WatchInbound: observeInbound error")
		}

		newInterval := ticker.SecondsFromUint64(ob.ChainParams().InboundTicker)
		t.SetInterval(newInterval)

		return nil
	}

	return ticker.Run(
		ctx,
		initialInterval,
		task,
		ticker.WithStopChan(ob.StopChannel()),
		ticker.WithLogger(ob.Logger().Inbound, "WatchInbound"),
	)
}

// observeGateway observes Gateway's account for new transactions.
// Due to TON architecture we have to scan for all net-new transactions.
// The main purpose is to observe inbounds from TON.
// Note that we might also have *outbounds* here (if a signer broadcasts a tx, it will be observed here).
func (ob *Observer) observeGateway(ctx context.Context) error {
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
	case len(txs) > MaxTransactionsPerTick:
		ob.Logger().Inbound.Info().
			Msgf("observeGateway: got %d transactions. Taking first %d", len(txs), MaxTransactionsPerTick)

		txs = txs[:MaxTransactionsPerTick]
	default:
		ob.Logger().Inbound.Info().Msgf("observeGateway: got %d transactions", len(txs))
	}

	for i := range txs {
		var skip bool

		tx, err := ob.gateway.ParseTransaction(txs[i])
		switch {
		case errors.Is(err, toncontracts.ErrParse), errors.Is(err, toncontracts.ErrUnknownOp):
			skip = true
		case err != nil:
			return errors.Wrap(err, "unable to parse tx")
		case tx.ExitCode != 0:
			skip = true
			ob.Logger().Inbound.Warn().Fields(txLogFields(tx)).Msg("observeGateway: observed a failed tx")
		}

		if skip {
			tx = &toncontracts.Transaction{Transaction: txs[i]}
			txHash := liteapi.TransactionToHashString(&tx.Transaction)
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

// addOutboundTracker publishes outbound tracker to Zetacore.
// In most cases will be a noop because the tracker is already published by the signer.
// See Signer{}.trackOutbound(...) for more details.
func (ob *Observer) addOutboundTracker(ctx context.Context, tx *toncontracts.Transaction) error {
	w, err := tx.Withdrawal()
	if err != nil {
		return errors.Wrap(err, "tx is not a withdrawal")
	}

	signer, err := w.Signer()
	switch {
	case err != nil:
		return errors.Wrap(err, "unable to get withdrawal signer")
	case signer != ob.TSS().EVMAddress():
		ob.Logger().Inbound.Warn().
			Fields(txLogFields(tx)).
			Str("transaction.ton.signer", signer.String()).
			Msg("observeGateway: addOutboundTracker: withdrawal signer is not TSS. Skipping")

		return nil
	}

	var (
		chainID = ob.Chain().ChainId
		nonce   = uint64(w.Seqno)
		hash    = liteapi.TransactionToHashString(&tx.Transaction)
	)

	// note it has a check for noop
	_, err = ob.
		ZetacoreClient().
		AddOutboundTracker(ctx, chainID, nonce, hash, nil, "", 0)

	return err
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
	txHash := liteapi.TransactionToHashString(&tx.Transaction)

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

func txLogFields(tx *toncontracts.Transaction) map[string]any {
	return map[string]any{
		"transaction.hash":           liteapi.TransactionToHashString(&tx.Transaction),
		"transaction.ton.lt":         tx.Lt,
		"transaction.ton.hash":       tx.Hash().Hex(),
		"transaction.ton.block_id":   tx.BlockID.BlockID.String(),
		"transaction.ton.is_inbound": tx.IsInbound(),
		"transaction.ton.op_code":    tx.Operation,
		"transaction.ton.exit_code":  tx.ExitCode,
	}
}
