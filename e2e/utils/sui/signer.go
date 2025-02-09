package sui

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	secp256k1_ecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"golang.org/x/crypto/blake2b"
)

type SignerSecp256k1 struct {
	privkey *secp256k1.PrivateKey
}

func NewSignerSecp256k1FromSecretKey(secret []byte) *SignerSecp256k1 {
	privKey := secp256k1.PrivKeyFromBytes(secret)

	return &SignerSecp256k1{
		privkey: privKey,
	}
}

// GetPublicKey returns the compressed public key bytes
func (s *SignerSecp256k1) GetPublicKey() []byte {
	pub := s.privkey.PubKey()

	// Create compressed public key format:
	// 0x02/0x03 | x-coordinate (32 bytes)
	x := pub.X().Bytes()

	// Ensure x coordinate is 32 bytes with leading zeros if needed
	paddedX := make([]byte, 32)
	copy(paddedX[32-len(x):], x)

	// Prefix with 0x02 for even Y or 0x03 for odd Y
	prefix := byte(0x02)
	if pub.Y().Bit(0) == 1 {
		prefix = 0x03
	}

	return append([]byte{prefix}, paddedX...)
}

// GetFlaggedPublicKey returns GetPublicKey flagged for use with wallets
// and command line tools
func (s *SignerSecp256k1) GetFlaggedPublicKey() []byte {
	return append([]byte{signer.SigntureFlagSecp256k1}, s.GetPublicKey()...)
}

func (s *SignerSecp256k1) Address() string {
	// Get the public key bytes
	pubKeyBytes := s.GetPublicKey()

	// Create BLAKE2b hash
	hash, err := blake2b.New256(nil)
	if err != nil {
		// This will never happen with nil key
		panic(err)
	}

	hash.Write([]byte{signer.SigntureFlagSecp256k1})
	hash.Write(pubKeyBytes)
	addrBytes := hash.Sum(nil)
	// convert to 0x hex
	return "0x" + hex.EncodeToString(addrBytes)
}

type SignedMessageSerializedSig struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

func ToSerializedSignature(signature, pubKey []byte) string {
	signatureLen := len(signature)
	pubKeyLen := len(pubKey)
	serializedSignature := make([]byte, 1+signatureLen+pubKeyLen)
	serializedSignature[0] = byte(signer.SigntureFlagSecp256k1)
	copy(serializedSignature[1:], signature)
	copy(serializedSignature[1+signatureLen:], pubKey)
	return base64.StdEncoding.EncodeToString(serializedSignature)
}

// SignTransactionBlock signs an encoded transaction block
func (s *SignerSecp256k1) SignTransactionBlock(txBytesEncoded string) (string, error) {
	txBytes, err := base64.StdEncoding.DecodeString(txBytesEncoded)
	if err != nil {
		return "", fmt.Errorf("decode tx bytes: %w", err)
	}
	message := messageWithIntent(txBytes)
	digest1 := blake2b.Sum256(message)
	// this additional hash is required for secp256k1 but not ed25519
	digest2 := sha256.Sum256(digest1[:])

	sigBytes := secp256k1_ecdsa.SignCompact(s.privkey, digest2[:], false)
	// Take R and S, skip recovery byte
	sigBytes = sigBytes[1:]

	signature := ToSerializedSignature(sigBytes, s.GetPublicKey())
	return signature, nil
}

func messageWithIntent(message []byte) []byte {
	intent := IntentBytes
	intentMessage := make([]byte, len(intent)+len(message))
	copy(intentMessage, intent)
	copy(intentMessage[len(intent):], message)
	return intentMessage
}

var IntentBytes = []byte{0, 0, 0}
