package integrationtests

import (
	"testing"

	"github.com/zeta-chain/zetacore/testutil/network"
)

func TestIntegrationTestSuite(t *testing.T) {
	_ = network.DefaultConfig()
	//suite.Run(t, NewIntegrationTestSuite(cfg))
}
