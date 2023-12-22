package v3

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

// MigrateStore adds a new TSSHistory store to crosschain module
func MigrateStore(
	ctx sdk.Context,
	crossChainStoreKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
) error {

	// Fetch existing TSS
	existingTss := observertypes.TSS{}
	store := prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(observertypes.TSSKey))
	b := store.Get([]byte{0})
	if b == nil {
		return errors.New("TSS not found")
	}

	// Add existing TSS to TSSHistory store
	cdc.MustUnmarshal(b, &existingTss)
	store = prefix.NewStore(ctx.KVStore(crossChainStoreKey), types.KeyPrefix(observertypes.TSSHistoryKey))
	b = cdc.MustMarshal(&existingTss)
	store.Set(types.KeyPrefix(fmt.Sprintf("%d", existingTss.FinalizedZetaHeight)), b)

	return nil
}
