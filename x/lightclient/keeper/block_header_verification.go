package keeper

import (
	cosmoserrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// SetVerificationFlags set the verification flags in the store. The key is the chain id
func (k Keeper) SetBlockHeaderVerification(ctx sdk.Context, bhv types.BlockHeaderVerification) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VerificationFlagsKey))
	b := k.cdc.MustMarshal(&bhv)
	store.Set([]byte{0}, b)
}

// GetBlockHeaderVerification returns the verification flags
func (k Keeper) GetBlockHeaderVerification(ctx sdk.Context) (bhv types.BlockHeaderVerification, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VerificationFlagsKey))
	b := store.Get([]byte{0})
	if b == nil {
		return bhv, false
	}

	k.cdc.MustUnmarshal(b, &bhv)
	return bhv, true
}

// CheckBlockHeaderVerificationEnabled checks for a specific chain if the verification flags are enabled
// It returns an error if the chain is not enabled or the verification flags are not for that chain
func (k Keeper) CheckBlockHeaderVerificationEnabled(ctx sdk.Context, chainID int64) error {
	bhv, found := k.GetBlockHeaderVerification(ctx)
	if !found {
		return cosmoserrors.Wrapf(
			types.ErrBlockHeaderVerificationDisabled,
			"proof verification is disabled for all chains",
		)
	}
	if !bhv.IsChainEnabled(chainID) {
		return cosmoserrors.Wrapf(
			types.ErrBlockHeaderVerificationDisabled,
			"proof verification is disabled for chain %d",
			chainID,
		)
	}
	return nil
}
