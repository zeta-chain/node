package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v5 "github.com/zeta-chain/node/x/crosschain/migrations/v5"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	crossChainKeeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{
		crossChainKeeper: keeper,
	}
}

// Migrate4to5 migrates the store from consensus version 4 to 5
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v5.MigrateStore(ctx, m.crossChainKeeper, m.crossChainKeeper.zetaObserverKeeper)
}
