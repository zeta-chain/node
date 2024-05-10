package observer

import (
	"math/big"
	"strconv"
	"sync"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/testutil/sample"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	clientcommon "github.com/zeta-chain/zetacore/zetaclient/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/context"
	"github.com/zeta-chain/zetacore/zetaclient/metrics"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/mocks"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	// tempSQLiteDbPath is the temporary SQLite database used for testing
	tempSQLiteDbPath = "file::memory:?cache=shared"
)

var (
	// the relative path to the testdata directory
	TestDataDir = "../../../"
)

// setupDBTxResults creates a new SQLite database and populates it with some transaction results.
func setupDBTxResults(t *testing.T) (*gorm.DB, map[string]btcjson.GetTransactionResult) {
	submittedTx := map[string]btcjson.GetTransactionResult{}

	db, err := gorm.Open(sqlite.Open(tempSQLiteDbPath), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&clienttypes.TransactionResultSQLType{})
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
		dbc := db.Create(&r)
		require.NoError(t, dbc.Error)
		submittedTx[strconv.Itoa(i)] = txResult
	}

	return db, submittedTx
}

func TestNewBitcoinObserver(t *testing.T) {
	t.Run("should return error because zetacore doesn't update core context", func(t *testing.T) {
		cfg := config.NewConfig()
		coreContext := context.NewZetaCoreContext(cfg)
		appContext := context.NewAppContext(coreContext, cfg)
		chain := chains.BtcMainnetChain
		zetacoreClient := mocks.NewMockZetaCoreClient()
		tss := mocks.NewMockTSS(chains.BtcTestNetChain, sample.EthAddress().String(), "")
		loggers := clientcommon.ClientLogger{}
		btcCfg := cfg.BitcoinConfig
		ts := metrics.NewTelemetryServer()

		client, err := NewObserver(appContext, chain, zetacoreClient, tss, tempSQLiteDbPath, loggers, btcCfg, ts)
		require.ErrorContains(t, err, "btc chains params not initialized")
		require.Nil(t, client)
	})
}

func TestConfirmationThreshold(t *testing.T) {
	ob := &Observer{Mu: &sync.Mutex{}}
	t.Run("should return confirmations in chain param", func(t *testing.T) {
		ob.SetChainParams(observertypes.ChainParams{ConfirmationCount: 3})
		require.Equal(t, int64(3), ob.ConfirmationsThreshold(big.NewInt(1000)))
	})

	t.Run("should return big value confirmations", func(t *testing.T) {
		ob.SetChainParams(observertypes.ChainParams{ConfirmationCount: 3})
		require.Equal(t, int64(bigValueConfirmationCount), ob.ConfirmationsThreshold(big.NewInt(bigValueSats)))
	})

	t.Run("big value confirmations is the upper cap", func(t *testing.T) {
		ob.SetChainParams(observertypes.ChainParams{ConfirmationCount: bigValueConfirmationCount + 1})
		require.Equal(t, int64(bigValueConfirmationCount), ob.ConfirmationsThreshold(big.NewInt(1000)))
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
