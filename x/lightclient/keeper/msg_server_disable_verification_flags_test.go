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

func TestMsgServer_DisableVerificationFlags(t *testing.T) {
	t.Run("emergency group can disable verification flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			ChainId: chains.EthChain.ChainId,
			Enabled: true,
		})
		k.SetVerificationFlags(ctx, types.VerificationFlags{
			ChainId: chains.BscMainnetChain.ChainId,
			Enabled: true,
		})

		// enable eth type chain
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)
		_, err := srv.DisableVerificationFlags(sdk.WrapSDKContext(ctx), &types.MsgDisableVerificationFlags{
			Creator:     admin,
			ChainIdList: []int64{chains.EthChain.ChainId, chains.BtcMainnetChain.ChainId},
		})
		require.NoError(t, err)
		vf, found := k.GetVerificationFlags(ctx, chains.EthChain.ChainId)
		require.True(t, found)
		require.False(t, vf.Enabled)
		vf, found = k.GetVerificationFlags(ctx, chains.BtcMainnetChain.ChainId)
		require.True(t, found)
		require.False(t, vf.Enabled)

	})

	t.Run("cannot update if not authorized group", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			ChainId: chains.EthChain.ChainId,
			Enabled: true,
		})

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)
		_, err := srv.DisableVerificationFlags(sdk.WrapSDKContext(ctx), &types.MsgDisableVerificationFlags{
			Creator:     admin,
			ChainIdList: []int64{chains.EthChain.ChainId},
		})
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("disable chain if even if the the chain has nto been set before", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		// enable eth type chain
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)
		_, err := srv.DisableVerificationFlags(sdk.WrapSDKContext(ctx), &types.MsgDisableVerificationFlags{
			Creator:     admin,
			ChainIdList: []int64{chains.EthChain.ChainId, chains.BtcMainnetChain.ChainId},
		})
		require.NoError(t, err)
		vf, found := k.GetVerificationFlags(ctx, chains.EthChain.ChainId)
		require.True(t, found)
		require.False(t, vf.Enabled)
		vf, found = k.GetVerificationFlags(ctx, chains.BtcMainnetChain.ChainId)
		require.True(t, found)
		require.False(t, vf.Enabled)
	})
}
