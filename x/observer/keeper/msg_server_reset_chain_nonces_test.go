package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_ResetChainNonces(t *testing.T) {
	t.Run("cannot reset chain nonces if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		chainId := chains.GoerliLocalnetChain().ChainId

		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, false)

		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        chainId,
			ChainNonceLow:  1,
			ChainNonceHigh: 5,
		})
		require.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)
	})

	t.Run("cannot reset chain nonces if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		chainId := chains.GoerliLocalnetChain().ChainId
		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        chainId,
			ChainNonceLow:  1,
			ChainNonceHigh: 5,
		})
		require.ErrorIs(t, err, types.ErrTssNotFound)
	})

	t.Run("cannot reset chain nonces if chain not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()
		k.SetTSS(ctx, tss)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        999,
			ChainNonceLow:  1,
			ChainNonceHigh: 5,
		})
		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("can reset chain nonces", func(t *testing.T) {
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})
		srv := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()
		k.SetTSS(ctx, tss)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupAdmin, true)

		chainId := chains.GoerliLocalnetChain().ChainId
		index := chains.GoerliLocalnetChain().ChainName.String()

		// check existing chain nonces
		_, found := k.GetChainNonces(ctx, index)
		require.False(t, found)
		_, found = k.GetPendingNonces(ctx, tss.TssPubkey, chainId)
		require.False(t, found)

		// reset chain nonces
		nonceLow := 1
		nonceHigh := 5
		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        chainId,
			ChainNonceLow:  int64(nonceLow),
			ChainNonceHigh: int64(nonceHigh),
		})
		require.NoError(t, err)

		// check updated chain nonces
		chainNonces, found := k.GetChainNonces(ctx, index)
		require.True(t, found)
		require.Equal(t, chainId, chainNonces.ChainId)
		require.Equal(t, index, chainNonces.Index)
		require.Equal(t, uint64(nonceHigh), chainNonces.Nonce)

		pendingNonces, found := k.GetPendingNonces(ctx, tss.TssPubkey, chainId)
		require.True(t, found)
		require.Equal(t, chainId, pendingNonces.ChainId)
		require.Equal(t, tss.TssPubkey, pendingNonces.Tss)
		require.Equal(t, int64(nonceLow), pendingNonces.NonceLow)
		require.Equal(t, int64(nonceHigh), pendingNonces.NonceHigh)

		// reset nonces back to 0
		_, err = srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        chainId,
			ChainNonceLow:  0,
			ChainNonceHigh: 0,
		})
		require.NoError(t, err)

		// check updated chain nonces
		chainNonces, found = k.GetChainNonces(ctx, index)
		require.True(t, found)
		require.Equal(t, chainId, chainNonces.ChainId)
		require.Equal(t, index, chainNonces.Index)
		require.Equal(t, uint64(0), chainNonces.Nonce)

		pendingNonces, found = k.GetPendingNonces(ctx, tss.TssPubkey, chainId)
		require.True(t, found)
		require.Equal(t, chainId, pendingNonces.ChainId)
		require.Equal(t, tss.TssPubkey, pendingNonces.Tss)
		require.Equal(t, int64(0), pendingNonces.NonceLow)
		require.Equal(t, int64(0), pendingNonces.NonceHigh)
	})
}
