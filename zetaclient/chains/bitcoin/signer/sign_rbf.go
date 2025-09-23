package signer

import (
	"context"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/client"
	"github.com/zeta-chain/node/zetaclient/logs"
)

// SignRBFTx signs a RBF (Replace-By-Fee) to unblock last stuck outbound transaction.
//
// The key points:
//   - It reuses the stuck tx's inputs and outputs but gives a higher fee to miners.
//   - Funding the last stuck outbound will be considered as CPFP (child-pays-for-parent) by miners.
func (signer *Signer) SignRBFTx(ctx context.Context, txData *OutboundData, lastTx *btcutil.Tx) (*wire.MsgTx, error) {
	var (
		logger = signer.Logger().Std.With().
			Uint64(logs.FieldNonce, txData.nonce).
			Str(logs.FieldTx, lastTx.MsgTx().TxID()).
			Logger()

		cctxRate = txData.feeRateLatest
	)

	// 1. for E2E test in regnet, hardcoded fee rate is used as we can't wait 40 minutes for zetacore to bump the fee rate
	// 2. for testnet and mainnet, we must wait for zetacore to bump the fee rate before signing the RBF tx
	if signer.isRegnet {
		cctxRate = client.FeeRateRegnetRBF
	} else if !txData.feeRateBumped {
		return nil, errors.New("fee rate is not bumped by zetacore yet, please hold on")
	}

	// create fee bumper
	fb, err := NewCPFPFeeBumper(
		ctx,
		signer.bitcoinClient,
		signer.Chain(),
		lastTx,
		cctxRate,
		txData.minRelayFee,
		logger,
	)
	if err != nil {
		return nil, errors.Wrap(err, "NewCPFPFeeBumper failed")
	}

	// bump tx fees
	result, err := fb.BumpTxFee()
	if err != nil {
		return nil, errors.Wrap(err, "BumpTxFee failed")
	}
	logger.Info().
		Uint64("old_fee_rate", fb.txsAndFees.AvgFeeRate).
		Uint64("new_fee_rate", result.NewFeeRate).
		Int64("additional_fees", result.AdditionalFees).
		Msg("call to BumpTxFee succeed")

	// collect input amounts for signing
	inAmounts := make([]int64, len(result.NewTx.TxIn))
	for i, input := range result.NewTx.TxIn {
		preOut := input.PreviousOutPoint
		preTx, err := signer.bitcoinClient.GetRawTransaction(ctx, &preOut.Hash)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to get previous tx %s", preOut.Hash)
		}
		inAmounts[i] = preTx.MsgTx().TxOut[preOut.Index].Value
	}

	// sign the RBF tx
	err = signer.SignTx(ctx, result.NewTx, inAmounts, txData.height, txData.nonce)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign tx")
	}

	return result.NewTx, nil
}
