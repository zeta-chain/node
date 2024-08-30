package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v2 "github.com/zeta-chain/node/x/authority/migrations/v2"
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
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.authorityKeeper)
}
