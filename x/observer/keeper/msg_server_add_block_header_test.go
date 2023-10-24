//go:build TESTNET
// +build TESTNET

package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_AddBlockHeader(t *testing.T) {
	header, err := sample.EthHeader()
	assert.NoError(t, err)
	observerChain := common.GoerliChain()
	observerAddress := sample.AccAddress()
	// Add tests for btc headers : https://github.com/zeta-chain/node/issues/1336
	tt := []struct {
		name                  string
		msg                   *types.MsgAddBlockHeader
		IsEthTypeChainEnabled bool
		IsBtcTypeChainEnabled bool
		wantErr               require.ErrorAssertionFunc
	}{
		{
			name: "success submit eth header",
			msg: &types.MsgAddBlockHeader{
				Creator:   observerAddress,
				ChainId:   common.GoerliChain().ChainId,
				BlockHash: sample.Bytes(),
				Height:    1,
				Header:    common.NewEthereumHeader(header),
			},
			IsEthTypeChainEnabled: true,
			IsBtcTypeChainEnabled: true,
			wantErr:               require.NoError,
		},
		{
			name: "failure submit eth header eth disabled",
			msg: &types.MsgAddBlockHeader{
				Creator:   observerAddress,
				ChainId:   common.GoerliChain().ChainId,
				BlockHash: sample.Bytes(),
				Height:    1,
				Header:    common.NewEthereumHeader(header),
			},
			IsEthTypeChainEnabled: false,
			IsBtcTypeChainEnabled: true,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				assert.ErrorIs(t, err, types.ErrBlockHeaderVerficationDisabled)
			},
		},
		{
			name: "failure submit eth header eth disabled",
			msg: &types.MsgAddBlockHeader{
				Creator:   sample.AccAddress(),
				ChainId:   common.GoerliChain().ChainId,
				BlockHash: sample.Bytes(),
				Height:    1,
				Header:    common.NewEthereumHeader(header),
			},
			IsEthTypeChainEnabled: false,
			IsBtcTypeChainEnabled: true,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				assert.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			k, ctx := keepertest.ObserverKeeper(t)
			srv := keeper.NewMsgServerImpl(*k)
			k.SetObserverMapper(ctx, &types.ObserverMapper{
				ObserverChain: &observerChain,
				ObserverList:  []string{observerAddress},
			})
			k.SetCrosschainFlags(ctx, types.CrosschainFlags{
				IsInboundEnabled:      true,
				IsOutboundEnabled:     true,
				GasPriceIncreaseFlags: nil,
				BlockHeaderVerificationFlags: &types.BlockHeaderVerificationFlags{
					IsEthTypeChainEnabled: tc.IsEthTypeChainEnabled,
					IsBtcTypeChainEnabled: tc.IsBtcTypeChainEnabled,
				},
			})
			_, err := srv.AddBlockHeader(ctx, tc.msg)
			tc.wantErr(t, err)
		})
	}
}
