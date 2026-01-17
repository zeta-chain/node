package runner

import (
	"context"
	"math/big"
	"testing"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/zetaclient/common"
)

func TestGetSuiGatewaySUIBalance(t *testing.T) {
	// live test, comment out to run
	if !common.LiveTestEnabled() {
		t.Skip("skipping live test")
	}

	// create a mainnet Sui client and use the gateway deployed on mainnet
	ctx := context.Background()
	client := sui.NewSuiClient("https://fullnode.mainnet.sui.io:443")
	gateway := "0xba477ad7b87a31fde3d29c4e4512329d7340ec23e61f130ebb4d0169ba37e189"

	balance, err := suiGetGatewaySUIBalance(ctx, client, gateway)
	require.NoError(t, err, "failed to get SUI balance for Sui gateway")

	// balance is bigger than zero
	require.NotNil(t, balance, "SUI balance should not be nil")
	require.Greater(t, balance.Cmp(big.NewInt(0)), 0, "SUI balance should be greater than zero")
}
