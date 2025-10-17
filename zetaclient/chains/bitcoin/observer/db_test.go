package observer

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
)

func Test_SaveBroadcastedTx(t *testing.T) {
	tests := []struct {
		name    string
		wantErr string
	}{
		{
			name:    "should be able to save broadcasted tx",
			wantErr: "",
		},
		{
			name:    "should fail on db error",
			wantErr: "failed to save broadcasted outbound hash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// test data
			nonce := uint64(1)
			txHash := sample.BtcHash().String()
			dbPath := sample.CreateTempDir(t)
			ob := newTestSuite(t, chains.BitcoinMainnet, withDatabasePath(dbPath))
			if tt.wantErr != "" {
				// delete db to simulate db error
				os.RemoveAll(dbPath)
			}

			// ACT
			// save a test tx
			err := ob.SaveBroadcastedTx(txHash, nonce)

			// ASSERT
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			// should always save broadcasted outbound to memory
			gotHash, found := ob.GetBroadcastedTx(nonce)
			require.True(t, found)
			require.Equal(t, txHash, gotHash)
			require.True(t, ob.IsTSSTransaction(txHash))
		})
	}
}

func Test_LoadLastBlockScanned(t *testing.T) {
	ctx := context.Background()

	// use Bitcoin mainnet chain for testing
	chain := chains.BitcoinMainnet

	t.Run("should load last block scanned", func(t *testing.T) {
		// create observer and write 199 as last block scanned
		ob := newTestSuite(t, chain)
		ob.WriteLastBlockScannedToDB(199)

		// load last block scanned
		err := ob.LoadLastBlockScanned(ctx)
		require.NoError(t, err)
		require.EqualValues(t, 199, ob.LastBlockScanned())
	})
	t.Run("should fail on invalid env var", func(t *testing.T) {
		// create observer
		ob := newTestSuite(t, chain)

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
		obOther := newTestSuite(t, chain)

		// reset last block scanned to 0 so that it will be loaded from RPC
		obOther.WithLastBlockScanned(0)

		// attach a mock btc client that returns rpc error
		obOther.client.ExpectedCalls = nil
		obOther.client.On("GetBlockCount", mock.Anything).Return(int64(0), errors.New("rpc error"))

		// load last block scanned
		err := obOther.LoadLastBlockScanned(ctx)
		require.ErrorContains(t, err, "unable to get block count")
	})
	t.Run("should use hardcode block 100 for regtest", func(t *testing.T) {
		// use regtest chain
		obRegnet := newTestSuite(t, chains.BitcoinRegtest)

		// load last block scanned
		err := obRegnet.LoadLastBlockScanned(ctx)
		require.NoError(t, err)
		require.EqualValues(t, RegnetStartBlock, obRegnet.LastBlockScanned())
	})
}

func Test_LoadBroadcastedTxMap(t *testing.T) {
	t.Run("should load broadcasted tx map", func(t *testing.T) {
		// test data
		nonce := uint64(1)
		txHash := sample.BtcHash().String()

		// create observer and save a test tx
		dbPath := sample.CreateTempDir(t)
		obOld := newTestSuite(t, chains.BitcoinMainnet, withDatabasePath(dbPath))
		obOld.SaveBroadcastedTx(txHash, nonce)

		// create new observer using same db path
		obNew := newTestSuite(t, chains.BitcoinMainnet, withDatabasePath(dbPath))

		// check if the txHash is a TSS outbound
		require.True(t, obNew.IsTSSTransaction(txHash))

		// get the broadcasted tx
		gotHash, found := obNew.GetBroadcastedTx(nonce)
		require.True(t, found)
		require.Equal(t, txHash, gotHash)
	})
}
