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
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
	"github.com/zeta-chain/node/zetaclient/config"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/mode"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
)

// signerTestSuite is a test suite for testing the signer
type signerTestSuite struct {
	*Signer
	tss *mocks.TSS
}

// newTestSuite creates a new test suite for testing
func newSignerTestSuite(t *testing.T) *signerTestSuite {
	// constructor parameters
	chain := chains.Ethereum
	tss := mocks.NewTSS(t)

	//logger := DefaultLogger()
	logger := Logger{}
	signer := NewSigner(chain, tss, logger, mode.StandardMode)

	suite := &signerTestSuite{
		Signer: signer,
		tss:    tss,
	}

	return suite
}

// createSigner creates a new signer for testing
func createSigner(t *testing.T) *Signer {
	// constructor parameters
	chain := chains.Ethereum
	tss := mocks.NewTSS(t)
	logger := DefaultLogger()

	// create signer
	return NewSigner(chain, tss, logger, mode.StandardMode)
}

func TestNewSigner(t *testing.T) {
	signer := createSigner(t)
	require.NotNil(t, signer)
}

func Test_BeingReportedFlag(t *testing.T) {
	signer := createSigner(t)

	// hash to be reported
	hash := "0x1234"
	alreadySet := signer.SetBeingReportedFlag(hash)
	require.False(t, alreadySet)

	// set reported outbound again and check
	alreadySet = signer.SetBeingReportedFlag(hash)
	require.True(t, alreadySet)

	// clear reported outbound and check again
	signer.ClearBeingReportedFlag(hash)
	alreadySet = signer.SetBeingReportedFlag(hash)
	require.False(t, alreadySet)
}

func Test_PassesCompliance(t *testing.T) {
	signer := createSigner(t)

	// create config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return false for restricted CCTX", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "abcd")
		cfg.ComplianceConfig.RestrictedAddresses = []string{cctx.InboundParams.Sender}
		config.SetRestrictedAddressesFromConfig(cfg)

		require.False(t, signer.PassesCompliance(cctx))
	})
	t.Run("should return true for non restricted CCTX", func(t *testing.T) {
		cctx := sample.CrossChainTxV2(t, "abcd")
		cfg.ComplianceConfig.RestrictedAddresses = []string{sample.EthAddress().Hex()}
		config.SetRestrictedAddressesFromConfig(cfg)

		require.True(t, signer.PassesCompliance(cctx))
	})
}

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
