package base

import (
	"context"
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
)

func Test_GetScanRangeInboundSafe(t *testing.T) {
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
			expectedBlockRange: [2]uint64{91, 91}, // [91, 91), nothing to scan
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

			start, end := ob.GetScanRangeInboundSafe(tt.blockLimit)
			require.Equal(t, tt.expectedBlockRange, [2]uint64{start, end})
		})
	}
}

func Test_GetScanRangeInboundFast(t *testing.T) {
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
			expectedBlockRange: [2]uint64{91, 91}, // [91, 91), nothing to scan
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

			start, end := ob.GetScanRangeInboundFast(tt.blockLimit)
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

func Test_IsInboundEligibleForFastConfirmation(t *testing.T) {
	chain := chains.Ethereum
	liquidityCap := sdkmath.NewUint(100_000)
	fastAmountCap := chains.CalcInboundFastConfirmationAmountCap(chain.ChainId, liquidityCap)
	confParamsEnabled := observertypes.ConfirmationParams{
		SafeInboundCount: 2,
		FastInboundCount: 1,
	}

	tests := []struct {
		name                string
		confParams          observertypes.ConfirmationParams
		msg                 *crosschaintypes.MsgVoteInbound
		failForeignCoinsRPC bool
		eligible            bool
		errMsg              string
	}{
		{
			name:       "eligible for fast confirmation",
			confParams: confParamsEnabled,
			msg: &crosschaintypes.MsgVoteInbound{
				SenderChainId:           chain.ChainId,
				Amount:                  sdkmath.NewUint(fastAmountCap.Uint64()),
				CoinType:                coin.CoinType_Gas,
				Asset:                   "",
				ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V2,
			},
			eligible: true,
		},
		{
			name: "not eligible if fast confirmation is disabled",
			confParams: observertypes.ConfirmationParams{
				SafeInboundCount: 2,
				FastInboundCount: 2, // equal to safe confirmation, effectively disabled
			},
			msg: &crosschaintypes.MsgVoteInbound{
				SenderChainId: chains.SolanaMainnet.ChainId, // not set for Solana
			},
			eligible: false,
		},
		{
			name:       "not eligible if protocol contract version V1 is used",
			confParams: confParamsEnabled,
			msg: &crosschaintypes.MsgVoteInbound{
				SenderChainId:           chain.ChainId,
				ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V1, // not eligible for V1
			},
			eligible: false,
		},
		{
			name:       "return error if foreign coins query RPC fails",
			confParams: confParamsEnabled,
			msg: &crosschaintypes.MsgVoteInbound{
				SenderChainId:           chain.ChainId,
				CoinType:                coin.CoinType_Gas,
				Asset:                   "",
				ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V2,
			},
			failForeignCoinsRPC: true,
			eligible:            false,
			errMsg:              zrepo.ErrClientGetForeignCoinsForAsset.Error(),
		},
		{
			name:       "not eligible if amount exceeds fast amount cap",
			confParams: confParamsEnabled,
			msg: &crosschaintypes.MsgVoteInbound{
				SenderChainId:           chain.ChainId,
				Amount:                  sdkmath.NewUint(fastAmountCap.Uint64() + 1), // +1 to exceed
				CoinType:                coin.CoinType_Gas,
				Asset:                   "",
				ProtocolContractVersion: crosschaintypes.ProtocolContractVersion_V2,
			},
			eligible: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// ARRANGE
			ob := newTestSuite(t, chain, withConfirmationParams(tt.confParams))

			// mock up the foreign coins RPC
			assetAddress := ethcommon.HexToAddress(tt.msg.Asset)
			if tt.failForeignCoinsRPC {
				ob.zetacore.
					On("GetForeignCoinsFromAsset", mock.Anything, chain.ChainId, assetAddress).
					Maybe().
					Return(fungibletypes.ForeignCoins{}, errors.New("rpc failed"))
			} else {
				ob.zetacore.
					On("GetForeignCoinsFromAsset", mock.Anything, chain.ChainId, assetAddress).
					Maybe().
					Return(fungibletypes.ForeignCoins{LiquidityCap: liquidityCap}, nil)
			}

			// ACT
			ctx := context.Background()
			eligible, err := ob.IsInboundEligibleForFastConfirmation(ctx, tt.msg)

			// ASSERT
			require.Equal(t, tt.eligible, eligible)
			if tt.errMsg != "" {
				require.Contains(t, err.Error(), tt.errMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}
