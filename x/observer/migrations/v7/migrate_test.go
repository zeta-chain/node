package v7_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	v7 "github.com/zeta-chain/zetacore/x/observer/migrations/v7"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMigrateStore(t *testing.T) {
	t.Run("Migrate store from v6 to v7", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		addr1 := sample.AccAddress()
		addr2 := sample.AccAddress()

		k.SetParams(ctx, types.Params{
			AdminPolicy: []*types.Admin_Policy{
				{
					PolicyType: types.Policy_Type_group1,
					Address:    addr1,
				},
				{
					PolicyType: types.Policy_Type_group2,
					Address:    addr2,
				},
			},
		})

		// Migrate store
		err := v7.MigrateStore(ctx, k)

		// Check if store is migrated
		require.NoError(t, err)
	})
}

func TestMigratePolicies(t *testing.T) {
	t.Run("Migrate policies from observer to authority with 2 types", func(t *testing.T) {
		k, ctx, _, zk := keepertest.ObserverKeeper(t)

		addr1 := sample.AccAddress()
		addr2 := sample.AccAddress()

		k.SetParams(ctx, types.Params{
			AdminPolicy: []*types.Admin_Policy{
				{
					PolicyType: types.Policy_Type_group1,
					Address:    addr1,
				},
				{
					PolicyType: types.Policy_Type_group2,
					Address:    addr2,
				},
			},
		})

		// Migrate policies
		err := v7.MigratePolicies(ctx, k)

		// Check if policies are migrated
		require.NoError(t, err)
		policies, found := zk.AuthorityKeeper.GetPolicies(ctx)
		require.True(t, found)
		items := policies.Items
		require.Len(t, items, 2)
		require.EqualValues(t, addr1, items[0].Address)
		require.EqualValues(t, addr2, items[1].Address)
		require.EqualValues(t, authoritytypes.PolicyType_groupEmergency, items[0].PolicyType)
		require.EqualValues(t, authoritytypes.PolicyType_groupAdmin, items[1].PolicyType)
	})

	t.Run("Can migrate with just emergency policy", func(t *testing.T) {
		k, ctx, _, zk := keepertest.ObserverKeeper(t)

		addr := sample.AccAddress()

		k.SetParams(ctx, types.Params{
			AdminPolicy: []*types.Admin_Policy{
				{
					PolicyType: types.Policy_Type_group1,
					Address:    addr,
				},
			},
		})

		// Migrate policies
		err := v7.MigratePolicies(ctx, k)

		// Check if policies are migrated
		require.NoError(t, err)
		policies, found := zk.AuthorityKeeper.GetPolicies(ctx)
		require.True(t, found)
		items := policies.Items
		require.Len(t, items, 1)
		require.EqualValues(t, addr, items[0].Address)
		require.EqualValues(t, authoritytypes.PolicyType_groupEmergency, items[0].PolicyType)
	})

	t.Run("Can migrate with just admin  policy", func(t *testing.T) {
		k, ctx, _, zk := keepertest.ObserverKeeper(t)

		addr := sample.AccAddress()

		k.SetParams(ctx, types.Params{
			AdminPolicy: []*types.Admin_Policy{
				{
					PolicyType: types.Policy_Type_group2,
					Address:    addr,
				},
			},
		})

		// Migrate policies
		err := v7.MigratePolicies(ctx, k)

		// Check if policies are migrated
		require.NoError(t, err)
		policies, found := zk.AuthorityKeeper.GetPolicies(ctx)
		require.True(t, found)
		items := policies.Items
		require.Len(t, items, 1)
		require.EqualValues(t, addr, items[0].Address)
		require.EqualValues(t, authoritytypes.PolicyType_groupAdmin, items[0].PolicyType)
	})

	t.Run("Can migrate with no policies", func(t *testing.T) {
		k, ctx, _, zk := keepertest.ObserverKeeper(t)

		k.SetParams(ctx, types.Params{})

		// Migrate policies
		err := v7.MigratePolicies(ctx, k)

		// Check if policies are migrated
		require.NoError(t, err)
		policies, found := zk.AuthorityKeeper.GetPolicies(ctx)
		require.True(t, found)
		items := policies.Items
		require.Len(t, items, 0)
	})

	t.Run("Fail to migrate if invalid policy", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeper(t)

		k.SetParams(ctx, types.Params{
			AdminPolicy: []*types.Admin_Policy{
				{
					PolicyType: types.Policy_Type_group1,
					Address:    "invalid",
				},
			},
		})

		// Migrate policies
		err := v7.MigratePolicies(ctx, k)
		require.Error(t, err)
	})
}
