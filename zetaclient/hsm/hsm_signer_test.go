//go:build hsm_test
// +build hsm_test

package hsm

import (
	"crypto/rand"
	"log"
	"testing"

	"github.com/frumioj/crypto11"

	btcsecp256k1 "github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/stretchr/testify/assert"
	keystone "github.com/zeta-chain/keystone/keys"
)

func TestSignSecp256k1(t *testing.T) {
	// PKCS11 configuration file
	config := &crypto11.Config{
		Path:       "/usr/local/lib/softhsm/libsofthsm2.so",
		TokenLabel: "My token 1",
		Pin:        "1234",
	}

	//Generate random label for key
	label, err := randomBytes(16)
	assert.NoError(t, err)

	//Generate key
	key, err := GenerateKey(string(label), keystone.KEYGEN_SECP256K1, config)
	assert.NoError(t, err)
	assert.NotNil(t, key)

	//Create sample message
	msg := []byte("Signing this plaintext tells me what exactly?")

	signature, err := Sign(config, msg, string(label))
	assert.NoError(t, err)
	assert.NotNil(t, signature)
	assert.Equal(t, key.KeyType(), keystone.KEYGEN_SECP256K1)

	pubkey := key.PubKey()
	secp256k1key := pubkey.(*secp256k1.PubKey)
	pub, err := btcsecp256k1.ParsePubKey(secp256k1key.Key, btcsecp256k1.S256())

	log.Printf("Pub: %v", pub)

	// Validate the signature made by the HSM key, but using the
	// BTC secp256k1 public key
	valid := secp256k1key.VerifySignature(msg, signature)
	log.Printf("Did the signature verify? True = yes: %v", valid)
	log.Printf("TM blockchain address from pubkey: %v", secp256k1key.Address())

	address, pubKey, err := GetHSMAddress(config, string(label))
	log.Printf("Address from HSM: %v, PubKey from HSM: %v", address, pubKey)

	err = key.Delete()
	assert.NoError(t, err)
}

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		log.Printf("Error reading random bytes: %s", err.Error())
		return nil, err
	}

	return b, nil
}
