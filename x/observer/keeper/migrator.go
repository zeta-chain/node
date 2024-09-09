package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v8 "github.com/zeta-chain/node/x/observer/migrations/v8"
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
