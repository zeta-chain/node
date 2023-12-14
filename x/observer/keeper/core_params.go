package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func (k Keeper) SetCoreParamsList(ctx sdk.Context, coreParams types.CoreParamsList) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&coreParams)
	key := types.KeyPrefix(fmt.Sprintf("%s", types.AllCoreParams))
	store.Set(key, b)
}

func (k Keeper) GetCoreParamsList(ctx sdk.Context) (val types.CoreParamsList, found bool) {
	found = false
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.KeyPrefix(fmt.Sprintf("%s", types.AllCoreParams)))
	if b == nil {
		return
	}
	found = true
	k.cdc.MustUnmarshal(b, &val)
	return
}

func (k Keeper) GetCoreParamsByChainID(ctx sdk.Context, chainID int64) (*types.CoreParams, bool) {
	allCoreParams, found := k.GetCoreParamsList(ctx)
	if !found {
		return &types.CoreParams{}, false
	}
	for _, coreParams := range allCoreParams.CoreParams {
		if coreParams.ChainId == chainID {
			return coreParams, true
		}
	}
	return &types.CoreParams{}, false
}

// GetSupportedChainFromChainID returns the chain from the chain id
// it returns nil if the chain doesn't exist or is not supported
// TODO: test this function
func (k Keeper) GetSupportedChainFromChainID(ctx sdk.Context, chainID int64) *common.Chain {
	cpl, found := k.GetCoreParamsList(ctx)
	if !found {
		return nil
	}

	for _, cp := range cpl.CoreParams {
		if cp.ChainId == chainID && cp.IsSupported {
			return common.GetChainFromChainID(chainID)
		}
	}
	return nil
}
