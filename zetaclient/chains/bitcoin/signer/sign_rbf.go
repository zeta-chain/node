package signer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/node/zetaclient/logs"
)

const (
	// rbfTxInSequenceNum is the sequence number used to signal an opt-in full-RBF (Replace-By-Fee) transaction
	// Setting sequenceNum to "1" effectively makes the transaction timelocks irrelevant.
	// See: https://github.com/bitcoin/bips/blob/master/bip-0125.mediawiki
	// See: https://github.com/BlockchainCommons/Learning-Bitcoin-from-the-Command-Line/blob/master/05_2_Resending_a_Transaction_with_RBF.md
	rbfTxInSequenceNum uint32 = 1
)

// SignRBFTx signs a RBF (Replace-By-Fee) to unblock last stuck outbound transaction.
//
// The key points:
//   - It reuses the stuck tx's inputs and outputs but gives a higher fee to miners.
//   - Funding the last stuck outbound will be considered as CPFP (child-pays-for-parent) by miners.
func (signer *Signer) SignRBFTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	height uint64,
	lastTx *btcutil.Tx,
	minRelayFee float64,
	memplTxsInfoFetcher MempoolTxsInfoFetcher,
) (*wire.MsgTx, error) {
	var (
		params = cctx.GetCurrentOutboundParam()
		lf     = map[string]any{
			logs.FieldMethod: "SignRBFTx",
			logs.FieldNonce:  params.TssNonce,
			logs.FieldTx:     lastTx.MsgTx().TxID(),
		}
		logger = signer.Logger().Std.With().Fields(lf).Logger()
	)

	var cctxRate int64
	switch signer.Chain().ChainId {
	case chains.BitcoinRegtest.ChainId:
		// hardcode for regnet E2E test, zetacore won't feed it to CCTX
		cctxRate = rpc.FeeRateRegnetRBF
	default:
		// parse recent fee rate from CCTX
		cctxRate, err := strconv.ParseInt(params.GasPriorityFee, 10, 64)
		if err != nil || cctxRate <= 0 {
			return nil, fmt.Errorf("invalid fee rate %s", params.GasPrice)
		}
	}

	// create fee bumper
	fb, err := NewCPFPFeeBumper(
		signer.Chain(),
		signer.client,
		memplTxsInfoFetcher,
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
		preTx, err := signer.client.GetRawTransaction(&preOut.Hash)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get previous tx %s", preOut.Hash)
		}
		inAmounts[i] = preTx.MsgTx().TxOut[preOut.Index].Value
	}

	// sign the RBF tx
	err = signer.SignTx(ctx, newTx, inAmounts, height, params.TssNonce)
	if err != nil {
		return nil, errors.Wrap(err, "SignTx failed")
	}

	return newTx, nil
}
