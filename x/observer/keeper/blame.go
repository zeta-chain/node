package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SetBlame(ctx sdk.Context, blame *types.Blame) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	b := k.cdc.MustMarshal(blame)
	store.Set([]byte(blame.Index), b)
}

func (k Keeper) GetBlame(ctx sdk.Context, index string) (val types.Blame, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllBlame(ctx sdk.Context) (BlameRecords []*types.Blame, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	found = false
	for ; iterator.Valid(); iterator.Next() {
		var val types.Blame
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		BlameRecords = append(BlameRecords, &val)
		found = true
	}
	return
}

func (k Keeper) GetBlameByChainAndNonce(ctx sdk.Context, chainID int64, nonce int64) (BlameRecords []*types.Blame, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.BlameKey))
	blamePrefix := fmt.Sprintf("%d-%d", chainID, nonce)
	iterator := sdk.KVStorePrefixIterator(store, []byte(blamePrefix))
	defer iterator.Close()
	found = false
	for ; iterator.Valid(); iterator.Next() {
		var val types.Blame
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		BlameRecords = append(BlameRecords, &val)
		found = true
	}
	return
}

// Query

func (k Keeper) BlameByIdentifier(goCtx context.Context, request *types.QueryBlameByIdentifierRequest) (*types.QueryBlameByIdentifierResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	blame, found := k.GetBlame(ctx, request.BlameIdentifier)
	if !found {
		return nil, status.Error(codes.NotFound, "blame info not found")
	}

	return &types.QueryBlameByIdentifierResponse{
		BlameInfo: &blame,
	}, nil
}

func (k Keeper) GetAllBlameRecords(goCtx context.Context, request *types.QueryAllBlameRecordsRequest) (*types.QueryAllBlameRecordsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	blameRecords, found := k.GetAllBlame(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "blame info not found")
	}

	return &types.QueryAllBlameRecordsResponse{
		BlameInfo: blameRecords,
	}, nil
}

func (k Keeper) BlameByChainAndNonce(goCtx context.Context, request *types.QueryBlameByChainAndNonceRequest) (*types.QueryBlameByChainAndNonceResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	blameRecords, found := k.GetBlameByChainAndNonce(ctx, request.ChainId, request.Nonce)
	if !found {
		return nil, status.Error(codes.NotFound, "blame info not found")
	}
	return &types.QueryBlameByChainAndNonceResponse{
		BlameInfo: blameRecords,
	}, nil
}
