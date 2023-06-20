//go:build TESTNET
// +build TESTNET

package testutil

import (
	"github.com/stretchr/testify/suite"
	"github.com/zeta-chain/zetacore/testutil/network"
	"testing"
)

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
