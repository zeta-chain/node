package base

import (
	goctx "context"
	"errors"
	"testing"
	"time"

	cometbft "github.com/cometbft/cometbft/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/scheduler"
	"github.com/zeta-chain/node/testutil/sample"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

func Test_IsStaleBlockEvent(t *testing.T) {
	tests := []struct {
		name                string
		eventHeight         int64
		zetaHeight          int64
		mockContext         bool
		mockZetaHeight      bool
		mockZetaHeightError error
		expectStale         bool
		expectHeight        int64
		errorMsg            string
	}{
		{
			name:           "stale block event",
			eventHeight:    100,
			zetaHeight:     101,
			mockContext:    true,
			mockZetaHeight: true,
			expectStale:    true,
			expectHeight:   101,
		},
		{
			name:           "not stale block event",
			eventHeight:    100,
			zetaHeight:     100,
			mockContext:    true,
			mockZetaHeight: true,
			expectStale:    false,
			expectHeight:   100,
		},
		{
			name:           "error getting block from context",
			eventHeight:    100,
			zetaHeight:     100,
			mockContext:    false,
			mockZetaHeight: false,
			errorMsg:       "unable to get block event from context",
		},
		{
			name:                "error getting zeta height",
			eventHeight:         100,
			zetaHeight:          0,
			mockContext:         true,
			mockZetaHeight:      true,
			mockZetaHeightError: errors.New("mock error"),
			errorMsg:            "unable to get zeta height",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			signer := newSignerTestSuite(t)
			zetacore := mocks.NewZetacoreClient(t)
			zetaRepo := zrepo.New(zetacore, chains.Ethereum, mode.StandardMode)

			ctx := goctx.Background()
			appCtx := zctx.New(config.New(false), nil, zerolog.Nop())
			ctx = zctx.WithAppContext(ctx, appCtx)

			// Mock context with block event
			if tc.mockContext {
				ctx = scheduler.WithBlockEvent(ctx, cometbft.EventDataNewBlock{
					Block: &cometbft.Block{
						Header: cometbft.Header{Height: tc.eventHeight, Time: time.Now()},
					},
				})
			}

			// Mock zeta height
			if tc.mockZetaHeight {
				zetacore.On("GetBlockHeight", mock.Anything).Return(tc.zetaHeight, tc.mockZetaHeightError).Once()
			}

			// ACT
			height, isStale, err := signer.IsStaleBlockEvent(ctx, zetaRepo)

			// ASSERT
			if tc.errorMsg != "" {
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectHeight, height)
				require.Equal(t, tc.expectStale, isStale)
			}
		})
	}
}

func Test_IsTimeToKeysign(t *testing.T) {
	tests := []struct {
		name               string
		nextTSSNonce       uint64
		zetaHeight         int64
		scheduleInterval   int64
		pendingNonces      observertypes.PendingNonces
		pendingNoncesError error
		staleNonces        []uint64
		shouldSign         bool
		errorMsg           string
	}{
		{
			name:             "it's time to sign",
			nextTSSNonce:     0,
			zetaHeight:       100,
			scheduleInterval: 10,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  5,
				NonceHigh: 10,
			},
			staleNonces: []uint64{1, 2, 3, 4},
			shouldSign:  true,
		},
		{
			name:             "not time to sign - height not multiple of interval",
			nextTSSNonce:     0,
			zetaHeight:       101,
			scheduleInterval: 10,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  5,
				NonceHigh: 10,
			},
			shouldSign: false,
		},
		{
			name:             "not time to sign - no pending cctx",
			nextTSSNonce:     0,
			zetaHeight:       100,
			scheduleInterval: 10,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  5,
				NonceHigh: 5,
			},
			staleNonces: []uint64{1, 2, 3, 4},
			shouldSign:  false,
		},
		{
			name:             "not time to sign - TSS nonce is ahead",
			nextTSSNonce:     10,
			zetaHeight:       100,
			scheduleInterval: 10,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  5,
				NonceHigh: 10,
			},
			staleNonces: []uint64{1, 2, 3, 4},
			shouldSign:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			signer := newSignerTestSuite(t)

			// mock pending keysign info in cache
			// #nosec G115 - always positive
			for nonce := uint64(tc.pendingNonces.NonceLow); nonce < uint64(tc.pendingNonces.NonceHigh); nonce++ {
				signer.GetSignatureOrAddDigest(nonce, sample.Digest32B(t))
			}

			// ACT
			shouldSign := signer.IsTimeToKeysign(tc.pendingNonces, tc.nextTSSNonce, tc.zetaHeight, tc.scheduleInterval)

			// ASSERT
			require.Equal(t, tc.shouldSign, shouldSign)
		})
	}
}

func Test_GetKeysignBatch(t *testing.T) {
	tests := []struct {
		name              string
		batchNumber       uint64
		pendingNonces     observertypes.PendingNonces
		keysignInfoNonces []uint64
		failRPC           bool
		expectBatch       bool
		expectLength      int
	}{
		{
			name:        "batch ready - successful collection",
			batchNumber: 0,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  0,
				NonceHigh: 5,
			},
			keysignInfoNonces: []uint64{0, 1, 2, 3, 4},
			expectBatch:       true,
			expectLength:      5,
		},
		{
			name:        "batch not ready - no overlap",
			batchNumber: 1,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  0,
				NonceHigh: 9,
			},
			expectBatch: false,
		},
		{
			name:        "batch not ready - no pending cctx",
			batchNumber: 0,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  5,
				NonceHigh: 5,
			},
			expectBatch: false,
		},
		{
			name:        "batch ready - waiting for digests",
			batchNumber: 0,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  0,
				NonceHigh: 5,
			},
			keysignInfoNonces: []uint64{}, // no info in cache
			expectBatch:       false,
		},
		{
			name:        "batch ready - waiting for gaps",
			batchNumber: 0,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  0,
				NonceHigh: 5,
			},
			keysignInfoNonces: []uint64{0, 1, 2, 4}, // with gap 3
			expectBatch:       false,
		},
		{
			name:        "batch ready - waiting until nonce",
			batchNumber: 0,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  0,
				NonceHigh: 5,
			},
			keysignInfoNonces: []uint64{0, 1, 2}, // no gap, but missing nonce 3, 4
			expectBatch:       false,
		},
		{
			name:        "error getting pending nonces",
			batchNumber: 0,
			pendingNonces: observertypes.PendingNonces{
				NonceLow:  0,
				NonceHigh: 5,
			},
			failRPC: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			signer := newSignerTestSuite(t)
			zetacore := mocks.NewZetacoreClient(t)
			zetaRepo := zrepo.New(zetacore, chains.Ethereum, mode.StandardMode)

			if tc.failRPC {
				zetacore.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).
					Return(observertypes.PendingNonces{}, errors.New("rpc error")).Maybe()
			} else {
				zetacore.On("GetPendingNoncesByChain", mock.Anything, mock.Anything).
					Return(tc.pendingNonces, nil).Maybe()
			}

			// Mock keysign info in cache
			for _, nonce := range tc.keysignInfoNonces {
				signer.GetSignatureOrAddDigest(nonce, sample.Digest32B(t))
			}

			// ACT
			ctx := goctx.Background()
			batch := signer.GetKeysignBatch(ctx, zetaRepo, tc.batchNumber)

			// ASSERT
			if tc.expectBatch {
				require.NotNil(t, batch)
				require.Len(t, batch.Digests(), tc.expectLength)
			} else {
				require.Nil(t, batch)
			}
		})
	}
}

func Test_SignBatch(t *testing.T) {
	// sample keysign infos
	infoList := createKeysignInfoList(t, 10)

	tests := []struct {
		name              string
		batchNumber       uint64
		keysignInfoNonces []uint64
		expectedError     bool
	}{
		{
			name:              "successful batch signing",
			batchNumber:       0,
			keysignInfoNonces: []uint64{0, 1, 2, 3, 4},
		},
		{
			name:              "TSS signing error",
			batchNumber:       0,
			keysignInfoNonces: []uint64{0, 1, 2, 3, 4},
			expectedError:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// ARRANGE
			signer := newSignerTestSuite(t)
			batch := NewTSSKeysignBatch()

			// Mock batch and keysign info in cache
			for _, nonce := range tc.keysignInfoNonces {
				info := infoList[nonce]
				batch.AddKeysignInfo(nonce, *info)

				_, found := signer.GetSignatureOrAddDigest(nonce, info.digest)
				require.False(t, found)
			}

			// Mock TSS SignBatch error if needed
			if tc.expectedError {
				signer.tss.Pause()
				defer signer.tss.Unpause()
			}

			// ACT
			ctx := goctx.Background()
			err := signer.SignBatch(ctx, *batch, 100)

			// ASSERT
			if tc.expectedError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify signatures were added to cache
			for _, nonce := range tc.keysignInfoNonces {
				info := infoList[nonce]
				sig, found := signer.GetSignatureOrAddDigest(nonce, info.digest)
				require.True(t, found)
				require.NotEqual(t, [65]byte{}, sig)
			}
		})
	}
}

func Test_GetSignatureOrAddDigest(t *testing.T) {
	// ARRANGE
	signer := newSignerTestSuite(t)
	batchNonces := []uint64{0, 1, 2, 3}

	// sample keysign infoList
	batch := NewTSSKeysignBatch()
	infos := createKeysignInfoList(t, 10)

	// sample keysign batch
	for _, nonce := range batchNonces {
		info := infos[nonce]
		batch.AddKeysignInfo(nonce, *info)
	}

	// ACT-1: add new digests and signatures
	for _, nonce := range batchNonces {
		sig, found := signer.GetSignatureOrAddDigest(nonce, infos[nonce].digest)
		require.False(t, found)
		require.Equal(t, [65]byte{}, sig)
	}

	sigs := genSignatures(t, len(batchNonces))
	signer.AddBatchSignatures(*batch, sigs)

	// ASSERT-1: verify signatures were added to cache
	for i, nonce := range batchNonces {
		sig, found := signer.GetSignatureOrAddDigest(nonce, infos[nonce].digest)
		require.True(t, found)
		require.Equal(t, sigs[i], sig)
	}

	// ACT-2: update existing digests, clearing signatures
	newBatch := NewTSSKeysignBatch()
	newInfos := createKeysignInfoList(t, 10)
	for _, nonce := range batchNonces {
		info := newInfos[nonce]
		newBatch.AddKeysignInfo(nonce, *info)
	}

	for _, nonce := range batchNonces {
		sig, found := signer.GetSignatureOrAddDigest(nonce, newInfos[nonce].digest)
		require.False(t, found)
		require.Equal(t, [65]byte{}, sig)
	}

	// ACT-3: add new signatures
	newSigs := genSignatures(t, len(batchNonces))
	signer.AddBatchSignatures(*newBatch, newSigs)

	// ASSERT-3: verify new signatures were added to cache
	for _, nonce := range batchNonces {
		sig, found := signer.GetSignatureOrAddDigest(nonce, newInfos[nonce].digest)
		require.True(t, found)
		require.Equal(t, newSigs[nonce], sig)
	}
}

func Test_RemoveKeysignInfo(t *testing.T) {
	// ARRANGE
	signer := newSignerTestSuite(t)

	// mock pending keysign info in cache
	for nonce := range uint64(10) {
		signer.GetSignatureOrAddDigest(nonce, sample.Digest32B(t))
	}

	// ACT
	signer.RemoveKeysignInfo(5)

	// ASSERT
	// verify cleanup happened
	for nonce := range uint64(10) {
		_, found := signer.tssKeysignInfoMap[nonce]
		if nonce < 5 {
			require.False(t, found)
		} else {
			require.True(t, found)
		}
	}
}

// genSignatures generates a list of sample signatures.
func genSignatures(t *testing.T, count int) [][65]byte {
	sigs := make([][65]byte, count)
	for i := range count {
		sigs[i] = sample.Signature65B(t)
	}
	return sigs
}
