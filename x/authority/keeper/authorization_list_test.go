package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	lightclienttypes "github.com/zeta-chain/zetacore/x/lightclient/types"

	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func TestKeeper_GetAuthorizationList(t *testing.T) {
	t.Run("successfully get authorizations list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		authorizationList := sample.AuthorizationList("sample")
		k.SetAuthorizationList(ctx, authorizationList)
		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, list)
	})

	t.Run("get authorizations list not found", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		list, found := k.GetAuthorizationList(ctx)
		require.False(t, found)
		require.Equal(t, types.AuthorizationList{}, list)
	})
}

func TestKeeper_SetAuthorizationList(t *testing.T) {
	t.Run("successfully set authorizations list when a list already exists", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		authorizationList := sample.AuthorizationList("sample")
		k.SetAuthorizationList(ctx, authorizationList)

		list, found := k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, authorizationList, list)

		newAuthorizationList := sample.AuthorizationList("sample2")
		require.NotEqual(t, authorizationList, newAuthorizationList)
		k.SetAuthorizationList(ctx, newAuthorizationList)

		list, found = k.GetAuthorizationList(ctx)
		require.True(t, found)
		require.Equal(t, newAuthorizationList, list)
	})
}

func TestKeeper_CheckAuthorization(t *testing.T) {
	t.Run("successfully check authorization", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		signer := sample.AccAddress()
		msg := lightclienttypes.MsgDisableHeaderVerification{
			Creator: signer,
		}
		authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           sdk.MsgTypeURL(&msg),
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
		},
		}
		policies := types.Policies{
			Items: []*types.Policy{
				{
					Address:    signer,
					PolicyType: types.PolicyType_groupOperational,
				},
			},
		}

		k.SetPolicies(ctx, policies)
		k.SetAuthorizationList(ctx, authorizationList)

		err := k.CheckAuthorization(ctx, &msg)
		require.NoError(t, err)
	})

	t.Run("unable to check authorization with multiple signers", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		signer := sample.AccAddress()
		msg := sample.MultipleSignerMessage()
		authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           sdk.MsgTypeURL(msg),
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
		},
		}
		policies := types.Policies{
			Items: []*types.Policy{
				{
					Address:    signer,
					PolicyType: types.PolicyType_groupOperational,
				},
			},
		}
		k.SetPolicies(ctx, policies)
		k.SetAuthorizationList(ctx, authorizationList)

		err := k.CheckAuthorization(ctx, msg)
		require.ErrorIs(t, err, types.ErrSigners)
	})

	t.Run("unable to check authorization with no authorization list", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		signer := sample.AccAddress()
		msg := lightclienttypes.MsgDisableHeaderVerification{
			Creator: signer,
		}

		policies := types.Policies{
			Items: []*types.Policy{
				{
					Address:    signer,
					PolicyType: types.PolicyType_groupOperational,
				},
			},
		}
		k.SetPolicies(ctx, policies)

		err := k.CheckAuthorization(ctx, &msg)
		require.ErrorIs(t, err, types.ErrAuthorizationListNotFound)
	})

	t.Run("unable to check authorization with no policies", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		signer := sample.AccAddress()
		msg := lightclienttypes.MsgDisableHeaderVerification{
			Creator: signer,
		}
		authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           sdk.MsgTypeURL(&msg),
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
		},
		}
		k.SetAuthorizationList(ctx, authorizationList)

		err := k.CheckAuthorization(ctx, &msg)
		require.ErrorIs(t, err, types.ErrPoliciesNotFound)
	})

	t.Run("unable to check authorization when the required authorization doesnt exist", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		signer := sample.AccAddress()
		msg := lightclienttypes.MsgDisableHeaderVerification{
			Creator: signer,
		}
		authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           "/zetachain.zetacore.observer.MsgDisableCCTX",
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
		},
		}
		policies := types.Policies{
			Items: []*types.Policy{
				{
					Address:    signer,
					PolicyType: types.PolicyType_groupOperational,
				},
			},
		}
		k.SetPolicies(ctx, policies)
		k.SetAuthorizationList(ctx, authorizationList)

		err := k.CheckAuthorization(ctx, &msg)
		require.ErrorIs(t, err, types.ErrAuthorizationNotFound)
	})

	t.Run("unable to check authorization when check signer fails", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		signer := sample.AccAddress()
		msg := lightclienttypes.MsgDisableHeaderVerification{
			Creator: signer,
		}
		authorizationList := types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           sdk.MsgTypeURL(&msg),
				AuthorizedPolicy: types.PolicyType_groupOperational,
			},
		},
		}
		policies := types.Policies{
			Items: []*types.Policy{
				{
					Address:    signer,
					PolicyType: types.PolicyType_groupAdmin,
				},
			},
		}
		k.SetPolicies(ctx, policies)
		k.SetAuthorizationList(ctx, authorizationList)

		err := k.CheckAuthorization(ctx, &msg)
		require.ErrorIs(t, err, types.ErrSignerDoesntMatch)
	})
}
