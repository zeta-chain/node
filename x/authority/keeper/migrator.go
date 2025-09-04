package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v4 "github.com/zeta-chain/node/x/authority/migrations/v4"
	v5 "github.com/zeta-chain/node/x/authority/migrations/v5"
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

// Migrate3to4 migrates the authority store from consensus version 3 to 4
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.MigrateStore(ctx, m.authorityKeeper)
}

// Migrate4to5 migrates the authority store from consensus version 4 to 5
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v5.MigrateStore(ctx, m.authorityKeeper)
}
