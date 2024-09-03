package mocks

import (
	"errors"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"

	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

// EvmClient interface
var _ interfaces.BTCRPCClient = &MockBTCRPCClient{}

// MockBTCRPCClient is a mock implementation of the BTCRPCClient interface
type MockBTCRPCClient struct {
	err            error
	blockCount     int64
	blockHash      *chainhash.Hash
	blockHeader    *wire.BlockHeader
	blockVerboseTx *btcjson.GetBlockVerboseTxResult
	Txs            []*btcutil.Tx
}

// NewMockBTCRPCClient creates a new mock BTC RPC client
func NewMockBTCRPCClient() *MockBTCRPCClient {
	client := &MockBTCRPCClient{}
	return client.Reset()
}

// Reset clears the mock data
func (c *MockBTCRPCClient) Reset() *MockBTCRPCClient {
	if c.err != nil {
		return nil
	}

	c.Txs = []*btcutil.Tx{}
	return c
}

func (c *MockBTCRPCClient) GetNetworkInfo() (*btcjson.GetNetworkInfoResult, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) CreateWallet(_ string, _ ...rpcclient.CreateWalletOpt) (*btcjson.CreateWalletResult, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) GetNewAddress(_ string) (btcutil.Address, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) GenerateToAddress(_ int64, _ btcutil.Address, _ *int64) ([]*chainhash.Hash, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) GetBalance(_ string) (btcutil.Amount, error) {
	return 0, errors.New("not implemented")
}

func (c *MockBTCRPCClient) SendRawTransaction(_ *wire.MsgTx, _ bool) (*chainhash.Hash, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) ListUnspent() ([]btcjson.ListUnspentResult, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) ListUnspentMinMaxAddresses(
	_ int,
	_ int,
	_ []btcutil.Address,
) ([]btcjson.ListUnspentResult, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) EstimateSmartFee(
	_ int64,
	_ *btcjson.EstimateSmartFeeMode,
) (*btcjson.EstimateSmartFeeResult, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) GetTransaction(_ *chainhash.Hash) (*btcjson.GetTransactionResult, error) {
	return nil, errors.New("not implemented")
}

// GetRawTransaction returns a pre-loaded transaction or nil
func (c *MockBTCRPCClient) GetRawTransaction(_ *chainhash.Hash) (*btcutil.Tx, error) {
	// pop a transaction from the list
	if len(c.Txs) > 0 {
		tx := c.Txs[len(c.Txs)-1]
		c.Txs = c.Txs[:len(c.Txs)-1]
		return tx, nil
	}
	return nil, errors.New("no transaction found")
}

func (c *MockBTCRPCClient) GetRawTransactionVerbose(_ *chainhash.Hash) (*btcjson.TxRawResult, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) GetBlockCount() (int64, error) {
	if c.err != nil {
		return 0, c.err
	}
	return c.blockCount, nil
}

func (c *MockBTCRPCClient) GetBlockHash(_ int64) (*chainhash.Hash, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.blockHash, nil
}

func (c *MockBTCRPCClient) GetBlockVerbose(_ *chainhash.Hash) (*btcjson.GetBlockVerboseResult, error) {
	return nil, errors.New("not implemented")
}

func (c *MockBTCRPCClient) GetBlockVerboseTx(_ *chainhash.Hash) (*btcjson.GetBlockVerboseTxResult, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.blockVerboseTx, nil
}

func (c *MockBTCRPCClient) GetBlockHeader(_ *chainhash.Hash) (*wire.BlockHeader, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.blockHeader, nil
}

// ----------------------------------------------------------------------------
// Feed data to the mock BTC RPC client for testing
// ----------------------------------------------------------------------------

func (c *MockBTCRPCClient) WithError(err error) *MockBTCRPCClient {
	c.err = err
	return c
}

func (c *MockBTCRPCClient) WithBlockCount(blkCnt int64) *MockBTCRPCClient {
	c.blockCount = blkCnt
	return c
}

func (c *MockBTCRPCClient) WithBlockHash(hash *chainhash.Hash) *MockBTCRPCClient {
	c.blockHash = hash
	return c
}

func (c *MockBTCRPCClient) WithBlockHeader(header *wire.BlockHeader) *MockBTCRPCClient {
	c.blockHeader = header
	return c
}

func (c *MockBTCRPCClient) WithBlockVerboseTx(block *btcjson.GetBlockVerboseTxResult) *MockBTCRPCClient {
	c.blockVerboseTx = block
	return c
}

func (c *MockBTCRPCClient) WithRawTransaction(tx *btcutil.Tx) *MockBTCRPCClient {
	c.Txs = append(c.Txs, tx)
	return c
}

func (c *MockBTCRPCClient) WithRawTransactions(txs []*btcutil.Tx) *MockBTCRPCClient {
	c.Txs = append(c.Txs, txs...)
	return c
}
