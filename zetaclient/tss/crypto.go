package tss

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/tss/go-tss/keysign"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/cosmos"
)

// PubKey represents TSS public key in various formats.
type PubKey struct {
	cosmosPubKey cryptotypes.PubKey
	ecdsaPubKey  *ecdsa.PublicKey
}

var (
	base64Decode       = base64.StdEncoding.Decode
	base64DecodeString = base64.StdEncoding.DecodeString
	base64EncodeString = base64.StdEncoding.EncodeToString
)

// NewPubKeyFromBech32 creates a new PubKey from a bech32 address.
// Example: `zetapub1addwnpepq2fdhcmfyv07s86djjca835l4f2n2ta0c7le6vnl508mseca2s9g6slj0gm`
func NewPubKeyFromBech32(bech32 string) (PubKey, error) {
	if bech32 == "" {
		return PubKey{}, errors.New("empty bech32 address")
	}

	cosmosPubKey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, bech32)
	if err != nil {
		return PubKey{}, errors.Wrap(err, "unable to GetPubKeyFromBech32")
	}

	pubKey, err := crypto.DecompressPubkey(cosmosPubKey.Bytes())
	if err != nil {
		return PubKey{}, errors.Wrap(err, "unable to DecompressPubkey")
	}

	return PubKey{
		cosmosPubKey: cosmosPubKey,
		ecdsaPubKey:  pubKey,
	}, nil
}

// NewPubKeyFromECDSA creates a new PubKey from an ECDSA public key.
func NewPubKeyFromECDSA(pk ecdsa.PublicKey) (PubKey, error) {
	compressed := elliptic.MarshalCompressed(pk.Curve, pk.X, pk.Y)

	return PubKey{
		cosmosPubKey: &secp256k1.PubKey{Key: compressed},
		ecdsaPubKey:  &pk,
	}, nil
}

// NewPubKeyFromECDSAHexString creates PubKey from 0xABC12...
func NewPubKeyFromECDSAHexString(raw string) (PubKey, error) {
	if strings.HasPrefix(raw, "0x") {
		raw = raw[2:]
	}

	b, err := hex.DecodeString(raw)
	if err != nil {
		return PubKey{}, errors.Wrap(err, "unable to decode hex string")
	}

	pk, err := crypto.UnmarshalPubkey(b)
	if err != nil {
		return PubKey{}, errors.Wrap(err, "unable to unmarshal pubkey")
	}

	return NewPubKeyFromECDSA(*pk)
}

// Bytes marshals pubKey to bytes either as compressed or uncompressed slice.
//
// In ECDSA, a compressed pubKey includes only the X and a parity bit for the Y,
// allowing the full Y to be reconstructed using the elliptic curve equation,
// thus reducing the key size while maintaining the ability to fully recover the pubKey.
func (k PubKey) Bytes(compress bool) []byte {
	pk := k.ecdsaPubKey
	if compress {
		return elliptic.MarshalCompressed(pk.Curve, pk.X, pk.Y)
	}

	return crypto.FromECDSAPub(pk)
}

// Bech32String returns the bech32 string of the public key.
// Example: `zetapub1addwnpepq2fdhcmfyv07s86djjca835l4f2n2ta0c7le6vnl508mseca2s9g6slj0gm`
func (k PubKey) Bech32String() string {
	v, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, k.cosmosPubKey)

	// should not happen as we only set k.cosmosPubKey from the constructor
	if err != nil {
		panic("PubKey.Bech32String: " + err.Error())
	}

	return v
}

// AddressBTC returns the bitcoin address of the public key.
func (k PubKey) AddressBTC(chainID int64) (*btcutil.AddressWitnessPubKeyHash, error) {
	return bitcoinP2WPKH(k.Bytes(true), chainID)
}

// AddressEVM returns the ethereum address of the public key.
func (k PubKey) AddressEVM() eth.Address {
	return crypto.PubkeyToAddress(*k.ecdsaPubKey)
}

// VerifySignature checks that keysign.Signature is valid and origins from expected TSS public key.
// Also returns signature as [65]byte (R, S, V)
func VerifySignature(sig keysign.Signature, pk PubKey, hash []byte) ([65]byte, error) {
	// Check that msg hash equals msg hash in the signature
	actualMsgHash, err := base64DecodeString(sig.Msg)
	switch {
	case err != nil:
		return [65]byte{}, errors.Wrap(err, "unable to decode message hash")
	case !bytes.Equal(hash, actualMsgHash):
		return [65]byte{}, errors.New("message hash mismatch")
	}

	sigBytes, err := SignatureToBytes(sig)
	if err != nil {
		return [65]byte{}, errors.Wrap(err, "unable to convert signature to bytes")
	}

	// Recover public key from signature
	actualPubKey, err := crypto.SigToPub(hash, sigBytes[:])
	switch {
	case err != nil:
		return [65]byte{}, errors.Wrap(err, "unable to recover public key from signature")
	case crypto.PubkeyToAddress(*actualPubKey) != pk.AddressEVM():
		return [65]byte{}, errors.New("public key mismatch")
	}

	return sigBytes, nil
}

// SignatureToBytes converts keysign.Signature to [65]byte (R, S, V)
func SignatureToBytes(input keysign.Signature) (sig [65]byte, err error) {
	if _, err = base64Decode(sig[:32], []byte(input.R)); err != nil {
		return sig, errors.Wrap(err, "unable to decode R")
	}

	if _, err = base64Decode(sig[32:64], []byte(input.S)); err != nil {
		return sig, errors.Wrap(err, "unable to decode S")
	}

	if _, err = base64Decode(sig[64:65], []byte(input.RecoveryID)); err != nil {
		return sig, errors.Wrap(err, "unable to decode RecoveryID (V)")
	}

	return sig, nil
}

// combineDigests combines the digests
func combineDigests(digestList []string) []byte {
	digestConcat := strings.Join(digestList, "")
	digestBytes := chainhash.DoubleHashH([]byte(digestConcat))
	return digestBytes.CloneBytes()
}

// bitcoinP2WPKH returns P2WPKH (pay to witness pub key hash) address from the compressed pub key.
func bitcoinP2WPKH(pkCompressed []byte, chainID int64) (*btcutil.AddressWitnessPubKeyHash, error) {
	params, err := chains.BitcoinNetParamsFromChainID(chainID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get btc net params")
	}

	hash := btcutil.Hash160(pkCompressed)

	return btcutil.NewAddressWitnessPubKeyHash(hash, params)
}
