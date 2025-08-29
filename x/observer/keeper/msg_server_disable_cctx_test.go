package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_MsgDisableCCTX(t *testing.T) {
	t.Run("can disable cctx flags if flags dont exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgDisableCCTX{
			Creator:         admin,
			DisableOutbound: true,
			DisableInbound:  true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.DisableCCTX(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
		require.Nil(t, flags.GasPriceIncreaseFlags)
	})

	t.Run("can disable cctx flags if flags set to true", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		gasPriceIncreaseFlags := sample.GasPriceIncreaseFlags()
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:      true,
			IsOutboundEnabled:     true,
			GasPriceIncreaseFlags: &gasPriceIncreaseFlags,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgDisableCCTX{
			Creator:         admin,
			DisableOutbound: true,
			DisableInbound:  true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.DisableCCTX(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
		require.Equal(t, gasPriceIncreaseFlags, *flags.GasPriceIncreaseFlags)
	})

	t.Run("can disable only outbound flag", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		gasPriceIncreaseFlags := sample.GasPriceIncreaseFlags()
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:      true,
			IsOutboundEnabled:     true,
			GasPriceIncreaseFlags: &gasPriceIncreaseFlags,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgDisableCCTX{
			Creator:         admin,
			DisableOutbound: true,
			DisableInbound:  false,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.DisableCCTX(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
		require.Equal(t, gasPriceIncreaseFlags, *flags.GasPriceIncreaseFlags)
	})

	t.Run("can disable only inbound flag", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		gasPriceIncreaseFlags := sample.GasPriceIncreaseFlags()
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:      true,
			IsOutboundEnabled:     true,
			GasPriceIncreaseFlags: &gasPriceIncreaseFlags,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgDisableCCTX{
			Creator:         admin,
			DisableOutbound: false,
			DisableInbound:  true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.DisableCCTX(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
		require.Equal(t, gasPriceIncreaseFlags, *flags.GasPriceIncreaseFlags)
	})

	t.Run("cannot disable cctx flags if not correct address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})

		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgDisableCCTX{
			Creator:         admin,
			DisableOutbound: true,
			DisableInbound:  false,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := srv.DisableCCTX(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, authoritytypes.ErrUnauthorized, err)

		_, found := k.GetCrosschainFlags(ctx)
		require.False(t, found)
	})
}
