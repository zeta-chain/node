package keeper_test

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	evmkeeper "github.com/evmos/ethermint/x/evm/keeper"
	"github.com/stretchr/testify/require"

	testkeeper "github.com/zeta-chain/zetacore/testutil/keeper"
	fungiblekeeper "github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

// setupGasCoin is a helper function to setup the gas coin for testing
func setupGasCoin(
	t *testing.T,
	ctx sdk.Context,
	k *fungiblekeeper.Keeper,
	evmk *evmkeeper.Keeper,
	chainID int64,
	assetName string,
	symbol string,
) (zrc20 common.Address) {
	addr, err := k.SetupChainGasCoinAndPool(
		ctx,
		chainID,
		assetName,
		symbol,
		8,
	)
	require.NoError(t, err)
	assertContractDeployment(t, evmk, ctx, addr)
	return addr
}

func TestKeeper_SetupChainGasCoinAndPool(t *testing.T) {
	t.Run("can setup a new chain gas coin", func(t *testing.T) {
		k, ctx, sdkk := testkeeper.FungibleKeeper(t)
		_ = k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// deploy the system contracts
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, 1, "foobar", "foobar")

		// can retrieve the gas coin
		found, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(1))
		require.NoError(t, err)
		require.Equal(t, zrc20, found)
	})
}
