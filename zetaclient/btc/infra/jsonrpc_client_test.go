package infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/zetaclient/btc/model"
)

type JSONRpcClientTestSuite struct {
	suite.Suite
	client *JSONRpcClient
}

func (suite *JSONRpcClientTestSuite) SetupTest() {
	endpoint := "https://nd-456-407-783.p2pify.com/cee81511fd724bdcc75021ae81a9b5c9"
	targetAddress := "tb1q9dlnu5dr254s8xvtzlhk5ttu0c923u623qup39"
	suite.client = NewJSONRpcClient(endpoint, targetAddress)
}

func (suite *JSONRpcClientTestSuite) TearDownSuite() {
}

func (suite *JSONRpcClientTestSuite) TestBlockHeight() {
	var min int64 = 2405366
	blockNumber, err := suite.client.GetBlockHeight()
	suite.Require().NoError(err)
	suite.Assert().Greater(blockNumber, min)
	suite.T().Logf("blockHeight: %v\n", blockNumber)
}

func (suite *JSONRpcClientTestSuite) TestBlockHash() {
	var block int64 = 2405366
	expected := "00000000000000388a41c619e631bcd034e22ec6ca12a8a40123b64677f62cdd"
	hash, err := suite.client.GetBlockHash(block)
	suite.Require().NoError(err)
	suite.Assert().Equal(expected, hash)
}

func (suite *JSONRpcClientTestSuite) TestBlockByHash() {
	hash := "00000000000000388a41c619e631bcd034e22ec6ca12a8a40123b64677f62cdd"
	expected := "Amount: 0.0002, Address: 0x6162303132333435363738393031323334353637, Message: 89012345678901Hello World!"
	rawEvents, err := suite.client.GetEventsByHash(hash)
	suite.Require().NoError(err)

	event, err := model.ParseRawEvents(rawEvents)
	suite.Require().NoError(err)

	suite.Assert().Equal(expected, event.String())
}

// TestJSONRpcClient is the entry point of this test suite
func TestJSONRpcClient(t *testing.T) {
	suite.Run(t, new(JSONRpcClientTestSuite))
}
