package testrpc

import (
	"fmt"
	"testing"
)

// EVMServer represents httptest for EVM RPC.
type EVMServer struct {
	*Server
	Endpoint string
}

// NewEVMServer creates a new EVMServer.
func NewEVMServer(t *testing.T) *EVMServer {
	rpc, endpoint := New(t, "EVM")

	return &EVMServer{Server: rpc, Endpoint: endpoint}
}

func (s *EVMServer) SetBlockNumber(n int) {
	s.On("eth_blockNumber", func(_ []any) (any, error) {
		return hex(n), nil
	})
}

func (s *EVMServer) SetBlockNumberFailure(err error) {
	s.On("eth_blockNumber", func(_ []any) (any, error) {
		return hex(0), err
	})
}

func (s *EVMServer) MockSendTransaction() {
	s.On("eth_sendRawTransaction", func(_ []any) (any, error) {
		return nil, nil
	})
}

func (s *EVMServer) MockNonceAt(nonce uint64) {
	s.On("eth_getTransactionCount", func(_ []any) (any, error) {
		return hex(nonce), nil
	})
}

func (s *EVMServer) SetChainID(n int) {
	s.On("eth_chainId", func(_ []any) (any, error) {
		return hex(n), nil
	})
}

func hex(v any) string {
	return fmt.Sprintf("%x", v)
}
