//go:build metacore_observer
// +build metacore_observer

package zetaclient

import (
	"os"
	"path/filepath"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/types"
	. "gopkg.in/check.v1"
)

type COSuite struct {
	bridge1      *ZetaCoreBridge
	bridge2      *ZetaCoreBridge
	signer       *EVMSigner
	coreObserver *CoreObserver
}

var _ = Suite(&COSuite{})

const (
	TEST_SENDER   = "0x566bF3b1993FFd4BA134c107A63bb2aebAcCdbA0"
	TEST_RECEIVER = "0x566bF3b1993FFd4BA134c107A63bb2aebAcCdbA0"
)

func (s *COSuite) SetUpTest(c *C) {
	types.SetupConfigForTest() // setup meta-prefix

	// setup 2 metabridges
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
		kb, _, err := GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
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

	// setup mock TSS signers:
	// The following PrivKey has address 0xE80B6467863EbF8865092544f441da8fD3cF6074
	privateKey, err := crypto.HexToECDSA(config.TssTestPrivkey)
	c.Assert(err, IsNil)
	tss := TestSigner{
		PrivKey: privateKey,
	}
	metaContractAddress := ethcommon.HexToAddress(config.ETH_MPI_ADDRESS)
	signer, err := NewEVMSigner(common.Chain("ETH"), config.GOERLI_RPC_ENDPOINT, tss.EVMAddress(), tss, config.META_TEST_GOERLI_ABI, metaContractAddress)
	c.Assert(err, IsNil)
	c.Logf("TSS EVMAddress %s", tss.EVMAddress().Hex())
	c.Logf("ETH MPI EVMAddress: %s", config.ETH_MPI_ADDRESS)

	s.signer = signer

	// setup zetacore observer
	co := &CoreObserver{
		bridge: s.bridge1,
		signer: signer,
	}
	s.coreObserver = co
	s.coreObserver.MonitorCore()
}

func (s *COSuite) TestSendFlow(c *C) {
	b1 := s.bridge1
	b2 := s.bridge2
	metaHash, err := b1.PostVoteInbound(TEST_SENDER, "Ethereum", TEST_SENDER, TEST_RECEIVER, "BSC", "1337", "0", "treat or trick",
		"0xtxhash", 123123, "0xtoken")
	c.Assert(err, IsNil)
	c.Logf("PostVoteInbound metaHash %s", metaHash)

	timer1 := time.NewTimer(2 * time.Second)
	<-timer1.C

	metaHash, err = b2.PostVoteInbound(TEST_SENDER, "Ethereum", TEST_SENDER, TEST_RECEIVER, "BSC", "1337", "0", "treat or trick",
		"0xtxhash", 123123, "0xtoken")
	c.Assert(err, IsNil)
	c.Logf("Second PostVoteInbound metaHash %s", metaHash)

	timer2 := time.NewTimer(2 * time.Second)
	<-timer2.C

	time.Sleep(15 * time.Second)
	//ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	//<-ch
	//c.Logf("stop signal received")
}
