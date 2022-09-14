package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

func (k Keeper) SetSupportedChain(ctx sdk.Context, chain types.SupportedChains) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SupportedChainsKey))
	b := k.cdc.MustMarshal(&chain)
	store.Set([]byte(types.AllSupportedChainsKey), b)
}

func (k Keeper) GetSupportedChains(ctx sdk.Context) (val types.SupportedChains, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SupportedChainsKey))
	b := store.Get([]byte(types.AllSupportedChainsKey))
	if b != nil {
		k.cdc.MustUnmarshal(b, &val)
		return val, true
	}
	return val, false
}

func (k Keeper) IsChainSupported(ctx sdk.Context, checkChain string) bool {
	chains, found := k.GetSupportedChains(ctx)
	if !found {
		return false
	}
	for _, chain := range chains.ChainList {
		if checkChain == chain {
			return true
		}
	}
	return false
}

func (k msgServer) SetSupportedChains(goCtx context.Context, msg *types.MsgSetSupportedChains) (*types.MsgSetSupportedChainsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	supportedChains := types.SupportedChains{ChainList: msg.GetChainlist()}
	k.SetSupportedChain(ctx, supportedChains)
	return &types.MsgSetSupportedChainsResponse{}, nil
}

func (k Keeper) SupportedChains(goCtx context.Context, req *types.QuerySupportedChains) (*types.QuerySupportedChainsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chains, found := k.GetSupportedChains(ctx)
	if !found {
		return nil, types.ErrSupportedChains.Wrap("Supported chains not set")
	}
	return &types.QuerySupportedChainsResponse{Chains: chains.ChainList}, nil
}
