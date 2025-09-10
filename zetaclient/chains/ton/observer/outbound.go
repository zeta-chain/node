package observer

import (
	"context"

	"cosmossdk.io/math"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

// We don't specify an effective gas limit for TON outbound, this value is used for the gas stability pool funding, which is not used for TON
const effectiveGasLimit = 0

// There's no sequential block height. Also, different txs might end up in different shards.
// tlb.BlockID is essentially a workchain+shard+seqno tuple. We can't use it as a block height, thus zero.
const tonBlockHeight = 0

type outbound struct {
	tx            *toncontracts.Transaction
	receiveStatus chains.ReceiveStatus
	nonce         uint64
}

// VoteOutboundIfConfirmed checks outbound status and returns (continueKeysign, error)
func (ob *Observer) VoteOutboundIfConfirmed(ctx context.Context, cctx *cctypes.CrossChainTx) (bool, error) {
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	res, exists := ob.getOutboundByNonce(nonce)
	if !exists {
		return true, nil
	}

	if err := ob.postVoteOutbound(ctx, cctx, res); err != nil {
		return false, errors.Wrap(err, "unable to post vote")
	}

	return false, nil
}

// ProcessOutboundTrackers pulls outbounds trackers from zetacore,
// fetches txs from TON and stores them in memory for further use.
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	var (
		chainID  = ob.Chain().ChainId
		zetacore = ob.ZetacoreClient()
	)

	trackers, err := zetacore.GetAllOutboundTrackerByChain(ctx, chainID, interfaces.Ascending)
	if err != nil {
		return errors.Wrap(err, "unable to get outbound trackers")
	}

	for _, tracker := range trackers {
		nonce := tracker.Nonce

		// If outbound is already in memory, skip.
		if _, ok := ob.getOutboundByNonce(nonce); ok {
			continue
		}

		// Let's not block other cctxs from being processed
		cctx, err := zetacore.GetCctxByNonce(ctx, chainID, nonce)
		if err != nil {
			ob.Logger().Outbound.
				Error().Err(err).
				Uint64(logs.FieldNonce, nonce).
				Msg("unable to get CCTX by nonce")
			continue
		}

		for _, txHash := range tracker.HashList {
			if err := ob.processOutboundTracker(ctx, cctx, txHash.TxHash); err != nil {
				ob.Logger().Outbound.
					Error().Err(err).
					Str(logs.FieldTx, txHash.TxHash).
					Uint64(logs.FieldNonce, nonce).
					Msg("unable to check CCTX by nonce")
			}
		}
	}

	return nil
}

// processOutboundTracker checks TON tx and stores it in memory for further processing
// by VoteOutboundIfConfirmed. Note that restricted txs (increase_seqno txs) are considered as success
// because they are committed on-chain.
func (ob *Observer) processOutboundTracker(ctx context.Context, cctx *cctypes.CrossChainTx, txHash string) error {
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	if cctx.InboundParams.CoinType != coin.CoinType_Gas {
		return errors.New("only gas cctxs are supported")
	}

	lt, hash, err := rpc.TransactionHashFromString(txHash)
	if err != nil {
		return errors.Wrap(err, "unable to parse tx hash")
	}

	rawTx, err := ob.rpc.GetTransaction(ctx, ob.gateway.AccountID(), lt, hash)
	if err != nil {
		return errors.Wrap(err, "unable to get tx")
	}

	tx, err := ob.gateway.ParseTransaction(rawTx)
	if err != nil {
		return errors.Wrap(err, "unable to parse tx")
	}

	receiveStatus := chains.ReceiveStatus_success

	out, err := tx.OutboundAuth()
	switch {
	case err != nil:
		return errors.Wrap(err, "unable to get outbound auth")
	case out.Signer != ob.TSS().PubKey().AddressEVM():
		return errors.Errorf("signer mismatch (got %s, want %s)", out.Signer, ob.TSS().PubKey().AddressEVM())
	case uint64(out.Seqno) != nonce:
		return errors.Errorf("nonce mismatch (got %d, want %d)", out.Seqno, nonce)
	case !tx.IsSuccess():
		receiveStatus = chains.ReceiveStatus_failed
	}

	// will be used by VoteOutboundIfConfirmed
	ob.setOutboundByNonce(outbound{tx, receiveStatus, nonce})

	return nil
}

// addOutboundTracker publishes outbound tracker to Zetacore.
// In most cases will be a noop because the tracker is already published by the signer.
// See Signer{}.trackOutbound(...) for more details.
func (ob *Observer) addOutboundTracker(ctx context.Context, tx *toncontracts.Transaction) error {
	auth, err := tx.OutboundAuth()
	switch {
	case err != nil:
		return err
	case auth.Signer != ob.TSS().PubKey().AddressEVM():
		ob.Logger().Inbound.Warn().
			Fields(txLogFields(tx)).
			Str("transaction_ton_signer", auth.Signer.String()).
			Msg("observe gateway: signer is not TSS; skipping")

		return nil
	}

	var (
		chainID = ob.Chain().ChainId
		nonce   = uint64(auth.Seqno)
		hash    = rpc.TransactionToHashString(tx.Transaction)
	)

	// note it has a check for noop
	_, err = ob.ZetacoreClient().PostOutboundTracker(ctx, chainID, nonce, hash)

	return err
}

// getOutboundByNonce returns outbound by nonce
func (ob *Observer) getOutboundByNonce(nonce uint64) (outbound, bool) {
	v, ok := ob.outbounds.Get(nonce)
	if !ok {
		return outbound{}, false
	}

	return v.(outbound), true
}

func (ob *Observer) setOutboundByNonce(entry outbound) {
	ob.outbounds.Add(entry.nonce, entry)
}

func (ob *Observer) postVoteOutbound(ctx context.Context, cctx *cctypes.CrossChainTx, res outbound) error {
	var (
		chainID       = ob.Chain().ChainId
		txHash        = rpc.TransactionToHashString(res.tx.Transaction)
		nonce         = cctx.GetCurrentOutboundParam().TssNonce
		signerAddress = ob.ZetacoreClient().GetKeys().GetOperatorAddress()
		coinType      = cctx.InboundParams.CoinType
	)

	receiveStatus, amount, err := receiveStatusWithAmount(res, cctx)
	if err != nil {
		return errors.Wrap(err, "unable to get status and amount")
	}

	gasPrice, err := ob.getLatestGasPrice()
	if err != nil {
		return errors.Wrap(err, "unable to get gas price")
	}

	// #nosec G115 len always in range
	gasPriceInt := math.NewInt(int64(gasPrice))

	msg := cctypes.NewMsgVoteOutbound(
		signerAddress.String(),
		cctx.Index,
		txHash,
		tonBlockHeight,
		res.tx.GasUsed().Uint64(),
		gasPriceInt,
		effectiveGasLimit,
		amount,
		receiveStatus,
		chainID,
		nonce,
		coinType,
		cctypes.ConfirmationMode_SAFE,
	)

	const gasLimit = zetacore.PostVoteOutboundGasLimit

	retryGasLimit := zetacore.PostVoteOutboundRetryGasLimit
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = zetacore.PostVoteOutboundRevertGasLimit
	}

	log := ob.Logger().Outbound.With().
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).
		Logger()

	zetaTxHash, ballot, err := ob.ZetacoreClient().PostVoteOutbound(ctx, gasLimit, retryGasLimit, msg)

	switch {
	case err != nil:
		log.Error().Err(err).Msg("unable to post outbound vote")
		return err
	case zetaTxHash != "":
		log.Info().
			Str(logs.FieldZetaTx, zetaTxHash).
			Str(logs.FieldBallotIndex, ballot).
			Msg("posted outbound vote")
	}

	return nil
}

func receiveStatusWithAmount(o outbound, cctx *cctypes.CrossChainTx) (chains.ReceiveStatus, math.Uint, error) {
	switch o.tx.Operation {
	case toncontracts.OpWithdraw:
		wd, err := o.tx.Withdrawal()
		if err != nil {
			return 0, math.Uint{}, errors.Wrap(err, "unable to get withdrawal")
		}

		return o.receiveStatus, wd.Amount, nil
	case toncontracts.OpIncreaseSeqno:
		// force failure to revert the CCTX in zetacore (similarly to SUI)
		return chains.ReceiveStatus_failed, cctx.GetCurrentOutboundParam().Amount, nil
	}

	return 0, math.Uint{}, errors.Errorf("unknown operation %d", o.tx.Operation)
}
