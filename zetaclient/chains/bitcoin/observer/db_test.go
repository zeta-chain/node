package observer_test

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
)

func Test_SaveBroadcastedTx(t *testing.T) {
	t.Run("should be able to save broadcasted tx", func(t *testing.T) {
		// test data
		nonce := uint64(1)
		txHash := sample.BtcHash().String()

		// create observer and open db
		ob := newTestSuite(t, chains.BitcoinMainnet, "")

		// save a test tx
		ob.SaveBroadcastedTx(txHash, nonce)

		// check if the txHash is a TSS outbound
		require.True(t, ob.IsTSSTransaction(txHash))

		// get the broadcasted tx
		gotHash, found := ob.GetBroadcastedTx(nonce)
		require.True(t, found)
		require.Equal(t, txHash, gotHash)
	})
}

func Test_LoadLastBlockScanned(t *testing.T) {
	ctx := context.Background()

	// use Bitcoin mainnet chain for testing
	chain := chains.BitcoinMainnet

	t.Run("should load last block scanned", func(t *testing.T) {
		// create observer and write 199 as last block scanned
		ob := newTestSuite(t, chain, "")
		ob.WriteLastBlockScannedToDB(199)

		// load last block scanned
		err := ob.LoadLastBlockScanned(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 199, ob.LastBlockScanned())
	})
	t.Run("should fail on invalid env var", func(t *testing.T) {
		// create observer
		ob := newTestSuite(t, chain, "")

		// set invalid environment variable
		envvar := base.EnvVarLatestBlockByChain(chain)
		os.Setenv(envvar, "invalid")
		defer os.Unsetenv(envvar)

		// load last block scanned
		err := ob.LoadLastBlockScanned(ctx)
		require.ErrorContains(t, err, "error LoadLastBlockScanned")
	})
	t.Run("should fail on RPC error", func(t *testing.T) {
		// create observer on separate path, as we need to reset last block scanned
		obOther := newTestSuite(t, chain, "")

		// reset last block scanned to 0 so that it will be loaded from RPC
		obOther.WithLastBlockScanned(0)

		// attach a mock btc client that returns rpc error
		obOther.client.ExpectedCalls = nil
		obOther.client.On("GetBlockCount", mock.Anything).Return(int64(0), errors.New("rpc error"))

		// load last block scanned
		err := obOther.LoadLastBlockScanned(ctx)
		require.ErrorContains(t, err, "rpc error")
	})
	t.Run("should use hardcode block 100 for regtest", func(t *testing.T) {
		// use regtest chain
		obRegnet := newTestSuite(t, chains.BitcoinRegtest, "")

		// load last block scanned
		err := obRegnet.LoadLastBlockScanned(ctx)
		require.NoError(t, err)
		require.EqualValues(t, observer.RegnetStartBlock, obRegnet.LastBlockScanned())
	})
}

func Test_LoadBroadcastedTxMap(t *testing.T) {
	t.Run("should load broadcasted tx map", func(t *testing.T) {
		// test data
		nonce := uint64(1)
		txHash := sample.BtcHash().String()

		// create observer and save a test tx
		dbPath := sample.CreateTempDir(t)
		obOld := newTestSuite(t, chains.BitcoinMainnet, dbPath)
		obOld.SaveBroadcastedTx(txHash, nonce)

		// create new observer using same db path
		obNew := newTestSuite(t, chains.BitcoinMainnet, dbPath)

		// load broadcasted tx map to new observer
		err := obNew.LoadBroadcastedTxMap()
		require.NoError(t, err)

		// check if the txHash is a TSS outbound
		require.True(t, obNew.IsTSSTransaction(txHash))

		// get the broadcasted tx
		gotHash, found := obNew.GetBroadcastedTx(nonce)
		require.True(t, found)
		require.Equal(t, txHash, gotHash)
	})
}
