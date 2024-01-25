//go:build integration
// +build integration

// this is integration test; must be run when a chain is running:
// starport chain serve

package zetaclient

import (
	"github.com/zeta-chain/zetacore/zetaclient/keys"
	"github.com/zeta-chain/zetacore/zetaclient/zetabridge"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	. "gopkg.in/check.v1"
)

type MySuite struct {
	bridge *zetabridge.ZetaCoreBridge
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpTest(c *C) {
	keys.SetupConfigForTest() // setup meta-prefix

	c.Logf("Settting up test...")
	homeDir, err := os.UserHomeDir()
	c.Logf("user home dir: %s", homeDir)
	chainHomeFoler := filepath.Join(homeDir, ".zetacored")
	c.Logf("chain home dir: %s", chainHomeFoler)

	// alice is the default user created by Starport chain serve
	signerName := "alice"
	signerPass := "password"
	kb, _, err := keys.GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get keyring keybase")
	}

	k := keys.NewKeysWithKeybase(kb, signerName, signerPass)
	//log.Info().Msgf("keybase: %s", k.GetSignerInfo().GetAddress())

	chainIP := os.Getenv("CHAIN_IP")
	if chainIP == "" {
		chainIP = "127.0.0.1"
	}
	bridge, err := zetabridge.NewZetaCoreBridge(k, chainIP, "alice")
	if err != nil {
		c.Fail()
	}
	s.bridge = bridge
}

func (s *MySuite) TestGetBlockHeight(c *C) {
	h, err := s.bridge.GetBlockHeight()
	c.Assert(err, IsNil)
	c.Logf("height %d", h)
}

func (s *MySuite) TestGetAccountNumberAndSeuqeuence(c *C) {
	an, as, err := s.bridge.GetAccountNumberAndSequenceNumber()
	c.Assert(err, IsNil)
	c.Logf("acc number %d acc sequence %d", an, as)
}
