package keeper

import (
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
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
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
