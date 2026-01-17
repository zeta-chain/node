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

	t.Run("can set new chain info if it doesnt exist", func(t *testing.T) {
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
		chain := sample.Chain(42)

		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		_, err := msgServer.UpdateChainInfo(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainInfo{
			Creator: admin,
			Chain:   chain,
		})
		require.NoError(t, err)

		// Check if the chain info is set
		storedChainInfo, found := k.GetChainInfo(ctx)
		require.True(t, found)
		require.Contains(t, storedChainInfo.Chains, chain)
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
		chainID := int64(42)
		chainInfo := sample.ChainInfo(1)
		chainInfo.Chains[0].ChainId = chainID

		chainInfo.Chains[0].Name = "name"
		k.SetChainInfo(ctx, chainInfo)
		chainInfo.Chains[0].Name = "updated name"
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		_, err := msgServer.UpdateChainInfo(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainInfo{
			Creator: admin,
			Chain:   chainInfo.Chains[0],
		})
		require.NoError(t, err)

		// Check if the chain info is set and updated
		storedChainInfo, found := k.GetChainInfo(ctx)
		require.True(t, found)
		for _, chain := range storedChainInfo.Chains {
			if chain.ChainId == chainID {
				require.Equal(t, "updated name", chain.Name)
			}
		}
	})

	t.Run("add chain to chain info if chain dos not exist", func(t *testing.T) {
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
		chainID := int64(103)
		newChain := sample.Chain(chainID)
		k.SetAuthorizationList(ctx, types.DefaultAuthorizationsList())

		_, err := msgServer.UpdateChainInfo(sdk.WrapSDKContext(ctx), &types.MsgUpdateChainInfo{
			Creator: admin,
			Chain:   newChain,
		})
		require.NoError(t, err)

		// Check if the chain info is set and updated
		storedChainInfo, found := k.GetChainInfo(ctx)
		require.True(t, found)
		require.Equal(t, 4, len(storedChainInfo.Chains))
		require.Contains(t, storedChainInfo.Chains, newChain)
	})
}
