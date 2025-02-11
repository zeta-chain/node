package sui

import (
	"context"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// This is a manual live test, uncomment the t.Skip to run it
// Localnet can currently be started and populated by running the instruction at:
// https://github.com/zeta-chain/protocol-contracts-sui?tab=readme-ov-file#integration-test
// packageID needs to be set to the value logged as moduleId when running `go run main.go`
func TestLiveGateway_ReadInbounds(t *testing.T) {
	//t.Skip("skipping live test")

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
	t.Logf("inbounds: %v", inbounds)
}
