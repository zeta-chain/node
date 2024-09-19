package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/x/emissions/exported"
	v3 "github.com/zeta-chain/node/x/emissions/migrations/v3"
	v4 "github.com/zeta-chain/node/x/emissions/migrations/v4"
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

// Migrate2to3 migrates the store from consensus version 2 to 3
func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	return v3.MigrateStore(ctx, m.keeper, m.legacySubspace)
}

// Migrate3to4 migrates the store from consensus version 3 to 4
func (m Migrator) Migrate3to4(ctx sdk.Context) error {
	return v4.MigrateStore(ctx, m.keeper)
}
