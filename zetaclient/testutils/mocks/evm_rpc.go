package mocks

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
	"golang.org/x/net/context"
)

const EVMRPCEnabled = "MockEVMRPCEnabled"

// Subscription interface
var _ ethereum.Subscription = subscription{}

type subscription struct {
}

func (s subscription) Unsubscribe() {
}

func (s subscription) Err() <-chan error {
	return nil
}

// EvmClient interface
var _ interfaces.EVMRPCClient = &MockEvmClient{}

type MockEvmClient struct {
	Receipts []*ethtypes.Receipt
}

func NewMockEvmClient() *MockEvmClient {
	client := &MockEvmClient{}
	return client.Reset()
}

func (e *MockEvmClient) SubscribeFilterLogs(_ context.Context, _ ethereum.FilterQuery, _ chan<- ethtypes.Log) (ethereum.Subscription, error) {
	return subscription{}, nil
}

func (e *MockEvmClient) CodeAt(_ context.Context, _ ethcommon.Address, _ *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (e *MockEvmClient) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (e *MockEvmClient) HeaderByNumber(_ context.Context, _ *big.Int) (*ethtypes.Header, error) {
	return &ethtypes.Header{}, nil
}

func (e *MockEvmClient) PendingCodeAt(_ context.Context, _ ethcommon.Address) ([]byte, error) {
	return []byte{}, nil
}

func (e *MockEvmClient) PendingNonceAt(_ context.Context, _ ethcommon.Address) (uint64, error) {
	return 0, nil
}

func (e *MockEvmClient) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (e *MockEvmClient) SuggestGasTipCap(_ context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (e *MockEvmClient) EstimateGas(_ context.Context, _ ethereum.CallMsg) (gas uint64, err error) {
	gas = 0
	err = nil
	return
}

func (e *MockEvmClient) SendTransaction(_ context.Context, _ *ethtypes.Transaction) error {
	return nil
}

func (e *MockEvmClient) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]ethtypes.Log, error) {
	return []ethtypes.Log{}, nil
}

func (e *MockEvmClient) BlockNumber(_ context.Context) (uint64, error) {
	return 0, nil
}

func (e *MockEvmClient) BlockByNumber(_ context.Context, _ *big.Int) (*ethtypes.Block, error) {
	return &ethtypes.Block{}, nil
}

func (e *MockEvmClient) TransactionByHash(_ context.Context, _ ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error) {
	return &ethtypes.Transaction{}, false, nil
}

func (e *MockEvmClient) TransactionReceipt(_ context.Context, _ ethcommon.Hash) (*ethtypes.Receipt, error) {
	// pop a receipt from the list
	if len(e.Receipts) > 0 {
		receipt := e.Receipts[len(e.Receipts)-1]
		e.Receipts = e.Receipts[:len(e.Receipts)-1]
		return receipt, nil
	}
	return nil, errors.New("no receipt found")
}

func (e *MockEvmClient) TransactionSender(_ context.Context, _ *ethtypes.Transaction, _ ethcommon.Hash, _ uint) (ethcommon.Address, error) {
	return ethcommon.Address{}, nil
}

func (e *MockEvmClient) Reset() *MockEvmClient {
	e.Receipts = []*ethtypes.Receipt{}
	return e
}

// ----------------------------------------------------------------------------
// Feed data to the mock evm client for testing
// ----------------------------------------------------------------------------
func (e *MockEvmClient) WithReceipt(receipt *ethtypes.Receipt) *MockEvmClient {
	e.Receipts = append(e.Receipts, receipt)
	return e
}

func (e *MockEvmClient) WithReceipts(receipts []*ethtypes.Receipt) *MockEvmClient {
	e.Receipts = append(e.Receipts, receipts...)
	return e
}
