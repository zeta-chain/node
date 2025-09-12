package observer

import (
	"context"
	"encoding/hex"
	"fmt"

	"cosmossdk.io/math"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	"github.com/zeta-chain/node/zetaclient/compliance"
	zetaclientconfig "github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

const paginationLimit = 100

// ObserveInbound observes Gateway's account for new transactions [INBOUND AND OUTBOUND]
//
// Due to TON's architecture we have to scan for all net-new transactions.
// The main purpose is to observe inbounds from TON.
// Note that we might also have *outbounds* here (if a signer broadcasts a tx, it will be observed here).
//
// The name `ObserveInbound` is used for consistency with other chains.
func (ob *Observer) ObserveInbound(ctx context.Context) error {
	lt, hashBits, err := ob.ensureLastScannedTx(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get last scanned tx")
	}

	txs, err := ob.rpc.GetTransactionsSince(ctx, ob.gateway.AccountID(), lt, hashBits)
	if err != nil {
		return errors.Wrap(err, "unable to get transactions")
	}

	switch {
	case len(txs) == 0:
		// noop
		return nil
	case len(txs) > paginationLimit:
		ob.Logger().Inbound.Info().
			Int("tx_count", len(txs)).
			Int("pagination_limit", paginationLimit).
			Msg("observe gateway: number of transactions exceeds pagination limit; taking only some")
		txs = txs[:paginationLimit]
	default:
		ob.Logger().Inbound.Info().
			Int("tx_count", len(txs)).
			Msg("observe gateway: got transactions")
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
			ob.Logger().Inbound.Warn().
				Fields(txLogFields(tx)).
				Msg("observe gateway: observed a failed tx")
		}

		if skip {
			tx = &toncontracts.Transaction{Transaction: txs[i]}
			txHash := rpc.TransactionToHashString(tx.Transaction)
			ob.Logger().Inbound.Warn().
				Str("transaction_hash", txHash).
				Msg("observe gateway: skipping tx")
			ob.setLastScannedTx(tx)
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
				ob.Logger().Inbound.Error().
					Err(err).
					Fields(txLogFields(tx)).
					Msg("observe gateway: unable to add outbound tracker")

				return errors.Wrap(err, "unable to add outbound tracker")
			}

			ob.setLastScannedTx(tx)
			continue
		}

		// Ok, let's process a new inbound tx
		if err := ob.voteInbound(ctx, tx); err != nil {
			ob.Logger().Inbound.Error().
				Err(err).
				Fields(txLogFields(tx)).
				Msg("observe gateway: unable to vote for inbound tx")

			return errors.Wrapf(err, "unable to vote for inbound tx %s", tx.Hash().Hex())
		}

		ob.setLastScannedTx(tx)
	}

	return nil
}

// ProcessInboundTrackers handles adhoc trackers that were somehow missed by
func (ob *Observer) ProcessInboundTrackers(ctx context.Context) error {
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

		lt, hash, err := rpc.TransactionHashFromString(txHash)
		if err != nil {
			ob.logSkippedTracker(txHash, "unable_to_parse_hash", err)
			continue
		}

		raw, err := ob.rpc.GetTransaction(ctx, gatewayAccountID, lt, hash)
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

		if err := ob.voteInbound(ctx, tx); err != nil {
			ob.logSkippedTracker(txHash, "vote_failed", err)
			continue
		}
	}

	return nil
}

// inboundData represents extract data from a TON inbound deposit
type inboundData struct {
	tx *toncontracts.Transaction

	sender   ton.AccountID
	receiver eth.Address

	coinType       coin.CoinType
	amount         math.Uint
	message        []byte
	isContractCall bool
}

// Sends PostVoteInbound to zetacore
func (ob *Observer) voteInbound(ctx context.Context, tx *toncontracts.Transaction) error {
	// noop
	if tx.Operation == toncontracts.OpDonate {
		ob.Logger().Inbound.Info().
			Fields(txLogFields(tx)).
			Msg("thank you rich folk for your donation")
		return nil
	}

	inbound, err := extractInboundData(tx)
	switch {
	case err != nil:
		return errors.Wrap(err, "unable to extract inbound data")
	case ob.inboundComplianceCheck(inbound):
		// do nothing
		return nil
	}

	const (
		seqno         = 0  // ton doesn't use sequential block numbers
		eventIndex    = 0  // not applicable for TON
		asset         = "" // empty for gas coin
		gasLimit      = zetacore.PostVoteInboundCallOptionsGasLimit
		retryGasLimit = zetacore.PostVoteInboundExecutionGasLimit
	)

	var (
		operatorAddress = ob.ZetacoreClient().GetKeys().GetOperatorAddress()
		inboundHash     = rpc.TransactionToHashString(inbound.tx.Transaction)
		sender          = inbound.sender.ToRaw()
		receiver        = inbound.receiver.Hex()
	)

	msg := types.NewMsgVoteInbound(
		operatorAddress.String(),
		sender,
		ob.Chain().ChainId,
		sender,
		receiver,
		ob.ZetacoreClient().Chain().ChainId,
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

	_, err = ob.PostVoteInbound(ctx, msg, retryGasLimit)
	if err != nil {
		return errors.Wrap(err, "unable to vote for inbound tx")
	}

	return nil
}

// extractInboundData parses Gateway tx into deposit (TON sender, amount, memo)
func extractInboundData(tx *toncontracts.Transaction) (inboundData, error) {
	in := inboundData{
		tx:             tx,
		sender:         ton.AccountID{},
		receiver:       eth.Address{},
		amount:         math.Uint{},
		coinType:       coin.CoinType_Gas,
		message:        []byte{},
		isContractCall: false,
	}

	switch tx.Operation {
	case toncontracts.OpDeposit:
		d, err := tx.Deposit()
		if err != nil {
			return inboundData{}, err
		}

		in.sender = d.Sender
		in.receiver = d.Recipient
		in.amount = d.Amount

		return in, nil
	case toncontracts.OpDepositAndCall:
		d, err := tx.DepositAndCall()
		if err != nil {
			return inboundData{}, err
		}

		in.sender = d.Sender
		in.receiver = d.Recipient
		in.amount = d.Amount
		in.message = d.CallData
		in.isContractCall = true

		return in, nil
	case toncontracts.OpCall:
		c, err := tx.Call()
		if err != nil {
			return inboundData{}, err
		}

		in.sender = c.Sender
		in.receiver = c.Recipient
		in.coinType = coin.CoinType_NoAssetCall
		in.amount = math.NewUint(0)
		in.message = c.CallData
		in.isContractCall = true

		return in, nil

	default:
		return inboundData{}, fmt.Errorf("unknown operation %d", tx.Operation)
	}
}

func (ob *Observer) inboundComplianceCheck(inbound inboundData) (restricted bool) {
	var addresses = []string{
		inbound.receiver.Hex(),
		inbound.sender.ToRaw(),
		inbound.sender.ToHuman(false, false),
		inbound.sender.ToHuman(true, false),
	}

	if !zetaclientconfig.ContainRestrictedAddress(addresses...) {
		return false
	}

	txHash := rpc.TransactionHashToString(inbound.tx.Lt, ton.Bits256(inbound.tx.Hash()))

	compliance.PrintComplianceLog(
		ob.Logger().Inbound,
		ob.Logger().Compliance,
		false,
		ob.Chain().ChainId,
		txHash,
		inbound.sender.ToRaw(),
		inbound.receiver.Hex(),
		inbound.coinType.String(),
	)

	return true
}

// ensureLastScannedTx or query the latest tx from RPC
func (ob *Observer) ensureLastScannedTx(ctx context.Context) (uint64, ton.Bits256, error) {
	// always expect init state.
	if txHash := ob.LastTxScanned(); txHash != "" {
		return rpc.TransactionHashFromString(txHash)
	}

	// get last txs from RPC and pick the oldest one
	const limit = 20

	txs, err := ob.rpc.GetTransactions(ctx, limit, ob.gateway.AccountID(), 0, ton.Bits256{})
	switch {
	case err != nil:
		return 0, ton.Bits256{}, errors.Wrap(err, "unable to get last scanned tx")
	case len(txs) == 0:
		return 0, ton.Bits256{}, errors.New("no transactions found")
	}

	tx := txs[len(txs)-1]

	// note this data is not persisted to DB unless real inbound is processed
	ob.WithLastTxScanned(rpc.TransactionToHashString(tx))

	return tx.Lt, ton.Bits256(tx.Hash()), nil
}

func (ob *Observer) setLastScannedTx(tx *toncontracts.Transaction) {
	logger := ob.Logger().Inbound.With().Fields(txLogFields(tx)).Logger()

	txHash := rpc.TransactionToHashString(tx.Transaction)
	ob.WithLastTxScanned(txHash)

	err := ob.WriteLastTxScannedToDB(txHash)
	if err != nil {
		logger.Error().Err(err).Msg("error calling WriteLastTxScannedToDB")
		return
	}

	logger.Info().Msg("call to WriteLastTxScannedToDB was successful")
}

func (ob *Observer) logSkippedTracker(hash string, reason string, err error) {
	ob.Logger().Inbound.Warn().
		Err(err).
		Str("transaction_hash", hash).
		Str("skip_reason", reason).
		Msg("skipping tracker")
}

func txLogFields(tx *toncontracts.Transaction) map[string]any {
	return map[string]any{
		"transaction_hash":           rpc.TransactionToHashString(tx.Transaction),
		"transaction_ton_is_inbound": tx.IsInbound(),
		"transaction_ton_op_code":    tx.Operation,
		"transaction_ton_exit_code":  tx.ExitCode,
	}
}

//nolint:unused // used for in tests
func castBlockID(id ton.BlockIDExt) rpc.BlockIDExt {
	return rpc.BlockIDExt{
		Workchain: int(id.Workchain),
		Seqno:     id.Seqno,
		Shard:     fmt.Sprintf("%d", id.Shard),
		RootHash:  id.RootHash.Base64(),
		FileHash:  id.FileHash.Base64(),
	}
}
