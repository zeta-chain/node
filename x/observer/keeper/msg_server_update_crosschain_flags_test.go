package keeper_test

import (
	"testing"
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func setAdminCrossChainFlags(ctx sdk.Context, k *keeper.Keeper, admin string) {
	k.SetParams(ctx, observertypes.Params{
		AdminPolicy: []*observertypes.Admin_Policy{
			{
				PolicyType: observertypes.Policy_Type_stop_inbound_cctx,
				Address:    admin,
			},
		},
	})
}

func TestMsgServer_UpdateCrosschainFlags(t *testing.T) {
	t.Run("can update crosschain flags", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin)

		_, err := srv.UpdateCrosschainFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateCrosschainFlags{
			Creator:           admin,
			IsInboundEnabled:  false,
			IsOutboundEnabled: false,
			GasPriceIncreaseFlags: &types.GasPriceIncreaseFlags{
				EpochLength:             42,
				RetryInterval:           time.Minute * 42,
				GasPriceIncreasePercent: 42,
			},
		})
		require.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		require.True(t, found)
		require.False(t, flags.IsInboundEnabled)
		require.False(t, flags.IsOutboundEnabled)
		require.Equal(t, int64(42), flags.GasPriceIncreaseFlags.EpochLength)
		require.Equal(t, time.Minute*42, flags.GasPriceIncreaseFlags.RetryInterval)
		require.Equal(t, uint32(42), flags.GasPriceIncreaseFlags.GasPriceIncreasePercent)
	})
}
