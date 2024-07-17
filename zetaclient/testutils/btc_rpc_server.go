package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zeta-chain/zetacore/zetaclient/config"
)

// BtcServer represents a BTC RPC mock with a "real" HTTP server allocated.
type BtcServer struct {
	t *testing.T
}

// NewBtcServer constructs BtcServer.
func NewBtcServer(t *testing.T) (*BtcServer, config.BTCConfig) {
	var (
		btcServer = &BtcServer{t: t}
		server    = httptest.NewUnstartedServer(http.HandlerFunc(btcServer.handler))
		cfg       = config.BTCConfig{
			RPCUsername: "btc-user",
			RPCPassword: "btc-password",
			RPCHost:     server.URL,
			RPCParams:   "",
		}
	)

	server.Start()
	t.Cleanup(server.Close)

	return btcServer, cfg
}

// handler is a simple HTTP handler that returns 200 OK.
// Later we can add any logic here.
func (s BtcServer) handler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
