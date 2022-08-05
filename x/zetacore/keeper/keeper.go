package keeper

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/zeta-chain/zetacore/x/zetacore/types"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc      codec.Codec
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey

		StakingKeeper types.StakingKeeper
		paramstore    paramtypes.Subspace

		// this line is used by starport scaffolding # ibc/keeper/attribute

	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey sdk.StoreKey,
	stakingKeeper types.StakingKeeper, // custom
	paramstore paramtypes.Subspace,

	// this line is used by starport scaffolding # ibc/keeper/parameter

) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		StakingKeeper: stakingKeeper,
		paramstore:    paramstore,

		// this line is used by starport scaffolding # ibc/keeper/return

	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
