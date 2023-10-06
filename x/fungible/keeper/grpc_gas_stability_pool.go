package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GasStabilityPoolAddress(
	_ context.Context,
	req *types.QueryGetGasStabilityPoolAddress,
) (*types.QueryGetGasStabilityPoolAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	return &types.QueryGetGasStabilityPoolAddressResponse{
		CosmosAddress: types.GasStabilityPoolAddress().String(),
		EvmAddress:    types.GasStabilityPoolAddressEVM().String(),
	}, nil
}

func (k Keeper) GasStabilityPoolBalance(
	c context.Context,
	req *types.QueryGetGasStabilityPoolBalance,
) (*types.QueryGetGasStabilityPoolBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	balance, err := k.GetGasStabilityPoolBalance(ctx, req.ChainId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if balance == nil {
		return nil, status.Error(codes.NotFound, "no balance for the gas stability pool")
	}

	return &types.QueryGetGasStabilityPoolBalanceResponse{Balance: balance.String()}, nil
}

func (k Keeper) GasStabilityPoolBalanceAll(
	c context.Context,
	req *types.QueryAllGasStabilityPoolBalance,
) (*types.QueryAllGasStabilityPoolBalanceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	// iterate supported chains
	chains := k.observerKeeper.GetParams(ctx).GetSupportedChains()
	balances := make([]types.QueryAllGasStabilityPoolBalanceResponse_Balance, 0, len(chains))
	for _, chain := range chains {
		if chain == nil {
			return nil, status.Error(codes.Internal, "invalid chain")
		}
		chainID := chain.ChainId

		balance, err := k.GetGasStabilityPoolBalance(ctx, chainID)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if balance == nil {
			return nil, status.Error(codes.NotFound, "no balance for the gas stability pool")
		}

		balances = append(balances, types.QueryAllGasStabilityPoolBalanceResponse_Balance{
			ChainId: chainID,
			Balance: balance.String(),
		})
	}

	return &types.QueryAllGasStabilityPoolBalanceResponse{
		Balances: balances,
	}, nil
}
