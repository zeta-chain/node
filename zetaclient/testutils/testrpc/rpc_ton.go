package testrpc

import (
	"testing"

	"github.com/zeta-chain/node/zetaclient/config"
)

// TONServer represents httptest for TON RPC.
type TONServer struct {
	*Server
	Endpoint string
}

// NewSuiServer creates a new SuiServer.
func NewTONServer(t *testing.T) (*TONServer, config.TONConfig) {
	rpc, endpoint := New(t, "TON")
	cfg := config.TONConfig{Endpoint: endpoint}

	return &TONServer{Server: rpc, Endpoint: endpoint}, cfg
}
