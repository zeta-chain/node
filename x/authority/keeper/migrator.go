package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v3 "github.com/zeta-chain/node/x/authority/migrations/v3"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	authorityKeeper Keeper
}

// NewMigrator returns a new Migrator for the authority module.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		authorityKeeper: keeper,
	}
}

// Migrate1to2 migrates the authority store from consensus version 1 to 2
func (m Migrator) Migrate1to2(_ sdk.Context) error {
	return nil
}

// Migrate2to3 migrates the authority store from consensus version 2 to 3
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.authorityKeeper)
}
