package zetaclient

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/zeta-chain/zetacore/common"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/chaincfg"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/btcsuite/btcutil"
)

type TSSSigner interface {
	Pubkey() []byte
	// Sign: Specify optionalPubkey to use a different pubkey than the current pubkey set during keygen
	Sign(data []byte, height uint64, nonce uint64, chain *common.Chain, optionalPubkey string) ([65]byte, error)
	EVMAddress() ethcommon.Address
	BTCAddress() string
	BTCAddressWitnessPubkeyHash() *btcutil.AddressWitnessPubKeyHash
	PubKeyCompressedBytes() []byte
}

var _ TSSSigner = (*TestSigner)(nil)

// a fake signer for testing
type TestSigner struct {
	PrivKey *ecdsa.PrivateKey
}

func (s TestSigner) Sign(digest []byte, _ uint64, _ uint64, _ *common.Chain, _ string) ([65]byte, error) {
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
	pk, err := btcec.ParsePubKey(pkBytes)
	if err != nil {
		panic(err)
	}
	return pk.SerializeCompressed()
}

func (s TestSigner) EVMAddress() ethcommon.Address {
	return crypto.PubkeyToAddress(s.PrivKey.PublicKey)
}

func (s TestSigner) BTCAddress() string {
	testnet3Addr := s.BTCAddressPubkey()
	if testnet3Addr == nil {
		return ""
	}
	return testnet3Addr.EncodeAddress()
}

func (s TestSigner) BTCAddressPubkey() *btcutil.AddressPubKey {
	pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	pk, err := btcec.ParsePubKey(pkBytes)
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

func (s TestSigner) BTCAddressWitnessPubkeyHash() *btcutil.AddressWitnessPubKeyHash {
	pkBytes := crypto.FromECDSAPub(&s.PrivKey.PublicKey)
	pk, err := btcec.ParsePubKey(pkBytes)
	if err != nil {
		fmt.Printf("error parsing pubkey: %v", err)
		return nil
	}
	// witness program: https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki#Witness_program
	// The HASH160 of the public key must match the 20-byte witness program.
	addrWPKH, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pk.SerializeCompressed()), &chaincfg.TestNet3Params)
	if err != nil {
		fmt.Printf("error NewAddressWitnessPubKeyHash: %v", err)
		return nil
	}

	return addrWPKH
}
