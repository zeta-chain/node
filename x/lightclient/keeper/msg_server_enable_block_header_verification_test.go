package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/lightclient/keeper"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

func TestMsgServer_EnableVerificationFlags(t *testing.T) {
	t.Run("operational group can enable verification flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			EnabledChains: []types.EnabledChain{
				{
					ChainId: chains.EthChain.ChainId,
					Enabled: false,
				},
				{
					ChainId: chains.BtcMainnetChain.ChainId,
					Enabled: false,
				},
			},
		})

		// enable both eth and btc type chain together
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)
		_, err := srv.EnableHeaderVerification(sdk.WrapSDKContext(ctx), &types.MsgEnableHeaderVerification{
			Creator:     admin,
			ChainIdList: []int64{chains.EthChain.ChainId, chains.BtcMainnetChain.ChainId},
		})
		require.NoError(t, err)
		bhv, found := k.GetBlockHeaderVerification(ctx)
		require.True(t, found)
		require.True(t, bhv.IsChainEnabled(chains.EthChain.ChainId))
		require.True(t, bhv.IsChainEnabled(chains.BtcMainnetChain.ChainId))
	})

	t.Run("enable verification flags even if the chain has not been set previously", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		// enable both eth and btc type chain together
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)
		_, err := srv.EnableHeaderVerification(sdk.WrapSDKContext(ctx), &types.MsgEnableHeaderVerification{
			Creator:     admin,
			ChainIdList: []int64{chains.EthChain.ChainId, chains.BtcMainnetChain.ChainId},
		})
		require.NoError(t, err)
		bhv, found := k.GetBlockHeaderVerification(ctx)
		require.True(t, found)
		require.True(t, bhv.IsChainEnabled(chains.EthChain.ChainId))
		require.True(t, bhv.IsChainEnabled(chains.BtcMainnetChain.ChainId))
	})

	t.Run("cannot update if not authorized group", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			EnabledChains: []types.EnabledChain{
				{
					ChainId: chains.EthChain.ChainId,
					Enabled: false,
				},
				{
					ChainId: chains.BtcMainnetChain.ChainId,
					Enabled: false,
				},
			},
		})

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, false)
		_, err := srv.EnableHeaderVerification(sdk.WrapSDKContext(ctx), &types.MsgEnableHeaderVerification{
			Creator:     admin,
			ChainIdList: []int64{chains.EthChain.ChainId},
		})
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
