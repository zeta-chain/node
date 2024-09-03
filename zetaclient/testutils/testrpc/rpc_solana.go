package testrpc

import (
	"testing"

	"github.com/zeta-chain/node/zetaclient/config"
)

// SolanaServer represents httptest for SOL RPC.
type SolanaServer struct {
	*Server
}

// NewSolanaServer creates a new SolanaServer.
func NewSolanaServer(t *testing.T) (*SolanaServer, config.SolanaConfig) {
	rpc, endpoint := New(t, "Solana")

	cfg := config.SolanaConfig{Endpoint: endpoint}

	return &SolanaServer{Server: rpc}, cfg
}
