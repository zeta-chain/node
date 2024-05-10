package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
)

func TestKeeper_SetPolicies(t *testing.T) {
	k, ctx := keepertest.AuthorityKeeper(t)
	policies := sample.Policies()

	_, found := k.GetPolicies(ctx)
	require.False(t, found)

	k.SetPolicies(ctx, policies)

	// Check policy is set
	got, found := k.GetPolicies(ctx)
	require.True(t, found)
	require.Equal(t, policies, got)

	// Can set policies again
	newPolicies := sample.Policies()
	require.NotEqual(t, policies, newPolicies)
	k.SetPolicies(ctx, newPolicies)
	got, found = k.GetPolicies(ctx)
	require.True(t, found)
	require.Equal(t, newPolicies, got)
}

func TestKeeper_IsAuthorized(t *testing.T) {
	t.Run("successfully authorized", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := sample.AccAddress()
		policies := sample.PoliciesWithAdmin(admin)
		k.SetPolicies(ctx, policies)
		msg := sample.AdminMessage(admin)
		ok, err := k.IsAuthorized(ctx, msg)
		require.True(t, ok)
		require.NoError(t, err)
	})

	t.Run("returns error if more than 1 signer", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msg := sample.MultipleSignerMessage()
		ok, err := k.IsAuthorized(ctx, msg)
		require.False(t, ok)
		require.ErrorIs(t, err, authoritytypes.ErrSigners)
	})

	t.Run("returns error if not admin message", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := sample.AccAddress()
		policies := sample.PoliciesWithAdmin(admin)
		k.SetPolicies(ctx, policies)
		msg := sample.NonAdminMessage(admin)
		ok, err := k.IsAuthorized(ctx, msg)
		require.False(t, ok)
		require.ErrorIs(t, err, authoritytypes.ErrMsgNotAuthorized)
	})

	t.Run("returns error if policies not found", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msg := sample.AdminMessage(sample.AccAddress())
		ok, err := k.IsAuthorized(ctx, msg)
		require.False(t, ok)
		require.ErrorIs(t, err, authoritytypes.ErrPoliciesNotFound)
	})

	t.Run("returns error if signer doesn't match", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		admin := sample.AccAddress()
		policies := sample.PoliciesWithAdmin(admin)
		k.SetPolicies(ctx, policies)
		msg := sample.AdminMessage(sample.AccAddress())
		ok, err := k.IsAuthorized(ctx, msg)
		require.False(t, ok)
		require.ErrorIs(t, err, authoritytypes.ErrSignerDoesntMatch)
	})
}
