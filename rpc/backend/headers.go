package backend

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	rpctypes "github.com/zeta-chain/node/rpc/types"
)

// GetBlockByNumber returns the JSON-RPC compatible Ethereum block identified by
// block number. Depending on fullTx it either returns the full transaction
// objects or if false only the hashes of the transactions.
func (b *Backend) GetHeaderByNumber(blockNum rpctypes.BlockNumber) (map[string]interface{}, error) {
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

	res, err := b.RPCHeaderFromCometBlock(resBlock, blockRes)
	if err != nil {
		b.Logger.Debug("RPCBlockFromCometBlock failed", "height", blockNum, "error", err.Error())
		return nil, err
	}

	return res, nil
}

// GetBlockByHash returns the JSON-RPC compatible Ethereum block identified by
// hash.
func (b *Backend) GetHeaderByHash(hash common.Hash) (map[string]interface{}, error) {
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

	res, err := b.RPCHeaderFromCometBlock(resBlock, blockRes)
	if err != nil {
		b.Logger.Debug("RPCBlockFromCometBlock failed", "hash", hash, "error", err.Error())
		return nil, err
	}

	return res, nil
}

// HeaderByNumber returns the block header identified by height.
func (b *Backend) HeaderByNumber(blockNum rpctypes.BlockNumber) (*ethtypes.Header, error) {
	resBlock, err := b.CometBlockByNumber(blockNum)
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		// block not found
		return nil, nil
	}

	blockRes, err := b.CometBlockResultByNumber(&resBlock.Block.Height)
	if err != nil {
		return nil, fmt.Errorf("header result not found for height %d", resBlock.Block.Height)
	}

	ethBlock, err := b.EthBlockFromCometBlock(resBlock, blockRes)
	if err != nil {
		return nil, fmt.Errorf("failed to get rpc block from comet block: %w", err)
	}

	return ethBlock.Header(), nil
}

// HeaderByHash returns the block header identified by hash.
func (b *Backend) HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error) {
	resBlock, err := b.CometBlockByHash(blockHash)
	if err != nil {
		return nil, err
	}

	if resBlock == nil {
		// block not found
		return nil, nil
	}

	blockRes, err := b.RPCClient.BlockResults(b.Ctx, &resBlock.Block.Height)
	if err != nil {
		b.Logger.Debug("failed to fetch block result from CometBFT", "block-hash", blockHash.String(), "error", err.Error())
		return nil, nil
	}

	ethBlock, err := b.EthBlockFromCometBlock(resBlock, blockRes)
	if err != nil {
		return nil, fmt.Errorf("failed to get rpc block from comet block: %w", err)
	}

	return ethBlock.Header(), nil
}
