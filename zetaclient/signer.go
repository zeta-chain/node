package zetaclient

import (
	"crypto/ecdsa"
	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/btcsuite/btcutil"
)

type TSSSigner interface {
	Pubkey() []byte
	Sign(data []byte) ([65]byte, error)
	EVMAddress() ethcommon.Address
	BTCAddress() string
}

// a fake signer for testing
type TestSigner struct {
	PrivKey *ecdsa.PrivateKey
}

func (s TestSigner) Sign(digest []byte) ([65]byte, error) {
	sig, err := crypto.Sign(digest, s.PrivKey)
	if err != nil {
		return [65]byte{}, err
	}
	var sigbyte [65]byte
	copy(sigbyte[:], sig[:65])
	return sigbyte, nil
}

func (s TestSigner) Pubkey() []byte {
	publicKeyBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	return publicKeyBytes
}

func (s TestSigner) EVMAddress() ethcommon.Address {
	return crypto.PubkeyToAddress(s.PrivKey.PublicKey)
}

func (s TestSigner) BTCAddress() string {
	pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	testnet3Addr, err := btcutil.NewAddressPubKey(pkBytes, &chaincfg.TestNet3Params)
	if err != nil {
		panic(err)
	}
	return testnet3Addr.EncodeAddress()
}

func (s TestSigner) BTCAddressPubkey() *btcutil.AddressPubKey {
	pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	testnet3Addr, err := btcutil.NewAddressPubKey(pkBytes, &chaincfg.TestNet3Params)
	if err != nil {
		panic(err)
	}
	return testnet3Addr
}
