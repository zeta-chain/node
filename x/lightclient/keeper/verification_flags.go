package keeper

import (
	"fmt"

	cosmoserrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// SetVerificationFlags set the verification flags in the store
func (k Keeper) SetVerificationFlags(ctx sdk.Context, vf types.VerificationFlags) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VerificationFlagsKey))
	b := k.cdc.MustMarshal(&vf)
	key := types.KeyPrefix(fmt.Sprintf("%d", vf.ChainId))
	store.Set(key, b)
}

// GetVerificationFlags returns the verification flags
func (k Keeper) GetVerificationFlags(ctx sdk.Context, chainID int64) (val types.VerificationFlags, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VerificationFlagsKey))
	key := types.KeyPrefix(fmt.Sprintf("%d", chainID))
	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetAllVerificationFlags(ctx sdk.Context) (verificationFlags []types.VerificationFlags) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VerificationFlagsKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.VerificationFlags
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		verificationFlags = append(verificationFlags, val)
	}

	return verificationFlags

}

// CheckVerificationFlagsEnabled checks for a specific chain if the verification flags are enabled
// It returns an error if the chain is not enabled
func (k Keeper) CheckVerificationFlagsEnabled(ctx sdk.Context, chainID int64) error {
	verificationFlags, found := k.GetVerificationFlags(ctx, chainID)
	if !found || verificationFlags.Enabled != true {
		return cosmoserrors.Wrapf(
			types.ErrBlockHeaderVerificationDisabled,
			"proof verification not enabled for,chain id: %d",
			chainID,
		)
	}
	return nil
}
