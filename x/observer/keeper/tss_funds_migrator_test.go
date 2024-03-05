package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

func TestKeeper_GetTssFundMigrator(t *testing.T) {
	t.Run("Successfully set funds migrator for chain", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)
		chain := sample.TssFundsMigrator(1)
		k.SetFundMigrator(ctx, chain)
		tfm, found := k.GetFundMigrator(ctx, chain.ChainId)
		require.True(t, found)
		require.Equal(t, chain, tfm)
	})
	t.Run("Verify only one migrator can be created for a chain", func(t *testing.T) {
		k, ctx, _ := keepertest.ObserverKeeper(t)
		tfm1 := sample.TssFundsMigrator(1)
		k.SetFundMigrator(ctx, tfm1)
		tfm2 := tfm1
		tfm2.MigrationCctxIndex = "sampleIndex2"
		k.SetFundMigrator(ctx, tfm2)
		migratorList := k.GetAllTssFundMigrators(ctx)
		require.Equal(t, 1, len(migratorList))
		require.Equal(t, tfm2, migratorList[0])
	})

}
