package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) SetChainParamsList(ctx sdk.Context, chainParams types.ChainParamsList) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&chainParams)
	key := types.KeyPrefix(fmt.Sprintf("%s", types.AllChainParamsKey))
	store.Set(key, b)
}

func (k Keeper) GetChainParamsList(ctx sdk.Context) (val types.ChainParamsList, found bool) {
	found = false
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.KeyPrefix(fmt.Sprintf("%s", types.AllChainParamsKey)))
	if b == nil {
		return
	}
	found = true
	k.cdc.MustUnmarshal(b, &val)
	return
}

func (k Keeper) GetChainParamsByChainID(ctx sdk.Context, chainID int64) (*types.ChainParams, bool) {
	allChainParams, found := k.GetChainParamsList(ctx)
	if !found {
		return &types.ChainParams{}, false
	}
	for _, chainParams := range allChainParams.ChainParams {
		if chainParams.ChainId == chainID {
			return chainParams, true
		}
	}
	return &types.ChainParams{}, false
}

// GetSupportedChainFromChainID returns the chain from the chain id
// it returns nil if the chain doesn't exist or is not supported
func (k Keeper) GetSupportedChainFromChainID(ctx sdk.Context, chainID int64) *chains.Chain {
	cpl, found := k.GetChainParamsList(ctx)
	if !found {
		return nil
	}

	for _, cp := range cpl.ChainParams {
		if cp.ChainId == chainID && cp.IsSupported {
			return chains.GetChainFromChainID(chainID)
		}
	}
	return nil
}

// GetSupportedChains returns the list of supported chains
func (k Keeper) GetSupportedChains(ctx sdk.Context) []*chains.Chain {
	cpl, found := k.GetChainParamsList(ctx)
	if !found {
		return []*chains.Chain{}
	}

	var c []*chains.Chain
	for _, cp := range cpl.ChainParams {
		if cp.IsSupported {
			c = append(c, chains.GetChainFromChainID(cp.ChainId))
		}
	}
	return c
}

// GetSupportedForeignChains returns the list of supported foreign chains
func (k Keeper) GetSupportedForeignChains(ctx sdk.Context) []*chains.Chain {
	allChains := k.GetSupportedChains(ctx)

	foreignChains := make([]*chains.Chain, 0)
	for _, chain := range allChains {
		if !chain.IsZetaChain() {
			foreignChains = append(foreignChains, chain)
		}
	}
	return foreignChains
}
