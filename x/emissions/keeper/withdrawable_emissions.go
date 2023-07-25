package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/emissions/types"
)

func (k Keeper) SetWithdrawableEmission(ctx sdk.Context, we types.WithdrawableEmissions) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WithdrawableEmissionsKey))
	b := k.cdc.MustMarshal(&we)
	store.Set([]byte(we.Address), b)
}

func (k Keeper) GetWithdrawableEmission(ctx sdk.Context, address string) (val types.WithdrawableEmissions, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.WithdrawableEmissionsKey))
	b := store.Get(types.KeyPrefix(address))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
