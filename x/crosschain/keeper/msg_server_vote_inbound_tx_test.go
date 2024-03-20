package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
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

func setObservers(t *testing.T, k *keeper.Keeper, ctx sdk.Context, zk keepertest.ZetaKeepers) []string {
	validators := k.GetStakingKeeper().GetAllValidators(ctx)

	validatorAddressListFormatted := make([]string, len(validators))
	for i, validator := range validators {
		valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
		require.NoError(t, err)
		addressTmp, err := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(valAddr.Bytes()))
		require.NoError(t, err)
		validatorAddressListFormatted[i] = addressTmp.String()
	}

	// Add validator to the observer list for voting
	zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{
		ObserverList: validatorAddressListFormatted,
	})
	return validatorAddressListFormatted
}

// TODO: Complete the test cases
// https://github.com/zeta-chain/node/issues/1542
func TestKeeper_VoteOnObservedInboundTx(t *testing.T) {
	t.Run("successfully vote on evm deposit", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		validatorList := setObservers(t, k, ctx, zk)
		to, from := int64(1337), int64(101)
		chains := zk.ObserverKeeper.GetSupportedChains(ctx)
		for _, chain := range chains {
			if common.IsEVMChain(chain.ChainId) {
				from = chain.ChainId
			}
			if common.IsZetaChain(chain.ChainId) {
				to = chain.ChainId
			}
		}
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())

		msg := sample.InboundVote(0, from, to)
		for _, validatorAddr := range validatorList {
			msg.Creator = validatorAddr
			_, err := msgServer.VoteOnObservedInboundTx(
				ctx,
				&msg,
			)
			require.NoError(t, err)
		}
		ballot, _, _ := zk.ObserverKeeper.FindBallot(
			ctx,
			msg.Digest(),
			zk.ObserverKeeper.GetSupportedChainFromChainID(ctx, msg.SenderChainId),
			observertypes.ObservationType_InBoundTx,
		)
		require.Equal(t, ballot.BallotStatus, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		cctx, found := k.GetCrossChainTx(ctx, msg.Digest())
		require.True(t, found)
		require.Equal(t, cctx.CctxStatus.Status, types.CctxStatus_OutboundMined)
		require.Equal(t, cctx.InboundTxParams.TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("prevent double event submission", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// MsgServer for the crosschain keeper
		msgServer := keeper.NewMsgServerImpl(*k)

		// Set the chain ids we want to use to be valid
		params := observertypes.DefaultParams()
		zk.ObserverKeeper.SetParams(
			ctx, params,
		)

		// Convert the validator address into a user address.
		validators := k.GetStakingKeeper().GetAllValidators(ctx)
		validatorAddress := validators[0].OperatorAddress
		valAddr, _ := sdk.ValAddressFromBech32(validatorAddress)
		addresstmp, _ := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(valAddr.Bytes()))
		validatorAddr := addresstmp.String()

		// Add validator to the observer list for voting
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{
			ObserverList: []string{validatorAddr},
		})

		// Add tss to the observer keeper
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())

		// Vote on the FIRST message.
		msg := &types.MsgVoteOnObservedInboundTx{
			Creator:       validatorAddr,
			Sender:        "0x954598965C2aCdA2885B037561526260764095B8",
			SenderChainId: 1337, // ETH
			Receiver:      "0x954598965C2aCdA2885B037561526260764095B8",
			ReceiverChain: 101, // zetachain
			Amount:        sdkmath.NewUintFromString("10000000"),
			Message:       "",
			InBlockHeight: 1,
			GasLimit:      1000000000,
			InTxHash:      "0x7a900ef978743f91f57ca47c6d1a1add75df4d3531da17671e9cf149e1aefe0b",
			CoinType:      0, // zeta
			TxOrigin:      "0x954598965C2aCdA2885B037561526260764095B8",
			Asset:         "",
			EventIndex:    1,
		}
		_, err := msgServer.VoteOnObservedInboundTx(
			ctx,
			msg,
		)
		require.NoError(t, err)

		// Check that the vote passed
		ballot, found := zk.ObserverKeeper.GetBallot(ctx, msg.Digest())
		require.True(t, found)
		require.Equal(t, ballot.BallotStatus, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		//Perform the SAME event. Except, this time, we resubmit the event.
		msg2 := &types.MsgVoteOnObservedInboundTx{
			Creator:       validatorAddr,
			Sender:        "0x954598965C2aCdA2885B037561526260764095B8",
			SenderChainId: 1337,
			Receiver:      "0x954598965C2aCdA2885B037561526260764095B8",
			ReceiverChain: 101,
			Amount:        sdkmath.NewUintFromString("10000000"),
			Message:       "",
			InBlockHeight: 1,
			GasLimit:      1000000001, // <---- Change here
			InTxHash:      "0x7a900ef978743f91f57ca47c6d1a1add75df4d3531da17671e9cf149e1aefe0b",
			CoinType:      0,
			TxOrigin:      "0x954598965C2aCdA2885B037561526260764095B8",
			Asset:         "",
			EventIndex:    1,
		}

		_, err = msgServer.VoteOnObservedInboundTx(
			ctx,
			msg2,
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrObservedTxAlreadyFinalized)
		_, found = zk.ObserverKeeper.GetBallot(ctx, msg2.Digest())
		require.False(t, found)
	})
}

func TestStatus_ChangeStatus(t *testing.T) {
	tt := []struct {
		Name         string
		Status       types.Status
		NonErrStatus types.CctxStatus
		Msg          string
		IsErr        bool
		ErrStatus    types.CctxStatus
	}{
		{
			Name: "Transition on finalize Inbound",
			Status: types.Status{
				Status:              types.CctxStatus_PendingInbound,
				StatusMessage:       "Getting InTX Votes",
				LastUpdateTimestamp: 0,
			},
			Msg:          "Got super majority and finalized Inbound",
			NonErrStatus: types.CctxStatus_PendingOutbound,
			ErrStatus:    types.CctxStatus_Aborted,
			IsErr:        false,
		},
		{
			Name: "Transition on finalize Inbound Fail",
			Status: types.Status{
				Status:              types.CctxStatus_PendingInbound,
				StatusMessage:       "Getting InTX Votes",
				LastUpdateTimestamp: 0,
			},
			Msg:          "Got super majority and finalized Inbound",
			NonErrStatus: types.CctxStatus_OutboundMined,
			ErrStatus:    types.CctxStatus_Aborted,
			IsErr:        false,
		},
	}
	for _, test := range tt {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			test.Status.ChangeStatus(test.NonErrStatus, test.Msg)
			if test.IsErr {
				require.Equal(t, test.ErrStatus, test.Status.Status)
			} else {
				require.Equal(t, test.NonErrStatus, test.Status.Status)
			}
		})
	}
}

func TestKeeper_ProcessZEVMDeposit(t *testing.T) {
	t.Run("process zevm deposit successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// expect DepositCoinZeta to be called
		fungibleMock.On("DepositCoinZeta", mock.Anything, receiver, amount).
			Return(nil)

		// call ProcessZEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = 0
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit returns err without reverting", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)

		// mock unsuccessful HandleEVMDeposit which does not revert
		fungibleMock.On("DepositCoinZeta", mock.Anything, receiver, amount).
			Return(fmt.Errorf("deposit error"), false)

		// call ProcessZEVMDeposit
		cctx := sample.CrossChainTx(t, "test")
		cctx.CctxStatus = &types.Status{Status: types.CctxStatus_PendingInbound}
		cctx.GetCurrentOutTxParam().Receiver = receiver.String()
		cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
		cctx.InboundTxParams.CoinType = common.CoinType_Zeta
		cctx.GetInboundTxParams().SenderChainId = 0
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, "deposit error", cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit reverts fails at GetSupportedChainFromChainID", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// mock unsuccessful GetSupportedChainFromChainID
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil)

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, "", amount)
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("invalid sender chain id %d", cctx.InboundTxParams.SenderChainId), cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails at and GetRevertGasLimit", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock unsuccessful GetRevertGasLimit for ERC20
		fungibleMock.On("GetForeignCoinFromAsset", mock.Anything, asset, senderChain.ChainId).
			Return(fungibletypes.ForeignCoins{}, false)

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("revert gas limit error: %s", types.ErrForeignCoinNotFound), cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert fails at PayGasInERC20AndUpdateCctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""

		// Setup expected calls
		errDeposit := fmt.Errorf("deposit failed")
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		observerMock := keepertest.GetCrosschainObserverMock(t, k)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// mock unsuccessful PayGasInERC20AndUpdateCctx
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).
			Return(nil).Once()

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, fmt.Sprintf("deposit revert message: %s err : %s", errDeposit, observertypes.ErrSupportedChains), cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit reverts fails at UpdateNonce", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)

		// Mock unsuccessful UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, senderChain.ChainName.String()).
			Return(observertypes.ChainNonces{}, false)

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Contains(t, cctx.CctxStatus.StatusMessage, "cannot find receiver chain nonce")
	})

	t.Run("unable to process zevm deposit HandleEVMDeposit revert successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		asset := ""
		errDeposit := fmt.Errorf("deposit failed")

		// Setup expected calls
		// mock unsuccessful HandleEVMDeposit which reverts , i.e returns err and isContractReverted = true
		keepertest.MockRevertForHandleEVMDeposit(fungibleMock, receiver, amount, senderChain.ChainId, errDeposit)

		// Mock successful GetSupportedChainFromChainID
		keepertest.MockGetSupportedChainFromChainID(observerMock, senderChain)

		// mock successful GetRevertGasLimit for ERC20
		keepertest.MockGetRevertGasLimitForERC20(fungibleMock, asset, *senderChain)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *senderChain, asset)
		// mock successful UpdateNonce
		updatedNonce := keepertest.MockUpdateNonce(observerMock, *senderChain)

		// call ProcessZEVMDeposit
		cctx := GetERC20Cctx(t, receiver, *senderChain, asset, amount)
		k.ProcessZEVMDeposit(ctx, cctx)
		require.Equal(t, types.CctxStatus_PendingRevert, cctx.CctxStatus.Status)
		require.Equal(t, errDeposit.Error(), cctx.CctxStatus.StatusMessage)
		require.Equal(t, updatedNonce, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
	})
}

func TestKeeper_ProcessCrosschainMsgPassing(t *testing.T) {
	t.Run("process crosschain msg passing successfully", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		receiverChain := getValidEthChain(t)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *receiverChain, "")

		// mock successful UpdateNonce
		updatedNonce := keepertest.MockUpdateNonce(observerMock, *receiverChain)

		// call ProcessCrosschainMsgPassing
		cctx := GetERC20Cctx(t, receiver, *receiverChain, "", amount)
		k.ProcessCrosschainMsgPassing(ctx, cctx)
		require.Equal(t, types.CctxStatus_PendingOutbound, cctx.CctxStatus.Status)
		require.Equal(t, updatedNonce, cctx.GetCurrentOutTxParam().OutboundTxTssNonce)
	})

	t.Run("unable to process crosschain msg passing PayGasAndUpdateCctx fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		receiverChain := getValidEthChain(t)

		// mock unsuccessful PayGasAndUpdateCctx
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, receiverChain.ChainId).
			Return(nil).Once()

		// call ProcessCrosschainMsgPassing
		cctx := GetERC20Cctx(t, receiver, *receiverChain, "", amount)
		k.ProcessCrosschainMsgPassing(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Equal(t, observertypes.ErrSupportedChains.Error(), cctx.CctxStatus.StatusMessage)
	})

	t.Run("unable to process crosschain msg passing UpdateNonce fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseFungibleMock: true,
			UseObserverMock: true,
		})

		// Setup mock data
		fungibleMock := keepertest.GetCrosschainFungibleMock(t, k)
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		receiverChain := getValidEthChain(t)

		// mock successful PayGasAndUpdateCctx
		keepertest.MockPayGasAndUpdateCCTX(fungibleMock, observerMock, ctx, *k, *receiverChain, "")

		// mock unsuccessful UpdateNonce
		observerMock.On("GetChainNonces", mock.Anything, receiverChain.ChainName.String()).
			Return(observertypes.ChainNonces{}, false)

		// call ProcessCrosschainMsgPassing
		cctx := GetERC20Cctx(t, receiver, *receiverChain, "", amount)
		k.ProcessCrosschainMsgPassing(ctx, cctx)
		require.Equal(t, types.CctxStatus_Aborted, cctx.CctxStatus.Status)
		require.Contains(t, cctx.CctxStatus.StatusMessage, "cannot find receiver chain nonce")
	})
}

func TestKeeper_SaveInbound(t *testing.T) {
	t.Run("should save the cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		cctx := GetERC20Cctx(t, receiver, *senderChain, "", amount)
		eventIndex := sample.Uint64InRange(1, 100)
		k.SaveInbound(ctx, cctx, eventIndex)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.InboundTxParams.TxFinalizationStatus)
		require.True(t, k.IsFinalizedInbound(ctx, cctx.GetInboundTxParams().InboundTxObservedHash, cctx.GetInboundTxParams().SenderChainId, eventIndex))
		_, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
	})

	t.Run("should save the cctx and remove tracker", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain(t)
		cctx := GetERC20Cctx(t, receiver, *senderChain, "", amount)
		hash := sample.Hash()
		cctx.InboundTxParams.InboundTxObservedHash = hash.String()
		k.SetInTxTracker(ctx, types.InTxTracker{
			ChainId:  senderChain.ChainId,
			TxHash:   hash.String(),
			CoinType: 0,
		})
		eventIndex := sample.Uint64InRange(1, 100)

		k.SaveInbound(ctx, cctx, eventIndex)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.InboundTxParams.TxFinalizationStatus)
		require.True(t, k.IsFinalizedInbound(ctx, cctx.GetInboundTxParams().InboundTxObservedHash, cctx.GetInboundTxParams().SenderChainId, eventIndex))
		_, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		_, found = k.GetInTxTracker(ctx, senderChain.ChainId, hash.String())
		require.False(t, found)
	})
}

// GetERC20Cctx returns a sample CrossChainTx with ERC20 params. This is used for testing Inbound and Outbound voting transactions
func GetERC20Cctx(t *testing.T, receiver ethcommon.Address, senderChain common.Chain, asset string, amount *big.Int) *types.CrossChainTx {
	r := sample.Rand()
	cctx := &types.CrossChainTx{
		Creator:          sample.AccAddress(),
		Index:            sample.ZetaIndex(t),
		ZetaFees:         sample.UintInRange(0, 100),
		RelayedMessage:   "",
		CctxStatus:       &types.Status{Status: types.CctxStatus_PendingInbound},
		InboundTxParams:  sample.InboundTxParams(r),
		OutboundTxParams: []*types.OutboundTxParams{sample.OutboundTxParams(r)},
	}

	cctx.GetInboundTxParams().Amount = sdkmath.NewUintFromBigInt(amount)
	cctx.GetInboundTxParams().SenderChainId = senderChain.ChainId
	cctx.GetInboundTxParams().InboundTxObservedHash = sample.Hash().String()
	cctx.GetInboundTxParams().InboundTxBallotIndex = sample.ZetaIndex(t)

	cctx.GetCurrentOutTxParam().ReceiverChainId = senderChain.ChainId
	cctx.GetCurrentOutTxParam().Receiver = receiver.String()
	cctx.GetCurrentOutTxParam().OutboundTxHash = sample.Hash().String()
	cctx.GetCurrentOutTxParam().OutboundTxBallotIndex = sample.ZetaIndex(t)

	cctx.InboundTxParams.CoinType = common.CoinType_ERC20
	for _, outboundTxParam := range cctx.OutboundTxParams {
		outboundTxParam.CoinType = common.CoinType_ERC20
	}

	cctx.GetInboundTxParams().Asset = asset
	cctx.GetInboundTxParams().Sender = sample.EthAddress().String()
	cctx.GetCurrentOutTxParam().OutboundTxTssNonce = 42
	cctx.GetCurrentOutTxParam().OutboundTxGasUsed = 100
	cctx.GetCurrentOutTxParam().OutboundTxEffectiveGasLimit = 100
	return cctx
}
