package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v2 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v2"
	v3 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v3"
	v4 "github.com/zeta-chain/zetacore/x/crosschain/migrations/v4"
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

// Migrate1to2 migrates the store from consensus version 1 to 2
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.crossChainKeeper.zetaObserverKeeper, m.crossChainKeeper.storeKey, m.crossChainKeeper.cdc)
}

// Migrate2to3 migrates the store from consensus version 2 to 3
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.crossChainKeeper.storeKey, m.crossChainKeeper.cdc)
}

// Migrate3to4 migrates the store from consensus version 3 to 4
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.MigrateStore(ctx, m.crossChainKeeper.zetaObserverKeeper, m.crossChainKeeper)
}
