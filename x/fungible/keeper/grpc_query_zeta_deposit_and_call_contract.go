package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ZetaDepositAndCallContract(c context.Context, req *types.QueryGetZetaDepositAndCallContractRequest) (*types.QueryGetZetaDepositAndCallContractResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetZetaDepositAndCallContract(ctx)
	if !found {
	    return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetZetaDepositAndCallContractResponse{ZetaDepositAndCallContract: val}, nil
}