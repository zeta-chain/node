//go:build voter
// +build voter

// this is integration test; must be run when a chain is running:
// starport chain serve

package zetaclient

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	. "gopkg.in/check.v1"
)

type VoterSuite struct {
	bridge1 *ZetaCoreBridge
	bridge2 *ZetaCoreBridge
}

var _ = Suite(&VoterSuite{})

func (s *VoterSuite) SetUpTest(c *C) {
	SetupConfigForTest() // setup meta-prefix

	c.Logf("Settting up test...")
	homeDir, err := os.UserHomeDir()
	if err != nil {
		c.Logf("UserHomeDir error")
		c.Fail()
	}
	c.Logf("user home dir: %s", homeDir)
	chainHomeFoler := filepath.Join(homeDir, ".zetacored")
	c.Logf("chain home dir: %s", chainHomeFoler)

	// first signer & zetaClient
	// alice is the default user created by Starport chain serve
	{
		signerName := "alice"
		signerPass := "password"
		kb, _, err := GetKeyringKeybase([]common.KeyType{common.ObserverGranteeKey}, chainHomeFoler, signerName, signerPass)
		if err != nil {
			log.Fatal().Err(err).Msg("fail to get keyring keybase")
		}

		k := NewKeysWithKeybase(kb, signerName, signerPass)

		chainIP := os.Getenv("CHAIN_IP")
		if chainIP == "" {
			chainIP = "127.0.0.1"
		}
		bridge, err := NewZetaCoreBridge(k, chainIP, "alice")
		if err != nil {
			c.Fail()
		}
		s.bridge1 = bridge
	}

	// second signer & zetaClient
	// alice is the default user created by Starport chain serve
	{
		signerName := "bob"
		signerPass := "password"
		kb, _, err := GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
		if err != nil {
			log.Fatal().Err(err).Msg("fail to get keyring keybase")
		}

		k := NewKeysWithKeybase(kb, signerName, signerPass)

		chainIP := os.Getenv("CHAIN_IP")
		if chainIP == "" {
			chainIP = "127.0.0.1"
		}
		bridge, err := NewZetaCoreBridge(k, chainIP, "bob")
		if err != nil {
			c.Fail()
		}
		s.bridge2 = bridge
	}
}

func (s *VoterSuite) TestSendVoter(c *C) {
	b1 := s.bridge1
	b2 := s.bridge2
	metaHash, err := b1.PostVoteInbound("0xfrom", "Ethereum", "0xfrom", "0xto", "BSC", "123456", "23245", "little message",
		"0xtxhash", 123123, "0xtoken")

	c.Assert(err, IsNil)
	log.Info().Msgf("PostVoteInbound metaHash %s", metaHash)

	// wait for the next block
	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C

	metaHash, err = b2.PostVoteInbound("0xfrom", "Ethereum", "0xfrom", "0xto", "BSC", "123456", "23245", "little message",
		"0xtxhash", 123123, "0xtoken")
	c.Assert(err, IsNil)
	log.Info().Msgf("Second PostVoteInbound metaHash %s", metaHash)

	// wait for the next block
	timer2 := time.NewTimer(2 * time.Second)
	<-timer2.C

	sends, err := b1.GetAllSend()
	c.Assert(err, IsNil)
	log.Info().Msgf("sends: %v", sends)
	c.Assert(len(sends) >= 1, Equals, true)

	send := sends[0]

	metaHash, err = b1.PostVoteOutbound(send.Index, "0xoutHash", 2123, "23245")
	c.Assert(err, IsNil)

	timer3 := time.NewTimer(2 * time.Second)
	<-timer3.C

	metaHash, err = b2.PostVoteOutbound(send.Index, "0xoutHash", 2123, "23245")
	c.Assert(err, IsNil)

	receives, err := b2.GetAllReceive()
	c.Assert(err, IsNil)
	log.Info().Msgf("receives: %v", receives)
	c.Assert(len(receives), Equals, 1)

	timer4 := time.NewTimer(2 * time.Second)
	<-timer4.C

	last, err := b1.GetLastBlockHeight()
	c.Assert(err, IsNil)
	c.Assert(len(last), Equals, 2)

}
