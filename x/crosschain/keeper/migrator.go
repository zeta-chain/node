package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v5 "github.com/zeta-chain/node/x/crosschain/migrations/v5"
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
func (m Migrator) Migrate1to2(_ sdk.Context) error {
	return nil
}

// Migrate2to3 migrates the store from consensus version 2 to 3
func (m Migrator) Migrate2to3(_ sdk.Context) error {
	return nil
}

// Migrate3to4 migrates the store from consensus version 3 to 4
func (m Migrator) Migrate3to4(_ sdk.Context) error {
	return nil
}

// Migrate4to5 migrates the store from consensus version 4 to 5
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v5.MigrateStore(ctx, m.crossChainKeeper, m.crossChainKeeper.zetaObserverKeeper)
}
