//go:build TESTNET
// +build TESTNET

package integrationtests

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/testutil/network"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
