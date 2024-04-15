package bitcoin

import (
	"errors"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
)

func setupBTCFallbackClient() *RPCClientFallback {
	client1 := stub.NewMockBTCRPCClient()
	client1.WithError(errors.New("rpc error"))

	client2 := stub.NewMockBTCRPCClient()
	client2.WithError(nil)

	clientq := common.NewClientQueue()
	clientq.Append(client1)
	clientq.Append(client2)

	return &RPCClientFallback{
		btcConfig:  config.BTCConfig{},
		rpcClients: clientq,
		logger:     zerolog.Logger{},
	}
}

func TestRPCClientFallback_CreateWallet(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.CreateWallet("testWallet")
	require.NoError(t, err)
}

func TestRPCClientFallback_EstimateSmartFee(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.EstimateSmartFee(0, &btcjson.EstimateModeUnset)
	require.NoError(t, err)
}

func TestRPCClientFallback_GenerateToAddress(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GenerateToAddress(0, &chains.AddressTaproot{}, nil)
	require.NoError(t, err)
}

func TestRPCClientFallback_GetBalance(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetBalance("")
	require.NoError(t, err)
}

func TestRPCClientFallback_GetBlockCount(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetBlockCount()
	require.NoError(t, err)
}

func TestRPCClientFallback_GetBlockHash(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetBlockHash(0)
	require.NoError(t, err)
}

func TestRPCClientFallback_GetBlockHeader(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetBlockHeader(&chainhash.Hash{})
	require.NoError(t, err)
}

func TestRPCClientFallback_GetBlockVerbose(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetBlockVerbose(&chainhash.Hash{})
	require.NoError(t, err)
}

func TestRPCClientFallback_GetBlockVerboseTx(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetBlockVerboseTx(&chainhash.Hash{})
	require.NoError(t, err)
}

func TestRPCClientFallback_GetNetworkInfo(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetNetworkInfo()
	require.NoError(t, err)
}

func TestRPCClientFallback_GetNewAddress(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetNewAddress("Test")
	require.NoError(t, err)
}

func TestRPCClientFallback_GetRawTransaction(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetRawTransaction(&chainhash.Hash{})
	require.NoError(t, err)
}

func TestRPCClientFallback_GetRawTransactionVerbose(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetRawTransactionVerbose(&chainhash.Hash{})
	require.NoError(t, err)
}

func TestRPCClientFallback_GetTransaction(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.GetTransaction(&chainhash.Hash{})
	require.NoError(t, err)
}

func TestRPCClientFallback_ListUnspentMinMaxAddresses(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.ListUnspentMinMaxAddresses(0, 0, []btcutil.Address{})
	require.NoError(t, err)
}

func TestRPCClientFallback_ListUnspent(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.ListUnspent()
	require.NoError(t, err)
}

func TestRPCClientFallback_SendRawTransaction(t *testing.T) {
	rpcClientFallback := setupBTCFallbackClient()
	_, err := rpcClientFallback.SendRawTransaction(&wire.MsgTx{}, false)
	require.NoError(t, err)
}
