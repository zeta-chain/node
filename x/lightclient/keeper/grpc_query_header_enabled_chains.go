package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/x/lightclient/types"
)

// HeaderSupportedChains implements the Query/HeaderEnabledChains gRPC method
// It returns a list for chains that support block header verification.
// Some chains in this list might be disabled which is indicated by the value of the `enabled` field.
func (k Keeper) HeaderSupportedChains(
	c context.Context,
	req *types.QueryHeaderSupportedChainsRequest,
) (*types.QueryHeaderSupportedChainsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetBlockHeaderVerification(ctx)
	if !found {
		return &types.QueryHeaderSupportedChainsResponse{}, types.ErrBlockHeaderVerificationDisabled.Wrapf(
			"proof verification is disabled for all chains",
		)
	}

	return &types.QueryHeaderSupportedChainsResponse{HeaderSupportedChains: val.GetHeaderSupportedChainsList()}, nil
}

// HeaderEnabledChains implements the Query/HeaderEnabledChains gRPC method
// It returns a list of chains that have block header verification enabled.
func (k Keeper) HeaderEnabledChains(
	c context.Context,
	req *types.QueryHeaderEnabledChainsRequest,
) (*types.QueryHeaderEnabledChainsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetBlockHeaderVerification(ctx)
	if !found {
		return &types.QueryHeaderEnabledChainsResponse{}, types.ErrBlockHeaderVerificationDisabled.Wrapf(
			"proof verification is disabled for all chains",
		)
	}

	return &types.QueryHeaderEnabledChainsResponse{HeaderEnabledChains: val.GetHeaderEnabledChains()}, nil
}
