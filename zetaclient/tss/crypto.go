package tss

import (
	"bytes"
	"encoding/base64"
	"strings"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"gitlab.com/thorchain/tss/go-tss/keysign"

	"github.com/zeta-chain/node/pkg/cosmos"
)

var (
	base64Decode       = base64.StdEncoding.Decode
	base64DecodeString = base64.StdEncoding.DecodeString
	base64EncodeString = base64.StdEncoding.EncodeToString
)

// VerifySignature checks that keysign.Signature is valid and origins from expected TSS public key.
// Also returns signature as [65]byte (R, S, V)
func VerifySignature(sig keysign.Signature, tssPubKey string, expectedMsgHash []byte) ([65]byte, error) {
	// Check that msg hash equals msg hash in the signature
	actualMsgHash, err := base64DecodeString(sig.Msg)
	switch {
	case err != nil:
		return [65]byte{}, errors.Wrap(err, "unable to decode message hash")
	case !bytes.Equal(expectedMsgHash, actualMsgHash):
		return [65]byte{}, errors.New("message hash mismatch")
	}

	// Prepare expected public key
	expectedPubKey, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, tssPubKey)
	if err != nil {
		return [65]byte{}, errors.Wrap(err, "unable to decode tss pub key from bech32")
	}

	sigBytes, err := SignatureToBytes(sig)
	if err != nil {
		return [65]byte{}, errors.Wrap(err, "unable to convert signature to bytes")
	}

	// Recover public key from signature
	actualPubKey, err := crypto.SigToPub(expectedMsgHash, sigBytes[:])
	if err != nil {
		return [65]byte{}, errors.Wrap(err, "unable to recover public key from signature")
	}

	if !bytes.Equal(expectedPubKey.Bytes(), crypto.CompressPubkey(actualPubKey)) {
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
