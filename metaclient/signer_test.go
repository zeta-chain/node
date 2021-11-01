package metaclient

import (
	"bytes"
	"github.com/Meta-Protocol/metacore/common"
	"github.com/Meta-Protocol/metacore/metaclient/config"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "gopkg.in/check.v1"
	"math/big"
)

type SignerSuite struct {
	signer *Signer
}

var _ = Suite(&SignerSuite{})


func (s *SignerSuite) SetUpTest(c *C) {
	// The following PrivKey has address 0xE80B6467863EbF8865092544f441da8fD3cF6074
	privateKey, err := crypto.HexToECDSA(config.TSS_TEST_PRIVKEY)
	// Uncomment the following code to generate new random private key pairs
	//privateKey, err := crypto.GenerateKey()
	//privkeyBytes := crypto.FromECDSA(privateKey)
	//c.Logf("privatekey %s", hex.EncodeToString(privkeyBytes))
	c.Assert(err, IsNil)
	tss := TestSigner{
		PrivKey: privateKey,
	}
	metaContractAddress := ethcommon.HexToAddress(config.BSC_ROUTER)
	signer, err := NewSigner(common.Chain("BSC"), config.BSC_ENDPOINT, tss.Address(), tss, config.META_ABI, metaContractAddress)
	c.Assert(err, IsNil)
	c.Logf("TSS Address %s", tss.Address().Hex())
	c.Logf("Contract on chain %s: %s", signer.chain, metaContractAddress.Hex())
	s.signer = signer

}

func (s *SignerSuite) TestSign(c *C) {
	data := []byte("1234")
	tx, sig, hash, err := s.signer.Sign(data, s.signer.tssSigner.Address(), 109, big.NewInt(2))
	_ = tx
	c.Assert(err, IsNil)
	pubkey, err := crypto.Ecrecover(hash, sig)
	c.Assert(err, IsNil)
	c.Assert(bytes.Equal(pubkey,s.signer.tssSigner.Pubkey()), Equals, true)
}

func (s *SignerSuite) TestMint(c *C) {
	txhash, err := s.signer.MMint(big.NewInt(1234), ethcommon.HexToAddress(config.TEST_RECEIVER), 80000, []byte{})
	c.Assert(err, IsNil)
	c.Logf("txhash %s", txhash)
}