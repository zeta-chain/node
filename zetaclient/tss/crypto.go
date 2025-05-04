package tss

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	eth "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/zeta-chain/go-tss/keysign"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/contracts/sui"
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
	b, err := hex.DecodeString(strings.TrimPrefix(raw, "0x"))
	if err != nil {
		return PubKey{}, errors.Wrap(err, "unable to decode hex string")
	}

	pk, err := crypto.UnmarshalPubkey(b)
	if err != nil {
		return PubKey{}, errors.Wrap(err, "unable to unmarshal pubkey")
	}

	return NewPubKeyFromECDSA(*pk)
}

func (k PubKey) AsECDSA() *ecdsa.PublicKey {
	return k.ecdsaPubKey
}

// Bytes marshals pubKey to bytes either as compressed or uncompressed slice.
//
// In ECDSA, a compressed pubKey includes only the X and a parity bit for the Y,
// allowing the full Y to be reconstructed using the elliptic curve equation,
// thus reducing the key size while maintaining the ability to fully recover the pubKey.
func (k PubKey) Bytes(compress bool) []byte {
	pk := k.ecdsaPubKey
	if compress {
		return crypto.CompressPubkey(pk)
	}

	return crypto.FromECDSAPub(pk)
}

// Bech32String returns the bech32 string of the public key.
// Example: `zetapub1addwnpepq2fdhcmfyv07s86djjca835l4f2n2ta0c7le6vnl508mseca2s9g6slj0gm`
func (k PubKey) Bech32String() string {
	v, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, k.cosmosPubKey)
	if err != nil {
		return "" // not possible
	}

	return v
}

// AddressBTC returns the Bitcoin address of the public key.
func (k PubKey) AddressBTC(chainID int64) (*btcutil.AddressWitnessPubKeyHash, error) {
	return bitcoinP2WPKH(k.Bytes(true), chainID)
}

// BTCPayToAddrScript returns the script for the Bitcoin TSS address.
func (k PubKey) BTCPayToAddrScript(chainID int64) ([]byte, error) {
	tssAddrP2WPKH, err := k.AddressBTC(chainID)
	if err != nil {
		return nil, err
	}
	return txscript.PayToAddrScript(tssAddrP2WPKH)
}

// AddressEVM returns the ethereum address of the public key.
func (k PubKey) AddressEVM() eth.Address {
	return crypto.PubkeyToAddress(*k.ecdsaPubKey)
}

// AddressSui returns Sui address of the public key.
func (k PubKey) AddressSui() string {
	return sui.AddressFromPubKeyECDSA(k.ecdsaPubKey)
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
		return [65]byte{}, errors.Errorf("message hash mismatch (got 0x%x, want 0x%x)", actualMsgHash, hash)
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
		return [65]byte{}, errors.Errorf(
			"public key mismatch (got %s, want %s)",
			crypto.PubkeyToAddress(*actualPubKey),
			pk.AddressEVM(),
		)
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

// apparently go-tss returns res.Signatures in a different order than digests,
// thus we need to ensure the order AND verify the signatures
func verifySignatures(digests [][]byte, res keysign.Response, pk PubKey) ([][65]byte, error) {
	switch {
	case len(digests) == 0:
		return nil, errors.New("empty digests list")
	case len(digests) != len(res.Signatures):
		return nil, errors.Errorf("length mismatch (got %d, want %d)", len(res.Signatures), len(digests))
	case len(digests) == 1:
		// most common case
		sig, err := VerifySignature(res.Signatures[0], pk, digests[0])
		if err != nil {
			return nil, err
		}

		return [][65]byte{sig}, nil
	}

	// map bas64(digest) => slice index
	cache := make(map[string]int, len(digests))
	for i, digest := range digests {
		cache[base64EncodeString(digest)] = i
	}

	signatures := make([][65]byte, len(res.Signatures))

	for _, sigResponse := range res.Signatures {
		i, ok := cache[sigResponse.Msg]
		if !ok {
			return nil, errors.Errorf("missing digest %s", sigResponse.Msg)
		}

		sig, err := VerifySignature(sigResponse, pk, digests[i])
		if err != nil {
			return nil, fmt.Errorf("unable to verify signature: %w (#%d)", err, i)
		}

		signatures[i] = sig
	}

	return signatures, nil
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
