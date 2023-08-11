package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v2 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	crossChainKeeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		crossChainKeeper: keeper,
	}
}

// Migrate2to3 migrates the store from consensus version 2 to 3
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.crossChainKeeper.ZetaObserverKeeper, m.crossChainKeeper.storeKey, m.crossChainKeeper.cdc)
}
