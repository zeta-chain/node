package zetaclient

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
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

// return 33B compressed pubkey
func (s TestSigner) PubKeyCompressedBytes() []byte {
	pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	pk, err := btcec.ParsePubKey(pkBytes, btcec.S256())
	if err != nil {
		panic(err)
	}
	return pk.SerializeCompressed()
}

func (s TestSigner) EVMAddress() ethcommon.Address {
	return crypto.PubkeyToAddress(s.PrivKey.PublicKey)
}

func (s TestSigner) BTCAddress() string {
	pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	pk, err := btcec.ParsePubKey(pkBytes, btcec.S256())
	if err != nil {
		return ""
	}
	testnet3Addr, err := btcutil.NewAddressPubKey(pk.SerializeCompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		return ""
	}
	return testnet3Addr.EncodeAddress()
}

func (s TestSigner) BTCAddressPubkey() *btcutil.AddressPubKey {
	pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	pk, err := btcec.ParsePubKey(pkBytes, btcec.S256())
	if err != nil {
		fmt.Printf("error parsing pubkey: %v", err)
		return nil
	}
	testnet3Addr, err := btcutil.NewAddressPubKey(pk.SerializeCompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Printf("error NewAddressPubKey: %v", err)
		return nil
	}
	return testnet3Addr
}

func (s TestSigner) BTCSegWitAddress() string {
	pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	pk, err := btcec.ParsePubKey(pkBytes, btcec.S256())
	if err != nil {
		return ""
	}

	testnet3Addr, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pk.SerializeCompressed()), &chaincfg.TestNet3Params)
	if err != nil {
		return ""
	}
	return testnet3Addr.EncodeAddress()
}
