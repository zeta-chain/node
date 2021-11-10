package metaclient

import (
	"context"
	"github.com/Meta-Protocol/metacore/metaclient/config"
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


}
