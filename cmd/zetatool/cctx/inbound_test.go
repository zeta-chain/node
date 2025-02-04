package cctx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/cmd/zetatool/cctx"
	zetatoolcontext "github.com/zeta-chain/node/cmd/zetatool/context"
	"github.com/zeta-chain/node/pkg/chains"
)

func Test_InboundBallotIdentifier(t *testing.T) {
	tt := []struct {
		name                     string
		inboundHash              string
		inboundChainID           int64
		expectedBallotIdentifier string
		expectError              bool
	}{
		{
			name:                     chains.Ethereum.Name,
			inboundHash:              "0x61008d7f79b2955a15e3cb95154a80e19c7385993fd0e083ff0cbe0b0f56cb9a",
			inboundChainID:           chains.Ethereum.ChainId,
			expectedBallotIdentifier: "0xae189ab5cd884af784835297ac43eb55deb8a7800023534c580f44ee2b3eb5ed",
			expectError:              false,
		},
		{
			name:                     chains.BaseMainnet.Name,
			inboundHash:              "0x88ee0943863fd8649546eb3affaf725f8caf09f44ebc5aa64de592b2edf378c8",
			inboundChainID:           chains.BaseMainnet.ChainId,
			expectedBallotIdentifier: "0xe2b4c3f5dbef8fb7feb14bdf0a3f63ca7018678ecb6ae99ff697ccd962932ca2",
			expectError:              false,
		},
		{
			name:                     chains.BscMainnet.Name,
			inboundHash:              "0xfa18cbcdbf70e987600647ee77a1a28f5ca707acf9b72462fada02fff2a94d2f",
			inboundChainID:           chains.BscMainnet.ChainId,
			expectedBallotIdentifier: "0xc7b289172db825b3c0490f263f35c8596b6f1fab8ec4c44db46de3020fe9e6e6",
			expectError:              false,
		},
		{
			name:                     chains.Polygon.Name,
			inboundHash:              "0x70b9b3ba89ff647257ab0085d90d60dc99b693c66931c4535e117b66a25236ce",
			inboundChainID:           chains.Polygon.ChainId,
			expectedBallotIdentifier: "0xf8ed419d9798aed83070763355628e2638ae9a4a47aa9c93ffc32f4b72c9fef4",
			expectError:              false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			ctx, err := zetatoolcontext.NewContext(context.Background(), tc.inboundChainID, tc.inboundHash, "")
			require.NoError(t, err)
			c := cctx.NewTrackingDetails()
			err = c.CheckInbound(ctx)
			require.NoError(t, err)
			if !tc.expectError && c.CCTXIdentifier != tc.expectedBallotIdentifier {
				t.Errorf("expected %s, got %s", tc.expectedBallotIdentifier, c.CCTXIdentifier)
			}
		})
	}

}
