package observer

import (
	"context"

	"cosmossdk.io/math"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/ton/encoder"
	"github.com/zeta-chain/node/zetaclient/logs"
	"github.com/zeta-chain/node/zetaclient/zetacore"
)

type Outbound struct {
	tx            *toncontracts.Transaction
	receiveStatus chains.ReceiveStatus
	nonce         uint64
}

// ------------------------------------------------------------------------------------------------
// VoteOutboundIfConfirmed
// ------------------------------------------------------------------------------------------------

// VoteOutboundIfConfirmed checks the outbound status and returns (continueKeysign, error).
func (ob *Observer) VoteOutboundIfConfirmed(ctx context.Context,
	cctx *cctypes.CrossChainTx,
) (bool, error) {
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	outbound := ob.getOutbound(nonce)
	if outbound == nil {
		return true, nil
	}

	err := ob.voteOutbound(ctx, cctx, outbound)
	if err != nil {
		return false, errors.Wrap(err, "unable to vote for outbound transaction")
	}

	return false, nil
}

func (ob *Observer) voteOutbound(ctx context.Context,
	cctx *cctypes.CrossChainTx,
	outbound *Outbound,
) error {
	receiveStatus, amount, err := receiveStatusWithAmount(cctx, outbound)
	if err != nil {
		return errors.Wrap(err, "unable to get status and amount")
	}

	gasPrice, err := ob.getLatestGasPrice()
	if err != nil {
		return err
	}

	// There is no sequential block height in TON.
	// Different txs might end up in different shards.
	// The ID in tlb.BlockID is essentially a workchain+shard+seqno tuple.
	// We cannot use it as a block height.
	const tonBlockHeight = 0

	// We do not specify an effective gas limit for TON outbounds.
	// This value is used for the gas stability pool funding, which is not applicable to TON.
	const effectiveGasLimit = 0

	txHash := encoder.EncodeTx(outbound.tx.Transaction)
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	msg := cctypes.NewMsgVoteOutbound(
		ob.ZetaRepo().GetOperatorAddress(),
		cctx.Index,
		txHash,
		tonBlockHeight,
		outbound.tx.GasUsed().Uint64(),
		math.NewIntFromUint64(gasPrice),
		effectiveGasLimit,
		amount,
		receiveStatus,
		ob.Chain().ChainId,
		nonce,
		cctx.InboundParams.CoinType,
		cctypes.ConfirmationMode_SAFE,
	)

	logger := ob.Logger().Outbound.With().
		Uint64(logs.FieldNonce, nonce).
		Str(logs.FieldTx, txHash).
		Logger()

	const gasLimit = zetacore.PostVoteOutboundGasLimit

	retryGasLimit := zetacore.PostVoteOutboundRetryGasLimit
	if msg.Status == chains.ReceiveStatus_failed {
		retryGasLimit = zetacore.PostVoteOutboundRevertGasLimit
	}

	_, _, err = ob.ZetaRepo().VoteOutbound(ctx, logger, gasLimit, retryGasLimit, msg)
	return err
}

func receiveStatusWithAmount(
	cctx *cctypes.CrossChainTx,
	outbound *Outbound,
) (chains.ReceiveStatus, math.Uint, error) {
	switch outbound.tx.Operation {
	case toncontracts.OpWithdraw:
		wd, err := outbound.tx.Withdrawal()
		if err != nil {
			return 0, math.Uint{}, errors.Wrap(err, "unable to get withdrawal")
		}

		return outbound.receiveStatus, wd.Amount, nil
	case toncontracts.OpIncreaseSeqno:
		// force failure to revert the CCTX in zetacore (similarly to SUI)
		return chains.ReceiveStatus_failed, cctx.GetCurrentOutboundParam().Amount, nil
	}

	return 0, math.Uint{}, errors.Errorf("unknown operation %d", outbound.tx.Operation)
}

// ------------------------------------------------------------------------------------------------
// ProcessOutboundTrackers
// ------------------------------------------------------------------------------------------------

// ProcessOutboundTrackers gets outbound trackers from zetacore, fetches their transaction data
// from TON, and stores it in memory for further use.
func (ob *Observer) ProcessOutboundTrackers(ctx context.Context) error {
	logger := ob.Logger().Outbound

	trackers, err := ob.ZetaRepo().GetOutboundTrackers(ctx)
	if err != nil {
		return err
	}

	for _, tracker := range trackers {
		nonce := tracker.Nonce

		// Skip outbound trackers that are already in the in-memory cache.
		if ob.getOutbound(nonce) != nil {
			continue
		}

		cctx, err := ob.ZetaRepo().GetCCTX(ctx, nonce)
		if err != nil {
			logger.Error().Err(err).Uint64(logs.FieldNonce, nonce).Send()
			continue // does not block other CCTXs from being processed
		}

		for _, txHash := range tracker.HashList {
			err := ob.processOutboundTracker(ctx, cctx, txHash.TxHash)
			if err != nil {
				logger.Error().
					Err(err).
					Str(logs.FieldTx, txHash.TxHash).
					Uint64(logs.FieldNonce, nonce).
					Msg("unable to process outbound tracker")
			}
		}
	}

	return nil
}

// processOutboundTracker validates the CCTX and stores it in the in-memory cache to be used by
// VoteOutboundIfConfirmed.
//
// NOTE: restricted transactions (increase-seqno transactions) are considered successful because
// they are committed on-chain.
func (ob *Observer) processOutboundTracker(ctx context.Context,
	cctx *cctypes.CrossChainTx,
	encodedHash string,
) error {
	nonce := cctx.GetCurrentOutboundParam().TssNonce

	if cctx.InboundParams.CoinType != coin.CoinType_Gas {
		return errors.New("only Gas CCTXs are supported")
	}

	rawTx, err := ob.tonRepo.GetTransactionByHash(ctx, encodedHash)
	if err != nil {
		return errors.Wrap(err, "unable to get transaction")
	}

	// TODO: why the different behavior from parseTransaction?
	tx, err := ob.gateway.ParseTransaction(*rawTx)
	if err != nil {
		return errors.Wrap(err, "unable to parse transaction")
	}

	receiveStatus := chains.ReceiveStatus_success

	out, err := tx.OutboundAuth()
	if err != nil {
		return errors.Wrap(err, "unable to get outbound auth")
	}

	tssSigner := ob.TSS().PubKey().AddressEVM()
	if out.Signer != tssSigner {
		return errors.Errorf("signer mismatch (got %s, want %s)", out.Signer, tssSigner)
	}

	if nonce != uint64(out.Seqno) {
		return errors.Errorf("nonce mismatch (got %d, want %d)", out.Seqno, nonce)
	}

	if !tx.IsSuccess() {
		receiveStatus = chains.ReceiveStatus_failed
	}

	// Will be used by VoteOutboundIfConfirmed.
	ob.addOutbound(Outbound{tx, receiveStatus, nonce})

	return nil
}
