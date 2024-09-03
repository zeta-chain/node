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

func TestKeeper_UpdateGasPriceIncreaseFlags(t *testing.T) {
	t.Run("can update gas price increase flags if crosschain flags dont exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		updatedFlags := sample.GasPriceIncreaseFlags()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// mock the authority keeper for authorization

		msg := types.MsgUpdateGasPriceIncreaseFlags{
			Creator:               admin,
			GasPriceIncreaseFlags: updatedFlags,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.UpdateGasPriceIncreaseFlags(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.Equal(t, updatedFlags, *flags.GasPriceIncreaseFlags)
		require.False(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
	})

	t.Run("can update gas price increase flags if crosschain flags exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		defaultCrosschainFlags := types.DefaultCrosschainFlags()
		k.SetCrosschainFlags(ctx, *defaultCrosschainFlags)
		updatedFlags := sample.GasPriceIncreaseFlags()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateGasPriceIncreaseFlags{
			Creator:               admin,
			GasPriceIncreaseFlags: updatedFlags,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.UpdateGasPriceIncreaseFlags(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.Equal(t, updatedFlags, *flags.GasPriceIncreaseFlags)
		require.Equal(t, defaultCrosschainFlags.IsInboundEnabled, flags.IsInboundEnabled)
		require.Equal(t, defaultCrosschainFlags.IsOutboundEnabled, flags.IsOutboundEnabled)
	})

	t.Run("cannot update invalid gas price increase flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// mock the authority keeper for authorization

		msg := types.MsgUpdateGasPriceIncreaseFlags{
			Creator: admin,
			GasPriceIncreaseFlags: types.GasPriceIncreaseFlags{
				EpochLength:             -1,
				RetryInterval:           1,
				GasPriceIncreasePercent: 1,
			},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.UpdateGasPriceIncreaseFlags(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorContains(t, err, "epoch length must be positive")

		_, found := k.GetCrosschainFlags(ctx)
		require.False(t, found)

	})

	t.Run("cannot update gas price increase flags if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		msg := types.MsgUpdateGasPriceIncreaseFlags{
			Creator:               admin,
			GasPriceIncreaseFlags: sample.GasPriceIncreaseFlags(),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := srv.UpdateGasPriceIncreaseFlags(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorContains(t, err, "sender not authorized")
	})
}
