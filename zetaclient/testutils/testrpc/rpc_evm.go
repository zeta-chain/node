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
	s.On("eth_blockNumber", func(_ map[string]any) (any, error) {
		return hex(n, true), nil
	})
}

func (s *EVMServer) SetBlockNumberFailure(err error) {
	s.On("eth_blockNumber", func(_ map[string]any) (any, error) {
		return hex(0, true), err
	})
}

func (s *EVMServer) MockSendTransaction() {
	s.On("eth_sendRawTransaction", func(_ map[string]any) (any, error) {
		return nil, nil
	})
}

func (s *EVMServer) MockNonceAt(nonce uint64) {
	s.On("eth_getTransactionCount", func(_ map[string]any) (any, error) {
		return hex(nonce, true), nil
	})
}

func (s *EVMServer) SetChainID(n int) {
	s.On("eth_chainId", func(_ map[string]any) (any, error) {
		return hex(n, true), nil
	})
}

func hex(v any, prefix bool) string {
	if prefix {
		return fmt.Sprintf("0x%x", v)
	}
	return fmt.Sprintf("%x", v)
}
