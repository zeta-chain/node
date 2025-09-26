package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v12 "github.com/zeta-chain/node/x/observer/migrations/v12"
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

// Migrate11to12 migrates the store from consensus version 11 to 12
func (m Migrator) Migrate11to12(ctx sdk.Context) error {
	return v12.MigrateStore(ctx, m.observerKeeper)
}
