package keeper_test

import (
	"math/big"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	zetacommon "github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_UpdateSystemContract(t *testing.T) {
	t.Run("can update the system contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		queryZRC20SystemContract := func(contract common.Address) string {
			abi, err := zrc20.ZRC20MetaData.GetAbi()
			require.NoError(t, err)
			res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, contract, keeper.BigIntZero, nil, false, false, "SYSTEM_CONTRACT_ADDRESS")
			require.NoError(t, err)
			unpacked, err := abi.Unpack("SYSTEM_CONTRACT_ADDRESS", res.Ret)
			require.NoError(t, err)
			address, ok := unpacked[0].(common.Address)
			require.True(t, ok)
			return address.Hex()
		}

		chains := zetacommon.DefaultChainsList()
		require.True(t, len(chains) > 1)
		require.NotNil(t, chains[0])
		require.NotNil(t, chains[1])
		chainID1 := chains[0].ChainId
		chainID2 := chains[1].ChainId

		wzeta, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		gas1 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID1, "foo", "foo")
		gas2 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID2, "bar", "bar")

		// deploy a new system contracts
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wzeta, factory, router)
		require.NoError(t, err)
		require.NotEqual(t, oldSystemContract, newSystemContract)

		// can update the system contract
		_, err = k.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex()))
		require.NoError(t, err)

		// can retrieve the system contract
		sc, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, newSystemContract.Hex(), sc.SystemContract)

		// check gas updated
		foundGas1, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID1))
		require.NoError(t, err)
		require.Equal(t, gas1, foundGas1)
		foundGas2, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID2))
		require.NoError(t, err)
		require.Equal(t, gas2, foundGas2)

		require.Equal(t, newSystemContract.Hex(), queryZRC20SystemContract(gas1))
		require.Equal(t, newSystemContract.Hex(), queryZRC20SystemContract(gas2))
	})

	t.Run("should not update the system contract if not admin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// deploy a new system contracts
		wzeta, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wzeta, factory, router)
		require.NoError(t, err)
		require.NotEqual(t, oldSystemContract, newSystemContract)

		// should not update the system contract if not admin
		_, err = k.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(sample.AccAddress(), newSystemContract.Hex()))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should not update the system contract if invalid address", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)

		// deploy a new system contracts
		wzeta, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wzeta, factory, router)
		require.NoError(t, err)
		require.NotEqual(t, oldSystemContract, newSystemContract)

		// should not update the system contract if invalid address
		_, err = k.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, "invalid"))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})
}
