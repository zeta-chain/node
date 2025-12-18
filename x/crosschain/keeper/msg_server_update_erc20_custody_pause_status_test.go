package keeper_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	testkeeper "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_UpdateERC20CustodyPauseStatus(t *testing.T) {
	t.Run("can create CCTX to pause ERC20 custody pause status", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, zk := testkeeper.CrosschainKeeperWithMocks(t, testkeeper.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		chainID := getValidEthChain().ChainId
		msgServer := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()

		msg := types.MsgUpdateERC20CustodyPauseStatus{
			Creator: sample.AccAddress(),
			ChainId: chainID,
			Pause:   true,
		}

		// mock authority calls
		authorityMock := testkeeper.GetCrosschainAuthorityMock(t, k)
		testkeeper.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// set necessary values in observer
		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{ChainId: chainID})
		zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{ChainId: chainID, Tss: tss.TssPubkey})
		zk.ObserverKeeper.SetTSS(ctx, tss)
		zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{
			ChainParams: []*observertypes.ChainParams{sample.ChainParamsSupported(chainID)},
		})
		k.SetGasPrice(ctx, sample.GasPriceWithChainID(t, chainID))
		medianGasPrice, priorityFee, isFound := k.GetMedianGasValues(ctx, msg.ChainId)
		require.True(t, isFound)

		// ACT
		res, err := msgServer.UpdateERC20CustodyPauseStatus(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.NoError(t, err)

		// check CCTX is created
		cctx, found := k.GetCrossChainTx(ctx, res.CctxIndex)
		require.True(t, found)
		require.Equal(t, coin.CoinType_Cmd, cctx.InboundParams.CoinType)
		require.Contains(t, cctx.RelayedMessage, constant.CmdUpdateERC20CustodyPauseStatus)
		require.Len(t, cctx.OutboundParams, 1)
		require.EqualValues(
			t,
			medianGasPrice.MulUint64(types.ERC20CustodyPausingGasMultiplierEVM).String(),
			cctx.OutboundParams[0].GasPrice,
		)
		require.EqualValues(
			t,
			priorityFee.MulUint64(types.ERC20CustodyPausingGasMultiplierEVM).String(),
			cctx.OutboundParams[0].GasPriorityFee,
		)
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := testkeeper.CrosschainKeeperWithMocks(t, testkeeper.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)

		msg := types.MsgUpdateERC20CustodyPauseStatus{
			Creator: sample.AccAddress(),
			ChainId: getValidEthChain().ChainId,
			Pause:   true,
		}

		// mock authority calls
		authorityMock := testkeeper.GetCrosschainAuthorityMock(t, k)
		testkeeper.MockCheckAuthorization(&authorityMock.Mock, &msg, errors.New("not authorized"))

		// ACT
		_, err := msgServer.UpdateERC20CustodyPauseStatus(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if can't find chain nonces", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, zk := testkeeper.CrosschainKeeperWithMocks(t, testkeeper.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		chainID := getValidEthChain().ChainId
		msgServer := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()

		msg := types.MsgUpdateERC20CustodyPauseStatus{
			Creator: sample.AccAddress(),
			ChainId: chainID,
			Pause:   true,
		}

		// mock authority calls
		authorityMock := testkeeper.GetCrosschainAuthorityMock(t, k)
		testkeeper.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// set necessary values in observer
		//zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{ChainId: chainID}) // not set
		zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{ChainId: chainID, Tss: tss.TssPubkey})
		zk.ObserverKeeper.SetTSS(ctx, tss)
		zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{
			ChainParams: []*observertypes.ChainParams{sample.ChainParamsSupported(chainID)},
		})
		k.SetGasPrice(ctx, sample.GasPriceWithChainID(t, chainID))

		// ACT
		_, err := msgServer.UpdateERC20CustodyPauseStatus(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.ErrorIs(t, err, types.ErrInvalidChainID)
	})

	t.Run("should fail if can't find current TSS", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, zk := testkeeper.CrosschainKeeperWithMocks(t, testkeeper.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		chainID := getValidEthChain().ChainId
		msgServer := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()

		msg := types.MsgUpdateERC20CustodyPauseStatus{
			Creator: sample.AccAddress(),
			ChainId: chainID,
			Pause:   true,
		}

		// mock authority calls
		authorityMock := testkeeper.GetCrosschainAuthorityMock(t, k)
		testkeeper.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// set necessary values in observer
		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{ChainId: chainID})
		zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{ChainId: chainID, Tss: tss.TssPubkey})
		//zk.ObserverKeeper.SetTSS(ctx, tss) // not set
		zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{
			ChainParams: []*observertypes.ChainParams{sample.ChainParamsSupported(chainID)},
		})
		k.SetGasPrice(ctx, sample.GasPriceWithChainID(t, chainID))

		// ACT
		_, err := msgServer.UpdateERC20CustodyPauseStatus(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
	})

	t.Run("should fail if can't find chain params", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, zk := testkeeper.CrosschainKeeperWithMocks(t, testkeeper.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		chainID := getValidEthChain().ChainId
		msgServer := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()

		msg := types.MsgUpdateERC20CustodyPauseStatus{
			Creator: sample.AccAddress(),
			ChainId: chainID,
			Pause:   true,
		}

		// mock authority calls
		authorityMock := testkeeper.GetCrosschainAuthorityMock(t, k)
		testkeeper.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// set necessary values in observer
		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{ChainId: chainID})
		zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{ChainId: chainID, Tss: tss.TssPubkey})
		zk.ObserverKeeper.SetTSS(ctx, tss)
		zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{}) // not set
		k.SetGasPrice(ctx, sample.GasPriceWithChainID(t, chainID))

		// ACT
		_, err := msgServer.UpdateERC20CustodyPauseStatus(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.ErrorIs(t, err, types.ErrInvalidChainID)
	})

	t.Run("should fail if can't find gas price", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, zk := testkeeper.CrosschainKeeperWithMocks(t, testkeeper.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		chainID := getValidEthChain().ChainId
		msgServer := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()

		msg := types.MsgUpdateERC20CustodyPauseStatus{
			Creator: sample.AccAddress(),
			ChainId: chainID,
			Pause:   true,
		}

		// mock authority calls
		authorityMock := testkeeper.GetCrosschainAuthorityMock(t, k)
		testkeeper.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// set necessary values in observer
		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{ChainId: chainID})
		zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{ChainId: chainID, Tss: tss.TssPubkey})
		zk.ObserverKeeper.SetTSS(ctx, tss)
		zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{
			ChainParams: []*observertypes.ChainParams{sample.ChainParamsSupported(chainID)},
		})
		//k.SetGasPrice(ctx, sample.GasPriceWithChainID(t, chainID)) // not set

		// ACT
		_, err := msgServer.UpdateERC20CustodyPauseStatus(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.ErrorIs(t, err, types.ErrUnableToGetGasPrice)
	})

	t.Run("should fail if priority fees higher than gas price", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, zk := testkeeper.CrosschainKeeperWithMocks(t, testkeeper.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		chainID := getValidEthChain().ChainId
		msgServer := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()

		msg := types.MsgUpdateERC20CustodyPauseStatus{
			Creator: sample.AccAddress(),
			ChainId: chainID,
			Pause:   true,
		}

		// mock authority calls
		authorityMock := testkeeper.GetCrosschainAuthorityMock(t, k)
		testkeeper.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// set necessary values in observer
		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{ChainId: chainID})
		zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{ChainId: chainID, Tss: tss.TssPubkey})
		zk.ObserverKeeper.SetTSS(ctx, tss)
		zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{
			ChainParams: []*observertypes.ChainParams{sample.ChainParamsSupported(chainID)},
		})
		k.SetGasPrice(ctx, types.GasPrice{
			Creator:      sample.AccAddress(),
			ChainId:      chainID,
			Signers:      []string{sample.AccAddress()},
			BlockNums:    []uint64{42},
			Prices:       []uint64{42},
			PriorityFees: []uint64{43},
			MedianIndex:  0,
		})

		// ACT
		_, err := msgServer.UpdateERC20CustodyPauseStatus(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.ErrorIs(t, err, types.ErrInvalidGasAmount)
	})

	t.Run("should fail if can't set outbound info", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, zk := testkeeper.CrosschainKeeperWithMocks(t, testkeeper.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		chainID := getValidEthChain().ChainId
		msgServer := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()

		msg := types.MsgUpdateERC20CustodyPauseStatus{
			Creator: sample.AccAddress(),
			ChainId: chainID,
			Pause:   true,
		}

		// mock authority calls
		authorityMock := testkeeper.GetCrosschainAuthorityMock(t, k)
		testkeeper.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// set necessary values in observer
		zk.ObserverKeeper.SetChainNonces(ctx, observertypes.ChainNonces{ChainId: chainID})
		zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{ChainId: chainID, Tss: tss.TssPubkey})
		zk.ObserverKeeper.SetTSS(ctx, tss)
		zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{
			ChainParams: []*observertypes.ChainParams{
				sample.ChainParams(chainID),
			}, // set non supported chain params to fail
		})
		k.SetGasPrice(ctx, sample.GasPriceWithChainID(t, chainID))

		// ACT
		_, err := msgServer.UpdateERC20CustodyPauseStatus(sdk.WrapSDKContext(ctx), &msg)

		// ASSERT
		require.ErrorIs(t, err, observertypes.ErrSupportedChains)
	})
}
