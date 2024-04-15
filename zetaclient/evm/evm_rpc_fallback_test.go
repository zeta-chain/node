package evm

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
)

func setupTestEVMClient() *EthClientFallback {
	client1 := stub.NewMockEvmClient()
	client1.WithError(errors.New("rpc error"))
	client2 := stub.NewMockEvmClient()

	clientQ := common.NewClientQueue()
	clientQ.Append(client1)
	clientQ.Append(client2)
	return &EthClientFallback{
		evmCfg:         config.EVMConfig{},
		ethClients:     clientQ,
		jsonRPCClients: clientQ,
		logger:         zerolog.Logger{},
	}
}

func TestEthClientFallback_BlockByNumber(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.BlockByNumber(context.Background(), big.NewInt(443))
	require.NoError(t, err)
	require.Equal(t, ethtypes.Block{}, *resp)
}

func TestEthClientFallback_CallContract(t *testing.T) {
	client := setupTestEVMClient()
	_, err := client.CallContract(context.Background(), ethereum.CallMsg{}, big.NewInt(54))
	require.NoError(t, err)
}

func TestEthClientFallback_BlockNumber(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.BlockNumber(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(88), resp)
}

func TestEthClientFallback_CodeAt(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.CodeAt(context.Background(), sample.EthAddress(), big.NewInt(88))
	require.NoError(t, err)
	require.Equal(t, []byte{}, resp)
}

func TestEthClientFallback_EstimateGas(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.EstimateGas(context.Background(), ethereum.CallMsg{})
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp)
}

func TestEthClientFallback_EthGetBlockByNumber(t *testing.T) {
	client := setupTestEVMClient()
	var expected *ethrpc.Block
	expected = nil

	resp, err := client.EthGetBlockByNumber(0, false)
	require.ErrorContains(t, err, "no block found")
	require.Equal(t, expected, resp)

}

func TestEthClientFallback_EthGetTransactionByHash(t *testing.T) {
	client := setupTestEVMClient()
	var expected *ethrpc.Transaction
	expected = nil

	resp, err := client.EthGetTransactionByHash("")
	require.ErrorContains(t, err, "no transaction found")
	require.Equal(t, expected, resp)
}

func TestEthClientFallback_FilterLogs(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{})
	require.NoError(t, err)
	require.Equal(t, []ethtypes.Log{}, resp)
}

func TestEthClientFallback_HeaderByNumber(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.HeaderByNumber(context.Background(), big.NewInt(88))
	require.NoError(t, err)
	require.Equal(t, ethtypes.Header{}, *resp)
}

func TestEthClientFallback_PendingCodeAt(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.PendingCodeAt(context.Background(), sample.EthAddress())
	require.NoError(t, err)
	require.Equal(t, []byte{}, resp)
}

func TestEthClientFallback_PendingNonceAt(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.PendingNonceAt(context.Background(), sample.EthAddress())
	require.NoError(t, err)
	require.Equal(t, uint64(0), resp)
}

func TestEthClientFallback_SendTransaction(t *testing.T) {
	client := setupTestEVMClient()
	err := client.SendTransaction(context.Background(), &ethtypes.Transaction{})
	require.NoError(t, err)
}

func TestEthClientFallback_SubscribeFilterLogs(t *testing.T) {
	client := setupTestEVMClient()
	channel := make(chan<- ethtypes.Log)

	resp, err := client.SubscribeFilterLogs(context.Background(), ethereum.FilterQuery{}, channel)
	require.NoError(t, err)
	require.Equal(t, stub.Subscription{}, resp)
}

func TestEthClientFallback_SuggestGasPrice(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0), resp)
}

func TestEthClientFallback_SuggestGasTipCap(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.SuggestGasTipCap(context.Background())
	require.NoError(t, err)
	require.Equal(t, big.NewInt(0), resp)
}

func TestEthClientFallback_TransactionByHash(t *testing.T) {
	client := setupTestEVMClient()
	resp, isPending, err := client.TransactionByHash(context.Background(), sample.Hash())
	require.NoError(t, err)
	require.Equal(t, false, isPending)
	require.Equal(t, ethtypes.Transaction{}, *resp)
}

func TestEthClientFallback_TransactionReceipt(t *testing.T) {
	client := setupTestEVMClient()
	var expected *ethtypes.Receipt
	expected = nil

	resp, err := client.TransactionReceipt(context.Background(), sample.Hash())
	require.ErrorContains(t, err, "no receipt found")
	require.Equal(t, expected, resp)
}

func TestEthClientFallback_TransactionSender(t *testing.T) {
	client := setupTestEVMClient()
	resp, err := client.TransactionSender(context.Background(), &ethtypes.Transaction{}, sample.Hash(), uint(0))
	require.NoError(t, err)
	require.Equal(t, ethcommon.Address{}, resp)
}
