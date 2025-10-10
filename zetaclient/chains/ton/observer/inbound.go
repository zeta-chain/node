package observer

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/ton/encoder"
	"github.com/zeta-chain/node/zetaclient/compliance"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// ------------------------------------------------------------------------------------------------
// ObserveInbound
// ------------------------------------------------------------------------------------------------

// ObserveInbounds observes the gateway for new transactions, inbounds and outbounds alike.
//
// The name "ObserveInbounds" is used for consistency with other chains.
// Also, the main purpose of this function is to indeed observe TON inbounds.
// However, when a signer broadcasts a transaction (an outbound), it also gets observed here.
// This happens because of TON's architecture (we have to scan for all net-new transactions).
func (ob *Observer) ObserveInbounds(ctx context.Context) error {
	logger := ob.Logger().Inbound

	lastScannedTx, err := ob.getLastScannedTransaction(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get the last scanned transaction")
	}

	rawTxs, err := ob.tonRepo.GetNextTransactions(ctx, logger, lastScannedTx)
	if err != nil {
		return errors.Wrap(err, "unable to get the next transactions to be processed")
	}

	for _, rawTx := range rawTxs {
		tx, err := ob.parseTransaction(rawTx)
		if err != nil {
			return err
		}

		// Skip unparseable transactions.
		if tx == nil {
			tx = &toncontracts.Transaction{Transaction: rawTx}
			encodedTx := encoder.EncodeTx(tx.Transaction)
			ob.saveLastScannedTransaction(tx)
			logger.Warn().Str(logs.FieldTx, encodedTx).Msg("skipping unparseable transaction")
			continue
		}

		txLogger := logger.With().Fields(txLogFields(tx)).Logger()

		// Skip failed transactions.
		if tx.ExitCode != 0 {
			ob.saveLastScannedTransaction(tx)
			txLogger.Warn().Msg("skipping failed transaction")
			continue
		}

		// Process outbounds and inbounds.
		switch {
		case tx.IsOutbound():
			err := ob.addOutboundTracker(ctx, tx)
			if err != nil {
				msg := "unable to add outbound tracker"
				txLogger.Error().Err(err).Msg(msg)
				return errors.Wrapf(err, "%s: %s", msg, tx.Hash().Hex())
			}
		case tx.IsInbound():
			err := ob.voteInbound(ctx, tx)
			if err != nil {
				msg := "unable to vote for inbound transaction"
				txLogger.Error().Err(err).Msg(msg)
				return errors.Wrapf(err, "%s: %s", msg, tx.Hash().Hex())
			}
		default:
			return errors.New("unreachable code (internal error)")
		}

		ob.saveLastScannedTransaction(tx)
	}

	return nil
}

// getLastScannedTransaction returns the last scanned transaction from the database.
//
// If there is no transaction in the database, it queries the blockchain for 20th most recent
// transaction and (arbitrarily) returns it as the last scanned transaction.
func (ob *Observer) getLastScannedTransaction(ctx context.Context) (string, error) {
	encodedTx := ob.LastTxScanned()
	if encodedTx != "" {
		return encodedTx, nil
	}

	const limit = 20 // arbitrary
	tx, err := ob.tonRepo.GetTransactionByIndex(ctx, limit)
	if err != nil {
		return "", errors.Wrap(err, "unable to query the blockchain")
	}

	encodedTx = encoder.EncodeTx(*tx)
	ob.WithLastTxScanned(encodedTx)
	return encodedTx, nil
}

// saveLastScannedTransaction sets the last scanned transaction and stores it in the database.
func (ob *Observer) saveLastScannedTransaction(tx *toncontracts.Transaction) {
	logger := ob.Logger().Inbound.With().Fields(txLogFields(tx)).Logger()

	encodedHash := encoder.EncodeTx(tx.Transaction)
	ob.WithLastTxScanned(encodedHash)

	err := ob.WriteLastTxScannedToDB(encodedHash)
	if err != nil {
		logger.Error().Err(err).Msg("error writing last scanned transaction to the database")
		return
	}

	logger.Info().Msg("last scanned transaction saved to the database")
}

// addOutboundTracker adds an outbound tracker to zetacore.
// In most cases it does nothing because signers usually add trackers first.
func (ob *Observer) addOutboundTracker(ctx context.Context, tx *toncontracts.Transaction) error {
	logger := ob.Logger().Inbound.With().Fields(txLogFields(tx)).Logger()

	auth, err := tx.OutboundAuth()
	if err != nil {
		return err
	}

	txSigner := auth.Signer
	tssSigner := ob.TSS().PubKey().AddressEVM()

	if txSigner != tssSigner {
		logger.Warn().Stringer("signer", txSigner).Msg("skipping transaction; signer is not TSS")
		return nil
	}

	nonce := uint64(auth.Seqno)
	hash := encoder.EncodeTx(tx.Transaction)

	_, err = ob.ZetaRepo().PostOutboundTracker(ctx, logger, nonce, hash)

	return err
}

// ------------------------------------------------------------------------------------------------
// ProcessInboundTrackers
// ------------------------------------------------------------------------------------------------

// ProcessInboundTrackers processes inbound trackers (ad hoc inbounds that were missed).
//
// It processes each tracker individually (continues executing despite any errors so it does not
// block other trackers).
func (ob *Observer) ProcessInboundTrackers(ctx context.Context) error {
	trackers, err := ob.ZetaRepo().GetInboundTrackers(ctx)
	if err != nil {
		return err
	}

	return ob.observeInboundTrackers(ctx, trackers, false)
}

// ProcessInternalTrackers processes internal inbound trackers
func (ob *Observer) ProcessInternalTrackers(ctx context.Context) error {
	trackers := ob.GetInboundInternalTrackers(ctx, time.Now())
	if len(trackers) > 0 {
		ob.Logger().Inbound.Info().Int("total_count", len(trackers)).Msg("processing internal trackers")
	}

	return ob.observeInboundTrackers(ctx, trackers, true)
}

// observeInboundTrackers observes given inbound trackers
func (ob *Observer) observeInboundTrackers(
	ctx context.Context,
	trackers []types.InboundTracker,
	isInternal bool,
) error {
	// take at most MaxInternalTrackersPerScan for each scan
	if len(trackers) > config.MaxInboundTrackersPerScan {
		trackers = trackers[:config.MaxInboundTrackersPerScan]
	}

	logSkippedTracker := func(hash string, reason string, err error) {
		ob.Logger().Inbound.Warn().
			Err(err).
			Str(logs.FieldTx, hash).
			Str("reason", reason).
			Bool("is_internal", isInternal).
			Msg("skipping inbound tracker")
	}

	for _, tracker := range trackers {
		encodedHash := tracker.TxHash

		raw, err := ob.tonRepo.GetTransactionByHash(ctx, encodedHash)
		if err != nil {
			logSkippedTracker(encodedHash, "query failed", err)
			continue
		}

		tx, err := ob.parseTransaction(*raw)
		if err != nil {
			logSkippedTracker(encodedHash, "unexpected parsing error", err)
			continue
		}
		if tx == nil {
			logSkippedTracker(encodedHash, "unparseable transaction", nil)
			continue
		}

		if tx.ExitCode != 0 {
			logSkippedTracker(encodedHash, "failed transaction", nil)
			continue
		}

		if tx.IsOutbound() {
			logSkippedTracker(encodedHash, "outbound transaction", nil)
			continue
		}

		err = ob.voteInbound(ctx, tx)
		if err != nil {
			logSkippedTracker(encodedHash, "vote failed", err)
			continue
		}
	}

	return nil
}

// ------------------------------------------------------------------------------------------------
// Shared auxiliary functions
// ------------------------------------------------------------------------------------------------

// parseTransaction parses a TON transaction.
//
// It only returns an error when it encounters unexpected parsing errors.
// Otherwise, it returns nil for both fields.
func (ob *Observer) parseTransaction(raw ton.Transaction) (*toncontracts.Transaction, error) {
	tx, err := ob.gateway.ParseTransaction(raw)
	if err != nil {
		switch {
		case errors.Is(err, toncontracts.ErrParse):
			fallthrough
		case errors.Is(err, toncontracts.ErrUnknownOp):
			return nil, nil
		default:
			return nil, errors.Wrap(err, "unexpected error")
		}
	}

	if tx == nil {
		return nil, errors.New("transaction is nil") // should not happen
	}

	return tx, nil
}

// voteInbound sends a VoteInbound message to zetacore.
// It handles donations and non-compliant inbounds.
func (ob *Observer) voteInbound(ctx context.Context, tx *toncontracts.Transaction) error {
	logger := ob.Logger().Inbound.With().Fields(txLogFields(tx)).Logger()

	// Skip donations.
	if tx.Operation == toncontracts.OpDonate {
		logger.Info().Msg("thank you rich folk for your donation!")
		return nil
	}

	inbound, err := newInbound(tx)
	if err != nil {
		return errors.Wrap(err, "unable to extract inbound data")
	}

	// Don't vote for inbounds that are not compliant.
	if !inbound.isCompliant() {
		compliance.PrintComplianceLog(
			ob.Logger().Inbound,
			ob.Logger().Compliance,
			false,
			ob.Chain().ChainId,
			encoder.EncodeTx(inbound.tx.Transaction),
			inbound.sender.ToRaw(),
			inbound.receiver.Hex(),
			&inbound.coinType,
		)
		return nil
	}

	operatorAddress := ob.ZetaRepo().GetOperatorAddress()
	senderChain := ob.Chain().ChainId
	zetaChain := ob.ZetaRepo().ZetaChain().ChainId

	logger = ob.Logger().Inbound
	msg := inbound.intoVoteMessage(operatorAddress, senderChain, zetaChain)
	_, err = ob.ZetaRepo().
		VoteInbound(ctx, logger, msg, zetacore.PostVoteInboundExecutionGasLimit, ob.WatchMonitoringError)
	if err != nil {
		return err
	}

	return nil
}

func txLogFields(tx *toncontracts.Transaction) map[string]any {
	return map[string]any{
		logs.FieldTx: encoder.EncodeTx(tx.Transaction),
		"is_inbound": tx.IsInbound(),
		"op_code":    tx.Operation,
		"exit_code":  tx.ExitCode,
	}
}

// ------------------------------------------------------------------------------------------------
// Inbound
// ------------------------------------------------------------------------------------------------

// Inbound represents a TON inbound deposit.
type Inbound struct {
	tx *toncontracts.Transaction

	sender   ton.AccountID
	receiver eth.Address

	amount   math.Uint
	coinType coin.CoinType

	message []byte

	isContractCall bool
}

// newInbound creates a new Inbound from a toncontracts.Transaction.
func newInbound(tx *toncontracts.Transaction) (*Inbound, error) {
	inbound := &Inbound{tx: tx}

	switch tx.Operation {
	case toncontracts.OpDeposit:
		castTx, err := tx.Deposit()
		if err != nil {
			return nil, err
		}
		inbound.sender = castTx.Sender
		inbound.receiver = castTx.Recipient
		inbound.amount = castTx.Amount
		inbound.coinType = coin.CoinType_Gas
		inbound.message = []byte{}
		inbound.isContractCall = false
	case toncontracts.OpDepositAndCall:
		castTx, err := tx.DepositAndCall()
		if err != nil {
			return nil, err
		}
		inbound.sender = castTx.Sender
		inbound.receiver = castTx.Recipient
		inbound.amount = castTx.Amount
		inbound.coinType = coin.CoinType_Gas
		inbound.message = castTx.CallData
		inbound.isContractCall = true
	case toncontracts.OpCall:
		castTx, err := tx.Call()
		if err != nil {
			return nil, err
		}
		inbound.sender = castTx.Sender
		inbound.receiver = castTx.Recipient
		inbound.amount = math.NewUint(0)
		inbound.coinType = coin.CoinType_NoAssetCall
		inbound.message = castTx.CallData
		inbound.isContractCall = true
	default:
		return nil, fmt.Errorf("unknown operation: %d", tx.Operation)
	}

	return inbound, nil
}

func (inbound *Inbound) isCompliant() bool {
	return !config.ContainRestrictedAddress(
		inbound.receiver.Hex(),
		inbound.sender.ToRaw(),
		inbound.sender.ToHuman(false, false),
		inbound.sender.ToHuman(true, false),
	)
}

func (inbound *Inbound) intoVoteMessage(
	operatorAddress string,
	senderChain int64,
	zetaChain int64,
) *types.MsgVoteInbound {
	const (
		seqno      = 0  // TON does not use sequential block numbers
		eventIndex = 0  // not applicable for TON
		asset      = "" // empty for gas coin
		gasLimit   = zetacore.PostVoteInboundCallOptionsGasLimit
	)

	var (
		inboundHash = encoder.EncodeTx(inbound.tx.Transaction)
		sender      = inbound.sender.ToRaw()
		receiver    = inbound.receiver.Hex()
	)

	return types.NewMsgVoteInbound(
		operatorAddress,
		sender,
		senderChain,
		sender,
		receiver,
		zetaChain,
		inbound.amount,
		hex.EncodeToString(inbound.message),
		inboundHash,
		seqno,
		gasLimit,
		inbound.coinType,
		asset,
		eventIndex,
		types.ProtocolContractVersion_V2,
		false, // not used
		types.InboundStatus_SUCCESS,
		types.ConfirmationMode_SAFE,
		types.WithCrossChainCall(inbound.isContractCall),
	)
}
