package keeper

import (
	"context"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
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

func (k Keeper) GetChainFromChainID(ctx sdk.Context, chainId int64) (*types.Chain, bool) {
	chains, found := k.GetSupportedChains(ctx)
	if !found {
		return nil, false
	}
	for _, chain := range chains.ChainList {
		if chain.ChainId == chainId {
			return chain, true
		}
	}
	return nil, false
}

func (k Keeper) GetChainFromChainName(ctx sdk.Context, name types.ChainName) (*types.Chain, bool) {
	chains, found := k.GetSupportedChains(ctx)
	if !found {
		return nil, false
	}
	for _, chain := range chains.ChainList {
		if chain.ChainName == name {
			return chain, true
		}
	}
	return nil, false
}

func (k Keeper) IsChainSupported(ctx sdk.Context, checkChain types.Chain) bool {
	chains, found := k.GetSupportedChains(ctx)
	if !found {
		return false
	}
	for _, chain := range chains.ChainList {
		if checkChain.IsEqual(chain) {
			return true
		}
	}
	return false
}

func (k Keeper) SetSupportedChains(goCtx context.Context, msg *types.MsgSetSupportedChains) (*types.MsgSetSupportedChainsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chain := []*types.Chain{{
		ChainName: msg.ChainName,
		ChainId:   msg.ChainId,
	},
	}
	chains, found := k.GetSupportedChains(ctx)
	if !found {
		supportedChains := types.SupportedChains{ChainList: chain}
		k.SetSupportedChain(ctx, supportedChains)
		return &types.MsgSetSupportedChainsResponse{}, nil
	}
	chains.ChainList = append(chains.ChainList, chain...)
	k.SetSupportedChain(ctx, chains)
	return &types.MsgSetSupportedChainsResponse{}, nil
}

func (k Keeper) SupportedChains(goCtx context.Context, _ *types.QuerySupportedChains) (*types.QuerySupportedChainsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	chains, found := k.GetSupportedChains(ctx)
	if !found {
		return nil, types.ErrSupportedChains.Wrap("Supported chains not set")
	}
	return &types.QuerySupportedChainsResponse{Chains: chains.ChainList}, nil
}
