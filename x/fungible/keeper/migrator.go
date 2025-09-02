package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v4 "github.com/zeta-chain/node/x/fungible/migrations/v4"
	v5 "github.com/zeta-chain/node/x/fungible/migrations/v5"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	fungibleKeeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		fungibleKeeper: keeper,
	}
}

// Migrate3to4 migrates the store from consensus version 3 to 4
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.MigrateStore(ctx, &m.fungibleKeeper)
}

// Migrate4to5 migrates the store from consensus version 4 to 5
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v5.MigrateStore(ctx, &m.fungibleKeeper)
}
