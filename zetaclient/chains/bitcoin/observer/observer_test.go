package observer_test

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	lru "github.com/hashicorp/golang-lru"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/db"
	"gorm.io/gorm"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

var (
	// the relative path to the testdata directory
	TestDataDir = "../../../"
)

// setupDBTxResults creates a new SQLite database and populates it with some transaction results.
func setupDBTxResults(t *testing.T) (*gorm.DB, map[string]btcjson.GetTransactionResult) {
	submittedTx := map[string]btcjson.GetTransactionResult{}

	database, err := db.NewFromSqliteInMemory(true)
	require.NoError(t, err)

	//Create some Transaction entries in the DB
	for i := 0; i < 2; i++ {
		txResult := btcjson.GetTransactionResult{
			Amount:          float64(i),
			Fee:             0,
			Confirmations:   0,
			BlockHash:       "",
			BlockIndex:      0,
			BlockTime:       0,
			TxID:            strconv.Itoa(i),
			WalletConflicts: nil,
			Time:            0,
			TimeReceived:    0,
			Details:         nil,
			Hex:             "",
		}
		r, _ := clienttypes.ToTransactionResultSQLType(txResult, strconv.Itoa(i))
		dbc := database.Client().Create(&r)
		require.NoError(t, dbc.Error)
		submittedTx[strconv.Itoa(i)] = txResult
	}

	return database.Client(), submittedTx
}

// MockBTCObserver creates a mock Bitcoin observer for testing
func MockBTCObserver(
	t *testing.T,
	chain chains.Chain,
	params observertypes.ChainParams,
	btcClient interfaces.BTCRPCClient,
) *observer.Observer {
	// use default mock btc client if not provided
	if btcClient == nil {
		btcClient = mocks.NewMockBTCRPCClient().WithBlockCount(100)
	}

	database, err := db.NewFromSqliteInMemory(true)

	// create observer
	ob, err := observer.NewObserver(
		chain,
		btcClient,
		params,
		nil,
		nil,
		60,
		database,
		base.Logger{},
		nil,
	)
	require.NoError(t, err)

	return ob
}

func Test_NewObserver(t *testing.T) {
	// use Bitcoin mainnet chain for testing
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)

	// test cases
	tests := []struct {
		name         string
		chain        chains.Chain
		btcClient    interfaces.BTCRPCClient
		chainParams  observertypes.ChainParams
		coreClient   interfaces.ZetacoreClient
		tss          interfaces.TSSSigner
		logger       base.Logger
		ts           *metrics.TelemetryServer
		errorMessage string
		before       func()
		after        func()
	}{
		{
			name:        "should be able to create observer",
			chain:       chain,
			btcClient:   mocks.NewMockBTCRPCClient().WithBlockCount(100),
			chainParams: params,
			coreClient:  nil,
			tss:         mocks.NewTSSMainnet(),
		},
		{
			name:         "should fail if net params is not found",
			chain:        chains.Chain{ChainId: 111}, // invalid chain id
			btcClient:    mocks.NewMockBTCRPCClient().WithBlockCount(100),
			chainParams:  params,
			coreClient:   nil,
			tss:          mocks.NewTSSMainnet(),
			errorMessage: "unable to get BTC net params for chain",
		},
		{
			name:        "should fail if env var us invalid",
			chain:       chain,
			btcClient:   mocks.NewMockBTCRPCClient().WithBlockCount(100),
			chainParams: params,
			coreClient:  nil,
			tss:         mocks.NewTSSMainnet(),
			before: func() {
				envVar := base.EnvVarLatestBlockByChain(chain)
				os.Setenv(envVar, "invalid")
			},
			after: func() {
				envVar := base.EnvVarLatestBlockByChain(chain)
				os.Unsetenv(envVar)
			},
			errorMessage: "unable to parse block number from ENV",
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create db
			database, err := db.NewFromSqliteInMemory(true)
			require.NoError(t, err)

			if tt.before != nil {
				tt.before()
			}

			if tt.after != nil {
				defer tt.after()
			}

			// create observer
			ob, err := observer.NewObserver(
				tt.chain,
				tt.btcClient,
				tt.chainParams,
				tt.coreClient,
				tt.tss,
				60,
				database,
				tt.logger,
				tt.ts,
			)

			if tt.errorMessage != "" {
				require.ErrorContains(t, err, tt.errorMessage)
				require.Nil(t, ob)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, ob)
		})
	}
}

func Test_BlockCache(t *testing.T) {
	t.Run("should add and get block from cache", func(t *testing.T) {
		// create observer
		ob := &observer.Observer{}
		blockCache, err := lru.New(100)
		require.NoError(t, err)
		ob.WithBlockCache(blockCache)

		// create mock btc client
		btcClient := mocks.NewMockBTCRPCClient()
		ob.WithBtcClient(btcClient)

		// feed block hash, header and block to btc client
		hash := sample.BtcHash()
		header := &wire.BlockHeader{Version: 1}
		block := &btcjson.GetBlockVerboseTxResult{Version: 1}
		btcClient.WithBlockHash(&hash)
		btcClient.WithBlockHeader(header)
		btcClient.WithBlockVerboseTx(block)

		// get block and header from observer, fallback to btc client
		result, err := ob.GetBlockByNumberCached(100)
		require.NoError(t, err)
		require.EqualValues(t, header, result.Header)
		require.EqualValues(t, block, result.Block)

		// get block header from cache
		result, err = ob.GetBlockByNumberCached(100)
		require.NoError(t, err)
		require.EqualValues(t, header, result.Header)
		require.EqualValues(t, block, result.Block)
	})
	t.Run("should fail if stored type is not BlockNHeader", func(t *testing.T) {
		// create observer
		ob := &observer.Observer{}
		blockCache, err := lru.New(100)
		require.NoError(t, err)
		ob.WithBlockCache(blockCache)

		// add a string to cache
		blockNumber := int64(100)
		blockCache.Add(blockNumber, "a string value")

		// get result from cache
		result, err := ob.GetBlockByNumberCached(blockNumber)
		require.ErrorContains(t, err, "cached value is not of type *BTCBlockNHeader")
		require.Nil(t, result)
	})
}

func Test_LoadLastBlockScanned(t *testing.T) {
	// use Bitcoin mainnet chain for testing
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)

	// create observer using mock btc client
	btcClient := mocks.NewMockBTCRPCClient().WithBlockCount(200)

	t.Run("should load last block scanned", func(t *testing.T) {
		// create observer and write 199 as last block scanned
		ob := MockBTCObserver(t, chain, params, btcClient)
		ob.WriteLastBlockScannedToDB(199)

		// load last block scanned
		err := ob.LoadLastBlockScanned()
		require.NoError(t, err)
		require.EqualValues(t, 199, ob.LastBlockScanned())
	})
	t.Run("should fail on invalid env var", func(t *testing.T) {
		// create observer
		ob := MockBTCObserver(t, chain, params, btcClient)

		// set invalid environment variable
		envvar := base.EnvVarLatestBlockByChain(chain)
		os.Setenv(envvar, "invalid")
		defer os.Unsetenv(envvar)

		// load last block scanned
		err := ob.LoadLastBlockScanned()
		require.ErrorContains(t, err, "error LoadLastBlockScanned")
	})
	t.Run("should fail on RPC error", func(t *testing.T) {
		// create observer on separate path, as we need to reset last block scanned
		obOther := MockBTCObserver(t, chain, params, btcClient)

		// reset last block scanned to 0 so that it will be loaded from RPC
		obOther.WithLastBlockScanned(0)

		// set RPC error
		btcClient.WithError(fmt.Errorf("error RPC"))

		// load last block scanned
		err := obOther.LoadLastBlockScanned()
		require.ErrorContains(t, err, "error RPC")
	})
	t.Run("should use hardcode block 100 for regtest", func(t *testing.T) {
		// use regtest chain
		regtest := chains.BitcoinRegtest
		obRegnet := MockBTCObserver(t, regtest, params, nil)

		// load last block scanned
		err := obRegnet.LoadLastBlockScanned()
		require.NoError(t, err)
		require.EqualValues(t, observer.RegnetStartBlock, obRegnet.LastBlockScanned())
	})
}

func TestConfirmationThreshold(t *testing.T) {
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)
	ob := MockBTCObserver(t, chain, params, nil)

	t.Run("should return confirmations in chain param", func(t *testing.T) {
		ob.SetChainParams(observertypes.ChainParams{ConfirmationCount: 3})
		require.Equal(t, int64(3), ob.ConfirmationsThreshold(big.NewInt(1000)))
	})

	t.Run("should return big value confirmations", func(t *testing.T) {
		ob.SetChainParams(observertypes.ChainParams{ConfirmationCount: 3})
		require.Equal(
			t,
			int64(observer.BigValueConfirmationCount),
			ob.ConfirmationsThreshold(big.NewInt(observer.BigValueSats)),
		)
	})

	t.Run("big value confirmations is the upper cap", func(t *testing.T) {
		ob.SetChainParams(observertypes.ChainParams{ConfirmationCount: observer.BigValueConfirmationCount + 1})
		require.Equal(t, int64(observer.BigValueConfirmationCount), ob.ConfirmationsThreshold(big.NewInt(1000)))
	})
}

func TestSubmittedTx(t *testing.T) {
	// setup db
	db, submittedTx := setupDBTxResults(t)

	var submittedTransactions []clienttypes.TransactionResultSQLType
	err := db.Find(&submittedTransactions).Error
	require.NoError(t, err)

	for _, txResult := range submittedTransactions {
		r, err := clienttypes.FromTransactionResultSQLType(txResult)
		require.NoError(t, err)
		want := submittedTx[txResult.Key]
		have := r

		require.Equal(t, want, have)
	}
}
