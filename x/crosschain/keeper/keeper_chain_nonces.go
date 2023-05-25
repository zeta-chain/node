package keeper

import (
	"context"
	"fmt"

	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetChainNonces set a specific chainNonces in the store from its index
func (k Keeper) SetChainNonces(ctx sdk.Context, chainNonces types.ChainNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))
	b := k.cdc.MustMarshal(&chainNonces)
	store.Set(types.KeyPrefix(chainNonces.Index), b)
}

// GetChainNonces returns a chainNonces from its index
func (k Keeper) GetChainNonces(ctx sdk.Context, index string) (val types.ChainNonces, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))

	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveChainNonces removes a chainNonces from the store
func (k Keeper) RemoveChainNonces(ctx sdk.Context, index string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))
	store.Delete(types.KeyPrefix(index))
}

// GetAllChainNonces returns all chainNonces
func (k Keeper) GetAllChainNonces(ctx sdk.Context) (list []types.ChainNonces) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ChainNoncesKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ChainNonces
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Queries

func (k Keeper) ChainNoncesAll(c context.Context, req *types.QueryAllChainNoncesRequest) (*types.QueryAllChainNoncesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var chainNoncess []*types.ChainNonces
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	chainNoncesStore := prefix.NewStore(store, types.KeyPrefix(types.ChainNoncesKey))

	pageRes, err := query.Paginate(chainNoncesStore, req.Pagination, func(key []byte, value []byte) error {
		var chainNonces types.ChainNonces
		if err := k.cdc.Unmarshal(value, &chainNonces); err != nil {
			return err
		}

		chainNoncess = append(chainNoncess, &chainNonces)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllChainNoncesResponse{ChainNonces: chainNoncess, Pagination: pageRes}, nil
}

func (k Keeper) ChainNonces(c context.Context, req *types.QueryGetChainNoncesRequest) (*types.QueryGetChainNoncesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetChainNonces(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "not found")
	}

	return &types.QueryGetChainNoncesResponse{ChainNonces: &val}, nil
}

// MESSAGES

// Should be removed
func (k msgServer) NonceVoter(goCtx context.Context, msg *types.MsgNonceVoter) (*types.MsgNonceVoterResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := k.zetaObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.ChainId)
	if chain == nil {
		return nil, zetaObserverTypes.ErrSupportedChains
	}

	ok, err := k.IsAuthorized(ctx, msg.Creator, chain)
	if !ok {
		return nil, err
	}
	chainNonce, isFound := k.GetChainNonces(ctx, chain.ChainName.String())

	if isFound {
		isExisting := false
		for _, signer := range chainNonce.Signers {
			if signer == msg.Creator {
				isExisting = true
			}
		}
		if !isExisting {
			chainNonce.Signers = append(chainNonce.Signers, msg.Creator)
		}
		chainNonce.Nonce = msg.Nonce
	} else if !isFound {
		chainNonce = types.ChainNonces{
			Creator: msg.Creator,
			Index:   chain.ChainName.String(),
			ChainId: chain.ChainId,
			Nonce:   msg.Nonce,
			Signers: []string{msg.Creator},
		}
	} else {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("chainNonce vote msg does not match state: %v vs %v", msg, chainNonce))
	}

	//if hasSuperMajorityValidators(len(chainNonce.Signers), validators) {
	//	chainNonce.FinalizedHeight = uint64(ctx.BlockHeader().Height)
	//}

	k.SetChainNonces(ctx, chainNonce)
	return &types.MsgNonceVoterResponse{}, nil
}
