// Copyright 2021 Evmos Foundation
// This file is part of Evmos' Ethermint library.
//
// The Ethermint library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Ethermint library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Ethermint library. If not, see https://github.com/evmos/ethermint/blob/main/LICENSE
package backend

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	tmrpcclient "github.com/cometbft/cometbft/rpc/client"
	tmrpctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/pkg/errors"

	rpctypes "github.com/zeta-chain/zetacore/rpc/types"
)

// GetTransactionByHash returns the Ethereum format transaction identified by Ethereum transaction hash
func (b *Backend) GetTransactionByHash(txHash common.Hash) (*rpctypes.RPCTransaction, error) {
	res, additional, err := b.GetTxByEthHash(txHash)
	hexTx := txHash.Hex()
	if err != nil {
		return b.getTransactionByHashPending(txHash)
	}

	resBlock, err := b.TendermintBlockByNumber(rpctypes.BlockNumber(res.Height))
	if err != nil {
		b.logger.Debug("block not found", "height", res.Height, "error", err.Error())
		return nil, err
	}

	blockRes, err := b.TendermintBlockResultByNumber(&res.Height)
	if err != nil {
		b.logger.Debug("failed to retrieve block results", "height", res.Height, "error", err.Error())
		return nil, nil
	}

	var ethMsg *evmtypes.MsgEthereumTx
	// if additional fields are empty we can try to get MsgEthereumTx from sdk.Msg array
	if additional == nil {
		// #nosec G115 always in range
		if int(res.TxIndex) >= len(resBlock.Block.Txs) {
			b.logger.Error("tx out of bounds")
			return nil, fmt.Errorf("tx out of bounds")
		}
		tx, err := b.clientCtx.TxConfig.TxDecoder()(resBlock.Block.Txs[res.TxIndex])
		if err != nil {
			b.logger.Debug("decoding failed", "error", err.Error())
			return nil, fmt.Errorf("failed to decode tx: %w", err)
		}
		ethMsg = tx.GetMsgs()[res.MsgIndex].(*evmtypes.MsgEthereumTx)
		if ethMsg == nil {
			b.logger.Error("failed to get eth msg from sdk.Msgs")
			return nil, fmt.Errorf("failed to get eth msg from sdk.Msgs")
		}
	} else {
		// if additional fields are not empty try to parse synthetic tx from them
		ethMsg = b.parseSyntethicTxFromAdditionalFields(additional)
		if ethMsg == nil {
			b.logger.Error("failed to get synthetic eth msg from additional fields")
			return nil, fmt.Errorf("failed to get synthetic eth msg from additional fields")
		}
	}

	if res.EthTxIndex == -1 {
		// Fallback to find tx index by iterating all valid eth transactions
		msgs, _ := b.EthMsgsFromTendermintBlock(resBlock, blockRes)
		for i := range msgs {
			if msgs[i].Hash == hexTx {
				// #nosec G115 always in range
				res.EthTxIndex = int32(i)
				break
			}
		}
	}

	// if we still unable to find the eth tx index, return error, shouldn't happen.
	if res.EthTxIndex == -1 && additional == nil {
		return nil, errors.New("can't find index of ethereum tx")
	}
	if res.EthTxIndex == -1 {
		res.EthTxIndex = 0
	}

	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.logger.Error(
			"failed to fetch Base Fee from prunned block. Check node prunning configuration",
			"height",
			blockRes.Height,
			"error",
			err,
		)
	}

	return rpctypes.NewTransactionFromMsg(
		ethMsg,
		common.BytesToHash(resBlock.BlockID.Hash.Bytes()),
		// #nosec G115 always positive
		uint64(res.Height),
		// #nosec G115 always positive
		uint64(res.EthTxIndex),
		baseFee,
		b.chainID,
		additional,
	)
}

// getTransactionByHashPending find pending tx from mempool
func (b *Backend) getTransactionByHashPending(txHash common.Hash) (*rpctypes.RPCTransaction, error) {
	hexTx := txHash.Hex()
	// try to find tx in mempool
	txs, err := b.PendingTransactions()
	if err != nil {
		b.logger.Debug("tx not found", "hash", hexTx, "error", err.Error())
		return nil, nil
	}

	for _, tx := range txs {
		msg, err := evmtypes.UnwrapEthereumMsg(tx, txHash)
		if err != nil {
			// not ethereum tx
			continue
		}

		if msg.Hash == hexTx {
			// use zero block values since it's not included in a block yet
			rpctx, err := rpctypes.NewTransactionFromMsg(
				msg,
				common.Hash{},
				uint64(0),
				uint64(0),
				nil,
				b.chainID,
				nil,
			)
			if err != nil {
				return nil, err
			}
			return rpctx, nil
		}
	}

	b.logger.Debug("tx not found", "hash", hexTx)
	return nil, nil
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (b *Backend) GetTransactionReceipt(hash common.Hash) (map[string]interface{}, error) {
	hexTx := hash.Hex()
	b.logger.Debug("eth_getTransactionReceipt", "hash", hexTx)
	res, additional, err := b.GetTxByEthHash(hash)
	if err != nil {
		b.logger.Debug("tx not found", "hash", hexTx, "error", err.Error())
		return nil, nil
	}

	resBlock, err := b.TendermintBlockByNumber(rpctypes.BlockNumber(res.Height))
	if err != nil {
		b.logger.Debug("block not found", "height", res.Height, "error", err.Error())
		return nil, nil
	}

	var txData evmtypes.TxData
	var ethMsg *evmtypes.MsgEthereumTx
	// if additional fields are empty we can try to get MsgEthereumTx from sdk.Msg array
	if additional == nil {
		// #nosec G115 always in range
		if int(res.TxIndex) >= len(resBlock.Block.Txs) {
			b.logger.Error("tx out of bounds")
			return nil, fmt.Errorf("tx out of bounds")
		}
		tx, err := b.clientCtx.TxConfig.TxDecoder()(resBlock.Block.Txs[res.TxIndex])
		if err != nil {
			b.logger.Debug("decoding failed", "error", err.Error())
			return nil, fmt.Errorf("failed to decode tx: %w", err)
		}
		ethMsg = tx.GetMsgs()[res.MsgIndex].(*evmtypes.MsgEthereumTx)
		if ethMsg == nil {
			b.logger.Error("failed to get eth msg")
			return nil, fmt.Errorf("failed to get eth msg")
		}

		txData, err = evmtypes.UnpackTxData(ethMsg.Data)
		if err != nil {
			b.logger.Error("failed to unpack tx data", "error", err.Error())
			return nil, err
		}
	} else {
		// if additional fields are not empty try to parse synthetic tx from them
		ethMsg = b.parseSyntethicTxFromAdditionalFields(additional)
		if ethMsg == nil {
			b.logger.Error("failed to parse tx")
			return nil, fmt.Errorf("failed to parse tx")
		}
	}

	cumulativeGasUsed := uint64(0)
	blockRes, err := b.TendermintBlockResultByNumber(&res.Height)
	if err != nil {
		b.logger.Debug("failed to retrieve block results", "height", res.Height, "error", err.Error())
		return nil, nil
	}
	for _, txResult := range blockRes.TxsResults[0:res.TxIndex] {
		// #nosec G115 always positive
		cumulativeGasUsed += uint64(txResult.GasUsed)
	}
	cumulativeGasUsed += res.CumulativeGasUsed

	var status hexutil.Uint
	if res.Failed {
		status = hexutil.Uint(ethtypes.ReceiptStatusFailed)
	} else {
		status = hexutil.Uint(ethtypes.ReceiptStatusSuccessful)
	}

	chainID, err := b.ChainID()
	if err != nil {
		return nil, err
	}

	var from common.Address
	if additional != nil {
		from = common.HexToAddress(ethMsg.From)
	} else if ethMsg.Data != nil {
		from, err = ethMsg.GetSender(chainID.ToInt())
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("failed to parse receipt")
	}

	// parse tx logs from events
	// #nosec G115 always in range
	logs, err := TxLogsFromEvents(blockRes.TxsResults[res.TxIndex].Events, int(res.MsgIndex))
	if err != nil {
		b.logger.Debug("failed to parse logs", "hash", hexTx, "error", err.Error())
	}

	if res.EthTxIndex == -1 {
		// Fallback to find tx index by iterating all valid eth transactions
		msgs, _ := b.EthMsgsFromTendermintBlock(resBlock, blockRes)
		for i := range msgs {
			if msgs[i].Hash == hexTx {
				// #nosec G115 always in range
				res.EthTxIndex = int32(i)
				break
			}
		}
	}

	// return error if still unable to find the eth tx index
	if res.EthTxIndex == -1 {
		if additional != nil {
			res.EthTxIndex = 0
		} else {
			return nil, errors.New("can't find index of ethereum tx")
		}
	}
	to := &common.Address{}
	var txType uint8

	if txData == nil {
		// #nosec G115 always in range
		txType = uint8(additional.Type)
		*to = additional.Recipient
	} else {
		txType = ethMsg.AsTransaction().Type()
		to = txData.GetTo()
	}

	receipt := map[string]interface{}{
		// Consensus fields: These fields are defined by the Yellow Paper
		"status":            status,
		"cumulativeGasUsed": hexutil.Uint64(cumulativeGasUsed),
		"logsBloom":         ethtypes.BytesToBloom(ethtypes.LogsBloom(logs)),
		"logs":              logs,

		// Implementation fields: These fields are added by geth when processing a transaction.
		// They are stored in the chain database.
		"transactionHash": hash,
		"contractAddress": nil,
		"gasUsed":         hexutil.Uint64(res.GasUsed),

		// Inclusion information: These fields provide information about the inclusion of the
		// transaction corresponding to this receipt.
		"blockHash":        common.BytesToHash(resBlock.Block.Header.Hash()).Hex(),
		"blockNumber":      hexutil.Uint64(res.Height),
		"transactionIndex": hexutil.Uint64(res.EthTxIndex),

		// sender and receiver (contract or EOA) addreses
		"from": from,
		"to":   to,
		// #nosec G115 uint8 -> uint false positive
		"type": hexutil.Uint(txType),
	}

	if logs == nil {
		receipt["logs"] = [][]*ethtypes.Log{}
	}

	if txData != nil {
		// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
		if txData.GetTo() == nil {
			receipt["contractAddress"] = crypto.CreateAddress(from, txData.GetNonce())
		}

		if dynamicTx, ok := txData.(*evmtypes.DynamicFeeTx); ok {
			baseFee, err := b.BaseFee(blockRes)
			if err != nil {
				// tolerate the error for pruned node.
				b.logger.Error("fetch basefee failed, node is pruned?", "height", res.Height, "error", err)
			} else {
				receipt["effectiveGasPrice"] = hexutil.Big(*dynamicTx.EffectiveGasPrice(baseFee))
			}
		}
	}
	return receipt, nil
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
func (b *Backend) GetTransactionByBlockHashAndIndex(
	hash common.Hash,
	idx hexutil.Uint,
) (*rpctypes.RPCTransaction, error) {
	b.logger.Debug("eth_getTransactionByBlockHashAndIndex", "hash", hash.Hex(), "index", idx)
	sc, ok := b.clientCtx.Client.(tmrpcclient.SignClient)
	if !ok {
		return nil, errors.New("invalid rpc client")
	}
	block, err := sc.BlockByHash(b.ctx, hash.Bytes())
	if err != nil {
		b.logger.Debug("block not found", "hash", hash.Hex(), "error", err.Error())
		return nil, nil
	}

	if block.Block == nil {
		b.logger.Debug("block not found", "hash", hash.Hex())
		return nil, nil
	}

	return b.GetTransactionByBlockAndIndex(block, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (b *Backend) GetTransactionByBlockNumberAndIndex(
	blockNum rpctypes.BlockNumber,
	idx hexutil.Uint,
) (*rpctypes.RPCTransaction, error) {
	b.logger.Debug("eth_getTransactionByBlockNumberAndIndex", "number", blockNum, "index", idx)

	block, err := b.TendermintBlockByNumber(blockNum)
	if err != nil {
		b.logger.Debug("block not found", "height", blockNum.Int64(), "error", err.Error())
		return nil, nil
	}

	if block.Block == nil {
		b.logger.Debug("block not found", "height", blockNum.Int64())
		return nil, nil
	}

	return b.GetTransactionByBlockAndIndex(block, idx)
}

// GetTxByEthHash uses `/tx_query` to find transaction by ethereum tx hash
// TODO: Don't need to convert once hashing is fixed on Tendermint
// https://github.com/tendermint/tendermint/issues/6539
func (b *Backend) GetTxByEthHash(hash common.Hash) (*ethermint.TxResult, *rpctypes.TxResultAdditionalFields, error) {
	if b.indexer != nil {
		txRes, err := b.indexer.GetByTxHash(hash)
		if err != nil {
			return nil, nil, err
		}
		return txRes, nil, nil
	}

	// fallback to tendermint tx indexer
	query := fmt.Sprintf("%s.%s='%s'", evmtypes.TypeMsgEthereumTx, evmtypes.AttributeKeyEthereumTxHash, hash.Hex())
	txResult, txAdditional, err := b.queryTendermintTxIndexer(query, func(txs *rpctypes.ParsedTxs) *rpctypes.ParsedTx {
		return txs.GetTxByHash(hash)
	})
	if err != nil {
		return nil, nil, errorsmod.Wrapf(err, "GetTxByEthHash %s", hash.Hex())
	}
	return txResult, txAdditional, nil
}

// GetTxByTxIndex uses `/tx_query` to find transaction by tx index of valid ethereum txs
func (b *Backend) GetTxByTxIndex(
	height int64,
	index uint,
) (*ethermint.TxResult, *rpctypes.TxResultAdditionalFields, error) {
	if b.indexer != nil {
		// #nosec G115 always in range
		txRes, err := b.indexer.GetByBlockAndIndex(height, int32(index))
		if err == nil {
			return txRes, nil, nil
		}
	}

	// fallback to tendermint tx indexer
	query := fmt.Sprintf("tx.height=%d AND %s.%s=%d",
		height, evmtypes.TypeMsgEthereumTx,
		evmtypes.AttributeKeyTxIndex, index,
	)
	txResult, txAdditional, err := b.queryTendermintTxIndexer(query, func(txs *rpctypes.ParsedTxs) *rpctypes.ParsedTx {
		// #nosec G115 always in range
		return txs.GetTxByTxIndex(int(index))
	})
	if err != nil {
		return nil, nil, errorsmod.Wrapf(err, "GetTxByTxIndex %d %d", height, index)
	}
	return txResult, txAdditional, nil
}

// queryTendermintTxIndexer query tx in tendermint tx indexer
func (b *Backend) queryTendermintTxIndexer(
	query string,
	txGetter func(*rpctypes.ParsedTxs) *rpctypes.ParsedTx,
) (*ethermint.TxResult, *rpctypes.TxResultAdditionalFields, error) {
	resTxs, err := b.clientCtx.Client.TxSearch(b.ctx, query, false, nil, nil, "")
	if err != nil {
		return nil, nil, err
	}
	if len(resTxs.Txs) == 0 {
		return nil, nil, errors.New("ethereum tx not found")
	}
	txResult := resTxs.Txs[0]
	if !rpctypes.TxSuccessOrExceedsBlockGasLimit(&txResult.TxResult) {
		return nil, nil, errors.New("invalid ethereum tx")
	}

	var tx sdk.Tx
	if txResult.TxResult.Code != 0 {
		// it's only needed when the tx exceeds block gas limit
		tx, err = b.clientCtx.TxConfig.TxDecoder()(txResult.Tx)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid ethereum tx")
		}
	}

	return rpctypes.ParseTxIndexerResult(txResult, tx, txGetter)
}

// GetTransactionByBlockAndIndex is the common code shared by `GetTransactionByBlockNumberAndIndex` and `GetTransactionByBlockHashAndIndex`.
func (b *Backend) GetTransactionByBlockAndIndex(
	block *tmrpctypes.ResultBlock,
	idx hexutil.Uint,
) (*rpctypes.RPCTransaction, error) {
	blockRes, err := b.TendermintBlockResultByNumber(&block.Block.Height)
	if err != nil {
		return nil, nil
	}

	// #nosec G115 always in range
	i := int(idx)
	ethMsgs, additionals := b.EthMsgsFromTendermintBlock(block, blockRes)
	if i >= len(ethMsgs) {
		b.logger.Debug("block txs index out of bound", "index", i)
		return nil, nil
	}

	msg := ethMsgs[i]
	additional := additionals[i]
	baseFee, err := b.BaseFee(blockRes)
	if err != nil {
		// handle the error for pruned node.
		b.logger.Error(
			"failed to fetch Base Fee from prunned block. Check node prunning configuration",
			"height",
			block.Block.Height,
			"error",
			err,
		)
	}

	return rpctypes.NewTransactionFromMsg(
		msg,
		common.BytesToHash(block.Block.Hash()),
		// #nosec G115 always positive
		uint64(block.Block.Height),
		// #nosec G115 always positive
		uint64(idx),
		baseFee,
		b.chainID,
		additional,
	)
}
