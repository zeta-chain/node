package testrpc

import (
	"bytes"
	"fmt"
	"net/url"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
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
		RPCParams:   "mainnet",
	}

	rpc.On("ping", func(_ map[string]any) (any, error) {
		return nil, nil
	})

	return &BtcServer{rpc}, cfg
}

func (s *BtcServer) SetBlockCount(count int) {
	s.On("getblockcount", func(_ map[string]any) (any, error) {
		return count, nil
	})
}

// OnSetRawTransaction mocks the raw transaction response.
func (s *BtcServer) OnSetRawTransaction(t *testing.T, msgTx wire.MsgTx, params ...any) {
	var buf bytes.Buffer
	err := msgTx.Serialize(&buf)
	require.NoError(t, err)

	// append the default 'verbose' parameter, otherwise the calculated params key won't match
	params = append(params, btcjson.Int(0))

	s.On("getrawtransaction", func(_ map[string]any) (any, error) {
		return hex(buf.Bytes(), false), nil
	}, params...)
}

func formatBitcoinRPCHost(serverURL string) (string, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", u.Hostname(), u.Port()), nil
}
