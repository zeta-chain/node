package metaclientd

import (
	"bytes"
	"crypto/ecdsa"
	"github.com/Meta-Protocol/metacore/common"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	. "gopkg.in/check.v1"
	"math/big"
)

type SignerSuite struct {
	signer *Signer
}

var _ = Suite(&SignerSuite{})

type testSigner struct {
	privkey *ecdsa.PrivateKey
}

func (s testSigner) Sign(digest []byte) [65]byte{
	sig, _ := crypto.Sign(digest, s.privkey)
	var sigbyte [65]byte
	copy(sigbyte[:], sig[:65])
	return sigbyte
}

func (s testSigner) Pubkey() []byte {
	publicKeyBytes := crypto.FromECDSAPub(&s.privkey.PublicKey)
	return publicKeyBytes
}


func (s testSigner) Address() ethcommon.Address {
	return crypto.PubkeyToAddress(s.privkey.PublicKey)
}

func (s *SignerSuite) SetUpTest(c *C) {
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	c.Assert(err, IsNil)
	tss := testSigner {
		privkey: privateKey,
	}
	signer, err := NewSigner(common.Chain("ETH"), ETH_ENDPOINT, tss.Address(), tss)
	c.Assert(err, IsNil)
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
