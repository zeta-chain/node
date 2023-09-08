package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func getOutTrackerKey(chainID int64, txHash string) string {
	return fmt.Sprintf("%d-%s", chainID, txHash)
}

// SetInTxTracker set a specific InTxTracker in the store from its index
func (k Keeper) SetInTxTracker(ctx sdk.Context, InTxTracker types.InTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	b := k.cdc.MustMarshal(&InTxTracker)
	key := types.KeyPrefix(getOutTrackerKey(InTxTracker.ChainId, InTxTracker.TxHash))
	store.Set(key, b)
}

// GetInTxTracker returns a InTxTracker from its index
func (k Keeper) GetInTxTracker(ctx sdk.Context, chainID int64, txHash string) (val types.InTxTracker, found bool) {
	key := getOutTrackerKey(chainID, txHash)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	b := store.Get(types.KeyPrefix(key))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllInTxTrackerForChain(ctx sdk.Context, chainId int64) (list []types.InTxTracker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InTxTrackerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte(fmt.Sprintf("%d-", chainId)))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.InTxTracker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return list
}

func (k msgServer) AddToInTxTracker(goCtx context.Context, msg *types.MsgAddToInTxTracker) (*types.MsgAddToInTxTrackerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, zetaObserverTypes.ErrSupportedChains
	}

	adminPolicyAccount := k.zetaObserverKeeper.GetParams(ctx).GetAdminPolicyAccount(zetaObserverTypes.Policy_Type_out_tx_tracker)
	isAdmin := msg.Creator == adminPolicyAccount

	isObserver, err := k.zetaObserverKeeper.IsAuthorized(ctx, msg.Creator, chain)
	if err != nil {
		ctx.Logger().Error("Error while checking if the account is an observer", err)
	}
	// Sender needs to be either the admin policy account or an observer
	if !(isAdmin || isObserver) {
		return nil, sdkerrors.Wrap(zetaObserverTypes.ErrNotAuthorized, fmt.Sprintf("Creator %s", msg.Creator))
	}

	k.Keeper.SetInTxTracker(ctx, types.InTxTracker{
		ChainId:  msg.ChainId,
		TxHash:   msg.TxHash,
		CoinType: msg.CoinType,
	})
	return &types.MsgAddToInTxTrackerResponse{}, nil
}
