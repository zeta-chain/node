package signer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/node/zetaclient/logs"
)

const (
	// rbfTxInSequenceNum is the sequence number used to signal an opt-in full-RBF (Replace-By-Fee) transaction
	// Setting sequenceNum to "1" effectively makes the transaction timelocks irrelevant.
	// See: https://github.com/bitcoin/bips/blob/master/bip-0125.mediawiki
	// Also see: https://github.com/BlockchainCommons/Learning-Bitcoin-from-the-Command-Line/blob/master/05_2_Resending_a_Transaction_with_RBF.md
	rbfTxInSequenceNum uint32 = 1
)

func (signer *Signer) SignRBFTx(
	ctx context.Context,
	cctx *types.CrossChainTx,
	oldTx *btcutil.Tx,
	minRelayFee float64,
) (*wire.MsgTx, error) {
	var (
		params = cctx.GetCurrentOutboundParam()
		lf     = map[string]any{
			logs.FieldMethod: "SignRBFTx",
			logs.FieldNonce:  params.TssNonce,
			logs.FieldTx:     oldTx.MsgTx().TxID(),
		}
		logger = signer.Logger().Std.With().Fields(lf).Logger()
	)

	// parse recent fee rate from CCTX
	cctxRate, err := strconv.ParseInt(params.GasPrice, 10, 64)
	if err != nil || cctxRate <= 0 {
		return nil, fmt.Errorf("cannot convert fee rate %s", params.GasPrice)
	}

	// initiate fee bumper
	fb := NewCPFPFeeBumper(signer.client, oldTx, cctxRate, minRelayFee)
	err = fb.FetchFeeBumpInfo(rpc.GetTotalMempoolParentsSizeNFees, logger)
	if err != nil {
		return nil, errors.Wrap(err, "FetchFeeBumpInfo failed")
	}

	// bump tx fees
	newTx, additionalFees, err := fb.BumpTxFee()
	if err != nil {
		return nil, errors.Wrap(err, "BumpTxFee failed")
	}
	logger.Info().Msgf("BumpTxFee success, additional fees: %d satoshis", additionalFees)

	// collect input amounts for later signing
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
	err = signer.SignTx(ctx, newTx, inAmounts, 0, params.TssNonce)
	if err != nil {
		return nil, errors.Wrap(err, "SignTx failed")
	}

	return newTx, nil
}
