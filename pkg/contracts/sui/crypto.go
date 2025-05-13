package sui

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	secp256k1signer "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

const (
	// flagSecp256k1 is the flag to indicate secp256k1 key scheme.
	flagSecp256k1 = 0x01

	// suiPrivateKeyPrefix is the prefix for Sui private key string.
	suiPrivateKeyPrefix = "suiprivkey"
)

// Digest calculates tx digest (hash) for further signing by TSS.
func Digest(tx models.TxnMetaData) ([32]byte, error) {
	txBytes, err := base64.StdEncoding.DecodeString(tx.TxBytes)
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "failed to decode tx bytes")
	}

	message := messageWithIntentPrefix(txBytes)

	// "When invoking the signing API, you must first hash the intent message of the tx
	// data to 32 bytes using Blake2b ... For ECDSA Secp256k1 and Secp256r1,
	// you must use SHA256 as the internal hash function"
	// https://docs.sui.io/concepts/cryptography/transaction-auth/signatures#signature-requirements
	return blake2b.Sum256(message), nil
}

// https://github.com/MystenLabs/sui/blob/0dc1a38f800fc2d8fabe11477fdef702058cf00d/crates/sui-types/src/intent.rs
// #1 = IntentScope(transactionData=0)
// #2 = Version(0)
// #3 = AppId(Sui=0)
var defaultIntent = []byte{0, 0, 0}

// Constructs binary message with intent prefix.
// https://docs.sui.io/concepts/cryptography/transaction-auth/intent-signing#structs
func messageWithIntentPrefix(message []byte) []byte {
	glued := make([]byte, len(defaultIntent)+len(message))
	copy(glued, defaultIntent)
	copy(glued[len(defaultIntent):], message)

	return glued
}

// AddressFromPubKeyECDSA converts ECDSA public key to Sui address.
// https://docs.sui.io/concepts/cryptography/transaction-auth/keys-addresses
// https://docs.sui.io/concepts/cryptography/transaction-auth/signatures
func AddressFromPubKeyECDSA(pk *ecdsa.PublicKey) string {
	pubBytes := elliptic.MarshalCompressed(pk.Curve, pk.X, pk.Y)

	raw := make([]byte, 1+len(pubBytes))
	raw[0] = flagSecp256k1
	copy(raw[1:], pubBytes)

	addrBytes := blake2b.Sum256(raw)

	return "0x" + hex.EncodeToString(addrBytes[:])
}

// SerializeSignatureECDSA serializes secp256k1 sig (R|S|V) and a publicKey into base64 string
// https://docs.sui.io/concepts/cryptography/transaction-auth/signatures
func SerializeSignatureECDSA(signature [65]byte, pubKey *ecdsa.PublicKey) (string, error) {
	// we don't need the last V byte for recovery
	const sigLen = 64
	const pubKeyLen = 33

	pubKeyBytes := crypto.CompressPubkey(pubKey)

	// should not happen
	if len(pubKeyBytes) != pubKeyLen {
		return "", errors.Errorf("invalid pubKey length (got %d, want %d)", len(pubKeyBytes), pubKeyLen)
	}

	serialized := make([]byte, 1+sigLen+pubKeyLen)
	serialized[0] = flagSecp256k1
	copy(serialized[1:], signature[:sigLen])
	copy(serialized[1+sigLen:], pubKeyBytes)

	return base64.StdEncoding.EncodeToString(serialized), nil
}

// DeserializeSignatureECDSA deserializes SUI secp256k1 signature.
// Returns ECDSA public key and signature.
// Sequence: `flag(1b) + sig(64b) + compressedPubKey(33b)`
func DeserializeSignatureECDSA(sigBase64 string) (*ecdsa.PublicKey, [64]byte, error) {
	const (
		flagLen     = 1
		sigLen      = 64
		pubLen      = 33
		expectedLen = flagLen + sigLen + pubLen
		pubOffset   = flagLen + sigLen
	)

	sigBytes, err := base64.StdEncoding.DecodeString(sigBase64)
	switch {
	case err != nil:
		return nil, [64]byte{}, errors.Wrap(err, "failed to decode signature")
	case len(sigBytes) != expectedLen:
		return nil, [64]byte{}, errors.Errorf("invalid sig length (got %d, want %d)", len(sigBytes), expectedLen)
	case sigBytes[0] != flagSecp256k1:
		return nil, [64]byte{}, errors.Errorf("invalid sig flag (got %d, want %d)", sigBytes[0], flagSecp256k1)
	case len(sigBytes[pubOffset:]) != pubLen:
		return nil, [64]byte{}, errors.Errorf(
			"invalid pubKey length (got %d, want %d)",
			len(sigBytes[pubOffset:]),
			pubLen,
		)
	}

	pk, err := crypto.DecompressPubkey(sigBytes[pubOffset:])
	if err != nil {
		return nil, [64]byte{}, errors.Wrap(err, "failed to decompress public key")
	}

	var sig [64]byte
	copy(sig[:], sigBytes[flagLen:pubOffset])

	return pk, sig, nil
}

// PrivateKeyBech32Secp256k1FromHex converts private key in hex into bech32 format using secp256k1 scheme.
func PrivateKeyBech32Secp256k1FromHex(privKeyHex string) (string, error) {
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode private key hex")
	}

	if len(privKeyBytes) != 32 {
		return "", errors.Errorf("invalid private key length %d)", len(privKeyBytes))
	}

	// add secp256k1 flag
	privKeyBytes = append([]byte{flagSecp256k1}, privKeyBytes...)

	return bech32.ConvertAndEncode(suiPrivateKeyPrefix, privKeyBytes)
}

// SignerSecp256k1 represents Sui Secp256k1 signer.
type SignerSecp256k1 struct {
	pk      *secp256k1.PrivateKey
	address string
}

// NewSignerSecp256k1 creates new SignerSecp256k1.
func NewSignerSecp256k1(privateKeyBytes []byte) *SignerSecp256k1 {
	pk := secp256k1.PrivKeyFromBytes(privateKeyBytes)
	address := AddressFromPubKeyECDSA(pk.PubKey().ToECDSA())

	return &SignerSecp256k1{pk: pk, address: address}
}

func (s *SignerSecp256k1) Address() string {
	return s.address
}

func (s *SignerSecp256k1) SignTxBlock(tx models.TxnMetaData) (string, error) {
	digest, err := Digest(tx)
	if err != nil {
		return "", errors.Wrap(err, "unable to get digest")
	}

	// Another hashing is required for ECDSA.
	// https://docs.sui.io/concepts/cryptography/transaction-auth/signatures#signature-requirements
	digestWrapped := sha256.Sum256(digest[:])

	// returns V[1b] | R[32b] | S[32b], But we need R | S | V
	sig := secp256k1signer.SignCompact(s.pk, digestWrapped[:], false)

	var sigReordered [65]byte
	copy(sigReordered[0:32], sig[1:33])   // Copy R[32]
	copy(sigReordered[32:64], sig[33:65]) // Copy S[32]
	sigReordered[64] = sig[0]             // Move V[1] to the end

	return SerializeSignatureECDSA(sigReordered, &s.pk.ToECDSA().PublicKey)
}
