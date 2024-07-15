package rpc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/pkg/errors"

	"github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/config"
)

const (
	// feeRateCountBackBlocks is the default number of blocks to look back for fee rate estimation
	feeRateCountBackBlocks = 2

	// defaultTestnetFeeRate is the default fee rate for testnet, 10 sat/byte
	defaultTestnetFeeRate = 10
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

// GetRecentFeeRate gets the highest fee rate from recent blocks
// Note: this method is only used for testnet
func GetRecentFeeRate(rpcClient interfaces.BTCRPCClient, netParams *chaincfg.Params) (uint64, error) {
	blockNumber, err := rpcClient.GetBlockCount()
	if err != nil {
		return 0, err
	}

	// get the highest fee rate among recent 'countBack' blocks to avoid underestimation
	highestRate := int64(0)
	for i := int64(0); i < feeRateCountBackBlocks; i++ {
		// get the block
		hash, err := rpcClient.GetBlockHash(blockNumber - i)
		if err != nil {
			return 0, err
		}
		block, err := rpcClient.GetBlockVerboseTx(hash)
		if err != nil {
			return 0, err
		}

		// computes the average fee rate of the block and take the higher rate
		avgFeeRate, err := bitcoin.CalcBlockAvgFeeRate(block, netParams)
		if err != nil {
			return 0, err
		}
		if avgFeeRate > highestRate {
			highestRate = avgFeeRate
		}
	}

	// use 10 sat/byte as default estimation if recent fee rate drops to 0
	if highestRate == 0 {
		highestRate = defaultTestnetFeeRate
	}

	// #nosec G115 always in range
	return uint64(highestRate), nil
}
