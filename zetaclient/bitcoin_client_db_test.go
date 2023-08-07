package zetaclient

import (
	"strconv"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/stretchr/testify/suite"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BitcoinClientDBTestSuite struct {
	suite.Suite
	db          *gorm.DB
	submittedTx map[string]btcjson.GetTransactionResult
}

func TestBitcoinClientDB(t *testing.T) {
	suite.Run(t, new(BitcoinClientDBTestSuite))
}

func (suite *BitcoinClientDBTestSuite) SetupTest() {
	suite.submittedTx = map[string]btcjson.GetTransactionResult{}

	db, err := gorm.Open(sqlite.Open(TempSQLiteDbPath), &gorm.Config{})
	suite.NoError(err)

	suite.db = db

	err = db.AutoMigrate(&clienttypes.TransactionResultSQLType{})
	suite.NoError(err)

	//Create some Transaction entries in the DB
	for i := 0; i < NumOfEntries; i++ {
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
		dbc := suite.db.Create(&r)
		suite.NoError(dbc.Error)
		suite.submittedTx[strconv.Itoa(i)] = txResult
	}
}

func (suite *BitcoinClientDBTestSuite) TearDownSuite() {
}

func (suite *BitcoinClientDBTestSuite) TestSubmittedTx() {
	var submittedTransactions []clienttypes.TransactionResultSQLType
	err := suite.db.Find(&submittedTransactions).Error
	suite.NoError(err)

	for _, txResult := range submittedTransactions {
		r, err := clienttypes.FromTransactionResultSQLType(txResult)
		suite.NoError(err)
		want := suite.submittedTx[txResult.Key]
		have := r

		suite.Equal(want, have)
	}
}
