package testrpc

import (
	"testing"

	"github.com/zeta-chain/node/zetaclient/config"
)

// SuiServer represents httptest for Sui RPC.
type SuiServer struct {
	*Server
	Endpoint string
}

// NewSuiServer creates a new SuiServer.
func NewSuiServer(t *testing.T) (*SuiServer, config.SuiConfig) {
	rpc, endpoint := New(t, "Sui")
	cfg := config.SuiConfig{Endpoint: endpoint}

	return &SuiServer{Server: rpc, Endpoint: endpoint}, cfg
}

// todo
