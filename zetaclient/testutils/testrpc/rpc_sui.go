package testrpc

import (
	"testing"

	"github.com/zeta-chain/node/zetaclient/config"
)

// SUIServer represents httptest for SUI RPC.
type SUIServer struct {
	*Server
	Endpoint string
}

// SUIServer creates a new SUIServer.
func NewSUIServer(t *testing.T) (*SUIServer, config.SUIConfig) {
	rpc, endpoint := New(t, "SUI")
	cfg := config.SUIConfig{Endpoint: endpoint}

	return &SUIServer{Server: rpc, Endpoint: endpoint}, cfg
}

// todo
