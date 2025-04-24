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

// newTestWACPTBArgs creates a withdrawAndCallPTBArgs struct for testing
func newTestWACPTBArgs(
	t *testing.T,
	gatewayObjRef, suiCoinObjRef, withdrawCapObjRef sui.ObjectRef,
	onCallObjectRefs []sui.ObjectRef,
) withdrawAndCallPTBArgs {
	return withdrawAndCallPTBArgs{
		withdrawAndCallObjRefs: withdrawAndCallObjRefs{
			gatewayObjRef:     gatewayObjRef,
			withdrawCapObjRef: withdrawCapObjRef,
			onCallObjectRefs:  onCallObjectRefs,
			suiCoinObjRef:     suiCoinObjRef,
		},
		coinType:  string(zetasui.SUI),
		amount:    1000000,
		nonce:     1,
		gasBudget: 2000000,
		receiver:  sample.SuiAddress(t),
		cp: zetasui.CallPayload{
			TypeArgs:  []string{string(zetasui.SUI)},
			ObjectIDs: []string{sample.SuiAddress(t)},
			Message:   []byte("test message"),
		},
	}
}

func Test_withdrawAndCallPTB(t *testing.T) {
	// Create a test suite
	ts := newTestSuite(t)

	// create test objects references
	gatewayObjRef := sampleObjectRef(t)
	suiCoinObjRef := sampleObjectRef(t)
	withdrawCapObjRef := sampleObjectRef(t)
	onCallObjRef := sampleObjectRef(t)

	tests := []struct {
		name   string
		args   withdrawAndCallPTBArgs
		errMsg string
	}{
		{
			name: "successful withdraw and call",
			args: newTestWACPTBArgs(t, gatewayObjRef, suiCoinObjRef, withdrawCapObjRef, []sui.ObjectRef{onCallObjRef}),
		},
		{
			name: "successful withdraw and call with empty payload",
			args: func() withdrawAndCallPTBArgs {
				args := newTestWACPTBArgs(
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
			name: "invalid coin type",
			args: func() withdrawAndCallPTBArgs {
				args := newTestWACPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.coinType = "invalid_coin_type"
				return args
			}(),
			errMsg: "invalid coin type",
		},
		{
			name: "invalid target package ID",
			args: func() withdrawAndCallPTBArgs {
				args := newTestWACPTBArgs(
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
			args: func() withdrawAndCallPTBArgs {
				args := newTestWACPTBArgs(
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
			got, err := ts.Signer.withdrawAndCallPTB(tt.args)

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
		expected        withdrawAndCallObjRefs
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
			expected: withdrawAndCallObjRefs{
				gatewayObjRef: sui.ObjectRef{
					ObjectId: gatewayID,
					Version:  1,
					Digest:   digest1,
				},
				withdrawCapObjRef: sui.ObjectRef{
					ObjectId: withdrawCapID,
					Version:  2,
					Digest:   digest2,
				},
				onCallObjectRefs: []sui.ObjectRef{
					{
						ObjectId: onCallObjectID,
						Version:  1,
						Digest:   digest3,
					},
				},
			},
		},
		{
			name:            "rpc call fails",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			onCallObjectIDs: []string{onCallObjectID.String()},
			mockError:       sample.ErrSample,
			errMsg:          "failed to get objects",
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
			got, err := ts.Signer.getWithdrawAndCallObjectRefs(ctx, tt.withdrawCapID, tt.onCallObjectIDs)

			// ASSERT
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
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
