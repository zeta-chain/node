package v3

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

// MigrateStore migrates the x/crosschain module state from the consensus version 1 to 2
// This migration moves some data from the cross chain store to the observer store.
// The data moved is the node accounts, permission flags and keygen.
func MigrateStore(
	ctx sdk.Context,
	crossChainStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) error {

	// Fetch existing TSS
	existingTss := types.TSS{}
	store := prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(types.TSSKey))
	b := store.Get([]byte{0})
	if b == nil {
		return errors.New("TSS not found")
	}

	// Add existing TSS to TSSHistory store
	cdc.MustUnmarshal(b, &existingTss)
	store = prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(types.TSSHistoryKey))
	b = cdc.MustMarshal(&existingTss)
	store.Set(types.KeyPrefix(fmt.Sprintf("%d", existingTss.FinalizedZetaHeight)), b)

	return nil
}
