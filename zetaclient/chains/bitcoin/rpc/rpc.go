package rpc

import (
	"fmt"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"

	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/config"
)

const (
	// RPCAlertLatency is the default threshold for RPC latency to be considered unhealthy and trigger an alert.
	// Bitcoin block time is 10 minutes, 1200s (20 minutes) is a reasonable threshold for Bitcoin
	RPCAlertLatency = time.Duration(1200) * time.Second
)

// NewRPCClient creates a new RPC client by the given config.
func NewRPCClient(btcConfig config.BTCConfig) (*rpcclient.Client, error) {
	connCfg := &rpcclient.ConnConfig{
		Host:         btcConfig.RPCHost,
		User:         btcConfig.RPCUsername,
		Pass:         btcConfig.RPCPassword,
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       btcConfig.RPCParams,
	}

	rpcClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating rpc client: %s", err)
	}

	err = rpcClient.Ping()
	if err != nil {
		return nil, fmt.Errorf("error ping the bitcoin server: %s", err)
	}
	return rpcClient, nil
}

// GetTxResultByHash gets the transaction result by hash
func GetTxResultByHash(
	rpcClient interfaces.BTCRPCClient,
	txID string,
) (*chainhash.Hash, *btcjson.GetTransactionResult, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetTxResultByHash: error NewHashFromStr: %s", txID)
	}

	// The Bitcoin node has to be configured to watch TSS address
	txResult, err := rpcClient.GetTransaction(hash)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "GetTxResultByHash: error GetTransaction %s", hash.String())
	}
	return hash, txResult, nil
}

// GetTXRawResultByHash gets the raw transaction by hash
func GetRawTxByHash(rpcClient interfaces.BTCRPCClient, txID string) (*btcutil.Tx, error) {
	hash, err := chainhash.NewHashFromStr(txID)
	if err != nil {
		return nil, errors.Wrapf(err, "GetRawTxByHash: error NewHashFromStr: %s", txID)
	}

	tx, err := rpcClient.GetRawTransaction(hash)
	if err != nil {
		return nil, errors.Wrapf(err, "GetRawTxByHash: error GetRawTransaction %s", txID)
	}
	return tx, nil
}

// GetBlockHeightByHash gets the block height by block hash
func GetBlockHeightByHash(
	rpcClient interfaces.BTCRPCClient,
	hash string,
) (int64, error) {
	// decode the block hash
	var blockHash chainhash.Hash
	err := chainhash.Decode(&blockHash, hash)
	if err != nil {
		return 0, errors.Wrapf(err, "GetBlockHeightByHash: error decoding block hash %s", hash)
	}

	// get block by hash
	block, err := rpcClient.GetBlockVerbose(&blockHash)
	if err != nil {
		return 0, errors.Wrapf(err, "GetBlockHeightByHash: error GetBlockVerbose %s", hash)
	}
	return block.Height, nil
}

// GetRawTxResult gets the raw tx result
func GetRawTxResult(
	rpcClient interfaces.BTCRPCClient,
	hash *chainhash.Hash,
	res *btcjson.GetTransactionResult,
) (btcjson.TxRawResult, error) {
	if res.Confirmations == 0 { // for pending tx, we query the raw tx directly
		rawResult, err := rpcClient.GetRawTransactionVerbose(hash) // for pending tx, we query the raw tx
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(
				err,
				"GetRawTxResult: error GetRawTransactionVerbose %s",
				res.TxID,
			)
		}
		return *rawResult, nil
	} else if res.Confirmations > 0 { // for confirmed tx, we query the block
		blkHash, err := chainhash.NewHashFromStr(res.BlockHash)
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "GetRawTxResult: error NewHashFromStr for block hash %s", res.BlockHash)
		}
		block, err := rpcClient.GetBlockVerboseTx(blkHash)
		if err != nil {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "GetRawTxResult: error GetBlockVerboseTx %s", res.BlockHash)
		}
		if res.BlockIndex < 0 || res.BlockIndex >= int64(len(block.Tx)) {
			return btcjson.TxRawResult{}, errors.Wrapf(err, "GetRawTxResult: invalid outbound with invalid block index, TxID %s, BlockIndex %d", res.TxID, res.BlockIndex)
		}
		return block.Tx[res.BlockIndex], nil
	}

	// res.Confirmations < 0 (meaning not included)
	return btcjson.TxRawResult{}, fmt.Errorf("GetRawTxResult: tx %s not included yet", hash)
}

// GetTransactionFeeAndRate gets the transaction fee and rate for a given tx result
func GetTransactionFeeAndRate(rpcClient interfaces.BTCRPCClient, rawResult *btcjson.TxRawResult) (int64, int64, error) {
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
		prevTx, err := GetRawTxByHash(rpcClient, vin.Txid)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "failed to get previous tx: %s", vin.Txid)
		}
		totalInputValue += prevTx.MsgTx().TxOut[vin.Vout].Value
	}

	// query the raw tx
	tx, err := GetRawTxByHash(rpcClient, rawResult.Txid)
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

// CheckRPCStatus checks the RPC status of the bitcoin chain
func CheckRPCStatus(client interfaces.BTCRPCClient, tssAddress btcutil.Address) (time.Time, error) {
	// query latest block number
	bn, err := client.GetBlockCount()
	if err != nil {
		return time.Time{}, errors.Wrap(err, "RPC failed on GetBlockCount, RPC down?")
	}

	// query latest block header
	hash, err := client.GetBlockHash(bn)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "RPC failed on GetBlockHash, RPC down?")
	}

	// query latest block header thru hash
	header, err := client.GetBlockHeader(hash)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "RPC failed on GetBlockHeader, RPC down?")
	}

	// should be able to list utxos owned by TSS address
	res, err := client.ListUnspentMinMaxAddresses(0, 1000000, []btcutil.Address{tssAddress})
	if err != nil {
		return time.Time{}, errors.Wrap(err, "can't list utxos of TSS address; TSS address is not imported?")
	}

	// TSS address should have utxos
	if len(res) == 0 {
		return time.Time{}, errors.New("TSS address has no utxos; TSS address is not imported?")
	}

	return header.Timestamp, nil
}
