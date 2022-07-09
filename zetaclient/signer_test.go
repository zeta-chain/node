package zetaclient

import (
	"bytes"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/zetaclient/config"
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
	metaContractAddress := ethcommon.HexToAddress(config.Chains[common.BSCTestnetChain.String()].ConnectorContractAddress)
	signer, err := NewSigner(common.Chain("BSC"), config.Chains[common.BSCTestnetChain.String()].Endpoint, tss, config.CONNECTOR_ABI_STRING, metaContractAddress)
	c.Assert(err, IsNil)
	c.Logf("TSS Address %s", tss.Address().Hex())
	c.Logf("Contract on chain %s: %s", signer.chain, metaContractAddress.Hex())
	s.signer = signer

}

func (s *SignerSuite) TestSign(c *C) {
	data := []byte("1234")
	tx, sig, hash, err := s.signer.Sign(data, s.signer.tssSigner.Address(), 109, big.NewInt(2), 23)
	_ = tx
	c.Assert(err, IsNil)
	pubkey, err := crypto.Ecrecover(hash, sig)
	c.Assert(err, IsNil)
	c.Assert(bytes.Equal(pubkey, s.signer.tssSigner.Pubkey()), Equals, true)
}

//func (s *SignerSuite) TestMint(c *C) {
//	sendHash, err := hex.DecodeString(config.TSS_TEST_PRIVKEY)
//	c.Assert(err, IsNil)
//	c.Assert(len(sendHash), Equals, 32)
//	var sendHashBytes [32]byte
//	copy(sendHashBytes[:32], sendHash[:32])
//	tssAddr := ethcommon.HexToAddress(config.TSS_TEST_ADDRESS)
//	nonce, err := s.signer.client.NonceAt(context.TODO(), tssAddr, nil)
//	c.Assert(err, IsNil)
//	txhash, err := s.signer.MMint(big.NewInt(1234), ethcommon.HexToAddress(config.TEST_RECEIVER), 80000, []byte{}, sendHashBytes, nonce, big.NewInt(10_000_000_000))
//	c.Assert(err, IsNil)
//	c.Logf("txhash %s", txhash)
//}
