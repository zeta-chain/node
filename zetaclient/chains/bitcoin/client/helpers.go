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
	FeeRateRegnetRBF = 5

	// maxBTCSupply is the maximum supply of Bitcoin
	maxBTCSupply = 21000000.0
)

// IsRegnet returns true if the chain is regnet
func (c *Client) IsRegnet() bool {
	return c.isRegnet
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

// GetTotalMempoolParentsSizeNFees returns the information of all pending parent txs of a given tx (inclusive)
//
// A parent tx is defined as:
//   - a tx that is also pending in the mempool
//   - a tx that has its first output spent by the child as first input
//
// Returns: (totalTxs, totalFees, totalVSize, error)
func (c *Client) GetTotalMempoolParentsSizeNFees(
	ctx context.Context,
	childHash string,
	timeout time.Duration,
) (int64, float64, int64, uint64, error) {
	var (
		totalTxs   int64
		totalFees  float64
		totalVSize int64
		avgFeeRate uint64
	)

	// loop through all parents
	startTime := time.Now()
	parentHash := childHash
	for {
		memplEntry, err := c.GetMempoolEntry(ctx, parentHash)
		if err != nil {
			if strings.Contains(err.Error(), "Transaction not in mempool") {
				// not a mempool tx, stop looking for parents
				break
			}
			return 0, 0, 0, 0, errors.Wrapf(err, "unable to get mempool entry for tx %s", parentHash)
		}

		// accumulate fees and vsize
		totalTxs++
		totalFees += memplEntry.Fee
		totalVSize += int64(memplEntry.VSize)

		// find the parent tx
		tx, err := c.GetRawTransactionByStr(ctx, parentHash)
		if err != nil {
			return 0, 0, 0, 0, errors.Wrapf(err, "unable to get tx %s", parentHash)
		}
		parentHash = tx.MsgTx().TxIn[0].PreviousOutPoint.Hash.String()

		// check timeout to avoid infinite loop
		if time.Since(startTime) > timeout {
			return 0, 0, 0, 0, errors.Errorf("timeout reached on %dth tx: %s", totalTxs, parentHash)
		}
	}

	// no pending tx found
	if totalTxs == 0 {
		return 0, 0, 0, 0, errors.Errorf("given tx is not pending: %s", childHash)
	}

	// sanity check, should never happen
	if totalFees < 0 || totalVSize <= 0 {
		return 0, 0, 0, 0, errors.Errorf("invalid result: totalFees %f, totalVSize %d", totalFees, totalVSize)
	}

	// calculate the average fee rate
	// #nosec G115 always positive
	avgFeeRate = uint64(math.Ceil(totalFees / float64(totalVSize)))

	return totalTxs, totalFees, totalVSize, avgFeeRate, nil
}

func strToHash(s string) (*chainhash.Hash, error) {
	hash, err := chainhash.NewHashFromStr(s)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create btc hash from string")
	}

	return hash, nil
}
