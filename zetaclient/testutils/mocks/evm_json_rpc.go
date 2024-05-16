package mocks

import (
	"errors"

	"github.com/onrik/ethrpc"
	"github.com/zeta-chain/zetacore/zetaclient/chains/interfaces"
)

// EvmClient interface
var _ interfaces.EVMJSONRPCClient = &MockJSONRPCClient{}

// MockJSONRPCClient is a mock implementation of the EVMJSONRPCClient interface
type MockJSONRPCClient struct {
	Blocks       []*ethrpc.Block
	Transactions []*ethrpc.Transaction
}

// NewMockJSONRPCClient creates a new mock JSON RPC client
func NewMockJSONRPCClient() *MockJSONRPCClient {
	client := &MockJSONRPCClient{}
	return client.Reset()
}

// EthGetBlockByNumber returns a pre-loaded block or nil
func (e *MockJSONRPCClient) EthGetBlockByNumber(_ int, _ bool) (*ethrpc.Block, error) {
	// pop a block from the list
	if len(e.Blocks) > 0 {
		block := e.Blocks[len(e.Blocks)-1]
		e.Blocks = e.Blocks[:len(e.Blocks)-1]
		return block, nil
	}
	return nil, errors.New("no block found")
}

// EthGetTransactionByHash returns a pre-loaded transaction or nil
func (e *MockJSONRPCClient) EthGetTransactionByHash(_ string) (*ethrpc.Transaction, error) {
	// pop a transaction from the list
	if len(e.Transactions) > 0 {
		tx := e.Transactions[len(e.Transactions)-1]
		e.Transactions = e.Transactions[:len(e.Transactions)-1]
		return tx, nil
	}
	return nil, errors.New("no transaction found")
}

// Reset clears the mock data
func (e *MockJSONRPCClient) Reset() *MockJSONRPCClient {
	e.Blocks = []*ethrpc.Block{}
	e.Transactions = []*ethrpc.Transaction{}
	return e
}

// ----------------------------------------------------------------------------
// Feed data to the mock JSON RPC client for testing
// ----------------------------------------------------------------------------
func (e *MockJSONRPCClient) WithBlock(block *ethrpc.Block) *MockJSONRPCClient {
	e.Blocks = append(e.Blocks, block)
	return e
}

func (e *MockJSONRPCClient) WithBlocks(blocks []*ethrpc.Block) *MockJSONRPCClient {
	e.Blocks = append(e.Blocks, blocks...)
	return e
}

func (e *MockJSONRPCClient) WithTransaction(tx *ethrpc.Transaction) *MockJSONRPCClient {
	e.Transactions = append(e.Transactions, tx)
	return e
}

func (e *MockJSONRPCClient) WithTransactions(txs []*ethrpc.Transaction) *MockJSONRPCClient {
	e.Transactions = append(e.Transactions, txs...)
	return e
}
