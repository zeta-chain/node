package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/lightclient/keeper"
	"github.com/zeta-chain/node/x/lightclient/types"
)

func TestMsgServer_DisableVerificationFlags(t *testing.T) {
	t.Run("emergency group can disable verification flags", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Ethereum.ChainId,
					Enabled: true,
				},
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: true,
				},
			},
		})

		// enable eth type chain
		msg := types.MsgDisableHeaderVerification{
			Creator:     admin,
			ChainIdList: []int64{chains.Ethereum.ChainId, chains.BitcoinMainnet.ChainId},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.DisableHeaderVerification(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		bhv, found := k.GetBlockHeaderVerification(ctx)
		require.True(t, found)
		require.False(t, bhv.IsChainEnabled(chains.Ethereum.ChainId))
		require.False(t, bhv.IsChainEnabled(chains.BitcoinMainnet.ChainId))

	})

	t.Run("cannot update if not authorized group", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Ethereum.ChainId,
					Enabled: true,
				},
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: true,
				},
			},
		})

		msg := types.MsgDisableHeaderVerification{
			Creator:     admin,
			ChainIdList: []int64{chains.Ethereum.ChainId},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := srv.DisableHeaderVerification(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("disable chain if even if the the chain has nto been set before", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		// enable eth type chain
		msg := types.MsgDisableHeaderVerification{
			Creator:     admin,
			ChainIdList: []int64{chains.Ethereum.ChainId, chains.BitcoinMainnet.ChainId},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.DisableHeaderVerification(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)
		bhv, found := k.GetBlockHeaderVerification(ctx)
		require.True(t, found)
		require.False(t, bhv.IsChainEnabled(chains.Ethereum.ChainId))
		require.False(t, bhv.IsChainEnabled(chains.BitcoinMainnet.ChainId))
	})
}
