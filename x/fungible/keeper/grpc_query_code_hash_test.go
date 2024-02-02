package keeper_test

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_CodeHash(t *testing.T) {
	t.Run("should return code hash", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		wzeta, _, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		acc := sdkk.EvmKeeper.GetAccount(ctx, wzeta)
		assert.NotNil(t, acc)
		assert.NotNil(t, acc.CodeHash)

		res, err := k.CodeHash(ctx, &types.QueryCodeHashRequest{
			Address: wzeta.Hex(),
		})
		assert.NoError(t, err)
		assert.Equal(t, ethcommon.BytesToHash(acc.CodeHash).Hex(), res.CodeHash)
	})

	t.Run("should return error if address is invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)

		_, err := k.CodeHash(ctx, &types.QueryCodeHashRequest{
			Address: "invalid",
		})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid address")
	})

	t.Run("should return error if account not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)

		_, err := k.CodeHash(ctx, &types.QueryCodeHashRequest{
			Address: sample.EthAddress().Hex(),
		})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "account not found")
	})

	t.Run("should return error if account is not a contract", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		_, err := k.CodeHash(ctx, &types.QueryCodeHashRequest{
			Address: types.ModuleAddressEVM.Hex(),
		})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "account is not a contract")
	})
}
