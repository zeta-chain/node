package sui

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/common"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

const rpcTestnet = "https://sui-testnet-endpoint.blockvision.org"

// The test used the gateway deployed on Sui testnet at
// https://suiscan.xyz/testnet/object/0xe88db37ef3dd9f8b334e3839fa277a8d0e37c329b74a965c2c8e802a737885db/tx-blocks
func TestLiveGateway_ReadInbounds(t *testing.T) {
	if !common.LiveTestEnabled() {
		t.Skip("skipping live test")
		return
	}

	client := sui.NewSuiClient(rpcTestnet)
	ctx := context.Background()
	now := time.Now()

	// query event from last 2 hours
	from := uint64(now.Add(-2 * time.Hour).UnixMilli())

	gateway := NewGateway(client, "0xe88db37ef3dd9f8b334e3839fa277a8d0e37c329b74a965c2c8e802a737885db")
	inbounds, err := gateway.QueryDepositInbounds(ctx, from, uint64(now.UnixMilli()))
	require.NoError(t, err)
	t.Logf("deposit:")
	for _, inbound := range inbounds {
		t.Logf("amount: %d, receiver: %s", inbound.Amount, inbound.Receiver.Hex())
		require.True(t, inbound.IsGasDeposit())
		require.False(t, inbound.IsCrossChainCall)
	}

	inbounds, err = gateway.QueryDepositAndCallInbounds(ctx, from, uint64(now.UnixMilli()))
	require.NoError(t, err)
	t.Logf("depositAndCall:")
	for _, inbound := range inbounds {
		t.Logf("amount: %d, receiver: %s, payload: %v", inbound.Amount, inbound.Receiver.Hex(), inbound.Payload)
		require.True(t, inbound.IsGasDeposit())
		require.True(t, inbound.IsCrossChainCall)
	}
}

func TestGateway_QueryDepositInbounds(t *testing.T) {
	clientMock := mocks.NewSUIClient(t)
	gateway := NewGateway(clientMock, "packageID")
	ctx := context.Background()

	ethAddr1 := sample.EthAddress()
	ethAddr2 := sample.EthAddress()

	tt := []struct {
		name             string
		suiQueryRes      models.PaginatedEventsResponse
		suiQueryErr      error
		expectedInbounds []Inbound
		errContains      string
	}{
		{
			name: "no events",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{},
			},
			expectedInbounds: []Inbound{},
		},
		{
			name:        "query error",
			suiQueryErr: errors.New("query error"),
			errContains: "query error",
		},
		{
			name: "valid events",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "1",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "100",
							"sender":    "0x123",
							"receiver":  ethAddr1.Hex(),
						},
					},
					{
						Id: models.EventId{
							TxDigest: "0xefg",
							EventSeq: "2",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "200",
							"sender":    "0x456",
							"receiver":  ethAddr2.Hex(),
						},
					},
				},
			},
			expectedInbounds: []Inbound{
				{
					TxHash:           "0xabc",
					EventIndex:       1,
					CoinType:         SUI,
					Amount:           100,
					Sender:           "0x123",
					Receiver:         ethAddr1,
					IsCrossChainCall: false,
				},
				{
					TxHash:           "0xefg",
					EventIndex:       2,
					CoinType:         SUI,
					Amount:           200,
					Sender:           "0x456",
					Receiver:         ethAddr2,
					IsCrossChainCall: false,
				},
			},
		},
		{
			name: "invalid event index",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "invalid",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "100",
							"sender":    "0x123",
							"receiver":  ethAddr1.Hex(),
						},
					},
				},
			},
			errContains: "failed to parse event index",
		},
		{
			name: "invalid coin type",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "1",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": 1,
							"amount":    "100",
							"sender":    "0x123",
							"receiver":  ethAddr1.Hex(),
						},
					},
				},
			},
			errContains: "invalid coin type",
		},
		{
			name: "invalid amount",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "1",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    100,
							"sender":    "0x123",
							"receiver":  ethAddr1.Hex(),
						},
					},
				},
			},
			errContains: "invalid amount",
		},
		{
			name: "invalid sender",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "1",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "100",
							"sender":    123,
							"receiver":  ethAddr1.Hex(),
						},
					},
				},
			},
			errContains: "invalid sender",
		},
		{
			name: "invalid receiver",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "1",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "100",
							"sender":    "0x123",
							"receiver":  123,
						},
					},
				},
			},
			errContains: "invalid receiver",
		},
		{
			name: "can't parse receiver as evm address",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "1",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "100",
							"sender":    "0x123",
							"receiver":  "invalid",
						},
					},
				},
			},
			errContains: "can't parse receiver address",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			onQueryEvent(clientMock, tc.suiQueryRes, tc.suiQueryErr)
			inbounds, err := gateway.QueryDepositInbounds(ctx, 0, 0)
			if tc.errContains != "" {
				require.ErrorContains(t, err, tc.errContains)
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tc.expectedInbounds, inbounds)
		})
	}
}

func TestGateway_QueryDepositAndCallInbounds(t *testing.T) {
	clientMock := mocks.NewSUIClient(t)
	gateway := NewGateway(clientMock, "packageID")
	ctx := context.Background()

	ethAddr1 := sample.EthAddress()
	ethAddr2 := sample.EthAddress()

	tt := []struct {
		name             string
		suiQueryRes      models.PaginatedEventsResponse
		suiQueryErr      error
		expectedInbounds []Inbound
		errContains      string
	}{
		{
			name: "no events",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{},
			},
			expectedInbounds: []Inbound{},
		},
		{
			name:        "query error",
			suiQueryErr: errors.New("query error"),
			errContains: "query error",
		},
		{
			name: "valid events",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "1",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "100",
							"sender":    "0x123",
							"receiver":  ethAddr1.Hex(),
							"payload": []interface{}{
								float64(1),
								float64(2),
								float64(3),
							},
						},
					},
					{
						Id: models.EventId{
							TxDigest: "0xefg",
							EventSeq: "2",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "200",
							"sender":    "0x456",
							"receiver":  ethAddr2.Hex(),
							"payload":   []interface{}{},
						},
					},
				},
			},
			expectedInbounds: []Inbound{
				{
					TxHash:           "0xabc",
					EventIndex:       1,
					CoinType:         SUI,
					Amount:           100,
					Sender:           "0x123",
					Receiver:         ethAddr1,
					Payload:          []byte{1, 2, 3},
					IsCrossChainCall: true,
				},
				{
					TxHash:           "0xefg",
					EventIndex:       2,
					CoinType:         SUI,
					Amount:           200,
					Sender:           "0x456",
					Receiver:         ethAddr2,
					Payload:          []byte{},
					IsCrossChainCall: true,
				},
			},
		},
		{
			name: "invalid payload",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "2",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "100",
							"sender":    "0x123",
							"receiver":  ethAddr1.Hex(),
							"payload": []interface{}{
								float64(1),
								uint64(2),
								float64(3),
							},
						},
					},
				},
			},
			errContains: "failed to convert payload",
		},
		{
			name: "invalid payload, not a byte",
			suiQueryRes: models.PaginatedEventsResponse{
				Data: []models.SuiEventResponse{
					{
						Id: models.EventId{
							TxDigest: "0xabc",
							EventSeq: "2",
						},
						ParsedJson: map[string]interface{}{
							"coin_type": string(SUI),
							"amount":    "100",
							"sender":    "0x123",
							"receiver":  ethAddr1.Hex(),
							"payload": []interface{}{
								float64(1),
								float64(256),
								float64(3),
							},
						},
					},
				},
			},
			errContains: "failed to convert payload",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			onQueryEvent(clientMock, tc.suiQueryRes, tc.suiQueryErr)
			inbounds, err := gateway.QueryDepositAndCallInbounds(ctx, 0, 0)
			if tc.errContains != "" {
				require.ErrorContains(t, err, tc.errContains)
				return
			}
			require.NoError(t, err)
			require.ElementsMatch(t, tc.expectedInbounds, inbounds)
		})
	}
}

func onQueryEvent(m *mocks.SUIClient, res models.PaginatedEventsResponse, err error) {
	m.On("SuiXQueryEvents", mock.Anything, mock.Anything).Return(res, err).Once()
}
