package infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type JSONRpcClientTestSuite struct {
	suite.Suite
	cli *JSONRpcClient
}

func (suite *JSONRpcClientTestSuite) SetupTest() {
	endpoint := "https://nd-456-407-783.p2pify.com/cee81511fd724bdcc75021ae81a9b5c9"
	targetAddress := "tb1q9dlnu5dr254s8xvtzlhk5ttu0c923u623qup39"
	suite.cli = NewJSONRpcClient(endpoint, targetAddress)
}

func (suite *JSONRpcClientTestSuite) TearDownSuite() {
}

func (suite *JSONRpcClientTestSuite) TestBlockHeight() {
	var min int64 = 2378200
	blockNumber, err := suite.cli.GetBlockHeight()
	suite.Require().NoError(err)
	suite.Assert().Greater(blockNumber, min)
	suite.T().Logf("blockHeight: %v\n", blockNumber)
}

func (suite *JSONRpcClientTestSuite) TestBlockHash() {
	var expected string = "0000000039364a85b7326558802bacef7390aed973e9f9b0627f29e8ac3e6676"
	var block int64 = 2378200
	hash, err := suite.cli.GetBlockHash(block)
	suite.Require().NoError(err)
	suite.Assert().Equal(expected, hash)
}

func (suite *JSONRpcClientTestSuite) TestBlockByHash() {
	var hash string = "0000000039364a85b7326558802bacef7390aed973e9f9b0627f29e8ac3e6676"
	events, err := suite.cli.GetEventsByHash(hash)
	for _, evt := range events {
		suite.T().Logf("EVENT=>%v\n", evt)
	}
	suite.Require().NoError(err)
}

// TestJSONRpcClient is the entry point of this test suite
func TestJSONRpcClient(t *testing.T) {
	suite.Run(t, new(JSONRpcClientTestSuite))
}
