package zetaclient

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/99designs/keyring"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	hd "github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	. "gopkg.in/check.v1"

	"github.com/zeta-chain/zetacore/cmd"
	"github.com/zeta-chain/zetacore/common/cosmos"
)

type KeysSuite struct{}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&KeysSuite{})

func (*KeysSuite) SetUpSuite(c *C) {
	SetupConfigForTest()
}

var (
	password = "password"
)

func SetupConfigForTest() {
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(cmd.Bech32PrefixAccAddr, cmd.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(cmd.Bech32PrefixValAddr, cmd.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(cmd.Bech32PrefixConsAddr, cmd.Bech32PrefixConsPub)
	//config.SetCoinType(cmd.MetaChainCoinType)
	config.SetFullFundraiserPath(cmd.ZetaChainHDPath)
	types.SetCoinDenomRegex(func() string {
		return cmd.DenomRegex
	})
}

const (
	signerNameForTest     = `jack`
	signerPasswordForTest = `password`
)

func (*KeysSuite) setupKeysForTest(c *C) string {
	ns := strconv.Itoa(time.Now().Nanosecond())
	metaCliDir := filepath.Join(os.TempDir(), ns, ".metacli")
	c.Logf("metacliDir:%s", metaCliDir)
	buf := bytes.NewBufferString(signerPasswordForTest)
	// the library used by keyring is using ReadLine , which expect a new line
	buf.WriteByte('\n')
	buf.WriteString(signerPasswordForTest)
	buf.WriteByte('\n')
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)
	kb, err := cKeys.New(cosmos.KeyringServiceName(), cKeys.BackendTest, metaCliDir, buf, cdc)
	c.Assert(err, IsNil)

	_, _, err = kb.NewMnemonic(GetGranteeKeyName(signerNameForTest), cKeys.English, cmd.ZetaChainHDPath, password, hd.Secp256k1)
	c.Assert(err, IsNil)
	return metaCliDir
}

func (ks *KeysSuite) TestGetKeyringKeybase(c *C) {
	keyring.Debug = true
	cfg := &config.Config{
		AuthzHotkey:  "bob",
		ZetaCoreHome: "/Users/test/.zetacored/",
	}
	_, _, err := GetKeyringKeybase(cfg)
	c.Assert(err, NotNil)
}

func (ks *KeysSuite) TestNewKeys(c *C) {
	oldStdIn := os.Stdin
	defer func() {
		os.Stdin = oldStdIn
	}()
	os.Stdin = nil
	folder := ks.setupKeysForTest(c)
	defer func() {
		err := os.RemoveAll(folder)
		c.Assert(err, IsNil)
	}()

	cfg := &config.Config{
		AuthzHotkey:  signerNameForTest,
		ZetaCoreHome: folder,
	}

	k, _, err := GetKeyringKeybase(cfg)
	c.Assert(err, IsNil)
	c.Assert(k, NotNil)
	granter := cosmos.AccAddress(crypto.AddressHash([]byte("granter")))
	ki := NewKeysWithKeybase(k, granter, signerNameForTest, signerPasswordForTest)
	kInfo := ki.GetSignerInfo()
	c.Assert(kInfo, NotNil)
	//c.Assert(kInfo.G, Equals, signerNameForTest)
	priKey, err := ki.GetPrivateKey()
	c.Assert(err, IsNil)
	c.Assert(priKey, NotNil)
	c.Assert(priKey.Bytes(), HasLen, 32)
	kb := ki.GetKeybase()
	c.Assert(kb, NotNil)

	msg := "hello"
	signedMsg, err := priKey.Sign([]byte(msg))
	c.Assert(err, IsNil)
	pubKey, err := ki.GetSignerInfo().GetPubKey()
	c.Assert(err, IsNil)
	c.Assert(pubKey.VerifySignature([]byte(msg), signedMsg), Equals, true)
}
