package keeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

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
	_, _ = setupKeeper(t)
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
