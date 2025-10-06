package base_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/context"
)

func Test_CheckSkipInbound(t *testing.T) {
	tests := []struct {
		name                string
		isInboundEnabled    bool
		isMempoolCongested  bool
		expectedSkip        bool
		expectedLogContains string
	}{
		{
			name:               "should not skip when all conditions are met",
			isInboundEnabled:   true,
			isMempoolCongested: false,
			expectedSkip:       false,
		},
		{
			name:               "should skip when inbound is disabled",
			isInboundEnabled:   false,
			isMempoolCongested: false,
			expectedSkip:       true,
		},
		{
			name:               "should skip when mempool is congested",
			isInboundEnabled:   true,
			isMempoolCongested: true,
			expectedSkip:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// create test suite
			chain := chains.Ethereum
			ethParams := observertypes.GetDefaultEthMainnetChainParams()
			ethParams.IsSupported = true
			ob := newTestSuite(t, chains.Ethereum)

			// mock app context
			appCtx := mockAppContext(t, chain, *ethParams, tt.isInboundEnabled, true, tt.isMempoolCongested)

			// ACT
			result := base.CheckSkipInbound(ob.Observer, appCtx)

			// ASSERT
			assert.Equal(t, tt.expectedSkip, result)
		})
	}
}

func Test_CheckSkipOutbound(t *testing.T) {
	tests := []struct {
		name               string
		isOutboundEnabled  bool
		isMempoolCongested bool
		expectedSkip       bool
	}{
		{
			name:               "should not skip when all conditions are met",
			isOutboundEnabled:  true,
			isMempoolCongested: false,
			expectedSkip:       false,
		},
		{
			name:               "should skip when outbound is disabled",
			isOutboundEnabled:  false,
			isMempoolCongested: false,
			expectedSkip:       true,
		},
		{
			name:               "should skip when mempool is congested",
			isOutboundEnabled:  true,
			isMempoolCongested: true,
			expectedSkip:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// create test suite
			chain := chains.Ethereum
			ethParams := observertypes.GetDefaultEthMainnetChainParams()
			ethParams.IsSupported = true
			ob := newTestSuite(t, chains.Ethereum)

			// mock app context
			appCtx := mockAppContext(t, chain, *ethParams, true, tt.isOutboundEnabled, tt.isMempoolCongested)

			// ACT
			result := base.CheckSkipOutbound(ob.Observer, appCtx)

			// ASSERT
			assert.Equal(t, tt.expectedSkip, result)
		})
	}
}

func Test_CheckSkipGasPrice(t *testing.T) {
	tests := []struct {
		name               string
		isMempoolCongested bool
		expectedSkip       bool
	}{
		{
			name:               "should not skip when chain is supported and mempool is not congested",
			isMempoolCongested: false,
			expectedSkip:       false,
		},
		{
			name:               "should skip when mempool is congested",
			isMempoolCongested: true,
			expectedSkip:       true,
		},
		{
			name:               "should skip when chain is not supported and mempool is congested",
			isMempoolCongested: true,
			expectedSkip:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			// create test suite
			chain := chains.Ethereum
			ethParams := observertypes.GetDefaultEthMainnetChainParams()
			ethParams.IsSupported = true
			ob := newTestSuite(t, chains.Ethereum)

			// mock app context
			appCtx := mockAppContext(t, chain, *ethParams, true, true, tt.isMempoolCongested)

			// ACT
			result := base.CheckSkipInbound(ob.Observer, appCtx)

			// ASSERT
			assert.Equal(t, tt.expectedSkip, result)
		})
	}
}

// mockAppContext creates a mock AppContext for testing
func mockAppContext(t *testing.T, chain chains.Chain, chainParams observertypes.ChainParams, isInboundEnabled, isOutboundEnabled, isMempoolCongested bool) *context.AppContext {
	// Create a mock config
	cfg := config.New(false)
	cfg.MempoolCongestionThreshold = 1

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t))

	// Create AppContext
	appCtx := context.New(cfg, nil, logger)

	// Create crosschain flags
	crosschainFlags := observertypes.CrosschainFlags{
		IsInboundEnabled:  isInboundEnabled,
		IsOutboundEnabled: isOutboundEnabled,
	}

	// Create operational flags
	operationalFlags := observertypes.OperationalFlags{}

	// Set unconfirmed tx count based on mempool congestion
	unconfirmedTxCount := 0
	if isMempoolCongested {
		unconfirmedTxCount = 2 // above threshold of 1
	}

	// Update the context
	err := appCtx.Update(
		[]chains.Chain{chain},
		nil,
		map[int64]*observertypes.ChainParams{chain.ChainId: &chainParams},
		crosschainFlags,
		operationalFlags,
		unconfirmedTxCount,
	)
	require.NoError(t, err)

	return appCtx
}
