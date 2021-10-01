package metaclientd

import (
	"encoding/json"
	"fmt"
	"github.com/Meta-Protocol/metacore/cmd/metaclientd/types"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	bridge *MetachainBridge
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpTest(c *C) {
	c.Logf("Settting up test...")
	homeDir, err := os.UserHomeDir()
	c.Logf("user home dir: %s", homeDir)
	//chainHomeFoler := homeDir + "/.meta-chain"
	chainHomeFoler := filepath.Join(homeDir, ".metacore")
	c.Logf("chain home dir: %s", chainHomeFoler)
	signerName := "alice"
	signerPass := "password"
	kb, _, err := GetKeyringKeybase(chainHomeFoler, signerName, signerPass)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to get keyring keybase")
	}

	k := NewKeysWithKeybase(kb, signerName, signerPass)
	log.Info().Msgf("keybase: %s", k.GetSignerInfo().GetAddress())

	chainIP := "127.0.0.1"
	bridge, err := NewMetachainBridge(k, chainIP)
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
	c.Logf("acc number %d acc sequence %d err %w", an, as, err)
}

//
//func (s *MySuite) Test1(c *C) {
//	b := s.bridge
//	//alice_address, err := sdk.AccAddressFromBech32("cosmos1svjyqzr90s28njlfuuf4hr7swxvr4gh5v9ytxc")
//	address := b.keys.GetSignerInfo().GetAddress()
//
//	bankClient := banktypes.NewQueryClient(b.grpcConn)
//	bankRes, err := bankClient.AllBalances(
//		context.Background(),
//		&banktypes.QueryAllBalancesRequest{Address: address.String()},
//	)
//	if err != nil {
//		c.Errorf("fail to get balance : %s", err)
//	}
//	// Prints the account balance
//	c.Log(bankRes.Balances)
//}

func (s *MySuite) TestSequenceNumber(c *C) {
	b := s.bridge
	path := fmt.Sprintf("%s/%s", AuthAccountEndpoint, b.GetKeys().GetSignerInfo().GetAddress())

	body, _, err := b.GetWithPath(path)
	if err != nil {
		//return 0, 0, fmt.Errorf("failed to get auth accounts: %w", err)
		c.Errorf("failed to get auth accounts: %s", err)
	}

	var resp types.AccountResp
	if err := json.Unmarshal(body, &resp); err != nil {
		//return 0, 0, fmt.Errorf("failed to unmarshal account resp: %w", err)
		c.Errorf("failed to unmarshal account resp: %s", err)
	}

	c.Logf("acct # %d, seq # %d\n", resp.Result.Value.AccountNumber, resp.Result.Value.Sequence)

	//return acc.AccountNumber, acc.Sequence, nil
}

//
//func (s *MySuite) Test3(c *C) {
//	b := s.bridge
//	address := b.keys.GetSignerInfo().GetAddress()
//
//	authClient := authtypes.NewQueryClient(b.grpcConn)
//	accRes, err := authClient.Account(
//		context.TODO(),
//		&authtypes.QueryAccountRequest{Address: address.String()},
//	)
//	if err != nil {
//		fmt.Printf("fail to get acct : %s", err)
//	}
//	// Prints the account balance
//	fmt.Println(accRes.Account.GetCachedValue())
//}

func (s *MySuite) TestBroadcast(c *C) {
	//creator := s.bridge.keys.GetSignerInfo().GetAddress()
	//msg := metatypes.NewMsgCreateTxin(creator.String(), "1", "2", "3",
	//	"4", "5", "8", 123)
	//_, err := s.bridge.Broadcast(msg)
	//if err != nil {
	//	fmt.Println(err)
	//}
}

//func (s *MySuite) TestObservedTxIn(c *C) {
//	b := s.bridge
//	err := b.PostTxIn("ETH.ETH", 2, 4, "ETH.BSC", "0xdeadbeef", "0x1234", 2345)
//	c.Assert(err, IsNil)
//	err = s.bridge.PostTxoutConfirmation(0, "0x4445", 23, 1794)
//	c.Assert(err, IsNil)
//
//	timer1 := time.NewTimer(6 * time.Second)
//	<-timer1.C
//
//	chain, _ := common.NewChain("ETH")
//	_, err = s.bridge.GetLastBlockObserved(chain)
//	c.Assert(err, IsNil)
//}
//
//func (s *MySuite) TestTxoutObserve(c *C) {
//	_, err := s.bridge.GetMetachainTxout()
//	c.Assert(err, IsNil)
//}

