package metaclientd

import (
	"bytes"
	"github.com/cosmos/cosmos-sdk/types"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	hd "github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	. "gopkg.in/check.v1"

	"github.com/Meta-Protocol/metacore/cmd"
	"github.com/Meta-Protocol/metacore/common/cosmos"
)

type KeysSuite struct{}
func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&KeysSuite{})

func (*KeysSuite) SetUpSuite(c *C) {
	SetupConfigForTest()
}

func SetupConfigForTest() {
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(cmd.Bech32PrefixAccAddr, cmd.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(cmd.Bech32PrefixValAddr, cmd.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(cmd.Bech32PrefixConsAddr, cmd.Bech32PrefixConsPub)
	//config.SetCoinType(cmd.MetaChainCoinType)
	config.SetFullFundraiserPath(cmd.METAChainHDPath)
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
	kb, err := cKeys.New(cosmos.KeyringServiceName(), cKeys.BackendTest, metaCliDir, buf)
	c.Assert(err, IsNil)
	_, _, err = kb.NewMnemonic(signerNameForTest, cKeys.English, cmd.METAChainHDPath, hd.Secp256k1)
	c.Assert(err, IsNil)
	return metaCliDir
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

	k, info, err := GetKeyringKeybase(folder, signerNameForTest, signerPasswordForTest)
	c.Assert(err, IsNil)
	c.Assert(k, NotNil)
	c.Assert(info, NotNil)
	ki := NewKeysWithKeybase(k, signerNameForTest, signerPasswordForTest)
	kInfo := ki.GetSignerInfo()
	c.Assert(kInfo, NotNil)
	c.Assert(kInfo.GetName(), Equals, signerNameForTest)
	priKey, err := ki.GetPrivateKey()
	c.Assert(err, IsNil)
	c.Assert(priKey, NotNil)
	c.Assert(priKey.Bytes(), HasLen, 32)
	kb := ki.GetKeybase()
	c.Assert(kb, NotNil)
}
