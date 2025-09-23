package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v11 "github.com/zeta-chain/node/x/observer/migrations/v11"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	observerKeeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		observerKeeper: keeper,
	}
}

// Migrate10to11 migrates the store from consensus version 10 to 11
func (m Migrator) Migrate10to11(ctx sdk.Context) error {
	return v11.MigrateStore(ctx, m.observerKeeper)
}
