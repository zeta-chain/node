package signer

import (
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
)

func Test_withdrawCapID(t *testing.T) {
	// ARRANGE
	withdrawCapID := sample.SuiAddress(t)
	originalPackageID := sample.SuiAddress(t)

	// create test suite and specify withdraw cap ID
	ts := newTestSuite(t, func(cfg *testSuiteConfig) {
		cfg.withdrawCapID = withdrawCapID
		cfg.originalPackageID = originalPackageID
	})

	// ACT
	got, err := ts.withdrawCapID(ts.Ctx)

	// ASSERT
	require.NoError(t, err)
	require.Equal(t, withdrawCapID, got)
}

func Test_getMessageContextID(t *testing.T) {
	tests := []struct {
		name       string
		mockObject *models.SuiObjectData
		rpcError   bool
		want       string
		errMsg     string
	}{
		{
			name: "success",
			mockObject: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: map[string]any{"value": "0x123"},
					},
				},
			},
			want: "0x123",
		},
		{
			name:     "rpc error",
			rpcError: true,
			errMsg:   "unable to get message context dynamic field object",
		},
		{
			name:       "Sui object data is nil",
			mockObject: nil,
			errMsg:     "dynamic field object data is nil",
		},
		{
			name: "Sui object data content is nil",
			mockObject: &models.SuiObjectData{
				Content: nil,
			},
			errMsg: "dynamic field object data is nil",
		},
		{
			name: "unable to parse message context ID",
			mockObject: &models.SuiObjectData{
				Content: &models.SuiParsedData{
					SuiMoveObject: models.SuiMoveObject{
						Fields: map[string]any{"value": 123},
					},
				},
			},
			errMsg: "unable to parse message context ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			ts := newTestSuite(t)

			// setup RPC mock
			if tt.rpcError {
				ts.SuiMock.On("SuiXGetDynamicFieldObject", ts.Ctx, mock.Anything).Return(models.SuiObjectResponse{
					Data: nil,
				}, errors.New("rpc error"))
			} else {
				ts.SuiMock.On("SuiXGetDynamicFieldObject", ts.Ctx, mock.Anything).Return(models.SuiObjectResponse{
					Data: tt.mockObject,
				}, nil)
			}

			// ACT
			got, err := ts.getMessageContextID(ts.Ctx)

			// ASSERT
			if tt.errMsg != "" {
				require.Empty(t, got)
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
