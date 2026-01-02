package base

import (
	"fmt"

	mathpkg "github.com/zeta-chain/node/pkg/math"
)

const (
	// batchSize is the maximum number of digests in a keysign batch.
	// signing a 10-digest batch takes about 3~4 seconds on average, and
	// increasing batch size won't result in significant performance boost
	batchSize = 10

	// zetaHeightDivisor is the divisor to map zeta height to keysign height bucket.
	// 10 blocks (apprx. 20 seconds) is used to avoid keysign height drifting too fast.
	zetaHeightDivisor = 10
)

// TSSKeysignInfo represents a record of TSS keysign information.
type TSSKeysignInfo struct {
	// digest is the digest of the outbound transaction
	digest []byte

	// signature is the TSS signature of the digest
	signature [65]byte
}

// NewTSSKeysignInfo creates a new TSS keysign information.
func NewTSSKeysignInfo(digest []byte, signature [65]byte) *TSSKeysignInfo {
	return &TSSKeysignInfo{
		digest:    digest,
		signature: signature,
	}
}

// TSSKeysignBatch contains a batch of TSS keysign information.
type TSSKeysignBatch struct {
	// digests is a list of digests to sign
	digests [][]byte

	// nonceLow is the lowest cctx nonce in the batch
	nonceLow uint64

	// nonceHigh is the highest cctx nonce in the batch
	nonceHigh uint64
}

// NewTSSKeysignBatch creates a new TSS keysign batch.
func NewTSSKeysignBatch() *TSSKeysignBatch {
	return &TSSKeysignBatch{
		digests: make([][]byte, 0),
	}
}

// BatchNumber returns the batch number of the keysign batch
func (b *TSSKeysignBatch) BatchNumber() uint64 {
	return NonceToBatchNumber(b.nonceLow)
}

// Digests returns the digests in the keysign batch
func (b *TSSKeysignBatch) Digests() [][]byte {
	return b.digests
}

// NonceLow returns the lowest nonce in the keysign batch
func (b *TSSKeysignBatch) NonceLow() uint64 {
	return b.nonceLow
}

// NonceHigh returns the highest nonce in the keysign batch
func (b *TSSKeysignBatch) NonceHigh() uint64 {
	return b.nonceHigh
}

// AddKeysignInfo adds one TSS keysign information to the batch.
func (b *TSSKeysignBatch) AddKeysignInfo(nonce uint64, info TSSKeysignInfo) {
	b.digests = append(b.digests, info.digest)

	// initialize on first record
	if len(b.digests) == 1 {
		b.nonceLow = nonce
		b.nonceHigh = nonce
		return
	}

	// update nonceLow and heightLow
	if nonce < b.nonceLow {
		b.nonceLow = nonce
	} else if nonce > b.nonceHigh {
		b.nonceHigh = nonce
	}
}

// IsEmpty returns true if the batch is empty.
func (b *TSSKeysignBatch) IsEmpty() bool {
	return len(b.digests) == 0
}

// IsSequential returns true if the batch is sequential (no gaps in between).
// To make TSS keysign deterministic:
// - we ALWAYS sign sequential nonces, e.g.: [0,1,2,3,4], [7,8,9], [10,11,12,13], [14], [15,16,17,18,19], [20,21], ...
// - we NEVER sign nonces with gaps,  e.g.: [0,1,3,4], [5,6,7,9], [10,12,13], [14,15,16,18,19], ...
func (b *TSSKeysignBatch) IsSequential() bool {
	// #nosec G115 - always positive
	return uint64(len(b.digests)) == b.nonceHigh-b.nonceLow+1
}

// IsEnd returns true if the batch hits the end of the batch.
// For example: [6,7,8,9] is end of batch 0, [18,19] is end of batch 1, ...
func (b *TSSKeysignBatch) IsEnd() bool {
	_, batchNonceHigh := BatchNumberToRange(b.BatchNumber())
	return b.nonceHigh == batchNonceHigh
}

// ContainsNonce returns true if the batch contains the given nonce.
func (b *TSSKeysignBatch) ContainsNonce(nonce uint64) bool {
	return nonce >= b.nonceLow && nonce <= b.nonceHigh
}

// KeysignHeight calculates an artificial keysign height tweaked by chainID.
// This is used to uniquely identify TSS keysign request across all chains.
func KeysignHeight(chainID int64, zetaHeight int64) (uint64, error) {
	heightBucket := zetaHeight / zetaHeightDivisor
	if zetaHeight <= 0 || heightBucket >= mathpkg.MaxPairValue {
		return 0, fmt.Errorf("invalid zeta height: %d", zetaHeight)
	}

	if chainID <= 0 || chainID > mathpkg.MaxPairValue {
		return 0, fmt.Errorf("invalid chain ID: %d", chainID)
	}

	// use height bucket as unique identifier, added 1 to avoid 0
	// #nosec G115 - checked in range
	identifier := uint32(heightBucket + 1)

	// #nosec G115 - checked in range
	chainID32 := uint32(chainID)

	// create artificial height using pairing function
	return mathpkg.CantorPair(identifier, chainID32), nil
}

// NonceToBatchNumber maps a nonce to a batch number.
// For example:
// - nonce 0 falls into batch 0
// - nonce 9 falls into batch 0
// - nonce 10 falls into batch 1
// - nonce 19 falls into batch 1
// - nonce 20 falls into batch 2
// - ...
func NonceToBatchNumber(nonce uint64) uint64 {
	return nonce / batchSize
}

// BatchNumberToRange returns the range of nonces for the given batch number.
// Example ranges are: [0, 9], [10, 19], [20, 29], [30, 39], ...
func BatchNumberToRange(batchNumber uint64) (uint64, uint64) {
	return batchNumber * batchSize, (batchNumber+1)*batchSize - 1
}
