package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestAppendFVsSprintf(t *testing.T) {
	// Test cases for tss_funds_migrator.go
	testCases := []struct {
		name   string
		values []interface{}
		format string
	}{
		{
			name:   "chain id format",
			values: []interface{}{int64(123)},
			format: "%d",
		},
		{
			name:   "large chain id",
			values: []interface{}{int64(9223372036854775807)}, // Max int64
			format: "%d",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Generate string using fmt.Sprintf
			sprintfResult := []byte(fmt.Sprintf(tc.format, tc.values...))
			
			// Generate string using fmt.Appendf
			appendResult := fmt.Appendf(nil, tc.format, tc.values...)
			
			// Assert that both methods produce the same byte slice
			assert.Equal(t, sprintfResult, appendResult, "[]byte(fmt.Sprintf) and fmt.Appendf should produce identical byte slices")
		})
	}
}

func TestKeeper_GetTssFundMigrator(t *testing.T) {
	t.Run("Successfully set funds migrator for chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
		chain := sample.TssFundsMigrator(1)
		_, found := k.GetFundMigrator(ctx, chain.ChainId)
		require.False(t, found)
		k.SetFundMigrator(ctx, chain)
		tfm, found := k.GetFundMigrator(ctx, chain.ChainId)
		require.True(t, found)
		require.Equal(t, chain, tfm)

		k.RemoveAllExistingMigrators(ctx)
		_, found = k.GetFundMigrator(ctx, chain.ChainId)
		require.False(t, found)
	})
	t.Run("Verify only one migrator can be created for a chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)
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
