package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc      codec.Codec
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
}

// NewKeeper creates new instances of the authority Keeper
func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetStoreKey returns the key to the store for authority
func (k Keeper) GetStoreKey() storetypes.StoreKey {
	return k.storeKey
}

// GetMemKey returns the mem key to the store for authority
func (k Keeper) GetMemKey() storetypes.StoreKey {
	return k.memKey
}

// GetCodec returns the codec for authority
func (k Keeper) GetCodec() codec.Codec {
	return k.cdc
}
