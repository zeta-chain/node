package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zeta-chain/node/x/fungible/types"
)

// CodeHash returns the code hash of an account if it exists
func (k Keeper) CodeHash(c context.Context, req *types.QueryCodeHashRequest) (*types.QueryCodeHashResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// convert address to hex
	if !ethcommon.IsHexAddress(req.Address) {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}
	address := ethcommon.HexToAddress(req.Address)

	// fetch account
	ctx := sdk.UnwrapSDKContext(c)
	acc := k.evmKeeper.GetAccount(ctx, address)
	if acc == nil {
		return nil, status.Error(codes.NotFound, "account not found")
	}
	if !k.evmKeeper.IsContract(ctx, address) {
		return nil, status.Error(codes.NotFound, "account is not a contract")
	}

	// convert code hash to hex
	codeHash := ethcommon.BytesToHash(acc.CodeHash)

	return &types.QueryCodeHashResponse{CodeHash: codeHash.Hex()}, nil
}
