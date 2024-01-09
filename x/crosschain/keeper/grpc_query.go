package keeper

import (
	"context"

	"github.com/zeta-chain/zetacore/x/crosschain/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

// GetTssAddress returns the tss address
// Deprecated: GetTssAddress returns the tss address
// TODO: remove after v12 once upgrade testing is no longer needed with v11
// https://github.com/zeta-chain/node/issues/1547
func (k Keeper) GetTssAddress(_ context.Context, _ *types.QueryGetTssAddressRequest) (*types.QueryGetTssAddressResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Deprecated")
}
