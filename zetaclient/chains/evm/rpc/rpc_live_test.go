package rpc_test

import (
	"context"
	"math"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/chains/evm/rpc"
	"github.com/zeta-chain/node/zetaclient/common"

	"testing"
)

const (
	URLEthMainnet     = "https://rpc.ankr.com/eth"
	URLEthSepolia     = "https://rpc.ankr.com/eth_sepolia"
	URLBscMainnet     = "https://rpc.ankr.com/bsc"
	URLPolygonMainnet = "https://rpc.ankr.com/polygon"
)

// Test_EVMRPCLive is a phony test to run each live test individually
func Test_EVMRPCLive(t *testing.T) {
	if !common.LiveTestEnabled() {
		return
	}

	LiveTest_IsTxConfirmed(t)
	LiveTest_CheckRPCStatus(t)
}

func LiveTest_IsTxConfirmed(t *testing.T) {
	client, err := ethclient.Dial(URLEthMainnet)
	require.NoError(t, err)

	// check if the transaction is confirmed
	ctx := context.Background()
	txHash := "0xd2eba7ac3da1b62800165414ea4bcaf69a3b0fb9b13a0fc32f4be11bfef79146"

	t.Run("should confirm tx", func(t *testing.T) {
		confirmed, err := rpc.IsTxConfirmed(ctx, client, txHash, 12)
		require.NoError(t, err)
		require.True(t, confirmed)
	})

	t.Run("should not confirm tx if confirmations is not enough", func(t *testing.T) {
		confirmed, err := rpc.IsTxConfirmed(ctx, client, txHash, math.MaxUint64)
		require.NoError(t, err)
		require.False(t, confirmed)
	})
}

func LiveTest_CheckRPCStatus(t *testing.T) {
	client, err := ethclient.Dial(URLEthMainnet)
	require.NoError(t, err)

	ctx := context.Background()
	_, err = rpc.CheckRPCStatus(ctx, client)
	require.NoError(t, err)
}
