package zetacore

import (
	"context"
	"errors"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func TestClient_GetNumberOfUnconfirmedTxs(t *testing.T) {
	tests := []struct {
		name          string
		numTxs        int
		clientErr     error
		expectedCount int
		expectedErr   string
	}{
		{
			name:          "success with zero unconfirmed txs",
			numTxs:        0,
			expectedCount: 0,
		},
		{
			name:          "success with multiple unconfirmed txs",
			numTxs:        42,
			expectedCount: 42,
		},
		{
			name:        "error from cometbft client",
			clientErr:   errors.New("connection refused"),
			expectedErr: "failed to get number of unconfirmed txs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			cometBFTClient := mocks.NewSDKClientWithErr(t, nil, 0).
				SetNumUnconfirmedTxs(tt.numTxs)
			if tt.clientErr != nil {
				cometBFTClient.SetError(tt.clientErr)
			}

			client := setupZetacoreClient(
				t,
				withDefaultObserverKeys(),
				withCometBFT(cometBFTClient),
				withAccountRetriever(t, 5, 4),
			)

			// ACT
			count, err := client.GetNumberOfUnconfirmedTxs(context.Background())

			// ASSERT
			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
				require.Equal(t, 0, count)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedCount, count)
		})
	}
}

func TestClient_GetSyncStatus(t *testing.T) {
	tests := []struct {
		name           string
		syncing        bool
		expectedErr    string
		expectedStatus bool
	}{
		{
			name:           "node is syncing",
			syncing:        true,
			expectedStatus: true,
		},
		{
			name:           "node is not syncing",
			syncing:        false,
			expectedStatus: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			method := "/cosmos.base.tendermint.v1beta1.Service/GetSyncing"
			mockResponse := &cmtservice.GetSyncingResponse{
				Syncing: tt.syncing,
			}
			setupMockServer(t, cmtservice.RegisterServiceServer, method, &cmtservice.GetSyncingRequest{}, mockResponse)

			client := setupZetacoreClient(
				t,
				withDefaultObserverKeys(),
				withCometBFT(mocks.NewSDKClientWithErr(t, nil, 0)),
				withAccountRetriever(t, 5, 4),
			)

			// ACT
			syncing, err := client.GetSyncStatus(context.Background())

			// ASSERT
			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
				require.False(t, syncing)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expectedStatus, syncing)
		})
	}
}
