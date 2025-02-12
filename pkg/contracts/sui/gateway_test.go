package sui

import (
	"context"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// This is a manual live test, uncomment the t.Skip to run it
// The test used the gateway deployed on Sui testnet at
// https://suiscan.xyz/testnet/object/0xe88db37ef3dd9f8b334e3839fa277a8d0e37c329b74a965c2c8e802a737885db/tx-blocks
func TestLiveGateway_ReadInbounds(t *testing.T) {
	t.Skip("skipping live test")

	client := sui.NewSuiClient("https://sui-testnet-endpoint.blockvision.org")
	ctx := context.Background()
	now := time.Now()

	// query event from last 2 hours
	from := uint64(now.Add(-2 * time.Hour).UnixMilli())

	gateway := NewGateway(
		client,
		"0xe88db37ef3dd9f8b334e3839fa277a8d0e37c329b74a965c2c8e802a737885db",
	)
	inbounds, err := gateway.QueryDepositInbounds(ctx, from, uint64(now.UnixMilli()))
	require.NoError(t, err)
	t.Logf("deposit:")
	for _, inbound := range inbounds {
		t.Logf("amount: %d, receiver: %s", inbound.Amount, inbound.Receiver.Hex())
	}

	inbounds, err = gateway.QueryDepositAndCallInbounds(ctx, from, uint64(now.UnixMilli()))
	require.NoError(t, err)
	t.Logf("depositAndCall:")
	for _, inbound := range inbounds {
		t.Logf("amount: %d, receiver: %s, payload: %v", inbound.Amount, inbound.Receiver.Hex(), inbound.Payload)
	}
}
