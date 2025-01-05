package signer

import (
	"fmt"
	"math"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

const (
	// minCPFPFeeBumpFactor is the minimum factor by which the CPFP average fee rate should be bumped.
	// This value 20% is a heuristic, not mandated by the Bitcoin protocol, designed to balance effectiveness
	// in replacing stuck transactions while avoiding excessive sensitivity to fee market fluctuations.
	minCPFPFeeBumpFactor = 1.2
)

// MempoolTxsInfoFetcher is a function type to fetch mempool txs information
type MempoolTxsInfoFetcher func(interfaces.BTCRPCClient, string) (int64, float64, int64, int64, error)

// CPFPFeeBumper is a helper struct to contain CPFP (child-pays-for-parent) fee bumping logic
type CPFPFeeBumper struct {
	// Client is the RPC Client to interact with the Bitcoin chain
	Client interfaces.BTCRPCClient

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
	client interfaces.BTCRPCClient,
	memplTxsInfoFetcher MempoolTxsInfoFetcher,
	tx *btcutil.Tx,
	cctxRate int64,
	minRelayFee float64,
	logger zerolog.Logger,
) (*CPFPFeeBumper, error) {
	fb := &CPFPFeeBumper{
		Client:      client,
		Tx:          tx,
		MinRelayFee: minRelayFee,
		CCTXRate:    cctxRate,
	}

	err := fb.FetchFeeBumpInfo(memplTxsInfoFetcher, logger)
	if err != nil {
		return nil, err
	}
	return fb, nil
}

// BumpTxFee bumps the fee of the stuck transaction using reserved bump fees
func (b *CPFPFeeBumper) BumpTxFee() (*wire.MsgTx, int64, error) {
	// reuse old tx body and clear witness data (e.g., signatures)
	newTx := b.Tx.MsgTx().Copy()
	for idx := range newTx.TxIn {
		newTx.TxIn[idx].Witness = wire.TxWitness{}
	}

	// check reserved bump fees amount in the original tx
	if len(newTx.TxOut) < 3 {
		return nil, 0, errors.New("original tx has no reserved bump fees")
	}

	// tx replacement is triggered only when market fee rate goes 20% higher than current paid fee rate.
	// zetacore updates the cctx fee rate evey 10 minutes, we could hold on and retry later.
	minBumpRate := int64(math.Ceil(float64(b.AvgFeeRate) * minCPFPFeeBumpFactor))
	if b.CCTXRate < minBumpRate {
		return nil, 0, fmt.Errorf(
			"hold on RBF: cctx rate %d is lower than the min bumped rate %d",
			b.CCTXRate,
			minBumpRate,
		)
	}

	// the live rate may continue increasing during network congestion, we should wait until it stabilizes a bit.
	// this is to ensure the live rate is not 20%+ higher than the cctx rate, otherwise, the replacement tx may
	// also get stuck and need another replacement.
	bumpedRate := int64(math.Ceil(float64(b.CCTXRate) * minCPFPFeeBumpFactor))
	if b.LiveRate > bumpedRate {
		return nil, 0, fmt.Errorf(
			"hold on RBF: live rate %d is much higher than the cctx rate %d",
			b.LiveRate,
			b.CCTXRate,
		)
	}

	// calculate minmimum relay fees of the new replacement tx
	// the new tx will have almost same size as the old one because the tx body stays the same
	txVSize := mempool.GetTxVirtualSize(b.Tx)
	minRelayFeeRate := rpc.FeeRateToSatPerByte(b.MinRelayFee)
	minRelayTxFees := txVSize * minRelayFeeRate.Int64()

	// calculate the RBF additional fees required by Bitcoin protocol
	// two conditions to satisfy:
	// 1. new txFees >= old txFees (already handled above)
	// 2. additionalFees >= minRelayTxFees
	//
	// see: https://github.com/bitcoin/bitcoin/blob/master/src/policy/rbf.cpp#L166-L183
	additionalFees := b.TotalVSize*b.CCTXRate - b.TotalFees
	if additionalFees < minRelayTxFees {
		additionalFees = minRelayTxFees
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

	return newTx, additionalFees, nil
}

// fetchFeeBumpInfo fetches all necessary information needed to bump the stuck tx
func (b *CPFPFeeBumper) FetchFeeBumpInfo(memplTxsInfoFetcher MempoolTxsInfoFetcher, logger zerolog.Logger) error {
	// query live network fee rate
	liveRate, err := rpc.GetEstimatedFeeRate(b.Client, 1)
	if err != nil {
		return errors.Wrap(err, "GetEstimatedFeeRate failed")
	}
	b.LiveRate = liveRate

	// query total fees and sizes of all pending parent TSS txs
	totalTxs, totalFees, totalVSize, avgFeeRate, err := memplTxsInfoFetcher(b.Client, b.Tx.MsgTx().TxID())
	if err != nil {
		return errors.Wrap(err, "unable to fetch mempool txs info")
	}
	totalFeesSats, err := bitcoin.GetSatoshis(totalFees)
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
