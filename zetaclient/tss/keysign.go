package tss

import (
	"bytes"
	"encoding/base64"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	tsscommon "gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keysign"

	"github.com/zeta-chain/node/pkg/cosmos"
	"github.com/zeta-chain/node/zetaclient/logs"
)

var (
	testKeySignData    = []byte("hello meta")
	base64Decode       = base64.StdEncoding.Decode
	base64DecodeString = base64.StdEncoding.DecodeString
)

// TestKeySign performs a TSS key-sign test of sample data.
func TestKeySign(keySigner KeySigner, tssPubKey string, logger zerolog.Logger) error {
	logger = logger.With().Str(logs.FieldModule, "tss_keysign").Logger()

	hashedData := crypto.Keccak256Hash(testKeySignData)

	logger.Info().
		Str("keysign.test_data", string(testKeySignData)).
		Str("keysign.test_data_hash", hashedData.String()).
		Msg("Performing TSS key-sign test")

	req := keysign.NewRequest(
		tssPubKey,
		[]string{base64.StdEncoding.EncodeToString(hashedData.Bytes())},
		10,
		nil,
		Version,
	)

	res, err := keySigner.KeySign(req)
	switch {
	case err != nil:
		return errors.Wrap(err, "key signing request error")
	case res.Status != tsscommon.Success:
		logger.Error().Interface("keysign.fail_blame", res.Blame).Msg("Keysign failed")
		return errors.Wrapf(err, "key signing is not successful (status %d)", res.Status)
	case len(res.Signatures) == 0:
		return errors.New("signatures list is empty")
	}

	// 32B msg hash, 32B R, 32B S, 1B RC
	signature := res.Signatures[0]

	logger.Info().Interface("keysign.signature", signature).Msg("Received signature from TSS")

	if _, err = VerifySignature(signature, tssPubKey, hashedData.Bytes()); err != nil {
		return errors.Wrap(err, "signature verification failed")
	}

	logger.Info().Msg("TSS key-sign test passed")

	return nil
}

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
