package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/x/observer/types"
)

// TssFundsMigratorInfo queries a tss fund migrator info
func (k Keeper) TssFundsMigratorInfo(
	goCtx context.Context,
	req *types.QueryTssFundsMigratorInfoRequest,
) (*types.QueryTssFundsMigratorInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	_, found := chains.GetChainFromChainID(req.ChainId, k.GetAuthorityKeeper().GetAdditionalChainList(ctx))
	if !found {
		return nil, status.Error(codes.InvalidArgument, "invalid chain id")
	}

	fm, found := k.GetFundMigrator(ctx, req.ChainId)
	if !found {
		return nil, status.Error(codes.NotFound, "tss fund migrator not found")
	}
	return &types.QueryTssFundsMigratorInfoResponse{
		TssFundsMigrator: fm,
	}, nil
}

// TssFundsMigratorInfoAll queries all tss fund migrator info for all chains
func (k Keeper) TssFundsMigratorInfoAll(
	goCtx context.Context,
	request *types.QueryTssFundsMigratorInfoAllRequest,
) (*types.QueryTssFundsMigratorInfoAllResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	migrators := k.GetAllTssFundMigrators(ctx)

	return &types.QueryTssFundsMigratorInfoAllResponse{TssFundsMigrators: migrators}, nil
}
