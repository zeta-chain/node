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

func hex(v any) string {
	return fmt.Sprintf("0x%x", v)
}
