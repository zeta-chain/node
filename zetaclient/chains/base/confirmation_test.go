package base_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func Test_CalcUnscannedBlockRangeInboundSafe(t *testing.T) {
	chain := chains.BitcoinMainnet

	tests := []struct {
		name               string
		lastBlock          uint64
		lastScanned        uint64
		blockLimit         uint64
		confParams         observertypes.ConfirmationParams
		expectedBlockRange [2]uint64
	}{
		{
			name:        "no unscanned blocks",
			lastBlock:   99,
			lastScanned: 90,
			blockLimit:  10,
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 10,
			},
			expectedBlockRange: [2]uint64{0, 0}, // [0, 0)
		},
		{
			name:        "1 unscanned blocks",
			lastBlock:   100,
			lastScanned: 90,
			blockLimit:  10,
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 10,
			},
			expectedBlockRange: [2]uint64{91, 92}, // [91, 92)
		},
		{
			name:        "10 unscanned blocks",
			lastBlock:   109,
			lastScanned: 90,
			blockLimit:  10,
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 10,
			},
			expectedBlockRange: [2]uint64{91, 101}, // [91, 101)
		},
		{
			name:        "block limit applied",
			lastBlock:   110,
			lastScanned: 90,
			blockLimit:  10,
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 10,
			},
			expectedBlockRange: [2]uint64{91, 101}, // [91, 101), 11 unscanned blocks, but capped to 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := newTestSuite(t, chain, withConfirmationParams(tt.confParams))
			ob.Observer.WithLastBlock(tt.lastBlock)
			ob.Observer.WithLastBlockScanned(tt.lastScanned)

			start, end := ob.CalcUnscannedBlockRangeInboundSafe(tt.blockLimit)
			require.Equal(t, tt.expectedBlockRange, [2]uint64{start, end})
		})
	}
}

func Test_CalcUnscannedBlockRangeInboundFast(t *testing.T) {
	chain := chains.BitcoinMainnet

	tests := []struct {
		name               string
		lastBlock          uint64
		lastScanned        uint64
		blockLimit         uint64
		confParams         observertypes.ConfirmationParams
		expectedBlockRange [2]uint64
	}{
		{
			name:        "no unscanned blocks",
			lastBlock:   99,
			lastScanned: 90,
			blockLimit:  10,
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 10,
				FastInboundCount: 10,
			},
			expectedBlockRange: [2]uint64{0, 0}, // [0, 0)
		},
		{
			name:        "1 unscanned blocks",
			lastBlock:   100,
			lastScanned: 90,
			blockLimit:  10,
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 10,
				FastInboundCount: 0, // should fall back to safe confirmation
			},
			expectedBlockRange: [2]uint64{91, 92}, // [91, 92)
		},
		{
			name:        "10 unscanned blocks",
			lastBlock:   109,
			lastScanned: 90,
			blockLimit:  10,
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 10,
				FastInboundCount: 10,
			},
			expectedBlockRange: [2]uint64{91, 101}, // [91, 101)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := newTestSuite(t, chain, withConfirmationParams(tt.confParams))
			ob.Observer.WithLastBlock(tt.lastBlock)
			ob.Observer.WithLastBlockScanned(tt.lastScanned)

			start, end := ob.CalcUnscannedBlockRangeInboundFast(tt.blockLimit)
			require.Equal(t, tt.expectedBlockRange, [2]uint64{start, end})
		})
	}
}

func Test_IsBlockConfirmedForInboundSafe(t *testing.T) {
	chain := chains.BitcoinMainnet

	tests := []struct {
		name        string
		blockNumber uint64
		lastBlock   uint64
		confParams  observertypes.ConfirmationParams
		expected    bool
	}{
		{
			name:        "should confirm block 100 when confirmation arrives 2",
			blockNumber: 100,
			lastBlock:   101, // got 2 confirmations
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 2,
			},
			expected: true,
		},
		{
			name:        "should not confirm block 100 when confirmation < 2",
			blockNumber: 100,
			lastBlock:   100, // got 1 confirmation, need one more
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 2,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := newTestSuite(t, chain, withConfirmationParams(tt.confParams))
			ob.Observer.WithLastBlock(tt.lastBlock)

			isConfirmed := ob.IsBlockConfirmedForInboundSafe(tt.blockNumber)
			require.Equal(t, tt.expected, isConfirmed)
		})
	}
}

func Test_IsBlockConfirmedForInboundFast(t *testing.T) {
	chain := chains.BitcoinMainnet

	tests := []struct {
		name        string
		blockNumber uint64
		lastBlock   uint64
		confParams  observertypes.ConfirmationParams
		expected    bool
	}{
		{
			name:        "should confirm block 100 when confirmation arrives 2",
			blockNumber: 100,
			lastBlock:   101, // got 2 confirmations
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 2,
				FastInboundCount: 0, // should fall back to safe confirmation
			},
			expected: true,
		},
		{
			name:        "should not confirm block 100 when confirmation < 2",
			blockNumber: 100,
			lastBlock:   100, // got 1 confirmation, need one more
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 2,
				FastInboundCount: 2,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := newTestSuite(t, chain, withConfirmationParams(tt.confParams))
			ob.Observer.WithLastBlock(tt.lastBlock)

			isConfirmed := ob.IsBlockConfirmedForInboundFast(tt.blockNumber)
			require.Equal(t, tt.expected, isConfirmed)
		})
	}
}

func Test_IsBlockConfirmedForOutboundSafe(t *testing.T) {
	chain := chains.BitcoinMainnet

	tests := []struct {
		name        string
		blockNumber uint64
		lastBlock   uint64
		confParams  observertypes.ConfirmationParams
		expected    bool
	}{
		{
			name:        "should confirm block 100 when confirmation arrives 2",
			blockNumber: 100,
			lastBlock:   101, // got 2 confirmations
			confParams: observertypes.ConfirmationParams{
				SafeOutboundCount: 2,
			},
			expected: true,
		},
		{
			name:        "should not confirm block 100 when confirmation < 2",
			blockNumber: 100,
			lastBlock:   100, // got 1 confirmation, need one more
			confParams: observertypes.ConfirmationParams{
				SafeOutboundCount: 2,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := newTestSuite(t, chain, withConfirmationParams(tt.confParams))
			ob.Observer.WithLastBlock(tt.lastBlock)

			isConfirmed := ob.IsBlockConfirmedForOutboundSafe(tt.blockNumber)
			require.Equal(t, tt.expected, isConfirmed)
		})
	}
}

func Test_IsBlockConfirmedForOutboundFast(t *testing.T) {
	chain := chains.BitcoinMainnet

	tests := []struct {
		name        string
		blockNumber uint64
		lastBlock   uint64
		confParams  observertypes.ConfirmationParams
		expected    bool
	}{
		{
			name:        "should confirm block 100 when confirmation arrives 2",
			blockNumber: 100,
			lastBlock:   101, // got 2 confirmations
			confParams: observertypes.ConfirmationParams{
				SafeOutboundCount: 2,
				FastOutboundCount: 0, // should fall back to safe confirmation
			},
			expected: true,
		},
		{
			name:        "should not confirm block 100 when confirmation < 2",
			blockNumber: 100,
			lastBlock:   100, // got 1 confirmation, need one more
			confParams: observertypes.ConfirmationParams{
				SafeOutboundCount: 2,
				FastOutboundCount: 2,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := newTestSuite(t, chain, withConfirmationParams(tt.confParams))
			ob.Observer.WithLastBlock(tt.lastBlock)

			isConfirmed := ob.IsBlockConfirmedForOutboundFast(tt.blockNumber)
			require.Equal(t, tt.expected, isConfirmed)
		})
	}
}
