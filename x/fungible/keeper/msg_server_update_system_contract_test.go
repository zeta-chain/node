package keeper_test

import (
	"math/big"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
	zetacommon "github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/keeper"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_UpdateSystemContract(t *testing.T) {
	t.Run("can update the system contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		queryZRC20SystemContract := func(contract common.Address) string {
			abi, err := zrc20.ZRC20MetaData.GetAbi()
			assert.NoError(t, err)
			res, err := k.CallEVM(ctx, *abi, types.ModuleAddressEVM, contract, keeper.BigIntZero, nil, false, false, "SYSTEM_CONTRACT_ADDRESS")
			assert.NoError(t, err)
			unpacked, err := abi.Unpack("SYSTEM_CONTRACT_ADDRESS", res.Ret)
			assert.NoError(t, err)
			address, ok := unpacked[0].(common.Address)
			assert.True(t, ok)
			return address.Hex()
		}

		chains := zetacommon.DefaultChainsList()
		assert.True(t, len(chains) > 1)
		assert.NotNil(t, chains[0])
		assert.NotNil(t, chains[1])
		chainID1 := chains[0].ChainId
		chainID2 := chains[1].ChainId

		wzeta, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		gas1 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID1, "foo", "foo")
		gas2 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID2, "bar", "bar")

		// deploy a new system contracts
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wzeta, factory, router)
		assert.NoError(t, err)
		assert.NotEqual(t, oldSystemContract, newSystemContract)

		// can update the system contract
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex()))
		assert.NoError(t, err)

		// can retrieve the system contract
		sc, found := k.GetSystemContract(ctx)
		assert.True(t, found)
		assert.Equal(t, newSystemContract.Hex(), sc.SystemContract)

		// check gas updated
		foundGas1, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID1))
		assert.NoError(t, err)
		assert.Equal(t, gas1, foundGas1)
		foundGas2, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID2))
		assert.NoError(t, err)
		assert.Equal(t, gas2, foundGas2)

		assert.Equal(t, newSystemContract.Hex(), queryZRC20SystemContract(gas1))
		assert.Equal(t, newSystemContract.Hex(), queryZRC20SystemContract(gas2))
	})

	t.Run("should not update the system contract if not admin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// deploy a new system contracts
		wzeta, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wzeta, factory, router)
		assert.NoError(t, err)
		assert.NotEqual(t, oldSystemContract, newSystemContract)

		// should not update the system contract if not admin
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(sample.AccAddress(), newSystemContract.Hex()))
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should not update the system contract if invalid address", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		// deploy a new system contracts
		wzeta, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wzeta, factory, router)
		assert.NoError(t, err)
		assert.NotEqual(t, oldSystemContract, newSystemContract)

		// should not update the system contract if invalid address
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, "invalid"))
		assert.Error(t, err)
		assert.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})
}
