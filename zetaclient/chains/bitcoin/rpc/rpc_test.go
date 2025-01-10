package rpc_test

import (
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/rpc"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_GetEstimatedFeeRate(t *testing.T) {
	tests := []struct {
		name         string
		rate         float64
		regnet       bool
		resultError  bool
		rpcError     bool
		expectedRate int64
		errMsg       string
	}{
		{
			name:         "normal",
			rate:         0.0001,
			regnet:       false,
			expectedRate: 10,
		},
		{
			name:         "should return 1 for regnet",
			rate:         0.0001,
			regnet:       true,
			expectedRate: 1,
		},
		{
			name:     "should return error on rpc error",
			rpcError: true,
			errMsg:   "unable to estimate smart fee",
		},
		{
			name:        "should return error on result error",
			rate:        0.0001,
			resultError: true,
			errMsg:      "fee result contains errors",
		},
		{
			name:         "should return error on negative rate",
			rate:         -0.0001,
			expectedRate: 0,
			errMsg:       "invalid fee rate",
		},
		{
			name:         "should return error if it's greater than max supply",
			rate:         21000000,
			expectedRate: 0,
			errMsg:       "invalid fee rate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := mocks.NewBTCRPCClient(t)

			switch {
			case tt.rpcError:
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
			case tt.resultError:
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).Return(&btcjson.EstimateSmartFeeResult{
					Errors: []string{"error"},
				}, nil)
			default:
				client.On("EstimateSmartFee", mock.Anything, mock.Anything).
					Maybe().
					Return(&btcjson.EstimateSmartFeeResult{
						Errors:  nil,
						FeeRate: &tt.rate,
					}, nil)
			}

			rate, err := rpc.GetEstimatedFeeRate(client, 1, tt.regnet)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				require.Zero(t, rate)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedRate, rate)
		})
	}
}
