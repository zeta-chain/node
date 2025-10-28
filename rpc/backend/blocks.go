package backend

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"

	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/zeta-chain/node/rpc/types"

	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
)

// BlockNumber returns the current block number in abci app state. Because abci
// app state could lag behind from cometbft latest block, it's more stable for
// the client to use the latest block number in abci app state than cometbft
// rpc.
func (b *Backend) BlockNumber() (hexutil.Uint64, error) {
	// do any grpc query, ignore the response and use the returned block height
	var header metadata.MD
	_, err := b.QueryClient.Params(b.Ctx, &evmtypes.QueryParamsRequest{}, grpc.Header(&header))
	if err != nil {
		return hexutil.Uint64(0), err
	}

	blockHeightHeader := header.Get(grpctypes.GRPCBlockHeightHeader)
	if headerLen := len(blockHeightHeader); headerLen != 1 {
		return 0, fmt.Errorf("unexpected '%s' gRPC header length; got %d, expected: %d", grpctypes.GRPCBlockHeightHeader, headerLen, 1)
	}

	height, err := strconv.ParseUint(blockHeightHeader[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block height: %w", err)
	}

	return hexutil.Uint64(height), nil
}

// GetBlockByNumber returns the JSON-RPC compatible Ethereum block identified by
// block number. Depending on fullTx it either returns the full transaction
// objects or if false only the hashes of the transactions.
func (b *Backend) GetBlockByNumber(blockNum types.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	resBlock, err := b.CometBlockByNumber(blockNum)
	if err != nil {
		return nil, nil
	}

	// return if requested block height is greater than the current one
	if resBlock == nil || resBlock.Block == nil {
		return nil, nil
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resBlock.Block.Height)
	if err != nil {
		b.Logger.Debug("failed to fetch block result from CometBFT", "height", blockNum, "error", err.Error())
		return nil, nil
	}

	res, err := b.RPCBlockFromCometBlock(resBlock, blockRes, fullTx)
	if err != nil {
		b.Logger.Debug("RPCBlockFromCometBlock failed", "height", blockNum, "error", err.Error())
		return nil, err
	}

	return res, nil
}

// GetBlockByHash returns the JSON-RPC compatible Ethereum block identified by
// hash.
func (b *Backend) GetBlockByHash(hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	resBlock, err := b.CometBlockByHash(hash)
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		// block not found
		return nil, nil
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resBlock.Block.Height)
	if err != nil {
		b.Logger.Debug("failed to fetch block result from CometBFT", "block-hash", hash.String(), "error", err.Error())
		return nil, nil
	}

	res, err := b.RPCBlockFromCometBlock(resBlock, blockRes, fullTx)
	if err != nil {
		b.Logger.Debug("RPCBlockFromCometBlock failed", "hash", hash, "error", err.Error())
		return nil, err
	}

	return res, nil
}

// GetBlockTransactionCountByHash returns the number of Ethereum transactions in
// the block identified by hash.
func (b *Backend) GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint {
	block, err := b.RPCClient.BlockByHash(b.Ctx, hash.Bytes())
	if err != nil {
		b.Logger.Debug("block not found", "hash", hash.Hex(), "error", err.Error())
		return nil
	}

	if block.Block == nil {
		b.Logger.Debug("block not found", "hash", hash.Hex())
		return nil
	}

	return b.getBlockTransactionCount(block)
}

// GetBlockTransactionCountByNumber returns the number of Ethereum transactions
// in the block identified by number.
func (b *Backend) GetBlockTransactionCountByNumber(blockNum types.BlockNumber) *hexutil.Uint {
	block, err := b.CometBlockByNumber(blockNum)
	if err != nil {
		b.Logger.Debug("block not found", "height", blockNum.Int64(), "error", err.Error())
		return nil
	}

	if block.Block == nil {
		b.Logger.Debug("block not found", "height", blockNum.Int64())
		return nil
	}

	return b.getBlockTransactionCount(block)
}

// getBlockTransactionCount returns the number of Ethereum transactions in a
// given block.
func (b *Backend) getBlockTransactionCount(block *cmtrpctypes.ResultBlock) *hexutil.Uint {
	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &block.Block.Height)
	if err != nil {
		return nil
	}

	ethMsgs, _ := b.EthMsgsFromCometBlock(block, blockRes)
	n := hexutil.Uint(len(ethMsgs))
	return &n
}

// EthBlockByNumber returns the Ethereum Block identified by number.
func (b *Backend) EthBlockByNumber(blockNum types.BlockNumber) (*ethtypes.Block, error) {
	resBlock, err := b.CometBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		// block not found
		return nil, fmt.Errorf("block not found for height %d", blockNum)
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resBlock.Block.Height)
	if err != nil {
		return nil, fmt.Errorf("block result not found for height %d", resBlock.Block.Height)
	}

	ethBlock, err := b.EthBlockFromCometBlock(resBlock, blockRes)
	if err != nil {
		return nil, fmt.Errorf("failed to get eth block from comet block: %w", err)
	}

	return ethBlock, nil
}

// // GetBlockReceipts returns the receipts for a given block number or hash.
// func (b *Backend) GetBlockReceipts(
// 	blockNrOrHash types.BlockNumberOrHash,
// ) ([]map[string]interface{}, error) {
// 	blockNum, err := b.BlockNumberFromComet(blockNrOrHash)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get block number from hash: %w", err)
// 	}

// 	resBlock, err := b.CometBlockByNumber(blockNum)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get block by number: %w", err)
// 	}

// 	if resBlock == nil {
// 		return nil, fmt.Errorf("block not found for height %d", *blockNum.CmtHeight())
// 	}

// 	blockRes, err := b.RPCClient.BlockResults(b.Ctx, blockNum.CmtHeight())
// 	if err != nil {
// 		return nil, fmt.Errorf("block result not found for height %d", resBlock.Block.Height)
// 	}

// 	msgs := b.EthMsgsFromCometBlock(resBlock, blockRes)

// 	receipts, err := b.ReceiptsFromCometBlock(resBlock, blockRes, msgs)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get receipts from comet block: %w, ", err)
// 	}

// 	result := make([]map[string]interface{}, len(msgs))
// 	for i, msg := range msgs {
// 		var signer ethtypes.Signer
// 		tx := msg.AsTransaction()
// 		if tx.Protected() {
// 			signer = ethtypes.LatestSignerForChainID(tx.ChainId())
// 		} else {
// 			signer = ethtypes.FrontierSigner{}
// 		}
// 		from, err := msg.GetSenderLegacy(signer)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get sender: %w", err)
// 		}

// 		result[i], err = types.RPCMarshalReceipt(receipts[i], tx, from)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to marshal receipt")
// 		}
// 	}
// 	return result, nil
// }
