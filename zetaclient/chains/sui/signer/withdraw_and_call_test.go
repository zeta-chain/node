package signer

import (
	"context"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pattonkan/sui-go/sui"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	zetasui "github.com/zeta-chain/node/pkg/contracts/sui"
	"github.com/zeta-chain/node/testutil/sample"
)

// testPTBArgs holds all the arguments needed for withdrawAndCallPTB
type testPTBArgs struct {
	signerAddrStr       string
	gatewayPackageIDStr string
	gatewayModule       string
	gatewayObjRef       sui.ObjectRef
	suiCoinObjRef       sui.ObjectRef
	withdrawCapObjRef   sui.ObjectRef
	onCallObjectRefs    []sui.ObjectRef
	coinTypeStr         string
	amountStr           string
	nonceStr            string
	gasBudget           uint64
	receiver            string
	cp                  zetasui.CallPayload
}

// newTestPTBArgs creates a testArgs struct with default values
func newTestPTBArgs(
	t *testing.T,
	gatewayObjRef, suiCoinObjRef, withdrawCapObjRef sui.ObjectRef,
	onCallObjectRefs []sui.ObjectRef,
) testPTBArgs {
	return testPTBArgs{
		signerAddrStr:       sample.SuiAddress(t),
		gatewayPackageIDStr: sample.SuiAddress(t),
		gatewayModule:       "gateway",
		gatewayObjRef:       gatewayObjRef,
		suiCoinObjRef:       suiCoinObjRef,
		withdrawCapObjRef:   withdrawCapObjRef,
		onCallObjectRefs:    onCallObjectRefs,
		coinTypeStr:         string(zetasui.SUI),
		amountStr:           "1000000",
		nonceStr:            "1",
		gasBudget:           2000000,
		receiver:            sample.SuiAddress(t),
		cp: zetasui.CallPayload{
			TypeArgs:  []string{string(zetasui.SUI)},
			ObjectIDs: []string{sample.SuiAddress(t)},
			Message:   []byte("test message"),
		},
	}
}

func Test_withdrawAndCallPTB(t *testing.T) {
	// create test objects references
	gatewayObjRef := sampleObjectRef(t)
	suiCoinObjRef := sampleObjectRef(t)
	withdrawCapObjRef := sampleObjectRef(t)
	onCallObjRef := sampleObjectRef(t)

	tests := []struct {
		name   string
		args   testPTBArgs
		errMsg string
	}{
		{
			name: "successful withdraw and call",
			args: newTestPTBArgs(t, gatewayObjRef, suiCoinObjRef, withdrawCapObjRef, []sui.ObjectRef{onCallObjRef}),
		},
		{
			name: "successful withdraw and call with empty payload",
			args: func() testPTBArgs {
				args := newTestPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.cp.Message = []byte{}
				return args
			}(),
		},
		{
			name: "invalid signer address",
			args: func() testPTBArgs {
				args := newTestPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.signerAddrStr = "invalid_address"
				return args
			}(),
			errMsg: "invalid signer address",
		},
		{
			name: "invalid gateway package ID",
			args: func() testPTBArgs {
				args := newTestPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.gatewayPackageIDStr = "invalid_package_id"
				return args
			}(),
			errMsg: "invalid gateway package ID",
		},
		{
			name: "invalid coin type",
			args: func() testPTBArgs {
				args := newTestPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.coinTypeStr = "invalid_coin_type"
				return args
			}(),
			errMsg: "invalid coin type",
		},
		{
			name: "unable to create amount argument",
			args: func() testPTBArgs {
				args := newTestPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.amountStr = "invalid_amount"
				return args
			}(),
			errMsg: "unable to create amount argument",
		},
		{
			name: "unable to create nonce argument",
			args: func() testPTBArgs {
				args := newTestPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.nonceStr = "invalid_nonce"
				return args
			}(),
			errMsg: "unable to create nonce argument",
		},
		{
			name: "invalid target package ID",
			args: func() testPTBArgs {
				args := newTestPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.receiver = "invalid_target_package_id"
				return args
			}(),
			errMsg: "invalid target package ID",
		},
		{
			name: "invalid type argument",
			args: func() testPTBArgs {
				args := newTestPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.cp.TypeArgs[0] = "invalid_type_argument"
				return args
			}(),
			errMsg: "invalid type argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := withdrawAndCallPTB(
				tt.args.signerAddrStr,
				tt.args.gatewayPackageIDStr,
				tt.args.gatewayModule,
				tt.args.gatewayObjRef,
				tt.args.suiCoinObjRef,
				tt.args.withdrawCapObjRef,
				tt.args.onCallObjectRefs,
				tt.args.coinTypeStr,
				tt.args.amountStr,
				tt.args.nonceStr,
				tt.args.gasBudget,
				tt.args.receiver,
				tt.args.cp,
			)

			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, got.TxBytes)
		})
	}
}

func Test_getWithdrawAndCallObjectRefs(t *testing.T) {
	// create test objects references
	gatewayID, err := sui.ObjectIdFromHex(sample.SuiAddress(t))
	require.NoError(t, err)
	withdrawCapID, err := sui.ObjectIdFromHex(sample.SuiAddress(t))
	require.NoError(t, err)
	onCallObjectID, err := sui.ObjectIdFromHex(sample.SuiAddress(t))
	require.NoError(t, err)

	// create test object digests
	digest1, err := sui.NewBase58(sample.SuiDigest(t))
	require.NoError(t, err)
	digest2, err := sui.NewBase58(sample.SuiDigest(t))
	require.NoError(t, err)
	digest3, err := sui.NewBase58(sample.SuiDigest(t))
	require.NoError(t, err)

	tests := []struct {
		name            string
		gatewayID       string
		withdrawCapID   string
		onCallObjectIDs []string
		mockObjects     []*models.SuiObjectResponse
		mockError       error
		expected        []sui.ObjectRef
		errMsg          string
	}{
		{
			name:            "successful get object refs",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			onCallObjectIDs: []string{onCallObjectID.String()},
			mockObjects: []*models.SuiObjectResponse{
				{
					Data: &models.SuiObjectData{
						ObjectId: gatewayID.String(),
						Version:  "3",
						Digest:   digest1.String(),
						Owner: map[string]any{
							"Shared": map[string]any{
								"initial_shared_version": float64(1),
							},
						},
					},
				},
				{
					Data: &models.SuiObjectData{
						ObjectId: withdrawCapID.String(),
						Version:  "2",
						Digest:   digest2.String(),
					},
				},
				{
					Data: &models.SuiObjectData{
						ObjectId: onCallObjectID.String(),
						Version:  "3",
						Digest:   digest3.String(),
						Owner: map[string]any{
							"Shared": map[string]any{
								"initial_shared_version": float64(1),
							},
						},
					},
				},
			},
			expected: []sui.ObjectRef{
				{
					ObjectId: gatewayID,
					Version:  1,
					Digest:   digest1,
				},
				{
					ObjectId: withdrawCapID,
					Version:  2,
					Digest:   digest2,
				},
				{
					ObjectId: onCallObjectID,
					Version:  1,
					Digest:   digest3,
				},
			},
		},
		{
			name:            "rpc call fails",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			onCallObjectIDs: []string{onCallObjectID.String()},
			mockError:       sample.ErrSample,
			errMsg:          "failed to get SUI objects",
		},
		{
			name:            "invalid object ID",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			onCallObjectIDs: []string{onCallObjectID.String()},
			mockObjects: []*models.SuiObjectResponse{
				{
					Data: &models.SuiObjectData{
						ObjectId: "invalid_id",
						Version:  "1",
						Digest:   digest1.String(),
					},
				},
			},
			errMsg: "failed to parse object ID",
		},
		{
			name:            "invalid object version",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			onCallObjectIDs: []string{onCallObjectID.String()},
			mockObjects: []*models.SuiObjectResponse{
				{
					Data: &models.SuiObjectData{
						ObjectId: gatewayID.String(),
						Version:  "invalid_version",
						Digest:   digest1.String(),
					},
				},
			},
			errMsg: "failed to parse object version",
		},
		{
			name:            "invalid initial shared version",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			onCallObjectIDs: []string{onCallObjectID.String()},
			mockObjects: []*models.SuiObjectResponse{
				{
					Data: &models.SuiObjectData{
						ObjectId: gatewayID.String(),
						Version:  "1",
						Digest:   digest1.String(),
						Owner: map[string]any{
							"Shared": map[string]any{
								"initial_shared_version": "invalid_version",
							},
						},
					},
				},
			},
			errMsg: "failed to extract initial shared version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			ts := newTestSuite(t)

			// setup mock
			ctx := context.Background()
			ts.SuiMock.On("SuiMultiGetObjects", ctx, mock.Anything).Return(tt.mockObjects, tt.mockError)

			// ACT
			gatewayObjRef, withdrawCapObjRef, onCallObjectRefs, err := getWithdrawAndCallObjectRefs(
				ctx,
				ts.SuiMock,
				tt.gatewayID,
				tt.withdrawCapID,
				tt.onCallObjectIDs,
			)

			// ASSERT
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.EqualValues(t, tt.expected[0], gatewayObjRef)
			require.EqualValues(t, tt.expected[1], withdrawCapObjRef)
			require.EqualValues(t, tt.expected[2:], onCallObjectRefs)
			require.Len(t, onCallObjectRefs, len(tt.onCallObjectIDs))
		})
	}
}

// sampleObjectRef creates a sample Sui object reference
func sampleObjectRef(t *testing.T) sui.ObjectRef {
	objectID := sui.MustObjectIdFromHex(sample.SuiAddress(t))
	digest, err := sui.NewBase58(sample.SuiDigest(t))
	require.NoError(t, err)

	return sui.ObjectRef{
		ObjectId: objectID,
		Version:  1,
		Digest:   digest,
	}
}
