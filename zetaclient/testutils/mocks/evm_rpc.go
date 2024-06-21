package mocks

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/net/context"

	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
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
	err         error
	blockNumber uint64
	Receipts    []*ethtypes.Receipt
}

func NewMockEvmClient() *MockEvmClient {
	client := &MockEvmClient{}
	return client.Reset()
}

func (e *MockEvmClient) SubscribeFilterLogs(
	_ context.Context,
	_ ethereum.FilterQuery,
	_ chan<- ethtypes.Log,
) (ethereum.Subscription, error) {
	if e.err != nil {
		return subscription{}, e.err
	}
	return subscription{}, nil
}

func (e *MockEvmClient) CodeAt(_ context.Context, _ ethcommon.Address, _ *big.Int) ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}
	return []byte{}, nil
}

func (e *MockEvmClient) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}
	return []byte{}, nil
}

func (e *MockEvmClient) HeaderByNumber(_ context.Context, _ *big.Int) (*ethtypes.Header, error) {
	if e.err != nil {
		return nil, e.err
	}
	return &ethtypes.Header{}, nil
}

func (e *MockEvmClient) PendingCodeAt(_ context.Context, _ ethcommon.Address) ([]byte, error) {
	if e.err != nil {
		return nil, e.err
	}
	return []byte{}, nil
}

func (e *MockEvmClient) PendingNonceAt(_ context.Context, _ ethcommon.Address) (uint64, error) {
	if e.err != nil {
		return 0, e.err
	}
	return 0, nil
}

func (e *MockEvmClient) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	if e.err != nil {
		return nil, e.err
	}
	return big.NewInt(0), nil
}

func (e *MockEvmClient) SuggestGasTipCap(_ context.Context) (*big.Int, error) {
	if e.err != nil {
		return nil, e.err
	}
	return big.NewInt(0), nil
}

func (e *MockEvmClient) EstimateGas(_ context.Context, _ ethereum.CallMsg) (gas uint64, err error) {
	if e.err != nil {
		return 0, e.err
	}
	gas = 0
	err = nil
	return
}

func (e *MockEvmClient) SendTransaction(_ context.Context, _ *ethtypes.Transaction) error {
	return e.err
}

func (e *MockEvmClient) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]ethtypes.Log, error) {
	if e.err != nil {
		return nil, e.err
	}
	return []ethtypes.Log{}, nil
}

func (e *MockEvmClient) BlockNumber(_ context.Context) (uint64, error) {
	if e.err != nil {
		return 0, e.err
	}
	return e.blockNumber, nil
}

func (e *MockEvmClient) BlockByNumber(_ context.Context, _ *big.Int) (*ethtypes.Block, error) {
	if e.err != nil {
		return nil, e.err
	}
	return &ethtypes.Block{}, nil
}

func (e *MockEvmClient) TransactionByHash(
	_ context.Context,
	_ ethcommon.Hash,
) (tx *ethtypes.Transaction, isPending bool, err error) {
	if e.err != nil {
		return nil, false, e.err
	}
	return &ethtypes.Transaction{}, false, nil
}

func (e *MockEvmClient) TransactionReceipt(_ context.Context, _ ethcommon.Hash) (*ethtypes.Receipt, error) {
	if e.err != nil {
		return nil, e.err
	}

	// pop a receipt from the list
	if len(e.Receipts) > 0 {
		receipt := e.Receipts[len(e.Receipts)-1]
		e.Receipts = e.Receipts[:len(e.Receipts)-1]
		return receipt, nil
	}
	return nil, errors.New("no receipt found")
}

func (e *MockEvmClient) TransactionSender(
	_ context.Context,
	_ *ethtypes.Transaction,
	_ ethcommon.Hash,
	_ uint,
) (ethcommon.Address, error) {
	if e.err != nil {
		return ethcommon.Address{}, e.err
	}
	return ethcommon.Address{}, nil
}

func (e *MockEvmClient) Reset() *MockEvmClient {
	e.Receipts = []*ethtypes.Receipt{}
	return e
}

// ----------------------------------------------------------------------------
// Feed data to the mock evm client for testing
// ----------------------------------------------------------------------------
func (e *MockEvmClient) WithError(err error) *MockEvmClient {
	e.err = err
	return e
}

func (e *MockEvmClient) WithBlockNumber(blockNumber uint64) *MockEvmClient {
	e.blockNumber = blockNumber
	return e
}

func (e *MockEvmClient) WithReceipt(receipt *ethtypes.Receipt) *MockEvmClient {
	e.Receipts = append(e.Receipts, receipt)
	return e
}

func (e *MockEvmClient) WithReceipts(receipts []*ethtypes.Receipt) *MockEvmClient {
	e.Receipts = append(e.Receipts, receipts...)
	return e
}
