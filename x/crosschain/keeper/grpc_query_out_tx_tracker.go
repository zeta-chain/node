package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OutTxTrackerAll(c context.Context, req *types.QueryAllOutTxTrackerRequest) (*types.QueryAllOutTxTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var outTxTrackers []types.OutTxTracker
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	outTxTrackerStore := prefix.NewStore(store, types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	pageRes, err := query.Paginate(outTxTrackerStore, req.Pagination, func(key []byte, value []byte) error {
		var outTxTracker types.OutTxTracker
		if err := k.cdc.Unmarshal(value, &outTxTracker); err != nil {
			return err
		}

		outTxTrackers = append(outTxTrackers, outTxTracker)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllOutTxTrackerResponse{OutTxTracker: outTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) OutTxTrackerAllByChain(c context.Context, req *types.QueryAllOutTxTrackerByChainRequest) (*types.QueryAllOutTxTrackerByChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var outTxTrackers []types.OutTxTracker
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutTxTrackerKeyPrefix))
	chainStore := prefix.NewStore(store, types.KeyPrefix(fmt.Sprintf("%d-", req.Chain)))

	pageRes, err := query.Paginate(chainStore, req.Pagination, func(key []byte, value []byte) error {
		var outTxTracker types.OutTxTracker
		if err := k.cdc.Unmarshal(value, &outTxTracker); err != nil {
			return err
		}
		outTxTrackers = append(outTxTrackers, outTxTracker)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOutTxTrackerByChainResponse{OutTxTracker: outTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) OutTxTracker(c context.Context, req *types.QueryGetOutTxTrackerRequest) (*types.QueryGetOutTxTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	val, found := k.GetOutTxTracker(
		ctx,
		req.ChainID,
		req.Nonce,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetOutTxTrackerResponse{OutTxTracker: val}, nil
}
