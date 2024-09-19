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

func TestMsgServer_EnableVerificationFlags(t *testing.T) {
	t.Run("operational group can enable verification flags", func(t *testing.T) {
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
					Enabled: false,
				},
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: false,
				},
			},
		})

		// enable both eth and btc type chain together
		msg := types.MsgEnableHeaderVerification{
			Creator:     admin,
			ChainIdList: []int64{chains.Ethereum.ChainId, chains.BitcoinMainnet.ChainId},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.EnableHeaderVerification(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)
		bhv, found := k.GetBlockHeaderVerification(ctx)
		require.True(t, found)
		require.True(t, bhv.IsChainEnabled(chains.Ethereum.ChainId))
		require.True(t, bhv.IsChainEnabled(chains.BitcoinMainnet.ChainId))
	})

	t.Run("enable verification flags even if the chain has not been set previously", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeperWithMocks(t, keepertest.LightclientMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		// mock the authority keeper for authorization
		authorityMock := keepertest.GetLightclientAuthorityMock(t, k)

		// enable both eth and btc type chain together
		msg := types.MsgEnableHeaderVerification{
			Creator:     admin,
			ChainIdList: []int64{chains.Ethereum.ChainId, chains.BitcoinMainnet.ChainId},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.EnableHeaderVerification(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)
		bhv, found := k.GetBlockHeaderVerification(ctx)
		require.True(t, found)
		require.True(t, bhv.IsChainEnabled(chains.Ethereum.ChainId))
		require.True(t, bhv.IsChainEnabled(chains.BitcoinMainnet.ChainId))
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
					Enabled: false,
				},
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: false,
				},
			},
		})

		msg := types.MsgEnableHeaderVerification{
			Creator:     admin,
			ChainIdList: []int64{chains.Ethereum.ChainId},
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := srv.EnableHeaderVerification(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})
}
