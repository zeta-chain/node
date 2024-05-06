package evm

import (
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const TempSQLiteDbPath = "file::memory:?cache=shared"
const NumOfEntries = 2

type EVMClientTestSuite struct {
	suite.Suite
	db                        *gorm.DB
	outboundConfirmedReceipts map[string]*ethtypes.Receipt
	outTXConfirmedTransaction map[string]*ethtypes.Transaction
}

func TestEVMClient(t *testing.T) {
	suite.Run(t, new(EVMClientTestSuite))
}

func (suite *EVMClientTestSuite) SetupTest() {
	suite.outboundConfirmedReceipts = map[string]*ethtypes.Receipt{}
	suite.outTXConfirmedTransaction = map[string]*ethtypes.Transaction{}

	db, err := gorm.Open(sqlite.Open(TempSQLiteDbPath), &gorm.Config{})
	suite.NoError(err)

	suite.db = db

	err = db.AutoMigrate(&clienttypes.ReceiptSQLType{},
		&clienttypes.TransactionSQLType{},
		&clienttypes.LastBlockSQLType{})
	suite.NoError(err)

	//Create some receipt entries in the DB
	for i := 0; i < NumOfEntries; i++ {
		receipt := &ethtypes.Receipt{
			Type:              0,
			PostState:         nil,
			Status:            0,
			CumulativeGasUsed: 0,
			Bloom:             ethtypes.Bloom{},
			Logs:              nil,
			TxHash:            crypto.Keccak256Hash([]byte{byte(i)}),
			ContractAddress:   common.Address{},
			GasUsed:           0,
			BlockHash:         common.Hash{},
			BlockNumber:       nil,
			TransactionIndex:  uint(i),
		}
		r, _ := clienttypes.ToReceiptSQLType(receipt, strconv.Itoa(i))
		dbc := suite.db.Create(r)
		suite.NoError(dbc.Error)
		suite.outboundConfirmedReceipts[strconv.Itoa(i)] = receipt
	}

	//Create some transaction entries in the DB
	for i := 0; i < NumOfEntries; i++ {
		transaction := legacyTx(i)
		trans, _ := clienttypes.ToTransactionSQLType(transaction, strconv.Itoa(i))
		dbc := suite.db.Create(trans)
		suite.NoError(dbc.Error)
		suite.outTXConfirmedTransaction[strconv.Itoa(i)] = transaction
	}
}

func (suite *EVMClientTestSuite) TearDownSuite() {
	dbInst, err := suite.db.DB()
	suite.NoError(err)
	err = dbInst.Close()
	suite.NoError(err)
}

func (suite *EVMClientTestSuite) TestEVMReceipts() {
	for key, value := range suite.outboundConfirmedReceipts {
		var receipt clienttypes.ReceiptSQLType
		suite.db.Where("Identifier = ?", key).First(&receipt)

		r, _ := clienttypes.FromReceiptDBType(receipt.Receipt)
		suite.Equal(*r, *value)
	}
}

func (suite *EVMClientTestSuite) TestEVMTransactions() {
	for key, value := range suite.outTXConfirmedTransaction {
		var transaction clienttypes.TransactionSQLType
		suite.db.Where("Identifier = ?", key).First(&transaction)

		trans, _ := clienttypes.FromTransactionDBType(transaction.Transaction)

		have := trans.Hash()
		want := value.Hash()

		suite.Equal(want, have)
	}
}

func (suite *EVMClientTestSuite) TestEVMLastBlock() {
	lastBlockNum := uint64(12345)
	dbc := suite.db.Create(clienttypes.ToLastBlockSQLType(lastBlockNum))
	suite.NoError(dbc.Error)

	var lastBlockDB clienttypes.LastBlockSQLType
	dbf := suite.db.First(&lastBlockDB)
	suite.NoError(dbf.Error)

	suite.Equal(lastBlockNum, lastBlockDB.Num)

	lastBlockNum++
	dbs := suite.db.Save(clienttypes.ToLastBlockSQLType(lastBlockNum))
	suite.NoError(dbs.Error)

	dbf = suite.db.First(&lastBlockDB)
	suite.NoError(dbf.Error)
	suite.Equal(lastBlockNum, lastBlockDB.Num)
}

func legacyTx(nonce int) *ethtypes.Transaction {
	gasPrice, err := hexutil.DecodeBig("0x2bd0875aed")
	if err != nil {
		panic(err)
	}
	gas, err := hexutil.DecodeUint64("0x5208")
	if err != nil {
		panic(err)
	}
	to := common.HexToAddress("0x2f14582947e292a2ecd20c430b46f2d27cfe213c")
	value, err := hexutil.DecodeBig("0x2386f26fc10000")
	if err != nil {
		panic(err)
	}
	data := common.Hex2Bytes("0x")
	v, err := hexutil.DecodeBig("0x1")
	if err != nil {
		panic(err)
	}
	r, err := hexutil.DecodeBig("0x56b5bf9222ce26c3239492173249696740bc7c28cd159ad083a0f4940baf6d03")
	if err != nil {
		panic(err)
	}
	s, err := hexutil.DecodeBig("0x5fcd608b3b638950d3fe007b19ca8c4ead37237eaf89a8426777a594fd245c2a")
	if err != nil {
		panic(err)
	}

	newLegacyTx := ethtypes.NewTx(&ethtypes.LegacyTx{
		Nonce:    uint64(nonce),
		GasPrice: gasPrice,
		Gas:      gas,
		To:       &to,
		Value:    value,
		Data:     data,
		V:        v,
		R:        r,
		S:        s,
	})

	return newLegacyTx
}
