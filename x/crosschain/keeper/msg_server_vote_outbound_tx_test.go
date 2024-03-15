package keeper_test

import (
	"errors"
	"math/big"
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
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
			name:                                  "should call fund stability pool with correct remaining fees",
			effectiveGasLimit:                     100,
			gasUsed:                               90,
			effectiveGasPrice:                     math.NewInt(100),
			fundStabilityPoolReturn:               nil,
			expectFundStabilityPoolCall:           true,
			fundStabilityPoolExpectedRemainingFee: big.NewInt(10 * keeper.RemainingFeesToStabilityPoolPercent), // (100-90)*100 = 1000 => statbilityPool% of 1000 = 10 * statbilityPool
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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			k, ctx := keepertest.CrosschainKeeperAllMocks(t)
			fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)

			// OutboundTxParams
			outbound := sample.OutboundTxParams(r)
			outbound.OutboundTxEffectiveGasLimit = tc.effectiveGasLimit
			outbound.OutboundTxGasUsed = tc.gasUsed
			outbound.OutboundTxEffectiveGasPrice = tc.effectiveGasPrice

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

func TestKeeper_VoteOnObservedOutboundTx(t *testing.T) {
	t.Run("successfully vote on outbound tx", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""
		observer := sample.AccAddress()
		tss := sample.Tss()
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{observer}})
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutTxParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundSuccessBallot(observerMock, ctx, cctx, *senderChain, observer)

		// Successfully mock GetOutBound
		keepertest.MockGetOutBound(observerMock, ctx)

		//Successfully mock SaveSuccessfulOutBound
		keepertest.MockSaveOutBound(observerMock, ctx, cctx, tss)

		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.VoteOnObservedOutboundTx(ctx, &types.MsgVoteOnObservedOutboundTx{
			CctxHash:                       cctx.Index,
			OutTxTssNonce:                  cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
			OutTxChain:                     cctx.GetCurrentOutTxParam().ReceiverChainId,
			Status:                         common.ReceiveStatus_Success,
			Creator:                        observer,
			ObservedOutTxHash:              sample.Hash().String(),
			ValueReceived:                  cctx.GetCurrentOutTxParam().Amount,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxEffectiveGasPrice: math.NewInt(21),
			ObservedOutTxGasUsed:           21,
			CoinType:                       cctx.CoinType,
		})
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_OutboundMined, c.CctxStatus.Status)
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
		senderChain := getValidEthChain(t)
		asset := ""
		observer := sample.AccAddress()
		tss := sample.Tss()
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{observer}})
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutTxParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundFailedBallot(observerMock, ctx, cctx, *senderChain, observer)

		// Successfully mock GetOutBound
		keepertest.MockGetOutBound(observerMock, ctx)

		// Fail ProcessOutbound so that changes are not committed to the state
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, mock.Anything, mock.Anything).Return(fungibletypes.ForeignCoins{}, false)

		//Successfully mock SaveFailedOutBound
		keepertest.MockSaveOutBound(observerMock, ctx, cctx, tss)

		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.VoteOnObservedOutboundTx(ctx, &types.MsgVoteOnObservedOutboundTx{
			CctxHash:                       cctx.Index,
			OutTxTssNonce:                  cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
			OutTxChain:                     cctx.GetCurrentOutTxParam().ReceiverChainId,
			Status:                         common.ReceiveStatus_Success,
			Creator:                        observer,
			ObservedOutTxHash:              sample.Hash().String(),
			ValueReceived:                  cctx.GetCurrentOutTxParam().Amount,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxEffectiveGasPrice: math.NewInt(21),
			ObservedOutTxGasUsed:           21,
			CoinType:                       cctx.CoinType,
		})
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		// Status would be CctxStatus_PendingRevert if process outbound did not fail
		require.Equal(t, types.CctxStatus_Aborted, c.CctxStatus.Status)
	})

	t.Run("fail to finalize outbound if not a finalizing vote", func(t *testing.T) {
		k, ctx, sk, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{})

		// Setup mock data
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		tss := sample.Tss()

		// set state to successfully vote on outbound tx
		accAddress, err := observertypes.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{accAddress.String(), sample.AccAddress(), sample.AccAddress()}})
		sk.StakingKeeper.SetValidator(ctx, validator)
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutTxParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		zk.ObserverKeeper.SetTSS(ctx, tss)

		msgServer := keeper.NewMsgServerImpl(*k)
		msg := &types.MsgVoteOnObservedOutboundTx{
			CctxHash:                       cctx.Index,
			OutTxTssNonce:                  cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
			OutTxChain:                     cctx.GetCurrentOutTxParam().ReceiverChainId,
			Status:                         common.ReceiveStatus_Success,
			Creator:                        accAddress.String(),
			ObservedOutTxHash:              sample.Hash().String(),
			ValueReceived:                  cctx.GetCurrentOutTxParam().Amount,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxEffectiveGasPrice: math.NewInt(21),
			ObservedOutTxGasUsed:           21,
			CoinType:                       cctx.CoinType,
		}
		_, err = msgServer.VoteOnObservedOutboundTx(ctx, msg)
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
		senderChain := getValidEthChain(t)
		asset := ""
		r := rand.New(rand.NewSource(42))
		validator := sample.Validator(t, r)
		tss := sample.Tss()

		// set state to successfully vote on outbound tx
		accAddress, err := observertypes.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
		require.NoError(t, err)
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{ObserverList: []string{accAddress.String()}})
		sk.StakingKeeper.SetValidator(ctx, validator)
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutTxParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)

		msgServer := keeper.NewMsgServerImpl(*k)
		msg := &types.MsgVoteOnObservedOutboundTx{
			CctxHash:                       cctx.Index,
			OutTxTssNonce:                  cctx.GetCurrentOutTxParam().OutboundTxTssNonce,
			OutTxChain:                     cctx.GetCurrentOutTxParam().ReceiverChainId,
			Status:                         common.ReceiveStatus_Success,
			Creator:                        accAddress.String(),
			ObservedOutTxHash:              sample.Hash().String(),
			ValueReceived:                  cctx.GetCurrentOutTxParam().Amount,
			ObservedOutTxBlockHeight:       10,
			ObservedOutTxEffectiveGasPrice: math.NewInt(21),
			ObservedOutTxGasUsed:           21,
			CoinType:                       cctx.CoinType,
		}
		_, err = msgServer.VoteOnObservedOutboundTx(ctx, msg)
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_PendingOutbound, c.CctxStatus.Status)
		_, found = zk.ObserverKeeper.GetBallot(ctx, msg.Digest())
		require.False(t, found)

	})
}
