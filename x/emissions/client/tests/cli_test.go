package querytests

import (
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/testutil/network"
	"testing"
)

func TestCLIQuerySuite(t *testing.T) {
	cfg := network.DefaultConfig()
	suite.Run(t, NewCLITestSuite(cfg))
}
