package keeper_test

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_FundGasStabilityPoolFromRemainingFees(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	tt := []struct {
		name                                  string
		gasUsed                               uint64
		effectiveGasPrice                     math.Int
		effectiveGasLimit                     uint64
		fundStabilityPoolReturn               error
		expectFundStabilityPoolCall           bool
		fundStabilityPoolExpectedRemainingFee *big.Int
		isError                               bool
	}{
		{
			name:                        "no call if gasLimit is equal to gasUsed",
			effectiveGasLimit:           42,
			gasUsed:                     42,
			effectiveGasPrice:           math.NewInt(42),
			expectFundStabilityPoolCall: false,
		},
		{
			name:                        "no call if gasLimit is 0",
			effectiveGasLimit:           0,
			gasUsed:                     42,
			effectiveGasPrice:           math.NewInt(42),
			expectFundStabilityPoolCall: false,
		},
		{
			name:                        "no call if gasUsed is 0",
			effectiveGasLimit:           42,
			gasUsed:                     0,
			effectiveGasPrice:           math.NewInt(42),
			expectFundStabilityPoolCall: false,
		},
		{
			name:                        "no call if effectiveGasPrice is 0",
			effectiveGasLimit:           42,
			gasUsed:                     42,
			effectiveGasPrice:           math.NewInt(0),
			expectFundStabilityPoolCall: false,
		},
		{
			name:              "should return error if gas limit is less than gas used",
			effectiveGasLimit: 41,
			gasUsed:           42,
			effectiveGasPrice: math.NewInt(42),
			isError:           true,
		},
		{
			name:                        "should call fund stability pool with correct remaining fees",
			effectiveGasLimit:           100,
			gasUsed:                     90,
			effectiveGasPrice:           math.NewInt(100),
			fundStabilityPoolReturn:     nil,
			expectFundStabilityPoolCall: true,
			fundStabilityPoolExpectedRemainingFee: big.NewInt(
				10 * keeper.RemainingFeesToStabilityPoolPercent,
			), // (100-90)*100 = 1000 => statbilityPool% of 1000 = 10 * statbilityPool
		},
		{
			name:                                  "should return error if fund stability pool returns error",
			effectiveGasLimit:                     100,
			gasUsed:                               90,
			effectiveGasPrice:                     math.NewInt(100),
			fundStabilityPoolReturn:               errors.New("fund stability pool error"),
			expectFundStabilityPoolCall:           true,
			fundStabilityPoolExpectedRemainingFee: big.NewInt(10 * keeper.RemainingFeesToStabilityPoolPercent),
			isError:                               true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx := keepertest.CrosschainKeeperAllMocks(t)
			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)

			// OutboundParams
			outbound := sample.OutboundParams(r)
			outbound.EffectiveGasLimit = tc.effectiveGasLimit
			outbound.GasUsed = tc.gasUsed
			outbound.EffectiveGasPrice = tc.effectiveGasPrice

			if tc.expectFundStabilityPoolCall {
				fungibleMock.On(
					"FundGasStabilityPool", ctx, int64(42), tc.fundStabilityPoolExpectedRemainingFee,
				).Return(tc.fundStabilityPoolReturn)
			}

			err := k.FundGasStabilityPoolFromRemainingFees(ctx, *outbound, 42)
			if tc.isError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			fungibleMock.AssertExpectations(t)
		})
	}
}

func TestKeeper_VoteOutbound(t *testing.T) {
	t.Run("successfully vote on outbound tx with status pending outbound ,vote-type success", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		observer := sample.AccAddress()
		tss := sample.Tss()
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{observer}})
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundSuccessBallot(observerMock, ctx, cctx, senderChain, observer)

		// Successfully mock HandleValidOutbound
		expectedNumberOfOutboundParams := 1
		keepertest.MockSaveOutbound(observerMock, ctx, cctx, tss, expectedNumberOfOutboundParams)

		msgServer := keeper.NewMsgServerImpl(*k)
		msg := types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            chains.ReceiveStatus_success,
			Creator:                           observer,
			ObservedOutboundHash:              sample.Hash().String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundEffectiveGasPrice: math.NewInt(21),
			ObservedOutboundGasUsed:           21,
			CoinType:                          cctx.InboundParams.CoinType,
			ConfirmationMode:                  cctx.GetCurrentOutboundParam().ConfirmationMode,
		}
		_, err := msgServer.VoteOutbound(ctx, &msg)
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_OutboundMined, c.CctxStatus.Status)
		require.Equal(t, msg.Digest(), c.GetCurrentOutboundParam().BallotIndex)
		require.Len(t, c.OutboundParams, expectedNumberOfOutboundParams)
	})

	t.Run("unable to finalize an outbound if the cctx has already been aborted ", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		observer := sample.AccAddress()
		tss := sample.Tss()
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{observer}})
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_Aborted
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundSuccessBallot(observerMock, ctx, cctx, senderChain, observer)

		msgServer := keeper.NewMsgServerImpl(*k)
		msg := types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            chains.ReceiveStatus_success,
			Creator:                           observer,
			ObservedOutboundHash:              sample.Hash().String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundEffectiveGasPrice: math.NewInt(21),
			ObservedOutboundGasUsed:           21,
			CoinType:                          cctx.InboundParams.CoinType,
		}
		_, err := msgServer.VoteOutbound(ctx, &msg)
		require.ErrorIs(t, err, types.ErrCCTXAlreadyFinalized)
	})

	t.Run("unable to finalize an outbound if the cctx has already been outboundmined", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		observer := sample.AccAddress()
		tss := sample.Tss()
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{observer}})
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_OutboundMined
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundSuccessBallot(observerMock, ctx, cctx, senderChain, observer)

		msgServer := keeper.NewMsgServerImpl(*k)
		msg := types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            chains.ReceiveStatus_success,
			Creator:                           observer,
			ObservedOutboundHash:              sample.Hash().String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundEffectiveGasPrice: math.NewInt(21),
			ObservedOutboundGasUsed:           21,
			CoinType:                          cctx.InboundParams.CoinType,
		}
		_, err := msgServer.VoteOutbound(ctx, &msg)
		require.ErrorIs(t, err, types.ErrCCTXAlreadyFinalized)
	})

	t.Run("vote on outbound tx fails if tss is not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		observer := sample.AccAddress()
		tss := sample.Tss()
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{observer}})
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, false).Once()

		msgServer := keeper.NewMsgServerImpl(*k)
		msg := types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            chains.ReceiveStatus_success,
			Creator:                           observer,
			ObservedOutboundHash:              sample.Hash().String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundEffectiveGasPrice: math.NewInt(21),
			ObservedOutboundGasUsed:           21,
			CoinType:                          cctx.InboundParams.CoinType,
		}
		_, err := msgServer.VoteOutbound(ctx, &msg)
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
	})

	t.Run("successfully vote on outbound tx, vote-type failed", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		observer := sample.AccAddress()
		tss := sample.Tss()
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{observer}})
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundFailedBallot(observerMock, ctx, cctx, senderChain, observer)

		// Successfully mock ProcessOutbound
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, senderChain, asset)
		_ = keepertest.MockUpdateNonce(observerMock, senderChain)

		//Successfully mock SaveOutbound
		expectedNumberOfOutboundParams := 2
		keepertest.MockSaveOutboundNewRevertCreated(observerMock, ctx, cctx, tss, expectedNumberOfOutboundParams)
		oldParamsLen := len(cctx.OutboundParams)
		msgServer := keeper.NewMsgServerImpl(*k)
		msg := types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            chains.ReceiveStatus_failed,
			Creator:                           observer,
			ObservedOutboundHash:              sample.Hash().String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundEffectiveGasPrice: math.NewInt(21),
			ObservedOutboundGasUsed:           21,
			CoinType:                          cctx.InboundParams.CoinType,
		}
		_, err := msgServer.VoteOutbound(ctx, &msg)
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_PendingRevert, c.CctxStatus.Status)
		require.Len(t, c.OutboundParams, expectedNumberOfOutboundParams)
		require.Equal(t, types.TxFinalizationStatus_Executed, c.OutboundParams[oldParamsLen-1].TxFinalizationStatus)
		require.Equal(t, types.TxFinalizationStatus_NotFinalized, cctx.GetCurrentOutboundParam().TxFinalizationStatus)
		require.Equal(t, msg.Digest(), c.OutboundParams[0].BallotIndex)
		require.Equal(t, "", c.GetCurrentOutboundParam().BallotIndex)
	})

	t.Run("unsuccessfully vote on outbound tx, vote-type failed", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		observer := sample.AccAddress()
		tss := sample.Tss()
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{observer}})
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundFailedBallot(observerMock, ctx, cctx, senderChain, observer)

		// Mock Failed ProcessOutbound
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, senderChain, 100)
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, senderChain, asset)
		observerMock.On("GetChainNonces", mock.Anything, senderChain.ChainId).
			Return(observertypes.ChainNonces{}, false)
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		//Successfully mock SaveOutBound
		expectedNumberOfOutboundParams := 2
		keepertest.MockSaveOutbound(observerMock, ctx, cctx, tss, expectedNumberOfOutboundParams)
		oldParamsLen := len(cctx.OutboundParams)
		msgServer := keeper.NewMsgServerImpl(*k)
		msg := types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            chains.ReceiveStatus_failed,
			Creator:                           observer,
			ObservedOutboundHash:              sample.Hash().String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundEffectiveGasPrice: math.NewInt(21),
			ObservedOutboundGasUsed:           21,
			CoinType:                          cctx.InboundParams.CoinType,
		}
		_, err := msgServer.VoteOutbound(ctx, &msg)
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_Aborted, c.CctxStatus.Status)
		require.Len(t, c.OutboundParams, expectedNumberOfOutboundParams)
		require.Equal(t, msg.Digest(), c.OutboundParams[0].BallotIndex)
		require.Equal(t, "", c.GetCurrentOutboundParam().BallotIndex)
		// The message processing fails during the creation of the revert tx
		// So the original outbound tx is executed and the revert tx is not finalized.
		// The cctx status is Aborted
		require.Equal(t, types.TxFinalizationStatus_NotFinalized, c.GetCurrentOutboundParam().TxFinalizationStatus)
		require.Equal(t, types.TxFinalizationStatus_Executed, c.OutboundParams[oldParamsLen-1].TxFinalizationStatus)
	})

	t.Run("failure in processing outbound tx", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
			UseFungibleMock: true,
		})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		observer := sample.AccAddress()
		tss := sample.Tss()
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{observer}})
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)

		// Successfully mock GetTSS
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundFailedBallot(observerMock, ctx, cctx, senderChain, observer)

		// Fail ProcessOutbound so that changes are not committed to the state
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, mock.Anything, mock.Anything).
			Return(fungibletypes.ForeignCoins{}, false)

		// Successfully mock GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		//Successfully mock HandleInvalidOutbound
		expectedNumberOfOutboundParams := 1
		keepertest.MockSaveOutbound(observerMock, ctx, cctx, tss, expectedNumberOfOutboundParams)

		msgServer := keeper.NewMsgServerImpl(*k)
		msg := types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            chains.ReceiveStatus_failed,
			Creator:                           observer,
			ObservedOutboundHash:              sample.Hash().String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundEffectiveGasPrice: math.NewInt(21),
			ObservedOutboundGasUsed:           21,
			CoinType:                          cctx.InboundParams.CoinType,
		}
		_, err := msgServer.VoteOutbound(ctx, &msg)
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, msg.Digest(), c.GetCurrentOutboundParam().BallotIndex)
		// Status would be CctxStatus_PendingRevert if process outbound did not fail
		require.Equal(t, types.CctxStatus_Aborted, c.CctxStatus.Status)
		require.Len(t, c.OutboundParams, expectedNumberOfOutboundParams)
	})

	t.Run("fail to finalize outbound if not a finalizing vote", func(t *testing.T) {
		k, ctx, sk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{})

		// Setup mock data
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		tss := sample.Tss()

		// set state to successfully vote on outbound tx
		accAddress, err := observertypes.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		zk.ObserverKeeper.SetObserverSet(
			ctx,
			observertypes.ObserverSet{
				ObserverList: []string{accAddress.String(), sample.AccAddress(), sample.AccAddress()},
			},
		)
		sk.StakingKeeper.SetValidator(ctx, validator)
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetTSS(ctx, tss)

		msgServer := keeper.NewMsgServerImpl(*k)
		msg := &types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            chains.ReceiveStatus_success,
			Creator:                           accAddress.String(),
			ObservedOutboundHash:              sample.Hash().String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundEffectiveGasPrice: math.NewInt(21),
			ObservedOutboundGasUsed:           21,
			CoinType:                          cctx.InboundParams.CoinType,
		}
		_, err = msgServer.VoteOutbound(ctx, msg)
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_PendingOutbound, c.CctxStatus.Status)
		ballot, found := zk.ObserverKeeper.GetBallot(ctx, msg.Digest())
		require.True(t, found)
		require.True(t, ballot.HasVoted(accAddress.String()))
	})

	t.Run("unable to add vote if tss is not present", func(t *testing.T) {
		k, ctx, sk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{})

		// Setup mock data
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		asset := ""
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		tss := sample.Tss()

		// set state to successfully vote on outbound tx
		accAddress, err := observertypes.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{accAddress.String()}})
		sk.StakingKeeper.SetValidator(ctx, validator)
		cctx := GetERC20Cctx(t, receiver, senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)

		msgServer := keeper.NewMsgServerImpl(*k)
		msg := &types.MsgVoteOutbound{
			CctxHash:                          cctx.Index,
			OutboundTssNonce:                  cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:                     cctx.GetCurrentOutboundParam().ReceiverChainId,
			Status:                            chains.ReceiveStatus_success,
			Creator:                           accAddress.String(),
			ObservedOutboundHash:              sample.Hash().String(),
			ValueReceived:                     cctx.GetCurrentOutboundParam().Amount,
			ObservedOutboundBlockHeight:       10,
			ObservedOutboundEffectiveGasPrice: math.NewInt(21),
			ObservedOutboundGasUsed:           21,
			CoinType:                          cctx.InboundParams.CoinType,
		}
		_, err = msgServer.VoteOutbound(ctx, msg)
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_PendingOutbound, c.CctxStatus.Status)
		_, found = zk.ObserverKeeper.GetBallot(ctx, msg.Digest())
		require.False(t, found)
	})
}

func TestKeeper_SaveOutbound(t *testing.T) {
	t.Run("successfully save outbound", func(t *testing.T) {
		//ARRANGE
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		// setup state for crosschain and observer modules
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			Index:    "",
			ChainId:  cctx.GetCurrentOutboundParam().ReceiverChainId,
			Nonce:    cctx.GetCurrentOutboundParam().TssNonce,
			HashList: nil,
		})
		zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
			NonceLow:  int64(cctx.GetCurrentOutboundParam().TssNonce) - 1,
			NonceHigh: int64(cctx.GetCurrentOutboundParam().TssNonce) + 1,
			ChainId:   cctx.GetCurrentOutboundParam().ReceiverChainId,
			Tss:       cctx.GetCurrentOutboundParam().TssPubkey,
		})
		zk.ObserverKeeper.SetTSS(ctx, observertypes.TSS{
			TssPubkey: cctx.GetCurrentOutboundParam().TssPubkey,
		})

		//ACT
		// Save outbound and assert all values are successfully saved
		k.SaveOutbound(ctx, cctx, cctx.GetCurrentOutboundParam().TssPubkey)

		//ASSERT
		_, found := k.GetOutboundTracker(
			ctx,
			cctx.GetCurrentOutboundParam().ReceiverChainId,
			cctx.GetCurrentOutboundParam().TssNonce,
		)
		require.False(t, found)
		pn, found := zk.ObserverKeeper.GetPendingNonces(
			ctx,
			cctx.GetCurrentOutboundParam().TssPubkey,
			cctx.GetCurrentOutboundParam().ReceiverChainId,
		)
		require.True(t, found)
		require.Equal(t, pn.NonceLow, int64(cctx.GetCurrentOutboundParam().TssNonce)+1)
		require.Equal(t, pn.NonceHigh, int64(cctx.GetCurrentOutboundParam().TssNonce)+1)
		_, found = k.GetInboundHashToCctx(ctx, cctx.InboundParams.ObservedHash)
		require.True(t, found)
		_, found = zk.ObserverKeeper.GetNonceToCctx(
			ctx,
			cctx.GetCurrentOutboundParam().TssPubkey,
			cctx.GetCurrentOutboundParam().ReceiverChainId,
			int64(cctx.GetCurrentOutboundParam().TssNonce),
		)
		require.True(t, found)
	})

	t.Run("successfully save outbound with multiple trackers", func(t *testing.T) {
		//ARRANGE
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		// setup state for crosschain and observer modules
		cctx := sample.CrossChainTx(t, "test")
		for _, outboundParams := range cctx.OutboundParams {
			k.SetOutboundTracker(ctx, types.OutboundTracker{
				Index:    "",
				ChainId:  outboundParams.ReceiverChainId,
				Nonce:    outboundParams.TssNonce,
				HashList: nil,
			})
			zk.ObserverKeeper.SetPendingNonces(ctx, observertypes.PendingNonces{
				NonceLow:  int64(cctx.GetCurrentOutboundParam().TssNonce) - 1,
				NonceHigh: int64(cctx.GetCurrentOutboundParam().TssNonce) + 1,
				ChainId:   outboundParams.ReceiverChainId,
				Tss:       outboundParams.TssPubkey,
			})
		}
		cctx.CctxStatus.Status = types.CctxStatus_PendingRevert
		tssPubkey := cctx.GetCurrentOutboundParam().TssPubkey
		zk.ObserverKeeper.SetTSS(ctx, observertypes.TSS{
			TssPubkey: tssPubkey,
		})

		//ACT
		// Save outbound and assert all values are successfully saved
		k.SaveOutbound(ctx, cctx, cctx.GetCurrentOutboundParam().TssPubkey)

		//ASSERT
		for _, outboundParams := range cctx.OutboundParams {
			_, found := k.GetOutboundTracker(
				ctx,
				outboundParams.ReceiverChainId,
				outboundParams.TssNonce,
			)
			require.False(t, found)
			_, found = k.GetInboundHashToCctx(ctx, cctx.InboundParams.ObservedHash)
			require.True(t, found)
		}

		// assert pending nonces
		pn, found := zk.ObserverKeeper.GetPendingNonces(
			ctx,
			cctx.GetCurrentOutboundParam().TssPubkey,
			cctx.GetCurrentOutboundParam().ReceiverChainId,
		)
		require.True(t, found)
		require.GreaterOrEqual(t, pn.NonceLow, int64(cctx.GetCurrentOutboundParam().TssNonce)+1)
		require.GreaterOrEqual(t, pn.NonceHigh, int64(cctx.GetCurrentOutboundParam().TssNonce)+1)

		// assert nonce to cctx mapping
		ncctx, found := zk.ObserverKeeper.GetNonceToCctx(ctx,
			tssPubkey,
			cctx.GetCurrentOutboundParam().ReceiverChainId,
			int64(cctx.GetCurrentOutboundParam().TssNonce))
		require.True(t, found)
		require.Equal(t, cctx.Index, ncctx.CctxIndex)
	})
}

func TestKeeper_ValidateOutboundMessage(t *testing.T) {
	t.Run("successfully validate outbound message", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())
		_, err := k.ValidateOutboundMessage(ctx, types.MsgVoteOutbound{
			CctxHash:         cctx.Index,
			OutboundTssNonce: cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:    cctx.GetCurrentOutboundParam().ReceiverChainId,
		})
		require.NoError(t, err)
	})

	t.Run("failed to validate outbound message if cctx not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		msg := types.MsgVoteOutbound{
			CctxHash:         sample.String(),
			OutboundTssNonce: 1,
		}
		_, err := k.ValidateOutboundMessage(ctx, msg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		require.ErrorContains(t, err, fmt.Sprintf("CCTX %s does not exist", msg.CctxHash))
	})

	t.Run("failed to validate outbound message if nonce does not match", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetCrossChainTx(ctx, *cctx)
		msg := types.MsgVoteOutbound{
			CctxHash:         cctx.Index,
			OutboundTssNonce: 2,
		}
		_, err := k.ValidateOutboundMessage(ctx, msg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		require.ErrorContains(
			t,
			err,
			fmt.Sprintf(
				"OutboundTssNonce %d does not match CCTX OutboundTssNonce %d",
				msg.OutboundTssNonce,
				cctx.GetCurrentOutboundParam().TssNonce,
			),
		)
	})

	t.Run("failed to validate outbound message if chain does not match", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())
		_, err := k.ValidateOutboundMessage(ctx, types.MsgVoteOutbound{
			CctxHash:         cctx.Index,
			OutboundTssNonce: cctx.GetCurrentOutboundParam().TssNonce,
			OutboundChain:    2,
		})
		require.ErrorIs(t, err, sdkerrors.ErrInvalidRequest)
		require.ErrorContains(
			t,
			err,
			fmt.Sprintf(
				"OutboundChain %d does not match CCTX OutboundChain %d",
				2,
				cctx.GetCurrentOutboundParam().ReceiverChainId,
			),
		)
	})
}
