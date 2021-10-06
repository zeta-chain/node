// +build voter

// this is integration test; must be run when a chain is running:
// starport chain serve

package metaclientd

import (
	"github.com/rs/zerolog/log"
	. "gopkg.in/check.v1"
	"os"
	"path/filepath"
	"time"
)

type VoterSuite struct {
	bridge1 *MetachainBridge
	bridge2 *MetachainBridge
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
	chainHomeFoler := filepath.Join(homeDir, ".metacore")
	c.Logf("chain home dir: %s", chainHomeFoler)

	// first signer & bridge
	// alice is the default user created by Starport chain serve
	{
		signerName := "alice"
		signerPass := "password"
		kb, _, err := GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
		if err != nil {
			log.Fatal().Err(err).Msg("fail to get keyring keybase")
		}

		k := NewKeysWithKeybase(kb, signerName, signerPass)

		chainIP := "127.0.0.1"
		bridge, err := NewMetachainBridge(k, chainIP, "alice")
		if err != nil {
			c.Fail()
		}
		s.bridge1 = bridge
	}

	// second signer & bridge
	// alice is the default user created by Starport chain serve
	{
		signerName := "bob"
		signerPass := "password"
		kb, _, err := GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
		if err != nil {
			log.Fatal().Err(err).Msg("fail to get keyring keybase")
		}

		k := NewKeysWithKeybase(kb, signerName, signerPass)

		chainIP := "127.0.0.1"
		bridge, err := NewMetachainBridge(k, chainIP, "bob")
		if err != nil {
			c.Fail()
		}
		s.bridge2 = bridge
	}
}

func (s *VoterSuite) TestObservedTxIn(c *C) {
	b1 := s.bridge1
	b2 := s.bridge2
	//err := b.PostTxIn("ETH.ETH", 2, 4, "ETH.BSC", "0xdeadbeef", "0x1234", 2345)
	metaHash, err := b1.PostTxIn("0xfrom", "0xto", "0xsource.ETH", 123456, 23245, "0xdest.BSC",
		"0xtxhash", 123123)

	c.Assert(err, IsNil)
	log.Info().Msgf("PostTxIn metaHash %s", metaHash)

	// wait for the next block
	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C

	metaHash, err = b2.PostTxIn("0xfrom", "0xto", "0xsource.ETH", 123456, 23245, "0xdest.BSC",
		"0xtxhash", 123123)
	c.Assert(err, IsNil)
	log.Info().Msgf("Second PostTxIn metaHash %s", metaHash)

	// wait for the next block
	timer2 := time.NewTimer(2 * time.Second)
	<-timer2.C

	txouts, err := b1.GetAllTxout()
	c.Assert(err, IsNil)
	log.Info().Msgf("txouts: %v", txouts)
	c.Assert(len(txouts) >=1, Equals, true)

	txout := txouts[0]
	tid := txout.Id
	metaHash, err = b1.PostTxoutConfirmation(tid, "0xhashtxout", 1337, "0xnicetoken", 1773, "0xmywallet", 12345)
	timer3 := time.NewTimer(2 * time.Second)
	<-timer3.C
	metaHash, err = b2.PostTxoutConfirmation(tid, "0xhashtxout", 1337, "0xnicetoken", 1773, "0xmywallet", 12345)

	timer4 := time.NewTimer(2 * time.Second)
	<-timer4.C

	txouts, err = b1.GetAllTxout()
	c.Assert(err, IsNil)
	log.Info().Msgf("txouts: %v", txouts)
}
