//go:build TESTNET
// +build TESTNET

package keeper_test

import (
	"math/rand"
	"testing"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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
	r := rand.New(rand.NewSource(9))
	validator := sample.Validator(t, r)
	observerAddress, err := types.GetAccAddressFromOperatorAddress(validator.OperatorAddress)
	assert.NoError(t, err)
	// Add tests for btc headers : https://github.com/zeta-chain/node/issues/1336
	tt := []struct {
		name                  string
		msg                   *types.MsgAddBlockHeader
		IsEthTypeChainEnabled bool
		IsBtcTypeChainEnabled bool
		validator             stakingtypes.Validator
		wantErr               require.ErrorAssertionFunc
	}{
		{
			name: "success submit eth header",
			msg: &types.MsgAddBlockHeader{
				Creator:   observerAddress.String(),
				ChainId:   common.GoerliChain().ChainId,
				BlockHash: sample.Bytes(),
				Height:    1,
				Header:    common.NewEthereumHeader(header),
			},
			IsEthTypeChainEnabled: true,
			IsBtcTypeChainEnabled: true,
			validator:             validator,
			wantErr:               require.NoError,
		},
		{
			name: "failure submit eth header eth disabled",
			msg: &types.MsgAddBlockHeader{
				Creator:   observerAddress.String(),
				ChainId:   common.GoerliChain().ChainId,
				BlockHash: sample.Bytes(),
				Height:    1,
				Header:    common.NewEthereumHeader(header),
			},
			IsEthTypeChainEnabled: false,
			IsBtcTypeChainEnabled: true,
			validator:             validator,
			wantErr: func(t require.TestingT, err error, i ...interface{}) {
				assert.ErrorIs(t, err, types.ErrBlockHeaderVerificationDisabled)
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
			validator:             validator,
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
				ObserverList:  []string{observerAddress.String()},
			})
			k.GetStakingKeeper().SetValidator(ctx, tc.validator)
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
