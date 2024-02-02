package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func setAdminCrossChainFlags(ctx sdk.Context, k *keeper.Keeper, admin string, group types.Policy_Type) {
	k.SetParams(ctx, observertypes.Params{
		AdminPolicy: []*observertypes.Admin_Policy{
			{
				PolicyType: group,
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

		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		_, err := srv.UpdateCrosschainFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateCrosschainFlags{
			Creator:           admin,
			IsInboundEnabled:  false,
			IsOutboundEnabled: false,
			GasPriceIncreaseFlags: &types.GasPriceIncreaseFlags{
				EpochLength:             42,
				RetryInterval:           time.Minute * 42,
				GasPriceIncreasePercent: 42,
			},
			BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
				IsEthTypeChainEnabled: true,
				IsBtcTypeChainEnabled: false,
			},
		})
		assert.NoError(t, err)

		flags, found := k.GetCrosschainFlags(ctx)
		assert.True(t, found)
		assert.False(t, flags.IsInboundEnabled)
		assert.False(t, flags.IsOutboundEnabled)
		assert.Equal(t, int64(42), flags.GasPriceIncreaseFlags.EpochLength)
		assert.Equal(t, time.Minute*42, flags.GasPriceIncreaseFlags.RetryInterval)
		assert.Equal(t, uint32(42), flags.GasPriceIncreaseFlags.GasPriceIncreasePercent)
		assert.True(t, flags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled)
		assert.False(t, flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled)

		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		// can update flags again
		_, err = srv.UpdateCrosschainFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateCrosschainFlags{
			Creator:           admin,
			IsInboundEnabled:  true,
			IsOutboundEnabled: true,
			GasPriceIncreaseFlags: &types.GasPriceIncreaseFlags{
				EpochLength:             43,
				RetryInterval:           time.Minute * 43,
				GasPriceIncreasePercent: 43,
			},
			BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
				IsEthTypeChainEnabled: false,
				IsBtcTypeChainEnabled: true,
			},
		})
		assert.NoError(t, err)

		flags, found = k.GetCrosschainFlags(ctx)
		assert.True(t, found)
		assert.True(t, flags.IsInboundEnabled)
		assert.True(t, flags.IsOutboundEnabled)
		assert.Equal(t, int64(43), flags.GasPriceIncreaseFlags.EpochLength)
		assert.Equal(t, time.Minute*43, flags.GasPriceIncreaseFlags.RetryInterval)
		assert.Equal(t, uint32(43), flags.GasPriceIncreaseFlags.GasPriceIncreasePercent)
		assert.False(t, flags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled)
		assert.True(t, flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled)

		// group 1 should be able to disable inbound and outbound
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group1)

		// if gas price increase flags is nil, it should not be updated
		_, err = srv.UpdateCrosschainFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateCrosschainFlags{
			Creator:           admin,
			IsInboundEnabled:  false,
			IsOutboundEnabled: false,
		})
		assert.NoError(t, err)

		flags, found = k.GetCrosschainFlags(ctx)
		assert.True(t, found)
		assert.False(t, flags.IsInboundEnabled)
		assert.False(t, flags.IsOutboundEnabled)
		assert.Equal(t, int64(43), flags.GasPriceIncreaseFlags.EpochLength)
		assert.Equal(t, time.Minute*43, flags.GasPriceIncreaseFlags.RetryInterval)
		assert.Equal(t, uint32(43), flags.GasPriceIncreaseFlags.GasPriceIncreasePercent)
		assert.False(t, flags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled)
		assert.True(t, flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled)

		// group 1 should be able to disable header verification
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group1)

		// if gas price increase flags is nil, it should not be updated
		_, err = srv.UpdateCrosschainFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateCrosschainFlags{
			Creator:           admin,
			IsInboundEnabled:  false,
			IsOutboundEnabled: false,
			BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
				IsEthTypeChainEnabled: false,
				IsBtcTypeChainEnabled: false,
			},
		})
		assert.NoError(t, err)

		flags, found = k.GetCrosschainFlags(ctx)
		assert.True(t, found)
		assert.False(t, flags.IsInboundEnabled)
		assert.False(t, flags.IsOutboundEnabled)
		assert.Equal(t, int64(43), flags.GasPriceIncreaseFlags.EpochLength)
		assert.Equal(t, time.Minute*43, flags.GasPriceIncreaseFlags.RetryInterval)
		assert.Equal(t, uint32(43), flags.GasPriceIncreaseFlags.GasPriceIncreasePercent)
		assert.False(t, flags.BlockHeaderVerificationFlags.IsEthTypeChainEnabled)
		assert.False(t, flags.BlockHeaderVerificationFlags.IsBtcTypeChainEnabled)

		// if flags are not defined, default should be used
		k.RemoveCrosschainFlags(ctx)
		_, found = k.GetCrosschainFlags(ctx)
		assert.False(t, found)

		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		_, err = srv.UpdateCrosschainFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateCrosschainFlags{
			Creator:           admin,
			IsInboundEnabled:  false,
			IsOutboundEnabled: true,
		})
		assert.NoError(t, err)

		flags, found = k.GetCrosschainFlags(ctx)
		assert.True(t, found)
		assert.False(t, flags.IsInboundEnabled)
		assert.True(t, flags.IsOutboundEnabled)
		assert.Equal(t, types.DefaultGasPriceIncreaseFlags.EpochLength, flags.GasPriceIncreaseFlags.EpochLength)
		assert.Equal(t, types.DefaultGasPriceIncreaseFlags.RetryInterval, flags.GasPriceIncreaseFlags.RetryInterval)
		assert.Equal(t, types.DefaultGasPriceIncreaseFlags.GasPriceIncreasePercent, flags.GasPriceIncreaseFlags.GasPriceIncreasePercent)
	})

	t.Run("cannot update crosschain flags if not authorized", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		_, err := srv.UpdateCrosschainFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateCrosschainFlags{
			Creator:           sample.AccAddress(),
			IsInboundEnabled:  false,
			IsOutboundEnabled: false,
		})
		assert.Error(t, err)
		assert.Equal(t, types.ErrNotAuthorizedPolicy, err)

		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group1)

		_, err = srv.UpdateCrosschainFlags(sdk.WrapSDKContext(ctx), &types.MsgUpdateCrosschainFlags{
			Creator:           admin,
			IsInboundEnabled:  false,
			IsOutboundEnabled: true,
		})
		assert.Error(t, err)
		assert.Equal(t, types.ErrNotAuthorizedPolicy, err)
	})
}
