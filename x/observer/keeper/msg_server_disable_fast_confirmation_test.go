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

// setupChainParamsList sets up the given chain params list in the keeper
func setupChainParamsList(
	t *testing.T,
	k *keeper.Keeper,
	ctx sdk.Context,
	admin string,
	chainParamsList []*types.ChainParams,
) {
	srv := keeper.NewMsgServerImpl(*k)

	// initial chain params list should be empty
	_, found := k.GetChainParamsList(ctx)
	require.False(t, found)

	for _, cp := range chainParamsList {
		// mock the authority keeper for authorization
		msg := types.MsgUpdateChainParams{Creator: admin, ChainParams: cp}
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)

		// add chain params to state
		_, err := srv.UpdateChainParams(ctx, &msg)
		require.NoError(t, err)
	}

	// chain params should have FAST confirmation enabled
	allChainParams, found := k.GetChainParamsList(ctx)
	require.True(t, found)
	require.Len(t, allChainParams.ChainParams, len(chainParamsList))

	for i, cp := range allChainParams.ChainParams {
		require.Equal(t, chainParamsList[i], cp)
		ensureFastConfirmationEnabled(t, cp)
	}
}

func TestMsgServer_DisableFastConfirmation(t *testing.T) {
	t.Run("emergency group can disable fast confirmation", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})

		// create mock admin, msg server and authority keeper
		admin := sample.AccAddress()
		srv := keeper.NewMsgServerImpl(*k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// create two chain params
		chainList := chains.ExternalChainList([]chains.Chain{})
		chainParams0 := sample.ChainParams(chainList[0].ChainId)
		chainParams1 := sample.ChainParams(chainList[1].ChainId)

		// setup chain params list
		setupChainParamsList(t, k, ctx, admin, []*types.ChainParams{chainParams0, chainParams1})

		// ACT
		// FAST confirmation can be disabled for the second chain
		msg := types.MsgDisableFastConfirmation{Creator: admin, ChainId: chainParams1.ChainId}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.DisableFastConfirmation(ctx, &msg)
		require.NoError(t, err)

		// ASSERT
		// check list has two chain params
		chainParamsList, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 2)

		// 1st chain params should be unchanged
		require.Equal(t, chainParams0, chainParamsList.ChainParams[0])

		// 2nd chain params should have FAST confirmation disabled
		// also, only the FAST confirmation counts were modified
		ensureFastConfirmationDisabled(t, chainParamsList.ChainParams[1])
		chainParamsList.ChainParams[1].ConfirmationParams.FastInboundCount = chainParams1.ConfirmationParams.FastInboundCount
		chainParamsList.ChainParams[1].ConfirmationParams.FastOutboundCount = chainParams1.ConfirmationParams.FastOutboundCount
		require.Equal(t, chainParams1, chainParamsList.ChainParams[1])
	})

	t.Run("cannot disable fast confirmation if not authorized", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})

		// create mock admin, msg server and authority keeper
		admin := sample.AccAddress()
		srv := keeper.NewMsgServerImpl(*k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// create chain params
		chainList := chains.ExternalChainList([]chains.Chain{})
		chainParams := sample.ChainParams(chainList[0].ChainId)

		// setup chain params list
		setupChainParamsList(t, k, ctx, admin, []*types.ChainParams{chainParams})

		// ACT
		// try to disable FAST confirmation for the chain
		msg := types.MsgDisableFastConfirmation{Creator: admin, ChainId: chainParams.ChainId}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, authoritytypes.ErrUnauthorized)
		_, err := srv.DisableFastConfirmation(ctx, &msg)

		// ASSERT
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)

		// chain params list should be unchanged
		chainParamsList, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 1)
		require.Equal(t, chainParams, chainParamsList.ChainParams[0])
	})

	t.Run("cannot disable fast confirmation if chain params not found", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})

		// create mock admin, msg server and authority keeper
		admin := sample.AccAddress()
		srv := keeper.NewMsgServerImpl(*k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// ACT
		msg := types.MsgDisableFastConfirmation{Creator: admin, ChainId: 1}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.DisableFastConfirmation(ctx, &msg)

		// ASSERT
		require.ErrorIs(t, err, types.ErrChainParamsNotFound)
	})

	t.Run("cannot disable fast confirmation if validation fails", func(t *testing.T) {
		// ARRANGE
		k, ctx, _, _ := keepertest.ObserverKeeperWithMocks(t, keepertest.ObserverMockOptions{
			UseAuthorityMock: true,
		})

		// create mock admin, msg server and authority keeper
		admin := sample.AccAddress()
		srv := keeper.NewMsgServerImpl(*k)
		authorityMock := keepertest.GetObserverAuthorityMock(t, k)

		// create invalid chain params
		chainList := chains.ExternalChainList([]chains.Chain{})
		chainParams := sample.ChainParams(chainList[0].ChainId)
		chainParams.ConfirmationParams.SafeInboundCount = 0 // <- invalid count

		// setup chain params list
		setupChainParamsList(t, k, ctx, admin, []*types.ChainParams{chainParams})

		// ACT
		// try to disable FAST confirmation for the chain
		msg := types.MsgDisableFastConfirmation{Creator: admin, ChainId: chainParams.ChainId}
		keepertest.MockCheckAuthorization(&authorityMock.Mock, &msg, nil)
		_, err := srv.DisableFastConfirmation(ctx, &msg)

		// ASSERT
		require.ErrorIs(t, err, types.ErrInvalidChainParams)

		// chain params list should be unchanged
		chainParamsList, found := k.GetChainParamsList(ctx)
		require.True(t, found)
		require.Len(t, chainParamsList.ChainParams, 1)
		require.Equal(t, chainParams, chainParamsList.ChainParams[0])
	})
}

// ensureFastConfirmationEnabled checks that the fast confirmation is enabled
func ensureFastConfirmationEnabled(t *testing.T, chainParams *types.ChainParams) {
	require.NotNil(t, chainParams.ConfirmationParams)
	require.Positive(t, chainParams.ConfirmationParams.FastInboundCount)
	require.Positive(t, chainParams.ConfirmationParams.FastOutboundCount)
}

// ensureFastConfirmationDisabled checks that the fast confirmation is disabled
func ensureFastConfirmationDisabled(t *testing.T, chainParams *types.ChainParams) {
	require.NotNil(t, chainParams.ConfirmationParams)
	require.Equal(t, chainParams.ConfirmationParams.SafeInboundCount, chainParams.ConfirmationParams.FastInboundCount)
	require.Equal(t, chainParams.ConfirmationParams.SafeOutboundCount, chainParams.ConfirmationParams.FastOutboundCount)
}
