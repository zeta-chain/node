package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v2 "github.com/zeta-chain/zetacore/x/observer/migrations/v2"
	v3 "github.com/zeta-chain/zetacore/x/observer/migrations/v3"
	v4 "github.com/zeta-chain/zetacore/x/observer/migrations/v4"
	v5 "github.com/zeta-chain/zetacore/x/observer/migrations/v5"
	v6 "github.com/zeta-chain/zetacore/x/observer/migrations/v6"
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

// Migrate1to2 migrates the store from consensus version 1 to 2
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return v2.MigrateStore(ctx, m.observerKeeper.storeKey, m.observerKeeper.cdc)
}

// Migrate2to3 migrates the store from consensus version 2 to 3
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.observerKeeper)
}

func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.MigrateStore(ctx, m.observerKeeper)
}

func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v5.MigrateStore(ctx, m.observerKeeper)
}

func (m Migrator) Migrate5to6(ctx sdk.Context) error {
	return v6.MigrateStore(ctx, m.observerKeeper)
}
