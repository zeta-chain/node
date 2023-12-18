package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/observer/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Chain nonces queries

func (k Keeper) ChainNoncesAll(c context.Context, req *types.QueryAllChainNoncesRequest) (*types.QueryAllChainNoncesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var chainNoncess []types.ChainNonces
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	chainNoncesStore := prefix.NewStore(store, types.KeyPrefix(types.ChainNoncesKey))

	pageRes, err := query.Paginate(chainNoncesStore, req.Pagination, func(key []byte, value []byte) error {
		var chainNonces types.ChainNonces
		if err := k.cdc.Unmarshal(value, &chainNonces); err != nil {
			return err
		}

		chainNoncess = append(chainNoncess, chainNonces)
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

	return &types.QueryGetChainNoncesResponse{ChainNonces: val}, nil
}

// Pending nonces queries

func (k Keeper) PendingNoncesAll(c context.Context, req *types.QueryAllPendingNoncesRequest) (*types.QueryAllPendingNoncesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	list, pageRes, err := k.GetAllPendingNoncesPaginated(ctx, req.Pagination)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllPendingNoncesResponse{
		PendingNonces: list,
		Pagination:    pageRes,
	}, nil
}

func (k Keeper) PendingNoncesByChain(c context.Context, req *types.QueryPendingNoncesByChainRequest) (*types.QueryPendingNoncesByChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	tss, found := k.GetTSS(ctx)
	if !found {
		return nil, status.Error(codes.NotFound, "tss not found")
	}
	list, found := k.GetPendingNonces(ctx, tss.TssPubkey, req.ChainId)
	if !found {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("pending nonces not found for chain id : %d", req.ChainId))
	}

	return &types.QueryPendingNoncesByChainResponse{
		PendingNonces: list,
	}, nil
}
