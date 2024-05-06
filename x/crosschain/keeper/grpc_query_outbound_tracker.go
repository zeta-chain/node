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

func (k Keeper) OutboundTrackerAll(c context.Context, req *types.QueryAllOutboundTrackerRequest) (*types.QueryAllOutboundTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var outTxTrackers []types.OutboundTracker
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	outTxTrackerStore := prefix.NewStore(store, types.KeyPrefix(types.OutboundTrackerKeyPrefix))
	pageRes, err := query.Paginate(outTxTrackerStore, req.Pagination, func(_ []byte, value []byte) error {
		var outTxTracker types.OutboundTracker
		if err := k.cdc.Unmarshal(value, &outTxTracker); err != nil {
			return err
		}

		outTxTrackers = append(outTxTrackers, outTxTracker)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryAllOutboundTrackerResponse{OutboundTracker: outTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) OutboundTrackerAllByChain(c context.Context, req *types.QueryAllOutboundTrackerByChainRequest) (*types.QueryAllOutboundTrackerByChainResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var outTxTrackers []types.OutboundTracker
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutboundTrackerKeyPrefix))
	chainStore := prefix.NewStore(store, types.KeyPrefix(fmt.Sprintf("%d-", req.Chain)))

	pageRes, err := query.Paginate(chainStore, req.Pagination, func(_ []byte, value []byte) error {
		var outTxTracker types.OutboundTracker
		if err := k.cdc.Unmarshal(value, &outTxTracker); err != nil {
			return err
		}
		outTxTrackers = append(outTxTrackers, outTxTracker)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOutboundTrackerByChainResponse{OutboundTracker: outTxTrackers, Pagination: pageRes}, nil
}

func (k Keeper) OutboundTracker(c context.Context, req *types.QueryGetOutboundTrackerRequest) (*types.QueryGetOutboundTrackerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	val, found := k.GetOutboundTracker(
		ctx,
		req.ChainID,
		req.Nonce,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetOutboundTrackerResponse{OutboundTracker: val}, nil
}
