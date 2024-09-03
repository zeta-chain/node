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

func TestMsgServer_EnableCCTX(t *testing.T) {
	t.Run("can enable cctx flags if flags dont exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := &types.MsgEnableCCTX{
			Creator:        admin,
			EnableInbound:  true,
			EnableOutbound: true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := srv.EnableCCTX(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
		require.Nil(t, flags.GasPriceIncreaseFlags)
	})

	t.Run("can enable cctx flags if flags set to false", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		gasPriceIncreaseFlags := sample.GasPriceIncreaseFlags()
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:      false,
			IsOutboundEnabled:     false,
			GasPriceIncreaseFlags: &gasPriceIncreaseFlags,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := &types.MsgEnableCCTX{
			Creator:        admin,
			EnableInbound:  true,
			EnableOutbound: true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := srv.EnableCCTX(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
		require.Equal(t, gasPriceIncreaseFlags, *flags.GasPriceIncreaseFlags)
	})

	t.Run("can enable cctx flags only inbound flag", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		gasPriceIncreaseFlags := sample.GasPriceIncreaseFlags()
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:      false,
			IsOutboundEnabled:     false,
			GasPriceIncreaseFlags: &gasPriceIncreaseFlags,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := &types.MsgEnableCCTX{
			Creator:        admin,
			EnableInbound:  true,
			EnableOutbound: false,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := srv.EnableCCTX(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.True(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
		require.Equal(t, gasPriceIncreaseFlags, *flags.GasPriceIncreaseFlags)
	})

	t.Run("can enable cctx flags only outbound flag", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		gasPriceIncreaseFlags := sample.GasPriceIncreaseFlags()
		k.SetCrosschainFlags(ctx, types.CrosschainFlags{
			IsInboundEnabled:      false,
			IsOutboundEnabled:     false,
			GasPriceIncreaseFlags: &gasPriceIncreaseFlags,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := &types.MsgEnableCCTX{
			Creator:        admin,
			EnableInbound:  false,
			EnableOutbound: true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := srv.EnableCCTX(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsInboundEnabled)
		require.True(t, flags.IsOutboundEnabled)
		require.Equal(t, gasPriceIncreaseFlags, *flags.GasPriceIncreaseFlags)
	})

	t.Run("cannot enable cctx flags if not correct address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})

		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := &types.MsgEnableCCTX{
			Creator:        admin,
			EnableInbound:  true,
			EnableOutbound: false,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)

		_, err := srv.EnableCCTX(sdk.WrapSDKContext(ctx), msg)
		require.ErrorIs(t, authoritytypes.ErrUnauthorized, err)

		_, found := k.GetCrosschainFlags(ctx)
		require.False(t, found)
	})
}
