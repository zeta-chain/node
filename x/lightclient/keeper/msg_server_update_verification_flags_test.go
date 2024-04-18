package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/lightclient/keeper"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

func TestMsgServer_UpdateVerificationFlags(t *testing.T) {
	t.Run("operational group can enable verification flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: false,
			BtcTypeChainEnabled: false,
		})

		// enable eth type chain
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)
		_, err := srv.UpdateVerificationFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateVerificationFlags{
			Creator: admin,
			VerificationFlags: types.VerificationFlags{
				EthTypeChainEnabled: true,
				BtcTypeChainEnabled: false,
			},
		})
		require.NoError(t, err)
		vf, found := k.GetVerificationFlags(ctx)
		require.True(t, found)
		require.True(t, vf.EthTypeChainEnabled)
		require.False(t, vf.BtcTypeChainEnabled)

		// enable btc type chain
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)
		_, err = srv.UpdateVerificationFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateVerificationFlags{
			Creator: admin,
			VerificationFlags: types.VerificationFlags{
				EthTypeChainEnabled: false,
				BtcTypeChainEnabled: true,
			},
		})
		require.NoError(t, err)
		vf, found = k.GetVerificationFlags(ctx)
		require.True(t, found)
		require.False(t, vf.EthTypeChainEnabled)
		require.True(t, vf.BtcTypeChainEnabled)

		// enable both eth and btc type chain
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)
		_, err = srv.UpdateVerificationFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateVerificationFlags{
			Creator: admin,
			VerificationFlags: types.VerificationFlags{
				EthTypeChainEnabled: true,
				BtcTypeChainEnabled: true,
			},
		})
		require.NoError(t, err)
		vf, found = k.GetVerificationFlags(ctx)
		require.True(t, found)
		require.True(t, vf.EthTypeChainEnabled)
		require.True(t, vf.BtcTypeChainEnabled)
	})

	t.Run("emergency group can disable verification flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: false,
			BtcTypeChainEnabled: false,
		})

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)
		_, err := srv.UpdateVerificationFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateVerificationFlags{
			Creator: admin,
			VerificationFlags: types.VerificationFlags{
				EthTypeChainEnabled: false,
				BtcTypeChainEnabled: false,
			},
		})
		require.NoError(t, err)
		vf, found := k.GetVerificationFlags(ctx)
		require.True(t, found)
		require.False(t, vf.EthTypeChainEnabled)
		require.False(t, vf.BtcTypeChainEnabled)
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
			EthTypeChainEnabled: false,
			BtcTypeChainEnabled: false,
		})

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, false)
		_, err := srv.UpdateVerificationFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateVerificationFlags{
			Creator: admin,
			VerificationFlags: types.VerificationFlags{
				EthTypeChainEnabled: true,
				BtcTypeChainEnabled: false,
			},
		})
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)

		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, false)
		_, err = srv.UpdateVerificationFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateVerificationFlags{
			Creator: admin,
			VerificationFlags: types.VerificationFlags{
				EthTypeChainEnabled: false,
				BtcTypeChainEnabled: false,
			},
		})
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
