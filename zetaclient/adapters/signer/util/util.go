package util

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/zeta-chain/go-tss-ctx/keysign"
	"github.com/zeta-chain/go-tss-ctx/tss"
	zcommon "github.com/zeta-chain/zetacore/common/cosmos"
)

func TestKeysign(tssPubkey string, tssServer tss.Server) error {
	log.Info().Msg("trying keysign...")
	data := []byte("hello meta")
	H := crypto.Keccak256Hash(data)
	log.Info().Msgf("hash of data (hello meta) is %s", H)

	keysignReq := keysign.NewRequest(tssPubkey, []string{base64.StdEncoding.EncodeToString(H.Bytes())}, 10, nil, "0.14.0")
	ksRes, err := tssServer.KeySign(keysignReq)
	if err != nil {
		log.Warn().Msg("keysign fail")
	}
	signature := ksRes.Signatures
	// [{cyP8i/UuCVfQKDsLr1kpg09/CeIHje1FU6GhfmyMD5Q= D4jXTH3/CSgCg+9kLjhhfnNo3ggy9DTQSlloe3bbKAs= eY++Z2LwsuKG1JcghChrsEJ4u9grLloaaFZNtXI3Ujk= AA==}]
	// 32B msg hash, 32B R, 32B S, 1B RC
	log.Info().Msgf("signature of helloworld... %v", signature)

	if len(signature) == 0 {
		log.Info().Msgf("signature has length 0, skipping verify")
		return fmt.Errorf("signature has length 0")
	}
	VerifySignature(tssPubkey, signature, H.Bytes())
	if VerifySignature(tssPubkey, signature, H.Bytes()) {
		return nil
	}
	return fmt.Errorf("verify signature fail")
}

func VerifySignature(tssPubkey string, signature []keysign.Signature, H []byte) bool {
	if len(signature) == 0 {
		log.Warn().Msg("verify_signature: empty signature array")
		return false
	}
	pubkey, err := zcommon.GetPubKeyFromBech32(zcommon.Bech32PubKeyTypeAccPub, tssPubkey)
	if err != nil {
		log.Error().Msg("get pubkey from bech32 fail")
	}
	// verify the signature of msg.
	var sigbyte [65]byte
	_, _ = base64.StdEncoding.Decode(sigbyte[:32], []byte(signature[0].R))
	_, _ = base64.StdEncoding.Decode(sigbyte[32:64], []byte(signature[0].S))
	_, _ = base64.StdEncoding.Decode(sigbyte[64:65], []byte(signature[0].RecoveryID))
	sigPublicKey, err := crypto.SigToPub(H, sigbyte[:])
	if err != nil {
		log.Error().Err(err).Msg("SigToPub error in verify_signature")
		return false
	}
	compressedPubkey := crypto.CompressPubkey(sigPublicKey)
	log.Info().Msgf("pubkey %s recovered pubkey %s", pubkey.String(), hex.EncodeToString(compressedPubkey))
	return bytes.Compare(pubkey.Bytes(), compressedPubkey) == 0
}
