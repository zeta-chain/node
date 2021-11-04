
package metaclient

import (
	. "gopkg.in/check.v1"
)

type ChainObSuite struct {
	chainob *ChainObserver
}

var _ = Suite(&ChainObSuite{})

func (s *ChainObSuite) SetUpTest(c *C) {

}