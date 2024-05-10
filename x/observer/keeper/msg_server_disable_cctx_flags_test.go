package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_DisableCCTXFlags(t *testing.T) {
	t.Run("can disable cctx flags if flags dont exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		msg := &types.MsgDisableCCTXFlags{
			Creator:         admin,
			DisableOutbound: true,
			DisableInbound:  true,
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		_, err := srv.DisableCCTXFlags(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
		require.Nil(t, flags.GasPriceIncreaseFlags)
		require.Nil(t, flags.BlockHeaderVerificationFlags)
	})

	t.Run("can disable cctx flags if flags set to true", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:  true,
			IsOutboundEnabled: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		msg := &types.MsgDisableCCTXFlags{
			Creator:         admin,
			DisableOutbound: true,
			DisableInbound:  true,
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		_, err := srv.DisableCCTXFlags(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
	})

	t.Run("can disable only 1 flag", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:  true,
			IsOutboundEnabled: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		msg := &types.MsgDisableCCTXFlags{
			Creator:         admin,
			DisableOutbound: true,
			DisableInbound:  false,
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		_, err := srv.DisableCCTXFlags(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
	})

	t.Run("cannot disable cctx flags if not correct address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})

		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		msg := &types.MsgDisableCCTXFlags{
			Creator:         admin,
			DisableOutbound: true,
			DisableInbound:  false,
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, false)

		_, err := srv.DisableCCTXFlags(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, authoritytypes.ErrUnauthorized, err)

		_, found := k.GetCrosschainFlags(ctx)
		require.False(t, found)
	})
}
