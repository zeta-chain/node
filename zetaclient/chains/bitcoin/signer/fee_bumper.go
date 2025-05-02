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
	Ctx context.Context

	Chain chains.Chain

	// RPC is the interface to interact with the Bitcoin chain
	RPC RPC

	// Tx is the stuck transaction to bump
	Tx *btcutil.Tx

	// MinRelayFee is the minimum relay fee in BTC
	MinRelayFee float64

	// CCTXRate is the most recent fee rate of the CCTX
	CCTXRate uint64

	// LiveRate is the most recent market fee rate
	LiveRate uint64

	// TotalTxs is the total number of stuck TSS txs
	TotalTxs int64

	// TotalFees is the total fees of all stuck TSS txs
	TotalFees int64

	// TotalVSize is the total vsize of all stuck TSS txs
	TotalVSize int64

	// AvgFeeRate is the average fee rate of all stuck TSS txs
	AvgFeeRate uint64

	Logger zerolog.Logger
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
		Ctx:         ctx,
		Chain:       chain,
		RPC:         rpc,
		Tx:          tx,
		MinRelayFee: minRelayFee,
		CCTXRate:    cctxRate,
		Logger:      logger,
	}

	err := fb.fetchFeeBumpInfo()
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// BumpTxFee bumps the fee of the stuck transaction using reserved bump fees
func (b *CPFPFeeBumper) BumpTxFee() (*wire.MsgTx, int64, uint64, error) {
	// reuse old tx body
	newTx := CopyMsgTxNoWitness(b.Tx.MsgTx())
	if len(newTx.TxOut) < 3 {
		return nil, 0, 0, errors.New("original tx has no reserved bump fees")
	}

	// the new fee rate is supposed to be much higher than current paid rate (old rate).
	// we print a warning message if it's not the case for monitoring purposes.
	// #nosec G115 always positive
	oldRateBumped, _ := mathpkg.IncreaseUintByPercent(sdkmath.NewUint(b.AvgFeeRate), decentFeeBumpPercent)
	if sdkmath.NewUint(b.CCTXRate).LT(oldRateBumped) {
		b.Logger.Warn().
			Uint64("old_fee_rate", b.AvgFeeRate).
			Uint64("new_fee_rate", b.CCTXRate).
			Msg("new fee rate is not much higher than the old fee rate")
	}

	// the live rate may continue increasing during network congestion, and the new fee rate is still not high enough.
	// but we should still continue with the tx replacement because zetacore had already bumped the fee rate.
	newRateBumped, _ := mathpkg.IncreaseUintByPercent(sdkmath.NewUint(b.CCTXRate), decentFeeBumpPercent)
	if sdkmath.NewUint(b.LiveRate).GT(newRateBumped) {
		b.Logger.Warn().
			Uint64("new_fee_rate", b.CCTXRate).
			Uint64("live_fee_rate", b.LiveRate).
			Msg("live fee rate is still much higher than the new fee rate")
	}

	// cap the fee rate to avoid excessive fees
	feeRateNew := min(b.CCTXRate, feeRateCap)

	// calculate minmimum relay fees of the new replacement tx
	// the new tx will have almost same size as the old one because the tx body stays the same
	txVSize := mempool.GetTxVirtualSize(b.Tx)
	minRelayFeeRate, err := common.FeeRateToSatPerByte(b.MinRelayFee)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, "unable to convert min relay fee rate")
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
	additionalFees := b.TotalVSize*int64(feeRateNew) - b.TotalFees
	if additionalFees < minRelayTxFees {
		return nil, 0, 0, fmt.Errorf(
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
	feeRateNew = uint64(math.Ceil(float64(b.TotalFees+additionalFees) / float64(b.TotalVSize)))

	return newTx, additionalFees, feeRateNew, nil
}

// fetchFeeBumpInfo fetches all necessary information needed to bump the stuck tx
func (b *CPFPFeeBumper) fetchFeeBumpInfo() error {
	// query live fee rate
	liveRate, err := b.RPC.GetEstimatedFeeRate(b.Ctx, 1)
	if err != nil {
		return errors.Wrap(err, "GetEstimatedFeeRate failed")
	}
	b.LiveRate = liveRate

	// query total fees and sizes of all pending parent TSS txs
	totalTxs, totalFees, totalVSize, avgFeeRate, err := b.RPC.GetTotalMempoolParentsSizeNFees(
		b.Ctx,
		b.Tx.MsgTx().TxID(),
		time.Minute,
	)
	if err != nil {
		return errors.Wrap(err, "unable to fetch mempool txs info")
	}
	totalFeesSats, err := common.GetSatoshis(totalFees)
	if err != nil {
		return errors.Wrapf(err, "cannot convert total fees %f", totalFees)
	}

	b.TotalTxs = totalTxs
	b.TotalFees = totalFeesSats
	b.TotalVSize = totalVSize
	b.AvgFeeRate = avgFeeRate

	b.Logger.Info().
		Int64("total_txs", totalTxs).
		Int64("total_fees", totalFeesSats).
		Int64("total_vsize", totalVSize).
		Uint64("avg_fee_rate", avgFeeRate).
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
