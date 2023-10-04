package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func getInTrackerKey(chainID int64, txHash string) string {
	return fmt.Sprintf("%d-%s", chainID, txHash)
}

// SetInTxTracker set a specific InTxTracker in the store from its index
func (k Keeper) SetInTxTracker(ctx sdk.Context, InTxTracker types.InTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	b := k.cdc.MustMarshal(&InTxTracker)
	key := types.KeyPrefix(getInTrackerKey(InTxTracker.ChainId, InTxTracker.TxHash))
	store.Set(key, b)
}

// GetInTxTracker returns a InTxTracker from its index
func (k Keeper) GetInTxTracker(ctx sdk.Context, chainID int64, txHash string) (val types.InTxTracker, found bool) {
	key := getInTrackerKey(chainID, txHash)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	b := store.Get(types.KeyPrefix(key))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveInTxTrackerIfExists(ctx sdk.Context, chainID int64, txHash string) {
	key := getInTrackerKey(chainID, txHash)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	if store.Has(types.KeyPrefix(key)) {
		store.Delete(types.KeyPrefix(key))
	}
}
func (k Keeper) GetAllInTxTracker(ctx sdk.Context) (list []types.InTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.InTxTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return list
}

func (k Keeper) GetAllInTxTrackerForChain(ctx sdk.Context, chainID int64) (list []types.InTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte(fmt.Sprintf("%d-", chainID)))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.InTxTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return list
}

func (k Keeper) GetAllInTxTrackerForChainPaginated(ctx sdk.Context, chainID int64, pagination *query.PageRequest) (inTxTrackers []types.InTxTracker, pageRes *query.PageResponse, err error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(fmt.Sprintf("%s", types.InTxTrackerKeyPrefix)))
	chainStore := prefix.NewStore(store, types.KeyPrefix(fmt.Sprintf("%d-", chainID)))
	pageRes, err = query.Paginate(chainStore, pagination, func(key []byte, value []byte) error {
		var inTxTracker types.InTxTracker
		if err := k.cdc.Unmarshal(value, &inTxTracker); err != nil {
			return err
		}
		inTxTrackers = append(inTxTrackers, inTxTracker)
		return nil
	})
	return
}

func (k msgServer) AddToInTxTracker(goCtx context.Context, msg *types.MsgAddToInTxTracker) (*types.MsgAddToInTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, observerTypes.ErrSupportedChains
	}

	adminPolicyAccount := k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(observerTypes.Policy_Type_group1)
	isAdmin := msg.Creator == adminPolicyAccount

	isObserver, err := k.zetaObserverKeeper.IsAuthorized(ctx, msg.Creator, chain)
	if err != nil {
		ctx.Logger().Error("Error while checking if the account is an observer", err)
	}
	isProven := false
	if msg.Proof != nil {
		isProven, err = k.VerifyInTxTrackerProof(ctx, msg.Proof, msg.BlockHash, msg.TxIndex, msg.ChainId, msg.CoinType)
		if err != nil {
			return nil, err
		}
	}

	// Sender needs to be either the admin policy account or an observer
	if !(isAdmin || isObserver || isProven) {
		return nil, errorsmod.Wrap(observerTypes.ErrNotAuthorized, fmt.Sprintf("Creator %s", msg.Creator))
	}

	k.Keeper.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})
	return &types.MsgAddToInTxTrackerResponse{}, nil
}

func (k Keeper) VerifyInTxTrackerProof(ctx sdk.Context, proof *common.Proof, hash string, txIndex int64, chainID int64, coinType common.CoinType) (bool, error) {

	senderChain := common.GetChainFromChainID(chainID)
	if senderChain == nil {
		return false, types.ErrUnsupportedChain
	}

	if !senderChain.IsProvable() {
		return false, types.ErrCannotVerifyProof.Wrapf("chain %d does not support block header verification", chainID)
	}

	blockHash := eth.HexToHash(hash)
	res, found := k.zetaObserverKeeper.GetBlockHeader(ctx, blockHash.Bytes())
	if !found {
		return false, errorsmod.Wrap(observerTypes.ErrBlockHeaderNotFound, fmt.Sprintf("block header not found %s", blockHash))
	}

	// verify and process the proof
	val, err := proof.Verify(res.Header, int(txIndex))
	if err != nil && !common.IsErrorInvalidProof(err) {
		return false, err
	}
	var txx ethtypes.Transaction
	err = txx.UnmarshalBinary(val)
	if err != nil {
		return false, err
	}

	coreParams, found := k.zetaObserverKeeper.GetCoreParamsByChainID(ctx, senderChain.ChainId)
	if !found {
		return false, types.ErrUnsupportedChain.Wrapf("core params not found for chain %d", senderChain.ChainId)
	}
	tssRes, err := k.GetTssAddress(ctx, &types.QueryGetTssAddressRequest{})
	if err != nil {
		return false, err
	}
	tssAddr := eth.HexToAddress(tssRes.Eth)
	if tssAddr == (eth.Address{}) {
		return false, fmt.Errorf("tss address not found")
	}
	if common.IsEVMChain(chainID) {
		switch coinType {
		case common.CoinType_Zeta:
			if txx.To().Hex() != coreParams.ConnectorContractAddress {
				return false, types.ErrCannotVerifyProof.Wrapf("receiver is not connector contract for coin type %s", coinType)
			}
			return true, nil
		case common.CoinType_ERC20:
			if txx.To().Hex() != coreParams.Erc20CustodyContractAddress {
				return false, types.ErrCannotVerifyProof.Wrapf("receiver is not erc20Custory contract for coin type %s", coinType)
			}
			return true, nil
		case common.CoinType_Gas:
			if txx.To().Hex() != tssAddr.Hex() {
				return false, types.ErrCannotVerifyProof.Wrapf("receiver is not tssAddress contract for coin type %s", coinType)
			}
			return true, nil
		}
	}
	if common.IsBitcoinChain(chainID) {
		return false, types.ErrCannotVerifyProof.Wrapf("cannot verify proof for bitcoin chain %d", chainID)
	}
	return false, fmt.Errorf("proof failed")
}
