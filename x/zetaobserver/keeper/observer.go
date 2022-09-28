package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
)

func (k Keeper) SetObserver(ctx sdk.Context, observer types.Observer) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverKey))
	observer.Index = fmt.Sprintf("%s-%s", observer.ObserverChain.String(), observer.ObserverAddress)
	b := k.cdc.MustMarshal(&observer)
	store.Set([]byte(observer.Index), b)
}

func (k Keeper) GetObserver(ctx sdk.Context, index string) (val types.Observer, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ObserverKey))
	b := store.Get(types.KeyPrefix(index))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
