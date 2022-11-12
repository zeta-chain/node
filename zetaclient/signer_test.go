package zetaclient

import (
	"github.com/ethereum/go-ethereum/crypto"
	. "gopkg.in/check.v1"
)

type SignerSuite struct {
	signer *TestSigner
}

var _ = Suite(&SignerSuite{})

func (s *SignerSuite) SetUpTest(c *C) {
	// test private key with EVM address
	//// EVM: 0x236C7f53a90493Bb423411fe4117Cb4c2De71DfB
	// BTC testnet3: muGe9prUBjQwEnX19zG26fVRHNi8z7kSPo
	skHex := "7b8507ba117e069f4a3f456f505276084f8c92aee86ac78ae37b4d1801d35fa8"
	privateKey, err := crypto.HexToECDSA(skHex)
	pkBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	c.Logf("pubkey: %d", len(pkBytes))
	// Uncomment the following code to generate new random private key pairs
	//privateKey, err := crypto.GenerateKey()
	//privkeyBytes := crypto.FromECDSA(privateKey)
	//c.Logf("privatekey %s", hex.EncodeToString(privkeyBytes))
	c.Assert(err, IsNil)
	tss := TestSigner{
		PrivKey: privateKey,
	}

	c.Logf("TSS EVMAddress %s", tss.EVMAddress().Hex())
	c.Logf("TSS BTCAddress %s", tss.BTCAddress())

	addr := tss.BTCAddressPubkey()
	c.Logf("TSS tx script: %x", addr.ScriptAddress())
}

func (s *SignerSuite) Test1(c *C) {
	//tss := s.signer
	//pkBytes := tss.Pubkey()
	//c.Logf("len of pkBytes %d", len(pkBytes))
}
