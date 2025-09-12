package keeper_test

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/x/vm/statedb"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/keeper"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func setObservers(t *testing.T, k *keeper.Keeper, ctx sdk.Context, zk keepertest.ZetaKeepers) []string {
	validators, err := k.GetStakingKeeper().GetAllValidators(ctx)
	require.NoError(t, err)

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
func TestKeeper_VoteInbound(t *testing.T) {
	t.Run("successfully vote on evm deposit", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		validatorList := setObservers(t, k, ctx, zk)

		to, from := chains.GoerliLocalnet.ChainId, chains.ZetaChainPrivnet.ChainId
		supportedChains := zk.ObserverKeeper.GetSupportedChains(ctx)
		for _, chain := range supportedChains {
			if chains.IsEthereumChain(chain.ChainId, []chains.Chain{}) {
				from = chain.ChainId
			}
			if chains.IsZetaChain(chain.ChainId, []chains.Chain{}) {
				to = chain.ChainId
			}
		}

		msg := sample.InboundVote(0, from, to)
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())

		err := sdkk.EvmKeeper.SetAccount(ctx, ethcommon.HexToAddress(msg.Receiver), statedb.Account{
			Nonce:    0,
			Balance:  uint256.NewInt(0),
			CodeHash: crypto.Keccak256(nil),
		})
		require.NoError(t, err)
		for _, validatorAddr := range validatorList {
			msg.Creator = validatorAddr
			_, err := msgServer.VoteInbound(
				ctx,
				&msg,
			)
			require.NoError(t, err)
		}

		chain, found := zk.ObserverKeeper.GetSupportedChainFromChainID(ctx, msg.SenderChainId)
		require.True(t, found)

		ballot, _, _ := zk.ObserverKeeper.FindBallot(
			ctx,
			msg.Digest(),
			chain,
			observertypes.ObservationType_InboundTx,
		)
		require.Equal(t, ballot.BallotStatus, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		cctx, found := k.GetCrossChainTx(ctx, msg.Digest())
		require.True(t, found)
		require.Equal(t, types.CctxStatus_OutboundMined, cctx.CctxStatus.Status)
		require.Equal(t, cctx.InboundParams.TxFinalizationStatus, types.TxFinalizationStatus_Executed)
	})

	t.Run("prevent double event submission", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// MsgServer for the crosschain keeper
		msgServer := keeper.NewMsgServerImpl(*k)

		// Convert the validator address into a user address.
		validators, err := k.GetStakingKeeper().GetAllValidators(ctx)
		require.NoError(t, err)
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
		msg := &types.MsgVoteInbound{
			Creator:            validatorAddr,
			Sender:             "0x954598965C2aCdA2885B037561526260764095B8",
			SenderChainId:      1337, // ETH
			Receiver:           "0x954598965C2aCdA2885B037561526260764095B8",
			ReceiverChain:      101, // zetachain
			Amount:             sdkmath.NewUintFromString("10000000"),
			Message:            "",
			InboundBlockHeight: 1,
			CallOptions: &types.CallOptions{
				GasLimit: 1000000000,
			},
			InboundHash:      "0x7a900ef978743f91f57ca47c6d1a1add75df4d3531da17671e9cf149e1aefe0b",
			CoinType:         0, // zeta
			TxOrigin:         "0x954598965C2aCdA2885B037561526260764095B8",
			Asset:            "",
			EventIndex:       1,
			Status:           types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE,
			ConfirmationMode: types.ConfirmationMode_FAST,
		}
		_, err = msgServer.VoteInbound(
			ctx,
			msg,
		)
		require.NoError(t, err)

		// Check that the vote passed
		ballot, found := zk.ObserverKeeper.GetBallot(ctx, msg.Digest())
		require.True(t, found)
		require.Equal(t, ballot.BallotStatus, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		//Perform the SAME event. Except, this time, we resubmit the event.
		msg = &types.MsgVoteInbound{
			Creator:            validatorAddr,
			Sender:             "0x954598965C2aCdA2885B037561526260764095B8",
			SenderChainId:      1337,
			Receiver:           "0x954598965C2aCdA2885B037561526260764095B8",
			ReceiverChain:      101,
			Amount:             sdkmath.NewUintFromString("10000000"),
			Message:            "",
			InboundBlockHeight: 1,
			CallOptions: &types.CallOptions{
				GasLimit: 1000000001, // <---- Change here
			},
			InboundHash:      "0x7a900ef978743f91f57ca47c6d1a1add75df4d3531da17671e9cf149e1aefe0b",
			CoinType:         0,
			TxOrigin:         "0x954598965C2aCdA2885B037561526260764095B8",
			Asset:            "",
			EventIndex:       1,
			Status:           types.InboundStatus_SUCCESS, // <---- Change here
			ConfirmationMode: types.ConfirmationMode_SAFE, // <---- Change here
		}

		_, err = msgServer.VoteInbound(
			ctx,
			msg,
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrObservedTxAlreadyFinalized)
		_, found = zk.ObserverKeeper.GetBallot(ctx, msg.Digest())
		require.False(t, found)
	})

	t.Run("prevent double event submission even if the second ballot is created before the first is finalized", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)

		// MsgServer for the crosschain keeper
		msgServer := keeper.NewMsgServerImpl(*k)

		// Convert the validator address into a user address.
		validators, err := k.GetStakingKeeper().GetAllValidators(ctx)
		require.NoError(t, err)

		fmt.Println("Validators:", len(validators))

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
		msg := &types.MsgVoteInbound{
			Creator:            validatorAddr,
			Sender:             "0x954598965C2aCdA2885B037561526260764095B8",
			SenderChainId:      1337, // ETH
			Receiver:           "0x954598965C2aCdA2885B037561526260764095B8",
			ReceiverChain:      101, // zetachain
			Amount:             sdkmath.NewUintFromString("10000000"),
			Message:            "",
			InboundBlockHeight: 1,
			CallOptions: &types.CallOptions{
				GasLimit: 1000000000,
			},
			InboundHash:      "0x7a900ef978743f91f57ca47c6d1a1add75df4d3531da17671e9cf149e1aefe0b",
			CoinType:         0, // zeta
			TxOrigin:         "0x954598965C2aCdA2885B037561526260764095B8",
			Asset:            "",
			EventIndex:       1,
			Status:           types.InboundStatus_INSUFFICIENT_DEPOSITOR_FEE,
			ConfirmationMode: types.ConfirmationMode_FAST,
		}
		_, err = msgServer.VoteInbound(
			ctx,
			msg,
		)
		require.NoError(t, err)

		// Check that the vote passed
		ballot, found := zk.ObserverKeeper.GetBallot(ctx, msg.Digest())
		require.True(t, found)
		require.Equal(t, ballot.BallotStatus, observertypes.BallotStatus_BallotFinalized_SuccessObservation)
		//Perform the SAME event. Except, this time, we resubmit the event.
		msg = &types.MsgVoteInbound{
			Creator:            validatorAddr,
			Sender:             "0x954598965C2aCdA2885B037561526260764095B8",
			SenderChainId:      1337,
			Receiver:           "0x954598965C2aCdA2885B037561526260764095B8",
			ReceiverChain:      101,
			Amount:             sdkmath.NewUintFromString("10000000"),
			Message:            "",
			InboundBlockHeight: 1,
			CallOptions: &types.CallOptions{
				GasLimit: 1000000001, // <---- Change here
			},
			InboundHash:      "0x7a900ef978743f91f57ca47c6d1a1add75df4d3531da17671e9cf149e1aefe0b",
			CoinType:         0,
			TxOrigin:         "0x954598965C2aCdA2885B037561526260764095B8",
			Asset:            "",
			EventIndex:       1,
			Status:           types.InboundStatus_SUCCESS, // <---- Change here
			ConfirmationMode: types.ConfirmationMode_SAFE, // <---- Change here
		}

		_, err = msgServer.VoteInbound(
			ctx,
			msg,
		)
		require.Error(t, err)
		require.ErrorIs(t, err, types.ErrObservedTxAlreadyFinalized)
		_, found = zk.ObserverKeeper.GetBallot(ctx, msg.Digest())
		require.False(t, found)
	})

	t.Run("should error if vote on inbound ballot fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("VoteOnInboundBallot", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, false, errors.New("err"))
		msgServer := keeper.NewMsgServerImpl(*k)
		to, from := chains.GoerliLocalnet.ChainId, chains.ZetaChainPrivnet.ChainId

		msg := sample.InboundVote(0, from, to)
		res, err := msgServer.VoteInbound(
			ctx,
			&msg,
		)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return if not finalized", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		validatorList := setObservers(t, k, ctx, zk)

		// add one more voter to make it not finalized
		r := rand.New(rand.NewSource(42))
		valAddr := sample.ValAddress(r)
		observerSet := append(validatorList, valAddr.String())
		zk.ObserverKeeper.SetObserverSet(ctx, observertypes.ObserverSet{
			ObserverList: observerSet,
		})
		to, from := chains.GoerliLocalnet.ChainId, chains.ZetaChainPrivnet.ChainId
		supportedChains := zk.ObserverKeeper.GetSupportedChains(ctx)
		for _, chain := range supportedChains {
			if chains.IsEthereumChain(chain.ChainId, []chains.Chain{}) {
				from = chain.ChainId
			}
			if chains.IsZetaChain(chain.ChainId, []chains.Chain{}) {
				to = chain.ChainId
			}
		}
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())

		msg := sample.InboundVote(0, from, to)
		for _, validatorAddr := range validatorList {
			msg.Creator = validatorAddr
			_, err := msgServer.VoteInbound(
				ctx,
				&msg,
			)
			require.NoError(t, err)
		}

		chain, found := zk.ObserverKeeper.GetSupportedChainFromChainID(ctx, msg.SenderChainId)
		require.True(t, found)

		ballot, _, _ := zk.ObserverKeeper.FindBallot(
			ctx,
			msg.Digest(),
			chain,
			observertypes.ObservationType_InboundTx,
		)
		require.Equal(t, ballot.BallotStatus, observertypes.BallotStatus_BallotInProgress)
		require.Equal(t, ballot.Votes[0], observertypes.VoteType_SuccessObservation)
		require.Equal(t, ballot.Votes[1], observertypes.VoteType_NotYetVoted)
		_, found = k.GetCrossChainTx(ctx, msg.Digest())
		require.False(t, found)
	})

	t.Run("should err if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("VoteOnInboundBallot", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(true, false, nil)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, false)
		msgServer := keeper.NewMsgServerImpl(*k)
		to, from := chains.GoerliLocalnet.ChainId, chains.ZetaChainPrivnet.ChainId

		msg := sample.InboundVote(0, from, to)
		res, err := msgServer.VoteInbound(
			ctx,
			&msg,
		)
		require.Error(t, err)
		require.Nil(t, res)
	})
}

func TestStatus_UpdateCctxStatus(t *testing.T) {
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
			test.Status.UpdateStatusAndErrorMessages(test.NonErrStatus, types.StatusMessages{StatusMessage: test.Msg})
			if test.IsErr {
				require.Equal(t, test.ErrStatus, test.Status.Status)
			} else {
				require.Equal(t, test.NonErrStatus, test.Status.Status)
			}
		})
	}
}

func TestKeeper_SaveObservedInboundInformation(t *testing.T) {
	t.Run("should save the cctx", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		cctx := GetERC20Cctx(t, receiver, senderChain, "", amount)
		eventIndex := sample.Uint64InRange(1, 100)
		k.SaveObservedInboundInformation(ctx, cctx, eventIndex)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.InboundParams.TxFinalizationStatus)
		require.True(
			t,
			k.IsFinalizedInbound(
				ctx,
				cctx.GetInboundParams().ObservedHash,
				cctx.GetInboundParams().SenderChainId,
				eventIndex,
			),
		)
		_, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
	})

	t.Run("should save the cctx and remove tracker", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		receiver := sample.EthAddress()
		amount := big.NewInt(42)
		senderChain := getValidEthChain()
		cctx := GetERC20Cctx(t, receiver, senderChain, "", amount)
		hash := sample.Hash()
		cctx.InboundParams.ObservedHash = hash.String()
		k.SetInboundTracker(ctx, types.InboundTracker{
			ChainId:  senderChain.ChainId,
			TxHash:   hash.String(),
			CoinType: 0,
		})
		eventIndex := sample.Uint64InRange(1, 100)
		zk.ObserverKeeper.SetTSS(ctx, sample.Tss())

		k.SaveObservedInboundInformation(ctx, cctx, eventIndex)
		require.Equal(t, types.TxFinalizationStatus_Executed, cctx.InboundParams.TxFinalizationStatus)
		require.True(
			t,
			k.IsFinalizedInbound(
				ctx,
				cctx.GetInboundParams().ObservedHash,
				cctx.GetInboundParams().SenderChainId,
				eventIndex,
			),
		)
		_, found := k.GetCrossChainTx(ctx, cctx.Index)
		require.True(t, found)
		_, found = k.GetInboundTracker(ctx, senderChain.ChainId, hash.String())
		require.False(t, found)
	})
}

// GetERC20Cctx returns a sample CrossChainTx with ERC20 params. This is used for testing Inbound and Outbound voting transactions
func GetERC20Cctx(
	t *testing.T,
	receiver ethcommon.Address,
	senderChain chains.Chain,
	asset string,
	amount *big.Int,
) *types.CrossChainTx {
	r := sample.Rand()
	cctx := &types.CrossChainTx{
		Creator:        sample.AccAddress(),
		Index:          sample.ZetaIndex(t),
		ZetaFees:       sample.UintInRange(0, 100),
		RelayedMessage: "",
		CctxStatus:     &types.Status{Status: types.CctxStatus_PendingInbound},
		InboundParams:  sample.InboundParams(r),
		OutboundParams: []*types.OutboundParams{sample.OutboundParams(r)},
	}

	cctx.GetInboundParams().Amount = sdkmath.NewUintFromBigInt(amount)
	cctx.GetInboundParams().SenderChainId = senderChain.ChainId
	cctx.GetInboundParams().ObservedHash = sample.Hash().String()
	cctx.GetInboundParams().BallotIndex = sample.ZetaIndex(t)

	cctx.GetCurrentOutboundParam().ReceiverChainId = senderChain.ChainId
	cctx.GetCurrentOutboundParam().Receiver = receiver.String()
	cctx.GetCurrentOutboundParam().Hash = sample.Hash().String()
	cctx.GetCurrentOutboundParam().BallotIndex = sample.ZetaIndex(t)

	cctx.InboundParams.CoinType = coin.CoinType_ERC20
	for _, outboundTxParam := range cctx.OutboundParams {
		outboundTxParam.CoinType = coin.CoinType_ERC20
	}

	cctx.GetInboundParams().Asset = asset
	cctx.GetInboundParams().Sender = sample.EthAddress().String()
	cctx.GetCurrentOutboundParam().TssNonce = 42
	cctx.GetCurrentOutboundParam().GasUsed = 100
	cctx.GetCurrentOutboundParam().EffectiveGasLimit = 100
	cctx.GetCurrentOutboundParam().ConfirmationMode = types.ConfirmationMode_SAFE
	return cctx
}
