package backend

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	cmtrpctypes "github.com/cometbft/cometbft/rpc/core/types"

	"github.com/cosmos/evm/utils"
	rpctypes "github.com/zeta-chain/node/rpc/types"
)

// CometBlockByNumber returns a CometBFT-formatted block for a given
// block number
func (b *Backend) CometBlockByNumber(blockNum rpctypes.BlockNumber) (*cmtrpctypes.ResultBlock, error) {
	height, err := b.getHeightByBlockNum(blockNum)
	if err != nil {
		return nil, err
	}
	resBlock, err := b.RPCClient.Block(b.Ctx, &height)
	if err != nil {
		b.Logger.Debug("cometbft client failed to get block", "height", height, "error", err.Error())
		return nil, err
	}

	if resBlock.Block == nil {
		b.Logger.Debug("CometBlockByNumber block not found", "height", height)
		return nil, nil
	}

	return resBlock, nil
}

// CometHeaderByNumber returns a CometBFT-formatted header for a given
// block number
func (b *Backend) CometHeaderByNumber(blockNum rpctypes.BlockNumber) (*cmtrpctypes.ResultHeader, error) {
	height, err := b.getHeightByBlockNum(blockNum)
	if err != nil {
		return nil, err
	}
	return b.RPCClient.Header(b.Ctx, &height)
}

// CometBlockResultByNumber returns a CometBFT-formatted block result
// by block number
func (b *Backend) CometBlockResultByNumber(height *int64) (*cmtrpctypes.ResultBlockResults, error) {
	if height != nil && *height == 0 {
		height = nil
	}
	res, err := b.RPCClient.BlockResults(b.Ctx, height)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block result from CometBFT %d: %w", *height, err)
	}

	return res, nil
}

// CometBlockByHash returns a CometBFT-formatted block by block number
func (b *Backend) CometBlockByHash(blockHash common.Hash) (*cmtrpctypes.ResultBlock, error) {
	resBlock, err := b.RPCClient.BlockByHash(b.Ctx, blockHash.Bytes())
	if err != nil {
		b.Logger.Debug("CometBFT client failed to get block", "blockHash", blockHash.Hex(), "error", err.Error())
		return nil, err
	}

	if resBlock == nil || resBlock.Block == nil {
		b.Logger.Debug("CometBlockByHash block not found", "blockHash", blockHash.Hex())
		return nil, fmt.Errorf("block not found for hash %s", blockHash.Hex())
	}

	return resBlock, nil
}

func (b *Backend) getHeightByBlockNum(blockNum rpctypes.BlockNumber) (int64, error) {
	if blockNum == rpctypes.EthEarliestBlockNumber {
		status, err := b.ClientCtx.Client.Status(b.Ctx)
		if err != nil {
			return 0, errors.New("failed to get earliest block height")
		}

		return status.SyncInfo.EarliestBlockHeight, nil
	}

	height := blockNum.Int64()
	if height <= 0 {
		// In cometBFT, LatestBlockNumber, FinalizedBlockNumber, SafeBlockNumber all map to the latest block height.
		// Fetch the latest block number from the app state, more accurate than the CometBFT block store state.
		//
		// For PendingBlockNumber, we alsoe returns the latest block height.
		// The reason is that CometBFT does not have the concept of pending block,
		// and the application state is only updated when a block is committed.
		n, err := b.BlockNumber()
		if err != nil {
			return 0, err
		}
		height, err = utils.SafeHexToInt64(n)
		if err != nil {
			return 0, err
		}
	}

	return height, nil
}
