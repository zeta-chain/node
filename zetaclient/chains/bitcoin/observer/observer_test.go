package observer

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/zetaclient/db"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"gorm.io/gorm"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/metrics"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/testutils/testlog"
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
	for i := range 2 {
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
	btcClient := mocks.NewBitcoinClient(t)
	btcClient.On("GetBlockCount", mock.Anything).Return(int64(100), nil)

	// test cases
	tests := []struct {
		name           string
		chain          chains.Chain
		btcClient      *mocks.BitcoinClient
		chainParams    observertypes.ChainParams
		zetacoreClient zrepo.ZetacoreClient
		tssSigner      interfaces.TSSSigner
		logger         base.Logger
		ts             *metrics.TelemetryServer
		errorMessage   string
		before         func()
		after          func()
	}{
		{
			name:           "should be able to create observer",
			chain:          chain,
			btcClient:      btcClient,
			chainParams:    params,
			zetacoreClient: nil,
			tssSigner:      mocks.NewTSS(t),
		},
		{
			name:           "should fail if net params is not found",
			chain:          chains.Chain{ChainId: 111}, // invalid chain id
			btcClient:      btcClient,
			chainParams:    params,
			zetacoreClient: nil,
			tssSigner:      mocks.NewTSS(t),
			errorMessage:   "unable to get BTC net params",
		},
		{
			name:           "should fail if env var us invalid",
			chain:          chain,
			btcClient:      btcClient,
			chainParams:    params,
			zetacoreClient: nil,
			tssSigner:      mocks.NewTSS(t),
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

			baseObserver, err := base.NewObserver(
				tt.chain,
				tt.chainParams,
				zrepo.New(tt.zetacoreClient, tt.chain, mode.StandardMode),
				tt.tssSigner,
				100,
				tt.ts,
				database,
				tt.logger,
			)
			require.NoError(t, err)

			// create observer
			ob, err := New(baseObserver, tt.btcClient, tt.chain)
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
		ob := newTestSuite(t, chains.BitcoinMainnet)

		// feed block hash, header and block to btc client
		hash := sample.BtcHash()
		header := &wire.BlockHeader{Version: 1}
		block := &btcjson.GetBlockVerboseTxResult{Version: 1}
		ob.client.On("GetBlockHash", mock.Anything, mock.Anything).Return(&hash, nil)
		ob.client.On("GetBlockHeader", mock.Anything, &hash).Return(header, nil)
		ob.client.On("GetBlockVerbose", mock.Anything, &hash).Return(block, nil)

		// get block and header from observer, fallback to btc client
		result, err := ob.GetBlockByNumberCached(ob.ctx, 100)
		require.NoError(t, err)
		require.EqualValues(t, header, result.Header)
		require.EqualValues(t, block, result.Block)

		// get block header from cache
		result, err = ob.GetBlockByNumberCached(ob.ctx, 100)
		require.NoError(t, err)
		require.EqualValues(t, header, result.Header)
		require.EqualValues(t, block, result.Block)
	})
	t.Run("should fail if stored type is not BlockNHeader", func(t *testing.T) {
		// create observer
		ob := newTestSuite(t, chains.BitcoinMainnet)

		// add a string to cache
		blockNumber := int64(100)
		ob.BlockCache().Add(blockNumber, "a string value")

		// get result from cache
		result, err := ob.GetBlockByNumberCached(ob.ctx, blockNumber)
		require.ErrorContains(t, err, "cached value is not of type *BTCBlockNHeader")
		require.Nil(t, result)
	})
}

func Test_SetLastStuckOutbound(t *testing.T) {
	// create observer and example stuck tx
	ob := newTestSuite(t, chains.BitcoinMainnet)
	btcTx := btcutil.NewTx(wire.NewMsgTx(wire.TxVersion))

	// STEP 1
	// initial stuck outbound does not exist
	_, found := ob.LastStuckOutbound()
	require.False(t, found)

	// STEP 2
	// set stuck outbound
	stuckTx := newLastStuckOutbound(100, btcTx, 30*time.Minute)
	ob.setLastStuckOutbound(stuckTx)

	// retrieve stuck outbound
	tx, found := ob.LastStuckOutbound()
	require.True(t, found)
	require.Equal(t, stuckTx, tx)

	// STEP 3
	// update stuck outbound
	stuckTxUpdate := newLastStuckOutbound(101, btcTx, 40*time.Minute)
	ob.setLastStuckOutbound(stuckTxUpdate)

	// retrieve updated stuck outbound
	tx, found = ob.LastStuckOutbound()
	require.True(t, found)
	require.Equal(t, stuckTxUpdate, tx)

	// STEP 4
	// clear stuck outbound
	ob.setLastStuckOutbound(nil)

	// stuck outbound should be nil
	tx, found = ob.LastStuckOutbound()
	require.False(t, found)
	require.Nil(t, tx)
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
	*Observer

	ctx      context.Context
	client   *mocks.BitcoinClient
	zetacore *mocks.ZetacoreClient
	db       *db.DB
}

type testSuiteOpts struct {
	dbPath string
}

type opt func(t *testSuiteOpts)

// withDatabasePath is an option to set custom db path
func withDatabasePath(dbPath string) opt {
	return func(t *testSuiteOpts) { t.dbPath = dbPath }
}

func newTestSuite(t *testing.T, chain chains.Chain, opts ...opt) *testSuite {
	// create test suite with options
	var testOpts testSuiteOpts
	for _, opt := range opts {
		opt(&testOpts)
	}

	require.True(t, chain.IsBitcoinChain())
	chainParams := mocks.MockChainParams(chain.ChainId, 10)

	client := mocks.NewBitcoinClient(t)
	zetacore := mocks.NewZetacoreClient(t)

	var tssSigner interfaces.TSSSigner
	if chains.IsBitcoinMainnet(chain.ChainId) {
		tssSigner = mocks.NewTSS(t).FakePubKey(testutils.TSSPubKeyMainnet)
	} else {
		tssSigner = mocks.NewTSS(t).FakePubKey(testutils.TSSPubkeyAthens3)
	}

	// create logger
	logger := testlog.New(t)
	baseLogger := base.Logger{Std: logger.Logger, Compliance: logger.Logger}

	var database *db.DB
	var err error
	if testOpts.dbPath == "" {
		database, err = db.NewFromSqliteInMemory(true)
		require.NoError(t, err)
	} else {
		database, err = db.NewFromSqlite(testOpts.dbPath, "test.db", true)
		require.NoError(t, err)
		t.Cleanup(func() { os.RemoveAll(testOpts.dbPath) })
	}

	client.On("GetBlockCount", mock.Anything).Maybe().Return(int64(100), nil).Maybe()

	baseObserver, err := base.NewObserver(
		chain,
		chainParams,
		zrepo.New(zetacore, chain, mode.StandardMode),
		tssSigner,
		100,
		&metrics.TelemetryServer{},
		database,
		baseLogger,
	)
	require.NoError(t, err)

	ob, err := New(baseObserver, client, chain)
	require.NoError(t, err)

	ts := &testSuite{
		ctx:      context.Background(),
		client:   client,
		zetacore: zetacore,
		db:       database,
		Observer: ob,
	}

	ts.zetacore.
		On("GetCctxByNonce", mock.Anything, mock.Anything, mock.Anything).
		Return(ts.mockGetCCTXByNonce).
		Maybe()

	return ts
}

func (ts *testSuite) mockGetCCTXByNonce(_ context.Context, chainID int64, nonce uint64) (*types.CrossChainTx, error) {
	// implement custom logic here if needed (e.g. mock)
	return nil, errors.New("not implemented")
}
