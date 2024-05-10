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

func TestMsgServer_EnableCCTXFlags(t *testing.T) {
	t.Run("can enable cctx flags if flags dont exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		msg := &types.MsgEnableCCTXFlags{
			Creator:        admin,
			EnableInbound:  true,
			EnableOutbound: true,
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		_, err := srv.EnableCCTXFlags(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
		require.Nil(t, flags.GasPriceIncreaseFlags)
		require.Nil(t, flags.BlockHeaderVerificationFlags)
	})

	t.Run("can enable cctx flags if flags set to false", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:  false,
			IsOutboundEnabled: false,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		msg := &types.MsgEnableCCTXFlags{
			Creator:        admin,
			EnableInbound:  true,
			EnableOutbound: true,
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		_, err := srv.EnableCCTXFlags(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
	})

	t.Run("can enable cctx flags one flag", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:  false,
			IsOutboundEnabled: false,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		msg := &types.MsgEnableCCTXFlags{
			Creator:        admin,
			EnableInbound:  true,
			EnableOutbound: false,
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		_, err := srv.EnableCCTXFlags(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
	})

	t.Run("cannot enable cctx flags if not correct address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})

		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		msg := &types.MsgEnableCCTXFlags{
			Creator:        admin,
			EnableInbound:  true,
			EnableOutbound: false,
		}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, false)

		_, err := srv.EnableCCTXFlags(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, authoritytypes.ErrUnauthorized, err)

		_, found := k.GetCrosschainFlags(ctx)
		require.False(t, found)
	})
}
