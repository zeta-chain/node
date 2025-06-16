package client

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	types "github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
)

const (
	// FeeRateRegnet is the hardcoded fee rate for regnet
	FeeRateRegnet = 1

	// FeeRateRegnetRBF is the hardcoded fee rate for regnet RBF
	// The zetacore bumps CCTX's fee rate every 40 minutes, and we can't wait that long in the E2E test.
	// For simplicity, zetaclient uses a constant fee rate (> above 1 sat/vB) to test RBF in the regnet.
	FeeRateRegnetRBF = 5

	// maxBTCSupply is the maximum supply of Bitcoin
	maxBTCSupply = 21000000.0
)

// MempoolTxsAndFees contains the information of pending mempool txs and fees
type MempoolTxsAndFees struct {
	TotalTxs   int64
	TotalFees  int64
	TotalVSize int64
	AvgFeeRate uint64
}

// GetBlockVerboseByStr alias for GetBlockVerbose
func (c *Client) GetBlockVerboseByStr(ctx context.Context, blockHash string) (*types.GetBlockVerboseTxResult, error) {
	h, err := strToHash(blockHash)
	if err != nil {
		return nil, err
	}

	return c.GetBlockVerbose(ctx, h)
}

// GetBlockHeightByStr alias for GetBlockVerbose
func (c *Client) GetBlockHeightByStr(ctx context.Context, blockHash string) (int64, error) {
	h, err := strToHash(blockHash)
	if err != nil {
		return 0, err
	}

	res, err := c.GetBlockVerbose(ctx, h)
	if err != nil {
		return 0, errors.Wrap(err, "unable to get block verbose")
	}

	return res.Height, nil
}

// GetTransactionByStr alias for GetTransaction
func (c *Client) GetTransactionByStr(
	ctx context.Context,
	hash string,
) (*chainhash.Hash, *types.GetTransactionResult, error) {
	h, err := strToHash(hash)
	if err != nil {
		return nil, nil, err
	}

	tx, err := c.GetTransaction(ctx, h)

	return h, tx, err
}

// GetRawTransactionByStr alias for GetRawTransaction
func (c *Client) GetRawTransactionByStr(ctx context.Context, hash string) (*btcutil.Tx, error) {
	h, err := strToHash(hash)
	if err != nil {
		return nil, err
	}

	return c.GetRawTransaction(ctx, h)
}

// GetRawTransactionResult gets the raw tx result
func (c *Client) GetRawTransactionResult(ctx context.Context,
	hash *chainhash.Hash,
	res *types.GetTransactionResult,
) (types.TxRawResult, error) {
	switch {
	case res.Confirmations == 0:
		// for pending tx, we query the raw tx
		rawResult, err := c.GetRawTransactionVerbose(ctx, hash)
		if err != nil {
			return types.TxRawResult{}, errors.Wrapf(err, "unable to get raw tx verbose %s", res.TxID)
		}

		return *rawResult, nil
	case res.Confirmations > 0:
		// for confirmed tx, we query the block

		blockHash, err := strToHash(res.BlockHash)
		if err != nil {
			return types.TxRawResult{}, err
		}

		block, err := c.GetBlockVerbose(ctx, blockHash)
		if err != nil {
			return types.TxRawResult{}, errors.Wrapf(err, "unable to get block versobse %s", res.BlockHash)
		}

		invalidRange := res.BlockIndex < 0 || res.BlockIndex >= int64(len(block.Tx))
		if invalidRange {
			return types.TxRawResult{}, errors.Errorf(
				"invalid block index: tx %s, block_index %d",
				res.TxID,
				res.BlockIndex,
			)
		}

		return block.Tx[res.BlockIndex], nil
	default:
		// res.Confirmations < 0 (meaning not included)
		return types.TxRawResult{}, fmt.Errorf("tx %s not included yet", hash)
	}
}

// GetEstimatedFeeRate gets estimated smart fee rate (sat/vB) targeting given block confirmation
func (c *Client) GetEstimatedFeeRate(ctx context.Context, confTarget int64) (satsPerByte uint64, err error) {
	// RPC 'EstimateSmartFee' is not available in regnet
	if c.isRegnet {
		return FeeRateRegnet, nil
	}

	feeResult, err := c.EstimateSmartFee(ctx, confTarget, &types.EstimateModeEconomical)
	switch {
	case err != nil:
		return 0, errors.Wrap(err, "unable to estimate smart fee")
	case feeResult.Errors != nil:
		return 0, fmt.Errorf("fee result contains errors: %s", feeResult.Errors)
	case feeResult.FeeRate == nil:
		return 0, errors.New("nil fee rate")
	}

	feeRate := *feeResult.FeeRate
	if feeRate <= 0 || feeRate >= maxBTCSupply {
		return 0, fmt.Errorf("invalid fee rate: %f", feeRate)
	}

	feeRateUint, err := common.FeeRateToSatPerByte(feeRate)
	if err != nil {
		return 0, errors.Wrapf(err, "invalid fee rate: %f", feeRate)
	}

	return feeRateUint, nil
}

// GetTransactionFeeAndRate gets the transaction fee and rate for a given tx result
func (c *Client) GetTransactionFeeAndRate(ctx context.Context, rawResult *types.TxRawResult) (int64, int64, error) {
	var (
		totalInputValue  int64
		totalOutputValue int64
	)

	// make sure the tx Vsize is not zero (should not happen)
	if rawResult.Vsize <= 0 {
		return 0, 0, fmt.Errorf("tx %s has non-positive Vsize: %d", rawResult.Txid, rawResult.Vsize)
	}

	// sum up total input value
	for _, vin := range rawResult.Vin {
		prevTx, err := c.GetRawTransactionByStr(ctx, vin.Txid)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "failed to get previous tx: %s", vin.Txid)
		}
		totalInputValue += prevTx.MsgTx().TxOut[vin.Vout].Value
	}

	// query the raw tx
	tx, err := c.GetRawTransactionByStr(ctx, rawResult.Txid)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to get tx: %s", rawResult.Txid)
	}

	// sum up total output value
	for _, vout := range tx.MsgTx().TxOut {
		totalOutputValue += vout.Value
	}

	// calculate the transaction fee in satoshis
	fee := totalInputValue - totalOutputValue
	if fee < 0 { // never happens
		return 0, 0, fmt.Errorf("got negative fee: %d", fee)
	}

	// Note: the calculation uses 'Vsize' returned by RPC to simplify dev experience:
	// 	- 1. the devs could use the same value returned by their RPC endpoints to estimate deposit fee.
	// 	- 2. the devs don't have to bother 'Vsize' calculation, even though there is more accurate formula.
	//		 Moreoever, the accurate 'Vsize' is usually an adjusted size (float value) by Bitcoin Core.
	//	- 3. the 'Vsize' calculation could depend on program language and the library used.
	//
	// calculate the fee rate in satoshis/vByte
	// #nosec G115 always in range
	feeRate := fee / int64(rawResult.Vsize)

	return fee, feeRate, nil
}

// GetTransactionInputSpender get the spender address of the given transaction input (vin)
func (c *Client) GetTransactionInputSpender(ctx context.Context, txid string, vout uint32) (string, error) {
	preTx, err := c.GetRawTransactionByStr(ctx, txid)
	if err != nil {
		return "", errors.Wrapf(err, "unable to get raw transaction %s", txid)
	}

	// #nosec G115 - always in range
	if len(preTx.MsgTx().TxOut) <= int(vout) {
		return "", fmt.Errorf("vout index %d out of range for tx %s", vout, txid)
	}

	// decode sender address from previous pkScript
	pkScript := preTx.MsgTx().TxOut[vout].PkScript

	return common.DecodeSenderFromScript(pkScript, c.NetParams())
}

// GetTransactionInitiator get the transaction initiator address of the given transaction
// The initiator is defined as the spender of the first input of the given transaction.
func (c *Client) GetTransactionInitiator(ctx context.Context, txid string) (string, error) {
	tx, err := c.GetRawTransactionByStr(ctx, txid)
	if err != nil {
		return "", errors.Wrapf(err, "unable to get raw transaction %s", txid)
	}

	if len(tx.MsgTx().TxIn) == 0 {
		return "", fmt.Errorf("tx %s has no inputs", txid)
	}

	// the first input
	preTxid := tx.MsgTx().TxIn[0].PreviousOutPoint.Hash.String()
	preVout := tx.MsgTx().TxIn[0].PreviousOutPoint.Index

	// get spender of the first input
	initiator, err := c.GetTransactionInputSpender(ctx, preTxid, preVout)
	if err != nil {
		return "", errors.Wrapf(err, "unable to get transaction input spender %s", txid)
	}

	return initiator, nil
}

// Healthcheck returns the latest block timestamp
func (c *Client) Healthcheck(ctx context.Context) (time.Time, error) {
	bn, err := c.GetBlockCount(ctx)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get block count")
	}

	hash, err := c.GetBlockHash(ctx, bn)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get block hash")
	}

	header, err := c.GetBlockHeader(ctx, hash)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "unable to get block header")
	}

	return header.Timestamp, nil
}

// GetMempoolTxsAndFees returns the information of all pending parent txs and fees of a given tx (inclusive)
//
// A parent tx is defined as:
//   - a tx that is also pending in the mempool
//   - a tx that has its first output spent by the child as first input
func (c *Client) GetMempoolTxsAndFees(
	ctx context.Context,
	childHash string,
) (txsAndFees MempoolTxsAndFees, err error) {
	totalFeesFloat := float64(0)

	// loop through all parents
	parentHash := childHash
	for {
		memplEntry, err := c.GetMempoolEntry(ctx, parentHash)
		if err != nil {
			if isTxNotInMempoolError(err) {
				// not a mempool tx, stop looking for parents
				break
			}
			return txsAndFees, errors.Wrapf(err, "unable to get mempool entry for tx %s", parentHash)
		}

		// accumulate fees and vsize
		txsAndFees.TotalTxs++
		totalFeesFloat += memplEntry.Fee
		txsAndFees.TotalVSize += int64(memplEntry.VSize)

		// find the parent tx
		tx, err := c.GetRawTransactionByStr(ctx, parentHash)
		if err != nil {
			return txsAndFees, errors.Wrapf(err, "unable to get tx %s", parentHash)
		}
		parentHash = tx.MsgTx().TxIn[0].PreviousOutPoint.Hash.String()

		// check timeout to avoid infinite loop
		if deadline, ok := ctx.Deadline(); ok {
			if time.Now().After(deadline) {
				return txsAndFees, errors.Errorf("timeout reached on %dth tx: %s", txsAndFees.TotalTxs, parentHash)
			}
		}
	}

	// no pending tx found
	if txsAndFees.TotalTxs == 0 {
		return txsAndFees, errors.Errorf("given tx is not pending: %s", childHash)
	}

	// convert total fees to satoshis
	txsAndFees.TotalFees, err = common.GetSatoshis(totalFeesFloat)
	if err != nil {
		return txsAndFees, errors.Wrapf(err, "invalid total fees: %f", totalFeesFloat)
	}

	// sanity check, should never happen
	if txsAndFees.TotalVSize <= 0 {
		return txsAndFees, errors.Errorf("invalid totalVSize %d", txsAndFees.TotalVSize)
	}

	// calculate the average fee rate
	// #nosec G115 always positive
	txsAndFees.AvgFeeRate = uint64(math.Ceil(totalFeesFloat / float64(txsAndFees.TotalVSize)))

	return txsAndFees, nil
}

// isTxNotInMempoolError checks if the given error is due to the transaction not being in the mempool.
func isTxNotInMempoolError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "Transaction not in mempool")
}

func strToHash(s string) (*chainhash.Hash, error) {
	hash, err := chainhash.NewHashFromStr(s)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create btc hash from string")
	}

	return hash, nil
}
