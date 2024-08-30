package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestKeeper_Policies(t *testing.T) {
	t.Run("invalid request", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)

		_, err := k.Policies(ctx, nil)
		require.ErrorContains(t, err, "invalid request")
	})

	t.Run("policies not found", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)

		_, err := k.Policies(ctx, &types.QueryGetPoliciesRequest{})
		require.ErrorContains(t, err, "policies not found")
	})

	t.Run("can retrieve policies", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)

		policies := sample.Policies()
		k.SetPolicies(ctx, policies)

		res, err := k.Policies(ctx, &types.QueryGetPoliciesRequest{})
		require.NoError(t, err)
		require.Equal(t, policies, res.Policies)
	})
}
