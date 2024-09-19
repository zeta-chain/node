package testrpc

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/zetaclient/config"
)

// BtcServer represents httptest for Bitcoin RPC.
type BtcServer struct {
	*Server
}

// NewBtcServer creates new BtcServer.
func NewBtcServer(t *testing.T) (*BtcServer, config.BTCConfig) {
	rpc, rpcURL := New(t, "bitcoin")

	host, err := formatBitcoinRPCHost(rpcURL)
	require.NoError(t, err)

	cfg := config.BTCConfig{
		RPCUsername: "btc-user",
		RPCPassword: "btc-password",
		RPCHost:     host,
		RPCParams:   "",
	}

	rpc.On("ping", func(_ []any) (any, error) {
		return nil, nil
	})

	return &BtcServer{rpc}, cfg
}

func (s *BtcServer) SetBlockCount(count int) {
	s.On("getblockcount", func(_ []any) (any, error) {
		return count, nil
	})
}

func formatBitcoinRPCHost(serverURL string) (string, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", u.Hostname(), u.Port()), nil
}
