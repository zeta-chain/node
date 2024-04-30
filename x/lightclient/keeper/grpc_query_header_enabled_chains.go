package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EnabledChain implements the Query/EnabledChain gRPC method
func (k Keeper) HeaderEnabledChains(c context.Context, req *types.QueryHeaderEnabledChainsRequest) (*types.QueryHeaderEnabledChainsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, _ := k.GetBlockHeaderVerification(ctx)

	return &types.QueryHeaderEnabledChainsResponse{EnabledChains: val.GetVerificationFlags()}, nil
}
