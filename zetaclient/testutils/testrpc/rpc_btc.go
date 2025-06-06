package testrpc

import (
	"fmt"
	"net/url"
	"path"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
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

// CreateBTCRPCAndLoadTx is a helper function to load raw txs and feed them to mock rpc client
func CreateBTCRPCAndLoadTx(t *testing.T, dir string, chainID int64, txHashes ...string) *mocks.BitcoinClient {
	// create mock rpc client
	rpcClient := mocks.NewBitcoinClient(t)

	// feed txs to mock rpc client
	for _, txHash := range txHashes {
		// file name for the archived MsgTx
		nameMsgTx := path.Join(dir, testutils.TestDataPathBTC, testutils.FileNameBTCMsgTx(chainID, txHash))

		// load archived MsgTx
		var msgTx wire.MsgTx
		testutils.LoadObjectFromJSONFile(t, &msgTx, nameMsgTx)

		// mock rpc response
		tx := btcutil.NewTx(&msgTx)
		rpcClient.On("GetTransactionInputSpender", mock.Anything, txHash, mock.Anything, mock.Anything).Return(tx, nil)
	}

	return rpcClient
}
