package keeper_test

import (
	"math/big"
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/systemcontract.sol"
	zrc20 "github.com/zeta-chain/protocol-contracts/pkg/contracts/zevm/zrc20.sol"
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
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex()))
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

	t.Run("can update the system contract if system contract not found", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		chains := zetacommon.DefaultChainsList()
		require.True(t, len(chains) > 1)
		require.NotNil(t, chains[0])
		require.NotNil(t, chains[1])

		wzeta, err := k.DeployWZETA(ctx)
		require.NoError(t, err)

		factory, err := k.DeployUniswapV2Factory(ctx)
		require.NoError(t, err)

		router, err := k.DeployUniswapV2Router02(ctx, factory, wzeta)
		require.NoError(t, err)

		// deploy a new system contracts
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wzeta, factory, router)
		require.NoError(t, err)

		// can update the system contract
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex()))
		require.NoError(t, err)

		// can retrieve the system contract
		sc, found := k.GetSystemContract(ctx)
		require.True(t, found)
		require.Equal(t, newSystemContract.Hex(), sc.SystemContract)
	})

	t.Run("should not update the system contract if not admin", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeper(t)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// deploy a new system contracts
		wzeta, factory, router, _, oldSystemContract := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wzeta, factory, router)
		require.NoError(t, err)
		require.NotEqual(t, oldSystemContract, newSystemContract)

		// should not update the system contract if not admin
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(sample.AccAddress(), newSystemContract.Hex()))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
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
		require.NoError(t, err)
		require.NotEqual(t, oldSystemContract, newSystemContract)

		// should not update the system contract if invalid address
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, "invalid"))
		require.Error(t, err)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should not update if any of 3 evm calls for foreign coin fail", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin, observertypes.Policy_Type_group2)

		chains := zetacommon.DefaultChainsList()
		require.True(t, len(chains) > 1)
		require.NotNil(t, chains[0])
		require.NotNil(t, chains[1])
		chainID1 := chains[0].ChainId

		wzeta, factory, router, _, _ := deploySystemContractsWithMockEvmKeeper(t, ctx, k, mockEVMKeeper)
		// setup mocks and setup gas coin
		var encodedAddress [32]byte
		copy(encodedAddress[12:], router[:])
		uniswapMock := &evmtypes.MsgEthereumTxResponse{
			Ret: encodedAddress[:],
		}
		mockEVMKeeper.MockEVMSuccessCallTimes(4)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(uniswapMock)
		mockEVMKeeper.MockEVMSuccessCallOnce()

		addLiqMockReturn := &evmtypes.MsgEthereumTxResponse{
			Ret: make([]byte, 3*32),
		}
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(addLiqMockReturn)

		setupGasCoin(t, ctx, k, mockEVMKeeper, chainID1, "foo", "foo")

		// deploy a new system contracts
		mockEVMKeeper.MockEVMSuccessCallOnce()
		newSystemContract, err := k.DeployContract(ctx, systemcontract.SystemContractMetaData, wzeta, factory, router)
		require.NoError(t, err)

		// fail on first evm call
		mockEVMKeeper.MockEVMFailCallOnce()

		// can update the system contract
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex()))
		require.ErrorIs(t, err, types.ErrContractCall)

		// fail on second evm call
		mockEVMKeeper.MockEVMSuccessCallOnce()
		mockEVMKeeper.MockEVMFailCallOnce()

		// can update the system contract
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex()))
		require.ErrorIs(t, err, types.ErrContractCall)

		// fail on third evm call
		mockEVMKeeper.MockEVMSuccessCallTimes(2)
		mockEVMKeeper.MockEVMFailCallOnce()

		// can update the system contract
		_, err = msgServer.UpdateSystemContract(ctx, types.NewMsgUpdateSystemContract(admin, newSystemContract.Hex()))
		require.ErrorIs(t, err, types.ErrContractCall)
	})
}
