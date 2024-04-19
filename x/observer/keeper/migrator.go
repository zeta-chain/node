package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	return nil
}

// Migrate2to3 migrates the store from consensus version 2 to 3
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return nil
}

func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return nil
}

func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return nil
}

func (m Migrator) Migrate5to6(ctx sdk.Context) error {
	return nil
}

// Migrate6to7 migrates the store from consensus version 6 to 7
func (m Migrator) Migrate6to7(ctx sdk.Context) error {
	return nil
}
