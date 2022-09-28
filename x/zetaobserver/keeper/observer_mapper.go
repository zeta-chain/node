package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
)

func (k Keeper) SetObserverMapper(ctx sdk.Context, om types.ObserverMapper) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	om.Index = fmt.Sprintf("%s-%s", om.ObserverChain.String(), om.ObservationType.String())
	b := k.cdc.MustMarshal(&om)
	store.Set([]byte(om.Index), b)
}

func (k Keeper) GetObserverMapper(ctx sdk.Context, index string) (val types.ObserverMapper, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverMapperKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
