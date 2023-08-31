package zetaclient

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
	. "gopkg.in/check.v1"
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
	sc, err := NewCctxScanner(nil, TempSQLiteDbPath, true, &logger)
	suite.NoError(err)
	suite.sc = sc
}

func (suite *CctxScannerTestSuite) SaveNonces(goerliNonce uint64, bsctestNonce uint64, mumbaiNonce uint64, btctestNonce uint64) {
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

func (suite *CctxScannerTestSuite) TestFirstNonceToScan() {
	// Check that the DB is empty
	suite.Equal(suite.sc.firstNonceToScan[5], uint64(0))
	suite.Equal(suite.sc.firstNonceToScan[97], uint64(0))
	suite.Equal(suite.sc.firstNonceToScan[80001], uint64(0))
	suite.Equal(suite.sc.firstNonceToScan[18332], uint64(0))

	suite.Equal(suite.sc.nextNonceToScan[5], uint64(0))
	suite.Equal(suite.sc.nextNonceToScan[97], uint64(0))
	suite.Equal(suite.sc.nextNonceToScan[80001], uint64(0))
	suite.Equal(suite.sc.nextNonceToScan[18332], uint64(0))

	// Create some entries in the DB
	suite.SaveNonces(1, 41806, 17490, 138)

	// Check the DB
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
	suite.SaveNonces(2349, 51570, 21086, 259)

	// Check the DB again
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
