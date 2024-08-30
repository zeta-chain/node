package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/x/emissions/types"
)

func TestKeeper_ListPoolAddresses(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		res, err := k.ListPoolAddresses(wctx, nil)
		require.Nil(t, res)
		require.Error(t, err)
	})

	t.Run("should not error if req is not nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.EmissionsKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		expectedRes := &types.QueryListPoolAddressesResponse{
			UndistributedObserverBalancesAddress: types.UndistributedObserverRewardsPoolAddress.String(),
			EmissionModuleAddress:                types.EmissionsModuleAddress.String(),
			UndistributedTssBalancesAddress:      types.UndistributedTssRewardsPoolAddress.String(),
		}
		res, err := k.ListPoolAddresses(wctx, &types.QueryListPoolAddressesRequest{})
		require.NoError(t, err)
		require.Equal(t, expectedRes, res)
	})
}
