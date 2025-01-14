package client

import (
	"context"
	"fmt"
	"time"

	types "github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/pkg/errors"
	"github.com/tendermint/btcd/btcjson"
)

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
func (c *Client) GetTransactionByStr(ctx context.Context, hash string) (*types.GetTransactionResult, error) {
	h, err := strToHash(hash)
	if err != nil {
		return nil, err
	}

	return c.GetTransaction(ctx, h)
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
	res *btcjson.GetTransactionResult,
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

// Healthcheck / checks the RPC status of the bitcoin chain. Returns the latest block timestamp
func (c *Client) Healthcheck(ctx context.Context, tssAddress btcutil.Address) (time.Time, error) {
	// 1. Query latest block header
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

	// 2. Query utxos owned by TSS address
	res, err := c.ListUnspentMinMaxAddresses(ctx, 0, 1000000, []btcutil.Address{tssAddress})
	switch {
	case err != nil:
		return time.Time{}, errors.Wrap(err, "unable to list TSS UTXOs")
	case len(res) == 0:
		return time.Time{}, errors.New("no UTXOs found for TSS")
	}

	return header.Timestamp, nil
}

func strToHash(s string) (*chainhash.Hash, error) {
	hash, err := chainhash.NewHashFromStr(s)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create btc hash from string")
	}

	return hash, nil
}
