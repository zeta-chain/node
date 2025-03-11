package base

import (
	"encoding/hex"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// TSSSignatureEntry represents a single entry in the TSS signature cache.
type TSSSignatureEntry struct {
	// signature is the 65-byte ECDSA TSS signature
	signature [65]byte

	// addedAt is the timestamp when the signature was added to the cache
	addedAt time.Time
}

// TSSSignatureCache stores cached TSS signatures
type TSSSignatureCache struct {
	// expiration is the expiration time for each cache entry
	expiration time.Duration

	// digestToSignature maps a digest to a TSS signature
	digestToSignature *lru.Cache[string, TSSSignatureEntry]
}

// NewTSSSignatureCache creates a new TSS signature cache
// Note: passing 0 expiration will effectively disable the cache
func NewTSSSignatureCache(size int, expiration time.Duration) (*TSSSignatureCache, error) {
	signatureCache, err := lru.New[string, TSSSignatureEntry](size)
	if err != nil {
		return nil, errors.Wrap(err, "error creating tss signature cache")
	}

	return &TSSSignatureCache{
		expiration:        expiration,
		digestToSignature: signatureCache,
	}, nil
}

// Add adds one TSS signature for given public key and digest to the cache.
func (c TSSSignatureCache) Add(pkBech32 string, digest []byte, signature [65]byte) {
	sigKey := signatureKey(pkBech32, digest)
	c.digestToSignature.Add(sigKey, TSSSignatureEntry{
		signature: signature,
		addedAt:   time.Now(),
	})
}

// Get fetches a TSS signature for given public key and digest from the cache.
func (c TSSSignatureCache) Get(pkBech32 string, digest []byte) (signature [65]byte, found bool) {
	sigKey := signatureKey(pkBech32, digest)
	entry, ok := c.digestToSignature.Get(sigKey)
	if !ok {
		return [65]byte{}, false
	}

	// early return if expired
	if c.expire(pkBech32, digest, entry) {
		return [65]byte{}, false
	}

	return entry.signature, true
}

// AddBatch adds multiple TSS signatures for given public key and digests to the cache.
func (c TSSSignatureCache) AddBatch(pkBech32 string, digests [][]byte, signatures [][65]byte) error {
	if len(digests) != len(signatures) {
		return fmt.Errorf("digests and signatures length mismatch: %d != %d", len(digests), len(signatures))
	}

	now := time.Now()
	for i, digest := range digests {
		sigKey := signatureKey(pkBech32, digest)
		c.digestToSignature.Add(sigKey, TSSSignatureEntry{
			signature: signatures[i],
			addedAt:   now,
		})
	}

	return nil
}

// GetBatch fetches TSS signatures for given public key and digests from the cache.
// Note: it returns true only if
//   - all the digests have signatures in the cache.
//   - none of the cached signatures has expired.
func (c TSSSignatureCache) GetBatch(pkBech32 string, digests [][]byte) (signatures [][65]byte, found bool) {
	signatures = make([][65]byte, len(digests))

	for i, digest := range digests {
		sigKey := signatureKey(pkBech32, digest)
		entry, ok := c.digestToSignature.Get(sigKey)
		if !ok {
			return nil, false
		}

		// early return if expired
		// if any one of the cached signatures has expired, all digests need a fresh TSS keysign
		if c.expire(pkBech32, digest, entry) {
			return nil, false
		}

		signatures[i] = entry.signature
	}

	return signatures, true
}

// expire expires the entry if expiration reached
func (c TSSSignatureCache) expire(pkBech32 string, digest []byte, entry TSSSignatureEntry) (expired bool) {
	sigKey := signatureKey(pkBech32, digest)
	if time.Since(entry.addedAt) >= c.expiration {
		c.digestToSignature.Remove(sigKey)
		log.Info().
			Str("pubkey", pkBech32).
			Str("digest", hex.EncodeToString(digest)).
			Str("signature", hex.EncodeToString(entry.signature[:])).
			Msg("tss signature has expired")
		return true
	}

	return false
}

// signatureKey builds the signature key for given public key and digest.
func signatureKey(pkBech32 string, digest []byte) string {
	digestHex := hex.EncodeToString(digest)
	return fmt.Sprintf("%s-%s", pkBech32, digestHex)
}
