package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/observer/keeper"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_ResetChainNonces(t *testing.T) {
	t.Run("cannot reset chain nonces if not authorized", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		chainId := common.GoerliLocalnetChain().ChainId
		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &types.MsgResetChainNonces{
			Creator:        sample.AccAddress(),
			ChainId:        chainId,
			ChainNonceLow:  1,
			ChainNonceHigh: 5,
		})
		require.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)

		// group 1 should not be able to reset chain nonces
		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group1)

		_, err = srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &types.MsgResetChainNonces{
			Creator:        sample.AccAddress(),
			ChainId:        chainId,
			ChainNonceLow:  1,
			ChainNonceHigh: 5,
		})
		require.ErrorIs(t, err, types.ErrNotAuthorizedPolicy)
	})

	t.Run("cannot reset chain nonces if tss not found", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		chainId := common.GoerliLocalnetChain().ChainId
		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        chainId,
			ChainNonceLow:  1,
			ChainNonceHigh: 5,
		})
		require.ErrorIs(t, err, types.ErrTssNotFound)
	})

	t.Run("cannot reset chain nonces if chain not supported", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()
		k.SetTSS(ctx, tss)

		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		_, err := srv.ResetChainNonces(sdk.WrapSDKContext(ctx), &types.MsgResetChainNonces{
			Creator:        admin,
			ChainId:        999,
			ChainNonceLow:  1,
			ChainNonceHigh: 5,
		})
		require.ErrorIs(t, err, types.ErrSupportedChains)
	})

	t.Run("can reset chain nonces", func(t *testing.T) {
		k, ctx := keepertest.ObserverKeeper(t)
		srv := keeper.NewMsgServerImpl(*k)
		tss := sample.Tss()
		k.SetTSS(ctx, tss)

		admin := sample.AccAddress()
		setAdminCrossChainFlags(ctx, k, admin, types.Policy_Type_group2)

		chainId := common.GoerliLocalnetChain().ChainId
		index := common.GoerliLocalnetChain().ChainName.String()

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
			ChainNonceLow:  uint64(nonceLow),
			ChainNonceHigh: uint64(nonceHigh),
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
			ChainNonceLow:  uint64(0),
			ChainNonceHigh: uint64(0),
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
