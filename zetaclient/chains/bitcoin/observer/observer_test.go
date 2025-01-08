package observer_test

import (
	"math/big"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/testutils"
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

func Test_NewObserver(t *testing.T) {
	// use Bitcoin mainnet chain for testing
	chain := chains.BitcoinMainnet
	params := mocks.MockChainParams(chain.ChainId, 10)

	// create mock btc client with block height 100
	btcClient := mocks.NewBTCRPCClient(t)
	btcClient.On("GetBlockCount").Return(int64(100), nil)

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
			btcClient:   btcClient,
			chainParams: params,
			coreClient:  nil,
			tss:         mocks.NewTSS(t),
		},
		{
			name:         "should fail if net params is not found",
			chain:        chains.Chain{ChainId: 111}, // invalid chain id
			btcClient:    btcClient,
			chainParams:  params,
			coreClient:   nil,
			tss:          mocks.NewTSS(t),
			errorMessage: "unable to get BTC net params for chain",
		},
		{
			name:        "should fail if env var us invalid",
			chain:       chain,
			btcClient:   btcClient,
			chainParams: params,
			coreClient:  nil,
			tss:         mocks.NewTSS(t),
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
		ob := newTestSuite(t, chains.BitcoinMainnet, "")

		// feed block hash, header and block to btc client
		hash := sample.BtcHash()
		header := &wire.BlockHeader{Version: 1}
		block := &btcjson.GetBlockVerboseTxResult{Version: 1}
		ob.client.On("GetBlockHash", mock.Anything).Return(&hash, nil)
		ob.client.On("GetBlockHeader", &hash).Return(header, nil)
		ob.client.On("GetBlockVerboseTx", &hash).Return(block, nil)

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
		ob := newTestSuite(t, chains.BitcoinMainnet, "")

		// add a string to cache
		blockNumber := int64(100)
		ob.BlockCache().Add(blockNumber, "a string value")

		// get result from cache
		result, err := ob.GetBlockByNumberCached(blockNumber)
		require.ErrorContains(t, err, "cached value is not of type *BTCBlockNHeader")
		require.Nil(t, result)
	})
}

func Test_SetPendingNonce(t *testing.T) {
	// create observer
	ob := newTestSuite(t, chains.BitcoinMainnet, "")

	// ensure pending nonce is 0
	require.Zero(t, ob.GetPendingNonce())

	// set and get pending nonce
	nonce := uint64(100)
	ob.SetPendingNonce(nonce)
	require.Equal(t, nonce, ob.GetPendingNonce())
}

func TestConfirmationThreshold(t *testing.T) {
	chain := chains.BitcoinMainnet
	ob := newTestSuite(t, chain, "")

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

func Test_SetLastStuckOutbound(t *testing.T) {
	// create observer and example stuck tx
	ob := newTestSuite(t, chains.BitcoinMainnet, "")
	tx := btcutil.NewTx(wire.NewMsgTx(wire.TxVersion))

	// STEP 1
	// initial stuck outbound is nil
	require.Nil(t, ob.GetLastStuckOutbound())

	// STEP 2
	// set stuck outbound
	stuckTx := observer.NewLastStuckOutbound(100, tx, 30*time.Minute)
	ob.SetLastStuckOutbound(stuckTx)

	// retrieve stuck outbound
	require.Equal(t, stuckTx, ob.GetLastStuckOutbound())

	// STEP 3
	// update stuck outbound
	stuckTxUpdate := observer.NewLastStuckOutbound(101, tx, 40*time.Minute)
	ob.SetLastStuckOutbound(stuckTxUpdate)

	// retrieve updated stuck outbound
	require.Equal(t, stuckTxUpdate, ob.GetLastStuckOutbound())

	// STEP 4
	// clear stuck outbound
	ob.SetLastStuckOutbound(nil)

	// stuck outbound should be nil
	require.Nil(t, ob.GetLastStuckOutbound())
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

type testSuite struct {
	*observer.Observer

	client   *mocks.BTCRPCClient
	zetacore *mocks.ZetacoreClient
	db       *db.DB
}

func newTestSuite(t *testing.T, chain chains.Chain, dbPath string) *testSuite {
	require.True(t, chain.IsBitcoinChain())

	chainParams := mocks.MockChainParams(chain.ChainId, 10)

	client := mocks.NewBTCRPCClient(t)
	client.On("GetBlockCount").Return(int64(100), nil).Maybe()

	zetacore := mocks.NewZetacoreClient(t)

	var tss interfaces.TSSSigner
	if chains.IsBitcoinMainnet(chain.ChainId) {
		tss = mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet)
	} else {
		tss = mocks.NewTSS(t).FakePubKey(testutils.TSSPubkeyAthens3)
	}

	// create test database
	var err error
	var database *db.DB
	if dbPath == "" {
		database, err = db.NewFromSqliteInMemory(true)
	} else {
		database, err = db.NewFromSqlite(dbPath, "test.db", true)
	}
	require.NoError(t, err)

	// create logger
	testLogger := zerolog.New(zerolog.NewTestWriter(t))
	logger := base.Logger{Std: testLogger, Compliance: testLogger}

	// create observer
	ob, err := observer.NewObserver(
		chain,
		client,
		chainParams,
		zetacore,
		tss,
		database,
		logger,
		&metrics.TelemetryServer{},
	)
	require.NoError(t, err)

	return &testSuite{
		Observer: ob,
		client:   client,
		zetacore: zetacore,
		db:       database,
	}
}
