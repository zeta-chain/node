package signer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// SignRBFTx signs a RBF (Replace-By-Fee) to unblock last stuck outbound transaction.
//
// The key points:
//   - It reuses the stuck tx's inputs and outputs but gives a higher fee to miners.
//   - Funding the last stuck outbound will be considered as CPFP (child-pays-for-parent) by miners.
func (signer *Signer) SignRBFTx(
	ctx context.Context,
	height uint64,
	nonce uint64,
	lastTx *btcutil.Tx,
	latestRateStr string,
	minRelayFee float64,
) (*wire.MsgTx, error) {
	var (
		lf = map[string]any{
			logs.FieldMethod: "SignRBFTx",
			logs.FieldNonce:  nonce,
			logs.FieldTx:     lastTx.MsgTx().TxID(),
		}
		logger = signer.Logger().Std.With().Fields(lf).Logger()
	)

	var cctxRate int64
	switch signer.Chain().ChainId {
	case chains.BitcoinRegtest.ChainId:
		// hardcode for regnet E2E test, zetacore won't feed it to CCTX
		cctxRate = client.FeeRateRegnetRBF
	default:
		// parse latest fee rate from CCTX
		latestRate, err := strconv.ParseInt(latestRateStr, 10, 64)
		if err != nil || latestRate <= 0 {
			return nil, fmt.Errorf("invalid fee rate %s", latestRateStr)
		}
		cctxRate = common.OutboundFeeRateFromCCTXRate(latestRate)
	}

	// create fee bumper
	fb, err := NewCPFPFeeBumper(
		ctx,
		signer.rpc,
		signer.Chain(),
		lastTx,
		cctxRate,
		minRelayFee,
		logger,
	)
	if err != nil {
		return nil, errors.Wrap(err, "NewCPFPFeeBumper failed")
	}

	// bump tx fees
	newTx, additionalFees, newRate, err := fb.BumpTxFee()
	if err != nil {
		return nil, errors.Wrap(err, "BumpTxFee failed")
	}
	logger.Info().
		Msgf("BumpTxFee succeed, additional fees: %d sats, rate: %d => %d sat/vB", additionalFees, fb.AvgFeeRate, newRate)

	// collect input amounts for signing
	inAmounts := make([]int64, len(newTx.TxIn))
	for i, input := range newTx.TxIn {
		preOut := input.PreviousOutPoint
		preTx, err := signer.rpc.GetRawTransaction(ctx, &preOut.Hash)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get previous tx %s", preOut.Hash)
		}
		inAmounts[i] = preTx.MsgTx().TxOut[preOut.Index].Value
	}

	// sign the RBF tx
	err = signer.SignTx(ctx, newTx, inAmounts, height, nonce)
	if err != nil {
		return nil, errors.Wrap(err, "SignTx failed")
	}

	return newTx, nil
}
