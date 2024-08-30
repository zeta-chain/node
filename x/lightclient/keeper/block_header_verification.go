package keeper

import (
	cosmoserrors "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/lightclient/types"
)

// SetBlockHeaderVerification sets BlockHeaderVerification settings for all chains
func (k Keeper) SetBlockHeaderVerification(ctx sdk.Context, bhv types.BlockHeaderVerification) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VerificationFlagsKey))
	b := k.cdc.MustMarshal(&bhv)
	store.Set([]byte{0}, b)
}

// GetBlockHeaderVerification returns the BlockHeaderVerification settings for all chains
func (k Keeper) GetBlockHeaderVerification(ctx sdk.Context) (bhv types.BlockHeaderVerification, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VerificationFlagsKey))
	b := store.Get([]byte{0})
	if b == nil {
		return bhv, false
	}

	k.cdc.MustUnmarshal(b, &bhv)
	return bhv, true
}

// CheckBlockHeaderVerificationEnabled checks for a specific chain if BlockHeaderVerification is enabled or not
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
