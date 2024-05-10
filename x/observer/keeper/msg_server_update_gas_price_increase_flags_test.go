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

func TestKeeper_UpdateGasPriceIncreaseFlags(t *testing.T) {
	t.Run("can update gas price increase flags if crosschain flags dont exist", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		updatedFlags := sample.GasPriceIncreaseFlags()
		msg := &types.MsgUpdateGasPriceIncreaseFlags{
			Creator:               admin,
			GasPriceIncreaseFlags: updatedFlags,
		}

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		_, err := srv.UpdateGasPriceIncreaseFlags(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.Equal(t, updatedFlags, *flags.GasPriceIncreaseFlags)
		require.False(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
	})

	t.Run("cannot update invalid gas price increase flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		msg := &types.MsgUpdateGasPriceIncreaseFlags{
			Creator: admin,
			GasPriceIncreaseFlags: types.GasPriceIncreaseFlags{
				EpochLength:             -1,
				RetryInterval:           1,
				GasPriceIncreasePercent: 1,
			},
		}

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, true)

		_, err := srv.UpdateGasPriceIncreaseFlags(sdk.WrapSDKContext(ctx), msg)
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
		msg := &types.MsgUpdateGasPriceIncreaseFlags{
			Creator:               admin,
			GasPriceIncreaseFlags: sample.GasPriceIncreaseFlags(),
		}

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupOperational, false)

		_, err := srv.UpdateGasPriceIncreaseFlags(sdk.WrapSDKContext(ctx), msg)
		require.ErrorContains(t, err, "sender not authorized")
	})
}
