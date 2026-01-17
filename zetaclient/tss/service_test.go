package tss_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/go-tss/blame"
	tsscommon "github.com/zeta-chain/go-tss/common"
	"github.com/zeta-chain/go-tss/keysign"

	"github.com/zeta-chain/node/pkg/cosmos"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	"github.com/zeta-chain/node/zetaclient/tss"
)

var (
	base64EncodeString = base64.StdEncoding.EncodeToString
	base64DecodeString = base64.StdEncoding.DecodeString
)

func TestService(t *testing.T) {
	t.Run("NewService", func(t *testing.T) {
		t.Run("Invalid pub key", func(t *testing.T) {
			s, err := tss.NewService(nil, "hello", nil, zerolog.Nop())
			require.ErrorContains(t, err, "invalid tss pub key")
			require.Empty(t, s)
		})

		t.Run("Creates new service", func(t *testing.T) {
			// ARRANGE
			ts := newTestSuite(t)

			// ACT
			s, err := tss.NewService(ts, ts.PubKeyBech32(), ts.zetacore, ts.logger)

			// ASSERT
			require.NoError(t, err)
			require.NotNil(t, s)
			assert.Regexp(t, regexp.MustCompile(`^zetapub.+$`), s.PubKey().Bech32String())
			assert.Equal(t, ts.PubKeyBech32(), s.PubKey().Bech32String())
		})
	})

	t.Run("Sign", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			// ARRANGE
			ts := newTestSuite(t)

			// Given tss service
			s, err := tss.NewService(ts, ts.PubKeyBech32(), ts.zetacore, ts.logger)
			require.NoError(t, err)

			// Given a sample msg to sign
			digest := ts.SampleDigest()

			// Given mock response
			const blockHeight = 123
			ts.keySignerMock.AddCall(ts.PubKeyBech32(), [][]byte{digest}, blockHeight, true, nil)

			// Ensure signature is not cached
			found := s.IsSignatureCached(3, [][]byte{digest})
			require.False(t, found)

			// ACT
			// Sign a message
			// - note that Sign() also contains sig verification
			// - note that Sign() is a wrapper for SignBatch()
			sig, err := s.Sign(ts.ctx, digest, blockHeight, 2, 3)

			// ASSERT
			require.NoError(t, err)
			require.NotEmpty(t, sig)

			// Ensure signature is cached
			found = s.IsSignatureCached(3, [][]byte{digest})
			require.True(t, found)
		})
	})

	t.Run("SignBatch", func(t *testing.T) {
		// ARRANGE
		ts := newTestSuite(t)

		// Given tss service
		s, err := tss.NewService(ts, ts.PubKeyBech32(), ts.zetacore, ts.logger)
		require.NoError(t, err)

		// Given several sample messages to sign
		digests := [][]byte{
			ts.SampleDigest(),
			ts.SampleDigest(),
			ts.SampleDigest(),
			ts.SampleDigest(),
			ts.SampleDigest(),
			ts.SampleDigest(),
			ts.SampleDigest(),
		}

		// Given mock response
		const blockHeight = 123
		ts.keySignerMock.AddCall(ts.PubKeyBech32(), digests, blockHeight, true, nil)

		// Ensure signature is not cached
		found := s.IsSignatureCached(3, digests)
		require.False(t, found)

		// ACT
		sig, err := s.SignBatch(ts.ctx, digests, blockHeight, 2, 3)

		// ASSERT
		require.NoError(t, err)
		require.NotEmpty(t, sig)

		// Ensure signature is cached
		found = s.IsSignatureCached(3, digests)
		require.True(t, found)
	})
}

type testSuite struct {
	*keySignerMock
	ctx      context.Context
	zetacore *mocks.ZetacoreClient
	logger   zerolog.Logger
}

func newTestSuite(t *testing.T) *testSuite {
	return &testSuite{
		keySignerMock: newKeySignerMock(t),
		ctx:           context.Background(),
		zetacore:      mocks.NewZetacoreClient(t),
		logger:        zerolog.New(zerolog.NewTestWriter(t)),
	}
}

func (ts *testSuite) SampleDigest() []byte {
	var digest [32]byte

	_, err := rand.Reader.Read(digest[:])
	require.NoError(ts.t, err)

	return digest[:]
}

type keySignerMock struct {
	t          *testing.T
	privateKey *ecdsa.PrivateKey
	mocks      map[string]lo.Tuple2[keysign.Response, error]
}

func newKeySignerMock(t *testing.T) *keySignerMock {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	return &keySignerMock{
		t:          t,
		privateKey: privateKey,
		mocks:      map[string]lo.Tuple2[keysign.Response, error]{},
	}
}

func (*keySignerMock) Stop() {}

func (m *keySignerMock) PubKeyBech32() string {
	cosmosPrivateKey := &secp256k1.PrivKey{Key: m.privateKey.D.Bytes()}
	pk := cosmosPrivateKey.PubKey()

	bech32, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pk)
	require.NoError(m.t, err)

	return bech32
}

// AddCall mimics TSS signature process (when called with provided arguments)
func (m *keySignerMock) AddCall(pk string, digests [][]byte, height int64, success bool, err error) {
	if success && err != nil {
		m.t.Fatalf("success and error are mutually exclusive")
	}

	var (
		msgs = lo.Map(digests, func(digest []byte, _ int) string {
			return base64EncodeString(digest)
		})

		req = keysign.NewRequest(pk, msgs, height, nil, tss.Version)
		key = m.key(req)

		res keysign.Response
	)

	if !success {
		res = keysign.Response{
			Status: tsscommon.Fail,
			Blame: blame.Blame{
				FailReason: "Ooopsie",
				BlameNodes: []blame.Node{{Pubkey: "some pub key"}},
			},
		}
		m.mocks[key] = lo.Tuple2[keysign.Response, error]{A: res}
		return
	}

	res = m.sign(req)
	m.mocks[key] = lo.Tuple2[keysign.Response, error]{A: res}
}

// sign actually signs the message using local private key instead of TSS
func (m *keySignerMock) sign(req keysign.Request) keysign.Response {
	var signatures []keysign.Signature

	for _, msg := range req.Messages {
		digest, err := base64DecodeString(msg)
		require.NoError(m.t, err)

		// [R || S || V]
		sig, err := crypto.Sign(digest, m.privateKey)
		require.NoError(m.t, err)

		signatures = append(signatures, keysign.Signature{
			Msg:        msg,
			R:          base64EncodeString(sig[:32]),
			S:          base64EncodeString(sig[32:64]),
			RecoveryID: base64EncodeString(sig[64:65]),
		})
	}

	// might be random... we should tolerate that
	signatures = lo.Shuffle(signatures)

	return keysign.Response{
		Signatures: signatures,
		Status:     tsscommon.Success,
	}
}

func (m *keySignerMock) KeySign(req keysign.Request) (keysign.Response, error) {
	key := m.key(req)
	v, ok := m.mocks[key]
	require.True(m.t, ok, "unexpected call KeySign(%+v)", req)

	return v.Unpack()
}

func (m *keySignerMock) key(req keysign.Request) string {
	return fmt.Sprintf("%s-%d:[%+v]", req.PoolPubKey, req.BlockHeight, req.Messages)
}
