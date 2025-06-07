package keeper

import (
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/proofs"
	"github.com/zeta-chain/node/x/lightclient/types"
)

// GetAllBlockHeaders returns all block headers
func (k Keeper) GetAllBlockHeaders(ctx sdk.Context) (list []proofs.BlockHeader) {
	p := types.KeyPrefix(types.BlockHeaderKey)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), p)

	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val proofs.BlockHeader
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

// SetBlockHeader set a specific block header in the store from its index
func (k Keeper) SetBlockHeader(ctx sdk.Context, header proofs.BlockHeader) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))
	b := k.cdc.MustMarshal(&header)
	store.Set(header.Hash, b)
}

// GetBlockHeader returns a block header from its hash
func (k Keeper) GetBlockHeader(ctx sdk.Context, hash []byte) (val proofs.BlockHeader, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))

	b := store.Get(hash)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveBlockHeader removes a block header from the store
func (k Keeper) RemoveBlockHeader(ctx sdk.Context, hash []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlockHeaderKey))
	store.Delete(hash)
}

// CheckNewBlockHeader checks if a new block header is valid and can be added to the store
// It checks that the parent block header exists and that the block height is valid
// It also checks that the block header does not already exist
// It returns an error if the block header is invalid
// Upon success, it returns the parent hash
func (k Keeper) CheckNewBlockHeader(
	ctx sdk.Context,
	chainID int64,
	blockHash []byte,
	height int64,
	header proofs.HeaderData,
) ([]byte, error) {
	// check verification flags are set
	if err := k.CheckBlockHeaderVerificationEnabled(ctx, chainID); err != nil {
		return nil, err
	}

	// check if the block header already exists
	if _, found := k.GetBlockHeader(ctx, blockHash); found {
		return nil, cosmoserrors.Wrap(types.ErrBlockAlreadyExist, fmt.Sprintf("block hash: %x", blockHash))
	}

	// NOTE: error is checked in BasicValidation in msg; check again for extra caution
	parentHash, err := header.ParentHash()
	if err != nil {
		return nil, cosmoserrors.Wrap(types.ErrNoParentHash, err.Error())
	}

	// if the chain state exists and parent block header is not found, returns error
	// the Earliest/Latest height with this block header (after voting, not here)
	// if ChainState is found, check if the block height is valid
	// validate block height as it's not part of the header itself
	chainState, found := k.GetChainState(ctx, chainID)
	if found && chainState.EarliestHeight > 0 && chainState.EarliestHeight < height {
		if height != chainState.LatestHeight+1 {
			return nil, cosmoserrors.Wrap(types.ErrInvalidHeight, fmt.Sprintf(
				"invalid block height: wanted %d, got %d",
				chainState.LatestHeight+1,
				height,
			))
		}
		_, found = k.GetBlockHeader(ctx, parentHash)
		if !found {
			return nil, cosmoserrors.Wrap(types.ErrNoParentHash, "parent block header not found")
		}
	}

	// Check timestamp
	if err := header.ValidateTimestamp(ctx.BlockTime()); err != nil {
		return nil, cosmoserrors.Wrap(types.ErrInvalidTimestamp, err.Error())
	}

	return parentHash, nil
}

// AddBlockHeader adds a new block header to the store and updates the chain state
func (k Keeper) AddBlockHeader(
	ctx sdk.Context,
	chainID int64,
	height int64,
	blockHash []byte,
	header proofs.HeaderData,
	parentHash []byte,
) {
	// update chain state
	chainState, found := k.GetChainState(ctx, chainID)
	if !found {
		// create a new chain state if it does not exist
		chainState = types.ChainState{
			ChainId:         chainID,
			LatestHeight:    height,
			EarliestHeight:  height,
			LatestBlockHash: blockHash,
		}
	} else {
		// update the chain state with the latest block header
		// TODO: these checks would need to be more sophisticated for production
		// We should investigate and implement the correct assumptions for adding new block header
		// https://github.com/zeta-chain/node/issues/1997
		if height > chainState.LatestHeight {
			chainState.LatestHeight = height
			chainState.LatestBlockHash = blockHash
		}
		if chainState.EarliestHeight == 0 {
			chainState.EarliestHeight = height
		}
	}
	k.SetChainState(ctx, chainState)

	// add the block header to the store
	blockHeader := proofs.BlockHeader{
		Header:     header,
		Height:     height,
		Hash:       blockHash,
		ParentHash: parentHash,
		ChainId:    chainID,
	}
	k.SetBlockHeader(ctx, blockHeader)
}
