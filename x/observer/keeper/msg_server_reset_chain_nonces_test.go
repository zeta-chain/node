package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/observer/keeper"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestMsgServer_ResetChainNonces(t *testing.T) {
	t.Run("cannot reset chain nonces if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		chainId := chains.GoerliLocalnet.ChainId
		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		msg := types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        chainId,
			ChainNonceLow:  1,
			ChainNonceHigh: 5,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("cannot reset chain nonces if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		chainId := chains.GoerliLocalnet.ChainId

		msg := types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        chainId,
			ChainNonceLow:  1,
			ChainNonceHigh: 5,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &msg)
		require.ErrorIs(t, err, types.ErrTssNotFound)
	})

	t.Run("can reset chain nonces", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()
		k.SetTSS(ctx, tss)

		admin := sample.AccAddress()
		chainId := chains.GoerliLocalnet.ChainId
		nonceLow := 1
		nonceHigh := 5
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// check existing chain nonces
		_, found := k.GetChainNonces(ctx, chainId)
		require.False(t, found)
		_, found = k.GetPendingNonces(ctx, tss.TssPubkey, chainId)
		require.False(t, found)

		// reset chain nonces
		// Reset nonces to nonceLow and nonceHigh
		msg := types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        chainId,
			ChainNonceLow:  int64(nonceLow),
			ChainNonceHigh: int64(nonceHigh),
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		//keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check updated chain nonces
		chainNonces, found := k.GetChainNonces(ctx, chainId)
		require.True(t, found)
		require.Equal(t, chainId, chainNonces.ChainId)
		require.Equal(t, uint64(nonceHigh), chainNonces.Nonce)

		pendingNonces, found := k.GetPendingNonces(ctx, tss.TssPubkey, chainId)
		require.True(t, found)
		require.Equal(t, chainId, pendingNonces.ChainId)
		require.Equal(t, tss.TssPubkey, pendingNonces.Tss)
		require.Equal(t, int64(nonceLow), pendingNonces.NonceLow)
		require.Equal(t, int64(nonceHigh), pendingNonces.NonceHigh)

		// reset nonces back to 0
		// Reset nonces back to 0
		msg = types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        chainId,
			ChainNonceLow:  0,
			ChainNonceHigh: 0,
		}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err = srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &msg)
		require.NoError(t, err)

		// check updated chain nonces
		chainNonces, found = k.GetChainNonces(ctx, chainId)
		require.True(t, found)
		require.Equal(t, chainId, chainNonces.ChainId)
		require.Equal(t, uint64(0), chainNonces.Nonce)

		pendingNonces, found = k.GetPendingNonces(ctx, tss.TssPubkey, chainId)
		require.True(t, found)
		require.Equal(t, chainId, pendingNonces.ChainId)
		require.Equal(t, tss.TssPubkey, pendingNonces.Tss)
		require.Equal(t, int64(0), pendingNonces.NonceLow)
		require.Equal(t, int64(0), pendingNonces.NonceHigh)
	})
}
