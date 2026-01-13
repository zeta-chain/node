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

func TestMsgServer_UpdateOperationalChainParams(t *testing.T) {
	t.Run("cannot update chain params if not authorized", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateOperationalChainParams{
			Creator:                   admin,
			ChainId:                   1,
			GasPriceTicker:            1,
			InboundTicker:             1,
			OutboundTicker:            1,
			WatchUtxoTicker:           1,
			OutboundScheduleInterval:  1,
			OutboundScheduleLookahead: 1,
			ConfirmationParams: types.ConfirmationParams{
				FastInboundCount:  1,
				FastOutboundCount: 1,
				SafeInboundCount:  1,
				SafeOutboundCount: 1,
			},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)

		// act
		_, err := srv.UpdateOperationalChainParams(sdk.WrapSDKContext(ctx), &msg)

		// assert
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("cannot update chain params if chain params list not found", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateOperationalChainParams{
			Creator:                   admin,
			ChainId:                   1,
			GasPriceTicker:            1,
			InboundTicker:             1,
			OutboundTicker:            1,
			WatchUtxoTicker:           1,
			OutboundScheduleInterval:  1,
			OutboundScheduleLookahead: 1,
			ConfirmationParams: types.ConfirmationParams{
				FastInboundCount:  1,
				FastOutboundCount: 1,
				SafeInboundCount:  1,
				SafeOutboundCount: 1,
			},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// act
		_, err := srv.UpdateOperationalChainParams(sdk.WrapSDKContext(ctx), &msg)

		// assert
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)
	})

	t.Run("cannot update chain params if chain params not found", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateOperationalChainParams{
			Creator:                   admin,
			ChainId:                   1,
			GasPriceTicker:            1,
			InboundTicker:             1,
			OutboundTicker:            1,
			WatchUtxoTicker:           1,
			OutboundScheduleInterval:  1,
			OutboundScheduleLookahead: 1,
			ConfirmationParams: types.ConfirmationParams{
				FastInboundCount:  1,
				FastOutboundCount: 1,
				SafeInboundCount:  1,
				SafeOutboundCount: 1,
			},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				sample.ChainParams(2),
			},
		})

		// act
		_, err := srv.UpdateOperationalChainParams(sdk.WrapSDKContext(ctx), &msg)

		// assert
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)
	})

	t.Run("can update chain params", func(t *testing.T) {
		// arrange
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgUpdateOperationalChainParams{
			Creator:                   admin,
			ChainId:                   1,
			GasPriceTicker:            1000,
			InboundTicker:             1001,
			OutboundTicker:            1002,
			WatchUtxoTicker:           1003,
			OutboundScheduleInterval:  1004,
			OutboundScheduleLookahead: 1005,
			ConfirmationParams: types.ConfirmationParams{
				FastInboundCount:  1006,
				FastOutboundCount: 1007,
				SafeInboundCount:  1008,
				SafeOutboundCount: 1009,
			},
			DisableTssBlockScan: true,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		originalChainParams := sample.ChainParams(1)
		otherChainParams1 := sample.ChainParams(2)
		otherChainParams2 := sample.ChainParams(3)
		otherChainParams3 := sample.ChainParams(4)
		otherChainParams4 := sample.ChainParams(5)

		// ensure original values are different
		require.NotEqualValues(t, originalChainParams.GasPriceTicker, msg.GasPriceTicker)
		require.NotEqualValues(t, originalChainParams.InboundTicker, msg.InboundTicker)
		require.NotEqualValues(t, originalChainParams.OutboundTicker, msg.OutboundTicker)
		require.NotEqualValues(t, originalChainParams.WatchUtxoTicker, msg.WatchUtxoTicker)
		require.NotEqualValues(t, originalChainParams.OutboundScheduleInterval, msg.OutboundScheduleInterval)
		require.NotEqualValues(t, originalChainParams.OutboundScheduleLookahead, msg.OutboundScheduleLookahead)
		require.NotEqualValues(t, *originalChainParams.ConfirmationParams, msg.ConfirmationParams)
		require.NotEqualValues(t, originalChainParams.DisableTssBlockScan, msg.DisableTssBlockScan)

		k.SetChainParamsList(ctx, types.ChainParamsList{
			ChainParams: []*types.ChainParams{
				otherChainParams1,
				otherChainParams2,
				originalChainParams,
				otherChainParams3,
				otherChainParams4,
			},
		})

		// act
		_, err := srv.UpdateOperationalChainParams(sdk.WrapSDKContext(ctx), &msg)

		// assert
		require.NoError(t, err)

		// check new chain param list
		chainParamsList, found := k.GetChainParamsList(ctx)
		require.True(t, found)

		require.Len(t, chainParamsList.ChainParams, 5)
		require.Contains(t, chainParamsList.ChainParams, otherChainParams1)
		require.Contains(t, chainParamsList.ChainParams, otherChainParams2)
		require.Contains(t, chainParamsList.ChainParams, otherChainParams3)
		require.Contains(t, chainParamsList.ChainParams, otherChainParams4)

		found = false
		for _, cp := range chainParamsList.ChainParams {
			if cp.ChainId == 1 {
				found = true

				// fields updated
				require.EqualValues(t, msg.GasPriceTicker, cp.GasPriceTicker)
				require.EqualValues(t, msg.InboundTicker, cp.InboundTicker)
				require.EqualValues(t, msg.OutboundTicker, cp.OutboundTicker)
				require.EqualValues(t, msg.WatchUtxoTicker, cp.WatchUtxoTicker)
				require.EqualValues(t, msg.OutboundScheduleInterval, cp.OutboundScheduleInterval)
				require.EqualValues(t, msg.OutboundScheduleLookahead, cp.OutboundScheduleLookahead)
				require.EqualValues(t, msg.ConfirmationParams, *cp.ConfirmationParams)

				// other chain params not changed
				require.EqualValues(t, originalChainParams.ZetaTokenContractAddress, cp.ZetaTokenContractAddress)
				require.EqualValues(t, originalChainParams.ConnectorContractAddress, cp.ConnectorContractAddress)
				require.EqualValues(t, originalChainParams.Erc20CustodyContractAddress, cp.Erc20CustodyContractAddress)
				require.EqualValues(t, originalChainParams.BallotThreshold.String(), cp.BallotThreshold.String())
				require.EqualValues(
					t,
					originalChainParams.MinObserverDelegation.String(),
					cp.MinObserverDelegation.String(),
				)
				require.EqualValues(t, originalChainParams.IsSupported, cp.IsSupported)
				require.EqualValues(t, originalChainParams.GatewayAddress, cp.GatewayAddress)

				return
			}
		}
		require.True(t, found)
	})
}
