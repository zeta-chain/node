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
	gatewayObjRef, suiCoinObjRef, withdrawCapObjRef, msgContextObjRef sui.ObjectRef,
	onCallObjectRefs []sui.ObjectRef,
) withdrawAndCallPTBArgs {
	return withdrawAndCallPTBArgs{
		withdrawAndCallObjRefs: withdrawAndCallObjRefs{
			gateway:       gatewayObjRef,
			withdrawCap:   withdrawCapObjRef,
			msgContextRef: msgContextObjRef,
			onCall:        onCallObjectRefs,
			suiCoins:      []*sui.ObjectRef{&suiCoinObjRef},
		},
		coinType:  string(zetasui.SUI),
		amount:    1000000,
		nonce:     1,
		gasBudget: 2000000,
		sender:    sample.EthAddress().Hex(),
		target:    sample.SuiAddress(t),
		payload: zetasui.CallPayload{
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
	msgContextObjRef := sampleObjectRef(t)
	onCallObjRef := sampleObjectRef(t)

	tests := []struct {
		name   string
		args   withdrawAndCallPTBArgs
		errMsg string
	}{
		{
			name: "successful withdraw and call",
			args: newTestWACPTBArgs(
				t,
				gatewayObjRef,
				suiCoinObjRef,
				withdrawCapObjRef,
				msgContextObjRef,
				[]sui.ObjectRef{onCallObjRef},
			),
		},
		{
			name: "successful withdraw and call with empty payload",
			args: func() withdrawAndCallPTBArgs {
				args := newTestWACPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					msgContextObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.payload.Message = []byte{}
				return args
			}(),
		},
		{
			name: "invalid target package ID",
			args: func() withdrawAndCallPTBArgs {
				args := newTestWACPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					msgContextObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.target = "invalid_target_package_id"
				return args
			}(),
			errMsg: "invalid target package ID",
		},
		{
			name: "invalid coin type",
			args: func() withdrawAndCallPTBArgs {
				args := newTestWACPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					msgContextObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.coinType = "invalid_coin_type"
				return args
			}(),
			errMsg: "invalid coin type",
		},
		{
			name: "invalid type argument",
			args: func() withdrawAndCallPTBArgs {
				args := newTestWACPTBArgs(
					t,
					gatewayObjRef,
					suiCoinObjRef,
					withdrawCapObjRef,
					msgContextObjRef,
					[]sui.ObjectRef{onCallObjRef},
				)
				args.payload.TypeArgs[0] = "invalid_type_argument"
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
	msgContextID, err := sui.ObjectIdFromHex(sample.SuiAddress(t))
	require.NoError(t, err)
	onCallObjectID, err := sui.ObjectIdFromHex(sample.SuiAddress(t))
	require.NoError(t, err)
	suiCoinID, err := sui.ObjectIdFromHex(sample.SuiAddress(t))
	require.NoError(t, err)

	// create test object digests
	digest1, err := sui.NewBase58(sample.SuiDigest(t))
	require.NoError(t, err)
	digest2, err := sui.NewBase58(sample.SuiDigest(t))
	require.NoError(t, err)
	digest3, err := sui.NewBase58(sample.SuiDigest(t))
	require.NoError(t, err)
	digest4, err := sui.NewBase58(sample.SuiDigest(t))
	require.NoError(t, err)
	digest5, err := sui.NewBase58(sample.SuiDigest(t))
	require.NoError(t, err)

	// create SUI coin object reference
	suiCoinObjRefs := []*sui.ObjectRef{
		{
			ObjectId: suiCoinID,
			Version:  1,
			Digest:   digest5,
		},
	}

	tests := []struct {
		name            string
		gatewayID       string
		withdrawCapID   string
		msgContextID    string
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
			msgContextID:    msgContextID.String(),
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
						ObjectId: msgContextID.String(),
						Version:  "7",
						Digest:   digest3.String(),
					},
				},
				{
					Data: &models.SuiObjectData{
						ObjectId: onCallObjectID.String(),
						Version:  "3",
						Digest:   digest4.String(),
						Owner: map[string]any{
							"Shared": map[string]any{
								"initial_shared_version": float64(1),
							},
						},
					},
				},
			},
			expected: withdrawAndCallObjRefs{
				gateway: sui.ObjectRef{
					ObjectId: gatewayID,
					Version:  1,
					Digest:   digest1,
				},
				withdrawCap: sui.ObjectRef{
					ObjectId: withdrawCapID,
					Version:  2,
					Digest:   digest2,
				},
				msgContextRef: sui.ObjectRef{
					ObjectId: msgContextID,
					Version:  7,
					Digest:   digest3,
				},
				onCall: []sui.ObjectRef{
					{
						ObjectId: onCallObjectID,
						Version:  1,
						Digest:   digest4,
					},
				},
				suiCoins: suiCoinObjRefs,
			},
		},
		{
			name:            "rpc call fails",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			msgContextID:    msgContextID.String(),
			onCallObjectIDs: []string{onCallObjectID.String()},
			mockError:       sample.ErrSample,
			errMsg:          "failed to get objects",
		},
		{
			name:            "invalid object ID",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			msgContextID:    msgContextID.String(),
			onCallObjectIDs: []string{onCallObjectID.String()},
			mockObjects: []*models.SuiObjectResponse{
				{
					Data: &models.SuiObjectData{
						ObjectId: "invalid_id",
						Version:  "1",
						Digest:   digest1.String(),
					},
				},
				{
					Data: sampleSharedObjectData(t),
				},
				{
					Data: sampleSharedObjectData(t),
				},
				{
					Data: sampleSharedObjectData(t),
				},
			},
			errMsg: "failed to parse object ID",
		},
		{
			name:            "invalid object version",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			msgContextID:    msgContextID.String(),
			onCallObjectIDs: []string{onCallObjectID.String()},
			mockObjects: []*models.SuiObjectResponse{
				{
					Data: &models.SuiObjectData{
						ObjectId: gatewayID.String(),
						Version:  "invalid_version",
						Digest:   digest1.String(),
					},
				},
				{
					Data: sampleSharedObjectData(t),
				},
				{
					Data: sampleSharedObjectData(t),
				},
				{
					Data: sampleSharedObjectData(t),
				},
			},
			errMsg: "failed to parse object version",
		},
		{
			name:            "invalid initial shared version",
			gatewayID:       gatewayID.String(),
			withdrawCapID:   withdrawCapID.String(),
			msgContextID:    msgContextID.String(),
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
				{
					Data: sampleSharedObjectData(t),
				},
				{
					Data: sampleSharedObjectData(t),
				},
				{
					Data: sampleSharedObjectData(t),
				},
			},
			errMsg: "failed to extract initial shared version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			ts := newTestSuite(t)

			// setup RPC mock
			ctx := context.Background()
			ts.SuiMock.On("SuiMultiGetObjects", ctx, mock.Anything).Return(tt.mockObjects, tt.mockError)
			ts.SuiMock.On("GetSuiCoinObjectRefs", ctx, mock.Anything, mock.Anything).Maybe().Return(suiCoinObjRefs, nil)

			// ACT
			got, err := ts.Signer.getWithdrawAndCallObjectRefs(
				ctx,
				tt.withdrawCapID,
				tt.msgContextID,
				tt.onCallObjectIDs,
				100,
			)

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

// sampleSharedObjectData creates a sample Sui object data for a shared object
func sampleSharedObjectData(t *testing.T) *models.SuiObjectData {
	return &models.SuiObjectData{
		ObjectId: sample.SuiAddress(t),
		Version:  "1",
		Digest:   sample.SuiDigest(t),
		Owner: map[string]any{
			"Shared": map[string]any{
				"initial_shared_version": float64(1),
			},
		},
	}
}
