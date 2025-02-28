package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v10 "github.com/zeta-chain/node/x/observer/migrations/v10"
	v11 "github.com/zeta-chain/node/x/observer/migrations/v11"
	v8 "github.com/zeta-chain/node/x/observer/migrations/v8"
	v9 "github.com/zeta-chain/node/x/observer/migrations/v9"
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
func (m Migrator) Migrate1to2(_ sdk.Context) error {
	return nil
}

// Migrate2to3 migrates the store from consensus version 2 to 3
func (m Migrator) Migrate2to3(_ sdk.Context) error {
	return nil
}

func (m Migrator) Migrate3to4(_ sdk.Context) error {
	return nil
}

func (m Migrator) Migrate4to5(_ sdk.Context) error {
	return nil
}

func (m Migrator) Migrate5to6(_ sdk.Context) error {
	return nil
}

// Migrate6to7 migrates the store from consensus version 6 to 7
func (m Migrator) Migrate6to7(_ sdk.Context) error {
	return nil
}

// Migrate7to8 migrates the store from consensus version 7 to 8
func (m Migrator) Migrate7to8(ctx sdk.Context) error {
	return v8.MigrateStore(ctx, m.observerKeeper)
}

// Migrate8to9 migrates the store from consensus version 8 to 9
func (m Migrator) Migrate8to9(ctx sdk.Context) error {
	return v9.MigrateStore(ctx, m.observerKeeper)
}

// Migrate9to10 migrates the store from consensus version 9 to 10
func (m Migrator) Migrate9to10(ctx sdk.Context) error {
	return v10.MigrateStore(ctx, m.observerKeeper)
}

// Migrate10to11 migrates the store from consensus version 10 to 11
func (m Migrator) Migrate10to11(ctx sdk.Context) error {
	return v11.MigrateStore(ctx, m.observerKeeper)
}
