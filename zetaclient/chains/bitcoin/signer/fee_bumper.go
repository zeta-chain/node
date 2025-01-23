package signer

import (
	"context"
	"fmt"
	"math"
	"time"

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

	// minCPFPFeeBumpPercent is the minimum percentage by which the CPFP average fee rate should be bumped.
	// This value 20% is a heuristic, not mandated by the Bitcoin protocol, designed to balance effectiveness
	// in replacing stuck transactions while avoiding excessive sensitivity to fee market fluctuations.
	minCPFPFeeBumpPercent = 20
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
	CCTXRate int64

	// LiveRate is the most recent market fee rate
	LiveRate int64

	// TotalTxs is the total number of stuck TSS txs
	TotalTxs int64

	// TotalFees is the total fees of all stuck TSS txs
	TotalFees int64

	// TotalVSize is the total vsize of all stuck TSS txs
	TotalVSize int64

	// AvgFeeRate is the average fee rate of all stuck TSS txs
	AvgFeeRate int64
}

// NewCPFPFeeBumper creates a new CPFPFeeBumper
func NewCPFPFeeBumper(
	ctx context.Context,
	rpc RPC,
	chain chains.Chain,
	tx *btcutil.Tx,
	cctxRate int64,
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
	}

	err := fb.FetchFeeBumpInfo(logger)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// BumpTxFee bumps the fee of the stuck transaction using reserved bump fees
func (b *CPFPFeeBumper) BumpTxFee() (*wire.MsgTx, int64, int64, error) {
	// reuse old tx body
	newTx := CopyMsgTxNoWitness(b.Tx.MsgTx())
	if len(newTx.TxOut) < 3 {
		return nil, 0, 0, errors.New("original tx has no reserved bump fees")
	}

	// tx replacement is triggered only when market fee rate goes 20% higher than current paid rate.
	// zetacore updates the cctx fee rate evey 10 minutes, we could hold on and retry later.
	minBumpRate := mathpkg.IncreaseIntByPercent(b.AvgFeeRate, minCPFPFeeBumpPercent)
	if b.CCTXRate < minBumpRate {
		return nil, 0, 0, fmt.Errorf(
			"hold on RBF: cctx rate %d is lower than the min bumped rate %d",
			b.CCTXRate,
			minBumpRate,
		)
	}

	// the live rate may continue increasing during network congestion, we should wait until it stabilizes a bit.
	// this is to ensure the live rate is not 20%+ higher than the cctx rate, otherwise, the replacement tx may
	// also get stuck and need another replacement.
	bumpedRate := mathpkg.IncreaseIntByPercent(b.CCTXRate, minCPFPFeeBumpPercent)
	if b.LiveRate > bumpedRate {
		return nil, 0, 0, fmt.Errorf(
			"hold on RBF: live rate %d is much higher than the cctx rate %d",
			b.LiveRate,
			b.CCTXRate,
		)
	}

	// cap the fee rate to avoid excessive fees
	feeRateNew := b.CCTXRate
	if b.CCTXRate > feeRateCap {
		feeRateNew = feeRateCap
	}

	// calculate minmimum relay fees of the new replacement tx
	// the new tx will have almost same size as the old one because the tx body stays the same
	txVSize := mempool.GetTxVirtualSize(b.Tx)
	minRelayFeeRate := common.FeeRateToSatPerByte(b.MinRelayFee)
	minRelayTxFees := txVSize * minRelayFeeRate

	// calculate the RBF additional fees required by Bitcoin protocol
	// two conditions to satisfy:
	// 1. new txFees >= old txFees (already handled above)
	// 2. additionalFees >= minRelayTxFees
	//
	// see: https://github.com/bitcoin/bitcoin/blob/master/src/policy/rbf.cpp#L166-L183
	additionalFees := b.TotalVSize*feeRateNew - b.TotalFees
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
	feeRateNew = int64(math.Ceil(float64(b.TotalFees+additionalFees) / float64(b.TotalVSize)))

	return newTx, additionalFees, feeRateNew, nil
}

// fetchFeeBumpInfo fetches all necessary information needed to bump the stuck tx
func (b *CPFPFeeBumper) FetchFeeBumpInfo(logger zerolog.Logger) error {
	// query live fee rate
	isRegnet := chains.IsBitcoinRegnet(b.Chain.ChainId)
	liveRate, err := b.RPC.GetEstimatedFeeRate(b.Ctx, 1, isRegnet)
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
	logger.Info().
		Msgf("totalTxs %d, totalFees %f, totalVSize %d, avgFeeRate %d", totalTxs, totalFees, totalVSize, avgFeeRate)

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
