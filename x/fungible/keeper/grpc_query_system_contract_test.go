package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_SystemContract(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		res, err := k.SystemContract(ctx, nil)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should error if system contract not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		res, err := k.SystemContract(ctx, &types.QueryGetSystemContractRequest{})
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("should return system contract if found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		sc := types.SystemContract{
			SystemContract:  sample.EthAddress().Hex(),
			ConnectorZevm:   sample.EthAddress().Hex(),
			GatewayGasLimit: sdkmath.NewIntFromBigInt(types.GatewayGasLimit),
		}
		k.SetSystemContract(ctx, sc)
		res, err := k.SystemContract(ctx, &types.QueryGetSystemContractRequest{})
		require.NoError(t, err)
		require.Equal(t, &types.QueryGetSystemContractResponse{
			SystemContract: sc,
		}, res)
	})
}
