package metaclient

import (
	"context"
	"github.com/Meta-Protocol/metacore/metaclient/config"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	. "gopkg.in/check.v1"
)

type ChainClientSuite struct {

}

var _ = Suite(&ChainClientSuite{})

func (s *ChainClientSuite) SetUpTest(c *C) {

}

func (s *ChainClientSuite) TestPolygonClient(c *C) {
	client, err := ethclient.Dial(config.POLY_ENDPOINT)
	c.Assert(err, IsNil)
	bn, err := client.BlockNumber(context.TODO())
	c.Assert(err, IsNil)
	c.Logf("blocknum %d", bn)

	gas, err := client.SuggestGasPrice(context.TODO())
	c.Assert(err, IsNil)
	c.Logf("gas price %d", gas)

	receipt,  err := client.TransactionReceipt(context.TODO(), ethcommon.HexToHash("0xa8ab7e7242ee1b00c7e4de581d9c87b2465bae76115bce086e7ff0e8d6a7e1ef"))
	c.Assert(err, IsNil)
	c.Log(receipt.Status, receipt.PostState, receipt.GasUsed, receipt.Logs[0], receipt.BlockNumber)

	// non-existent txhash
	_, _, err = client.TransactionByHash(context.TODO(), ethcommon.HexToHash("0x2c5d00aa638f04e49eb5d86499d1ea25d6ed1c62279008b93963d09b70fba270"))
	c.Assert(err, NotNil)

	_, err = client.TransactionReceipt(context.TODO(), ethcommon.HexToHash("0x2c5d00aa638f04e49eb5d86499d1ea25d6ed1c62279008b93963d09b70fba270"))
	c.Assert(err, NotNil)

	x := ethcommon.HexToHash("0x33")
	c.Log(x)
}

func (s *ChainClientSuite) TestBSCClient(c *C) {
	client, err := ethclient.Dial(config.BSC_ENDPOINT)
	c.Assert(err, IsNil)
	bn, err := client.BlockNumber(context.TODO())
	c.Assert(err, IsNil)
	c.Logf("blocknum %d", bn)

	gas, err := client.SuggestGasPrice(context.TODO())
	c.Assert(err, IsNil)
	c.Logf("gas price %d", gas)

	receipt,  err := client.TransactionReceipt(context.TODO(), ethcommon.HexToHash("0x63326995eb00cc49df7d2aa249ec473dc351cea30f230001c4d310d6e6763490"))
	c.Assert(err, IsNil)
	c.Log(receipt.Status, receipt.PostState, receipt.GasUsed, receipt.Logs[0], receipt.BlockNumber)
}
