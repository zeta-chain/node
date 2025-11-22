package mocks

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	_ "github.com/zeta-chain/node/pkg/sdkconfig/default" //nolint:blank-imports this is a test package
	zetatss "github.com/zeta-chain/node/zetaclient/tss"
)

type test = require.TestingT

type TSS struct {
	t          test
	privateKey *ecdsa.PrivateKey
	fakePubKey *zetatss.PubKey
	paused     bool
}

func NewTSS(t *testing.T) *TSS {
	pk, err := crypto.GenerateKey()
	require.NoError(t, err)

	return &TSS{t: t, privateKey: pk}
}

func NewTSSFromPrivateKey(t test, pk *ecdsa.PrivateKey) *TSS {
	return &TSS{t: t, privateKey: pk}
}

func (tss *TSS) PubKey() zetatss.PubKey {
	if tss.fakePubKey != nil {
		return *tss.fakePubKey
	}

	pubKey, err := zetatss.NewPubKeyFromECDSA(tss.privateKey.PublicKey)
	require.NoError(tss.t, err)

	return pubKey
}

func (tss *TSS) FakePubKey(pk any) *TSS {
	if pk == nil {
		tss.fakePubKey = nil
		return tss
	}

	if zpk, ok := pk.(zetatss.PubKey); ok {
		tss.fakePubKey = &zpk
		return tss
	}

	raw, ok := pk.(string)
	require.True(tss.t, ok, "invalid type for fake pub key (%v)", pk)

	if strings.HasPrefix(raw, "zetapub") {
		zpk, err := zetatss.NewPubKeyFromBech32(raw)
		require.NoError(tss.t, err)
		tss.fakePubKey = &zpk
		return tss
	}

	if strings.HasPrefix(raw, "0x") {
		zpk, err := zetatss.NewPubKeyFromECDSAHexString(raw)
		require.NoError(tss.t, err)
		tss.fakePubKey = &zpk
		return tss
	}

	tss.t.Errorf("invalid fake pub key format: %s", raw)
	tss.t.FailNow()

	return nil
}

func (tss *TSS) Sign(_ context.Context, digest []byte, _, _ uint64, _ int64) ([65]byte, error) {
	sigs, err := tss.SignBatch(context.Background(), [][]byte{digest}, 0, 0, 0)
	if err != nil {
		return [65]byte{}, err
	}

	return sigs[0], nil
}

func (tss *TSS) SignBatch(_ context.Context, digests [][]byte, _, _ uint64, _ int64) ([][65]byte, error) {
	// just for backwards compatibility (ideally we should remove this)
	if tss.paused {
		return nil, errors.New("tss is paused")
	}

	sigs := [][65]byte{}

	for _, digest := range digests {
		sigBytes, err := crypto.Sign(digest, tss.privateKey)
		require.NoError(tss.t, err)
		require.Len(tss.t, sigBytes, 65)

		var sig [65]byte
		copy(sig[:], sigBytes)

		sigs = append(sigs, sig)
	}

	return sigs, nil
}

func (tss *TSS) IsSignatureCached(_ int64, _ [][]byte) bool {
	return false
}

func (tss *TSS) UpdatePrivateKey(pk *ecdsa.PrivateKey) {
	tss.privateKey = pk
}

func (tss *TSS) Pause() {
	tss.paused = true
}

func (tss *TSS) Unpause() {
	tss.paused = false
}
