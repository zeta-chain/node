package keeper_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	zetacommon "github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_UpdateContractBytecode(t *testing.T) {
	k, ctx, sdkk := keepertest.FungibleKeeper(t)
	k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

	// Sample chainIDs and addresses
	chainList := zetacommon.DefaultChainsList()
	require.True(t, len(chainList) > 1)
	require.NotNil(t, chainList[0])
	require.NotNil(t, chainList[1])
	chainID1 := chainList[0].ChainId
	chainID2 := chainList[1].ChainId

	addr1 := sample.EthAddress()
	addr2 := sample.EthAddress()

	// Deploy the system contract and a ZRC20 contract
	deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
	zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID1, "alpha", "alpha")

	// Do some operation to populate the state
	_, err := k.DepositZRC20(ctx, zrc20, addr1, big.NewInt(100))
	require.NoError(t, err)
	_, err = k.DepositZRC20(ctx, zrc20, addr2, big.NewInt(200))
	require.NoError(t, err)

	// Check the state
	checkState := func() {
		// state that should not change
		balance, err := k.BalanceOfZRC4(ctx, zrc20, addr1)
		require.NoError(t, err)
		require.Equal(t, int64(100), balance.Int64())
		balance, err = k.BalanceOfZRC4(ctx, zrc20, addr2)
		require.NoError(t, err)
		require.Equal(t, int64(200), balance.Int64())
		totalSupply, err := k.TotalSupplyZRC4(ctx, zrc20)
		require.NoError(t, err)
		require.Equal(t, int64(10000300), totalSupply.Int64()) // 10000000 minted on deploy
	}

	checkState()
	chainID, err := k.QueryChainIDFromContract(ctx, zrc20)
	require.NoError(t, err)
	require.Equal(t, chainID1, chainID.Int64())

	// Deploy new zrc20
	newCodeAddress, err := k.DeployZRC20Contract(
		ctx,
		"beta",
		"BETA",
		18,
		chainID2,
		zetacommon.CoinType_ERC20,
		"beta",
		big.NewInt(90_000),
	)
	require.NoError(t, err)

	// Update the bytecode
	_, err = k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
		sample.AccAddress(),
		zrc20,
		newCodeAddress,
	))
	require.NoError(t, err)

	// Check the state
	// balances and total supply should remain
	// BYTECODE value is immutable and therefore part of the code, this value should change
	checkState()
	chainID, err = k.QueryChainIDFromContract(ctx, zrc20)
	require.NoError(t, err)
	require.Equal(t, chainID2, chainID.Int64())
}
