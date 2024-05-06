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
	"github.com/zeta-chain/zetacore/pkg/chains"
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

func TestKeeper_VoteOnObservedOutboundTx(t *testing.T) {
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
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundSuccessBallot(observerMock, ctx, cctx, *senderChain, observer)

		// Successfully mock GetOutBound
		keepertest.MockGetOutBound(observerMock, ctx)

		// Successfully mock SaveSuccessfulOutbound
		keepertest.MockSaveOutBound(observerMock, ctx, cctx, tss)

		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.VoteOutbound(ctx, &types.MsgVoteOutbound{
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
		})
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_OutboundMined, c.CctxStatus.Status)
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
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundFailedBallot(observerMock, ctx, cctx, *senderChain, observer)

		// Successfully mock GetOutBound
		keepertest.MockGetOutBound(observerMock, ctx)

		// Successfully mock ProcessOutbound
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)
		_ = keepertest.MockUpdateNonce(observerMock, *senderChain)

		//Successfully mock SaveOutBound
		keepertest.MockSaveOutBoundNewRevertCreated(observerMock, ctx, cctx, tss)
		oldParamsLen := len(cctx.OutboundParams)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.VoteOutbound(ctx, &types.MsgVoteOutbound{
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
		})
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_PendingRevert, c.CctxStatus.Status)
		require.Equal(t, oldParamsLen+1, len(c.OutboundParams))
		require.Equal(t, types.TxFinalizationStatus_Executed, c.OutboundParams[oldParamsLen-1].TxFinalizationStatus)
		require.Equal(t, types.TxFinalizationStatus_NotFinalized, cctx.GetCurrentOutboundParam().TxFinalizationStatus)
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
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundFailedBallot(observerMock, ctx, cctx, *senderChain, observer)

		// Successfully mock GetOutBound
		keepertest.MockGetOutBound(observerMock, ctx)

		// Mock Failed ProcessOutbound
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain, 100)
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)
		observerMock.On("GetChainNonces", mock.Anything, senderChain.ChainName.String()).
			Return(observertypes.ChainNonces{}, false)

		//Successfully mock SaveOutBound
		keepertest.MockSaveOutBound(observerMock, ctx, cctx, tss)
		oldParamsLen := len(cctx.OutboundParams)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.VoteOutbound(ctx, &types.MsgVoteOutbound{
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
		})
		require.NoError(t, err)
		c, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		require.Equal(t, types.CctxStatus_Aborted, c.CctxStatus.Status)
		require.Equal(t, oldParamsLen+1, len(c.OutboundParams))
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
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		cctx.GetCurrentOutboundParam().TssPubkey = tss.TssPubkey
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SetCrossChainTx(ctx, *cctx)

		// Successfully mock GetTSS
		observerMock.On("GetTSS", ctx).Return(observertypes.TSS{}, true).Once()

		// Successfully mock VoteOnOutboundBallot
		keepertest.MockVoteOnOutboundFailedBallot(observerMock, ctx, cctx, *senderChain, observer)

		// Successfully mock GetOutBound
		keepertest.MockGetOutBound(observerMock, ctx)

		// Fail ProcessOutbound so that changes are not committed to the state
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, mock.Anything, mock.Anything).Return(fungibletypes.ForeignCoins{}, false)

		//Successfully mock SaveFailedOutbound
		keepertest.MockSaveOutBound(observerMock, ctx, cctx, tss)

		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.VoteOutbound(ctx, &types.MsgVoteOutbound{
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
		senderChain := getValidEthChain()
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
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
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

func TestKeeper_SaveFailedOutBound(t *testing.T) {
	t.Run("successfully save failed outbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			Index:    "",
			ChainId:  cctx.GetCurrentOutboundParam().ReceiverChainId,
			Nonce:    cctx.GetCurrentOutboundParam().TssNonce,
			HashList: nil,
		})
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SaveFailedOutbound(ctx, cctx, sample.String(), sample.ZetaIndex(t))
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_Aborted)
		_, found := k.GetOutboundTracker(ctx, cctx.GetCurrentOutboundParam().ReceiverChainId, cctx.GetCurrentOutboundParam().TssNonce)
		require.False(t, found)
	})
}

func TestKeeper_SaveSuccessfulOutBound(t *testing.T) {
	t.Run("successfully save successful outbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetOutboundTracker(ctx, types.OutboundTracker{
			Index:    "",
			ChainId:  cctx.GetCurrentOutboundParam().ReceiverChainId,
			Nonce:    cctx.GetCurrentOutboundParam().TssNonce,
			HashList: nil,
		})
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		k.SaveSuccessfulOutbound(ctx, cctx, sample.String())
		require.Equal(t, cctx.GetCurrentOutboundParam().BallotIndex, sample.String())
		_, found := k.GetOutboundTracker(ctx, cctx.GetCurrentOutboundParam().ReceiverChainId, cctx.GetCurrentOutboundParam().TssNonce)
		require.False(t, found)
	})
}

func TestKeeper_SaveOutbound(t *testing.T) {
	t.Run("successfully save outbound", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// setup state for crosschain and observer modules
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus.Status = types.CctxStatus_PendingOutbound
		ballotIndex := sample.String()
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

		// Save outbound and assert all values are successfully saved
		k.SaveOutbound(ctx, cctx, ballotIndex)
		require.Equal(t, cctx.GetCurrentOutboundParam().BallotIndex, ballotIndex)
		_, found := k.GetOutboundTracker(ctx, cctx.GetCurrentOutboundParam().ReceiverChainId, cctx.GetCurrentOutboundParam().TssNonce)
		require.False(t, found)
		pn, found := zk.ObserverKeeper.GetPendingNonces(ctx, cctx.GetCurrentOutboundParam().TssPubkey, cctx.GetCurrentOutboundParam().ReceiverChainId)
		require.True(t, found)
		require.Equal(t, pn.NonceLow, int64(cctx.GetCurrentOutboundParam().TssNonce)+1)
		require.Equal(t, pn.NonceHigh, int64(cctx.GetCurrentOutboundParam().TssNonce)+1)
		_, found = k.GetInboundHashToCctx(ctx, cctx.InboundParams.ObservedHash)
		require.True(t, found)
		_, found = zk.ObserverKeeper.GetNonceToCctx(ctx, cctx.GetCurrentOutboundParam().TssPubkey, cctx.GetCurrentOutboundParam().ReceiverChainId, int64(cctx.GetCurrentOutboundParam().TssNonce))
		require.True(t, found)
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
		require.ErrorContains(t, err, fmt.Sprintf("OutboundTssNonce %d does not match CCTX OutboundTssNonce %d", msg.OutboundTssNonce, cctx.GetCurrentOutboundParam().TssNonce))
	})

	t.Run("failed to validate outbound message if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		cctx := sample.CrossChainTx(t, "test")
		k.SetCrossChainTx(ctx, *cctx)
		_, err := k.ValidateOutboundMessage(ctx, types.MsgVoteOutbound{
			CctxHash:         cctx.Index,
			OutboundTssNonce: cctx.GetCurrentOutboundParam().TssNonce,
		})
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
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
		require.ErrorContains(t, err, fmt.Sprintf("OutboundChain %d does not match CCTX OutboundChain %d", 2, cctx.GetCurrentOutboundParam().ReceiverChainId))
	})
}
