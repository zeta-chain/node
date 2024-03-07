package stub

import (
	"math/big"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
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
var _ interfaces.EVMRPCClient = EvmClient{}

type EvmClient struct {
}

func (e EvmClient) SubscribeFilterLogs(_ context.Context, _ ethereum.FilterQuery, _ chan<- ethtypes.Log) (ethereum.Subscription, error) {
	return subscription{}, nil
}

func (e EvmClient) CodeAt(_ context.Context, _ ethcommon.Address, _ *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (e EvmClient) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return []byte{}, nil
}

func (e EvmClient) HeaderByNumber(_ context.Context, _ *big.Int) (*ethtypes.Header, error) {
	return &ethtypes.Header{}, nil
}

func (e EvmClient) PendingCodeAt(_ context.Context, _ ethcommon.Address) ([]byte, error) {
	return []byte{}, nil
}

func (e EvmClient) PendingNonceAt(_ context.Context, _ ethcommon.Address) (uint64, error) {
	return 0, nil
}

func (e EvmClient) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (e EvmClient) SuggestGasTipCap(_ context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (e EvmClient) EstimateGas(_ context.Context, _ ethereum.CallMsg) (gas uint64, err error) {
	gas = 0
	err = nil
	return
}

func (e EvmClient) SendTransaction(_ context.Context, _ *ethtypes.Transaction) error {
	return nil
}

func (e EvmClient) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]ethtypes.Log, error) {
	return []ethtypes.Log{}, nil
}

func (e EvmClient) BlockNumber(_ context.Context) (uint64, error) {
	return 0, nil
}

func (e EvmClient) BlockByNumber(_ context.Context, _ *big.Int) (*ethtypes.Block, error) {
	return &ethtypes.Block{}, nil
}

func (e EvmClient) TransactionByHash(_ context.Context, _ ethcommon.Hash) (tx *ethtypes.Transaction, isPending bool, err error) {
	return &ethtypes.Transaction{}, false, nil
}

func (e EvmClient) TransactionReceipt(_ context.Context, _ ethcommon.Hash) (*ethtypes.Receipt, error) {
	return &ethtypes.Receipt{}, nil
}

func (e EvmClient) TransactionSender(_ context.Context, _ *ethtypes.Transaction, _ ethcommon.Hash, _ uint) (ethcommon.Address, error) {
	return ethcommon.Address{}, nil
}
