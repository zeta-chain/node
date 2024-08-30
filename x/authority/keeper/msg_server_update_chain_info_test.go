package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/keeper"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestMsgServer_UpdateChainInfo(t *testing.T) {
	t.Run("can't update chain info if not authorized", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		msg := types.MsgUpdateChainInfo{
			Creator: sample.AccAddress(),
		}
		k.SetAuthorizationList(ctx, types.AuthorizationList{Authorizations: []types.Authorization{
			{
				MsgUrl:           sdk.MsgTypeURL(&msg),
				AuthorizedPolicy: types.PolicyType_groupAdmin,
			},
		}})

		_, err := msgServer.UpdateChainInfo(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, types.ErrUnauthorized)
	})

	t.Run("can set chain info when it doesn't exist", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		_, found := k.GetChainInfo(ctx)
		require.False(t, found)

		// Set group admin policy
		admin := sample.AccAddress()
		k.SetPolicies(ctx, types.Policies{
			Items: []*types.Policy{
				{
					PolicyType: types.PolicyType_groupAdmin,
					Address:    admin,
				},
			},
		})
		chainInfo := sample.ChainInfo(42)

		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		_, err := msgServer.UpdateChainInfo(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainInfo{
			Creator:   admin,
			ChainInfo: chainInfo,
		})
		require.NoError(t, err)

		// Check if the chain info is set
		storedChainInfo, found := k.GetChainInfo(ctx)
		require.True(t, found)
		require.Equal(t, chainInfo, storedChainInfo)
	})

	t.Run("can update existing chain info", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		k.SetChainInfo(ctx, sample.ChainInfo(42))

		// Set group admin policy
		admin := sample.AccAddress()
		k.SetPolicies(ctx, types.Policies{
			Items: []*types.Policy{
				{
					PolicyType: types.PolicyType_groupAdmin,
					Address:    admin,
				},
			},
		})
		chainInfo := sample.ChainInfo(84)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		_, err := msgServer.UpdateChainInfo(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainInfo{
			Creator:   admin,
			ChainInfo: chainInfo,
		})
		require.NoError(t, err)

		// Check if the chain info is set
		storedChainInfo, found := k.GetChainInfo(ctx)
		require.True(t, found)
		require.Equal(t, chainInfo, storedChainInfo)
	})

	t.Run("can remove chain info", func(t *testing.T) {
		k, ctx := keepertest.AuthorityKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)

		k.SetChainInfo(ctx, sample.ChainInfo(42))

		// Set group admin policy
		admin := sample.AccAddress()
		k.SetPolicies(ctx, types.Policies{
			Items: []*types.Policy{
				{
					PolicyType: types.PolicyType_groupAdmin,
					Address:    admin,
				},
			},
		})
		chainInfo := types.ChainInfo{}
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		_, err := msgServer.UpdateChainInfo(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainInfo{
			Creator:   admin,
			ChainInfo: chainInfo,
		})
		require.NoError(t, err)

		// The structure should still exist but be empty
		storedChainInfo, found := k.GetChainInfo(ctx)
		require.True(t, found)
		require.Equal(t, chainInfo, storedChainInfo)
	})

}
