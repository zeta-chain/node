package zetaclient

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	. "gopkg.in/check.v1"
)

const (
	tssPubkey1 = "zetapub1addwnpepqde0ztz2agdt0ss47dhdj2867ad63ju82f87a7h97memasegvnr3xehkryd"
	tssPubkey2 = "zetapub1addwnpepqfapt52wqw6k2kv0kvkuf8u0e8l37q57ntau7qu5ppz9sh690cs9cg0yxzs"
)

type CctxScannerTestSuite struct {
	suite.Suite
	sc *CctxScanner
}

var _ = Suite(&CctxScannerTestSuite{})

func TestCctxScanner(t *testing.T) {
	suite.Run(t, new(CctxScannerTestSuite))
}

func (suite *CctxScannerTestSuite) SetupTest() {
	logger := zerolog.New(os.Stdout)
	sc, err := NewCctxScanner(nil, TempSQLiteDbPath, true, tssPubkey1, &logger)
	suite.NoError(err)
	suite.sc = sc
}

func (suite *CctxScannerTestSuite) SaveNLoadNonces(goerliNonce uint64, bsctestNonce uint64, mumbaiNonce uint64, btctestNonce uint64) {
	goerli := clienttypes.ToFirstNonceToScanSQLType(5, goerliNonce)
	bsctest := clienttypes.ToFirstNonceToScanSQLType(97, bsctestNonce)
	mumbai := clienttypes.ToFirstNonceToScanSQLType(80001, mumbaiNonce)
	btctest := clienttypes.ToFirstNonceToScanSQLType(18332, btctestNonce)
	firstNonces := []*clienttypes.FirstNonceToScanSQLType{goerli, bsctest, mumbai, btctest}
	for _, firstNonce := range firstNonces {
		dbc := suite.sc.db.Save(firstNonce)
		suite.NoError(dbc.Error)
	}
	err := suite.sc.LoadDB(TempSQLiteDbPath, true)
	suite.NoError(err)
}

func (suite *CctxScannerTestSuite) CheckEmptyNonces() {
	suite.Equal(suite.sc.firstNonceToScan[5], uint64(0))
	suite.Equal(suite.sc.firstNonceToScan[97], uint64(0))
	suite.Equal(suite.sc.firstNonceToScan[80001], uint64(0))
	suite.Equal(suite.sc.firstNonceToScan[18332], uint64(0))

	suite.Equal(suite.sc.nextNonceToScan[5], uint64(0))
	suite.Equal(suite.sc.nextNonceToScan[97], uint64(0))
	suite.Equal(suite.sc.nextNonceToScan[80001], uint64(0))
	suite.Equal(suite.sc.nextNonceToScan[18332], uint64(0))
}

func (suite *CctxScannerTestSuite) TestFirstNonceToScan() {
	// Make sure all maps are empty
	suite.CheckEmptyNonces()

	// Create some entries in the DB
	suite.SaveNLoadNonces(1, 41806, 17490, 138)

	// Check the DB nonces
	var firstNonces1 []clienttypes.FirstNonceToScanSQLType
	err := suite.sc.db.Find(&firstNonces1).Error
	suite.NoError(err)
	for _, firstNonce := range firstNonces1 {
		want1 := suite.sc.firstNonceToScan[firstNonce.ID]
		want2 := suite.sc.nextNonceToScan[firstNonce.ID]
		have := firstNonce.FirstNonce
		suite.Equal(want1, have)
		suite.Equal(want2, have)
	}

	// Update entries in the DB
	suite.SaveNLoadNonces(2349, 51570, 21086, 259)

	// Check the DB nonces again
	var firstNonces2 []clienttypes.FirstNonceToScanSQLType
	err = suite.sc.db.Find(&firstNonces2).Error
	suite.NoError(err)
	for _, firstNonce := range firstNonces2 {
		want1 := suite.sc.firstNonceToScan[firstNonce.ID]
		want2 := suite.sc.nextNonceToScan[firstNonce.ID]
		have := firstNonce.FirstNonce
		suite.Equal(want1, have)
		suite.Equal(want2, have)
	}
}

func (suite *CctxScannerTestSuite) TestReset() {
	// Create some entries in the DB
	suite.SaveNLoadNonces(1, 41806, 17490, 138)

	// create another scanner with different tss pubkey
	logger := zerolog.New(os.Stdout)
	sc, err := NewCctxScanner(nil, TempSQLiteDbPath, true, tssPubkey2, &logger)
	suite.NoError(err)
	suite.sc = sc

	// Make sure all maps are empty again
	suite.CheckEmptyNonces()
}
