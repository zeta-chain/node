package signer

import (
	"context"
	"fmt"
	"math"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/constant"
	mathpkg "github.com/zeta-chain/node/pkg/math"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

const (
	// feeRateCap is the maximum average fee rate for CPFP fee bumping
	// 100 sat/vB is a heuristic based on Bitcoin mempool statistics to avoid excessive fees
	// see: https://mempool.space/graphs/mempool#3y
	feeRateCap = 100

	// decentFeeBumpPercent is a decent percentage for a fee rate bump.
	// The value20% is a heuristic, not mandated by the Bitcoin protocol. It is used to measure the gap between
	// the old fee rate, new fee rate and live fee rate to emit warning messages during RBF fee bumping.
	decentFeeBumpPercent = 20
)

// CPFPFeeBumper is a helper struct to contain CPFP (child-pays-for-parent) fee bumping logic
type CPFPFeeBumper struct {
	ctx context.Context

	chain chains.Chain

	// rpc is the interface to interact with the Bitcoin chain
	rpc RPC

	// tx is the stuck transaction to bump
	tx *btcutil.Tx

	// minRelayFee is the minimum relay fee in BTC
	minRelayFee float64

	// cctxRate is the most recent fee rate of the CCTX
	cctxRate uint64

	// liveRate is the most recent market fee rate
	liveRate uint64

	// txsAndFees contains the information of all pending txs and fees
	txsAndFees client.MempoolTxsAndFees

	logger zerolog.Logger
}

// BumpResult contains the result of the fee bump
type BumpResult struct {
	NewTx          *wire.MsgTx
	AdditionalFees int64
	NewFeeRate     uint64
}

// NewCPFPFeeBumper creates a new CPFPFeeBumper
func NewCPFPFeeBumper(
	ctx context.Context,
	rpc RPC,
	chain chains.Chain,
	tx *btcutil.Tx,
	cctxRate uint64,
	minRelayFee float64,
	logger zerolog.Logger,
) (*CPFPFeeBumper, error) {
	fb := &CPFPFeeBumper{
		ctx:         ctx,
		chain:       chain,
		rpc:         rpc,
		tx:          tx,
		minRelayFee: minRelayFee,
		cctxRate:    cctxRate,
		logger:      logger,
	}

	err := fb.fetchFeeBumpInfo()
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// BumpTxFee bumps the fee of the stuck transaction using reserved bump fees
func (b *CPFPFeeBumper) BumpTxFee() (result BumpResult, err error) {
	// reuse old tx body
	newTx := CopyMsgTxNoWitness(b.tx.MsgTx())
	if len(newTx.TxOut) < 3 {
		return result, errors.New("original tx has no reserved bump fees")
	}

	// the new fee rate is supposed to be much higher than current paid rate (old rate).
	// we print a warning message if it's not the case for monitoring purposes.
	// #nosec G115 always positive
	oldRateBumped, _ := mathpkg.IncreaseUintByPercent(sdkmath.NewUint(b.txsAndFees.AvgFeeRate), decentFeeBumpPercent)
	if sdkmath.NewUint(b.cctxRate).LT(oldRateBumped) {
		b.logger.Warn().
			Uint64("old_fee_rate", b.txsAndFees.AvgFeeRate).
			Uint64("new_fee_rate", b.cctxRate).
			Msg("new fee rate is not much higher than the old fee rate")
	}

	// the live rate may continue increasing during network congestion, and the new fee rate is still not high enough.
	// but we should still continue with the tx replacement because zetacore had already bumped the fee rate.
	newRateBumped, _ := mathpkg.IncreaseUintByPercent(sdkmath.NewUint(b.cctxRate), decentFeeBumpPercent)
	if sdkmath.NewUint(b.liveRate).GT(newRateBumped) {
		b.logger.Warn().
			Uint64("new_fee_rate", b.cctxRate).
			Uint64("live_fee_rate", b.liveRate).
			Msg("live fee rate is still much higher than the new fee rate")
	}

	// cap the fee rate to avoid excessive fees
	feeRateNew := min(b.cctxRate, feeRateCap)

	// calculate minmimum relay fees of the new replacement tx
	// the new tx will have almost same size as the old one because the tx body stays the same
	txVSize := mempool.GetTxVirtualSize(b.tx)
	minRelayFeeRate, err := common.FeeRateToSatPerByte(b.minRelayFee)
	if err != nil {
		return result, errors.Wrapf(err, "unable to convert min relay fee rate")
	}
	// #nosec G115 always in range
	minRelayTxFees := txVSize * int64(minRelayFeeRate)

	// calculate the RBF additional fees required by Bitcoin protocol
	// two conditions to satisfy:
	// 1. new txFees >= old txFees (already handled above)
	// 2. additionalFees >= minRelayTxFees
	//
	// see: https://github.com/bitcoin/bitcoin/blob/5b8046a6e893b7fad5a93631e6d1e70db31878af/src/policy/rbf.cpp#L166-L183
	// #nosec G115 always in range
	additionalFees := b.txsAndFees.TotalVSize*int64(feeRateNew) - b.txsAndFees.TotalFees
	if additionalFees < minRelayTxFees {
		return result, fmt.Errorf(
			"hold on RBF: additional fees %d is lower than min relay fees %d",
			additionalFees,
			minRelayTxFees,
		)
	}

	// bump fees in two ways:
	// 1. deduct additional fees from the change amount
	// 2. give up the whole change amount if it's not enough
	if newTx.TxOut[2].Value >= additionalFees+constant.BTCWithdrawalDustAmount {
		newTx.TxOut[2].Value -= additionalFees
	} else {
		additionalFees = newTx.TxOut[2].Value
		newTx.TxOut = newTx.TxOut[:2]
	}

	// effective fee rate
	// #nosec G115 always positive
	feeRateNew = uint64(math.Ceil(float64(b.txsAndFees.TotalFees+additionalFees) / float64(b.txsAndFees.TotalVSize)))

	return BumpResult{
		NewTx:          newTx,
		AdditionalFees: additionalFees,
		NewFeeRate:     feeRateNew,
	}, nil
}

// fetchFeeBumpInfo fetches all necessary information needed to bump the stuck tx
func (b *CPFPFeeBumper) fetchFeeBumpInfo() error {
	// query live fee rate
	liveRate, err := b.rpc.GetEstimatedFeeRate(b.ctx, 1)
	if err != nil {
		return errors.Wrap(err, "GetEstimatedFeeRate failed")
	}
	b.liveRate = liveRate

	// create a new context with timeout
	ctx, cancel := context.WithTimeout(b.ctx, time.Minute)
	defer cancel()

	// query total fees and sizes of all pending parent TSS txs
	txsAndFees, err := b.rpc.GetMempoolTxsAndFees(ctx, b.tx.MsgTx().TxID())
	if err != nil {
		return errors.Wrap(err, "unable to fetch mempool txs info")
	}
	b.txsAndFees = txsAndFees

	b.logger.Info().
		Int64("total_txs", b.txsAndFees.TotalTxs).
		Int64("total_fees", b.txsAndFees.TotalFees).
		Int64("total_vsize", b.txsAndFees.TotalVSize).
		Uint64("avg_fee_rate", b.txsAndFees.AvgFeeRate).
		Msg("fetched fee bump information")

	return nil
}

// CopyMsgTxNoWitness creates a deep copy of the given MsgTx and clears the witness data
func CopyMsgTxNoWitness(tx *wire.MsgTx) *wire.MsgTx {
	copyTx := tx.Copy()
	for idx := range copyTx.TxIn {
		copyTx.TxIn[idx].Witness = wire.TxWitness{}
	}
	return copyTx
}
