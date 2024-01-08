package keeper_test

import (
	"encoding/hex"
	"fmt"
	"testing"

	//"github.com/zeta-chain/zetacore/common"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

/*
Potential Double Event Submission
*/
func TestNoDoubleEventProtections(t *testing.T) {
	k, ctx, _, zk := keepertest.CrosschainKeeper(t)

	// MsgServer for the crosschain keeper
	msgServer := keeper.NewMsgServerImpl(*k)

	// Set the chain ids we want to use to be valid
	params := observertypes.DefaultParams()
	zk.ObserverKeeper.SetParams(
		ctx, params,
	)

	// Convert the validator address into a user address.
	validators := k.StakingKeeper.GetAllValidators(ctx)
	validatorAddress := validators[0].OperatorAddress
	valAddr, _ := sdk.ValAddressFromBech32(validatorAddress)
	addresstmp, _ := sdk.AccAddressFromHexUnsafe(hex.EncodeToString(valAddr.Bytes()))
	validatorAddr := addresstmp.String()

	// Add validator to the observer list for voting
	chains := zk.ObserverKeeper.GetParams(ctx).GetSupportedChains()
	for _, chain := range chains {
		zk.ObserverKeeper.SetObserverMapper(ctx, &observertypes.ObserverMapper{
			ObserverChain: chain,
			ObserverList:  []string{validatorAddr},
		})
	}

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
	assert.NoError(t, err)

	// Check that the vote passed
	ballot, _, _ := zk.ObserverKeeper.FindBallot(ctx, msg.Digest(), zk.ObserverKeeper.GetParams(ctx).GetChainFromChainID(msg.SenderChainId), observerTypes.ObservationType_InBoundTx)
	assert.Equal(t, ballot.BallotStatus, observerTypes.BallotStatus_BallotFinalized_SuccessObservation)
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

	fmt.Println("Vote again with the same TxHash")
	_, err = msgServer.VoteOnObservedInboundTx(
		ctx,
		msg2,
	)

	assert.ErrorIs(t, err, types.ErrObservedTxAlreadyFinalized)
}

// FIMXE: make it work
//func Test_CalculateGasFee(t *testing.T) {
//
//	tt := []struct {
//		name        string
//		gasPrice    sdk.Uint // Sample gasPrice posted by zeta-client based on observed value and posted to core using PostGasPriceVoter
//		gasLimit    sdk.Uint // Sample gasLimit used in smartContract call
//		rate        sdk.Uint // Sample Rate obtained from UniSwapV2 / V3 and posted to core using PostGasPriceVoter
//		expectedFee sdk.Uint // ExpectedFee in Zeta Tokens
//	}{
//		{
//			name:        "Test Price1",
//			gasPrice:    sdk.NewUintFromString("20000000000"),
//			gasLimit:    sdk.NewUintFromString("90000"),
//			rate:        sdk.NewUintFromString("1000000000000000000"),
//			expectedFee: sdk.NewUintFromString("1001800000000000000"),
//		},
//	}
//	for _, test := range tt {
//		test := test
//		t.Run(test.name, func(t *testing.T) {
//			assert.Equal(t, test.expectedFee, CalculateFee(test.gasPrice, test.gasLimit, test.rate))
//		})
//	}
//}

// FIXME: make it work
//func Test_UpdateGasFees(t *testing.T) {
//	keeper, ctx := setupKeeper(t)
//	cctx := createNCctx(keeper, ctx, 1)
//	cctx[0].Amount = sdk.NewUintFromString("8000000000000000000")
//	keeper.SetGasPrice(ctx, types.GasPrice{
//		Creator:     cctx[0].Creator,
//		Index:       cctx[0].OutboundTxParams.ReceiverChain,
//		Chain:       cctx[0].OutboundTxParams.ReceiverChain,
//		Signers:     []string{cctx[0].Creator},
//		BlockNums:   nil,
//		Prices:      []uint64{20000000000, 20000000000, 20000000000, 20000000000},
//		MedianIndex: 0,
//	})
//	//keeper.SetZetaConversionRate(ctx, types.ZetaConversionRate{
//	//	Index:               cctx[0].OutboundTxParams.ReceiverChain,
//	//	Chain:               cctx[0].OutboundTxParams.ReceiverChain,
//	//	Signers:             []string{cctx[0].Creator},
//	//	BlockNums:           nil,
//	//	ZetaConversionRates: []string{"1000000000000000000", "1000000000000000000", "1000000000000000000", "1000000000000000000"},
//	//	NativeTokenSymbol:   "",
//	//	MedianIndex:         0,
//	//})
//	err := keeper.PayGasInZetaAndUpdateCctx(ctx, cctx[0].OutboundTxParams.ReceiverChain, &cctx[0])
//	assert.NoError(t, err)
//	fmt.Println(cctx[0].String())
//}

func TestStatus_StatusTransition(t *testing.T) {
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
				assert.Equal(t, test.ErrStatus, test.Status.Status)
			} else {
				assert.Equal(t, test.NonErrStatus, test.Status.Status)
			}
		})
	}
}
