package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/exported"
	v5 "github.com/zeta-chain/node/x/emissions/migrations/v5"
	v6 "github.com/zeta-chain/node/x/emissions/migrations/v6"
	v7 "github.com/zeta-chain/node/x/emissions/migrations/v7"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(k Keeper, ss exported.Subspace) Migrator {
	return Migrator{
		keeper:         k,
		legacySubspace: ss,
	}
}

// Migrate4to5 migrates the store from consensus version 4 to 5
func (m Migrator) Migrate4to5(ctx sdk.Context) error {
	return v5.MigrateStore(ctx, m.keeper)
}

// Migrate5to6 migrates the store from consensus version 5 to 6
func (m Migrator) Migrate5to6(ctx sdk.Context) error {
	return v6.MigrateStore(ctx, m.keeper)
}

// Migrate6to7 migrates the store from consensus version 6 to 7
func (m Migrator) Migrate6to7(ctx sdk.Context) error {
	return v7.MigrateStore(ctx, m.keeper)
}
