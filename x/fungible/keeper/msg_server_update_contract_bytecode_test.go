package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/x/evm/statedb"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	zetacommon "github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func setAdminDeployFungibleCoin(ctx sdk.Context, zk keepertest.ZetaKeepers, admin string) {
	zk.ObserverKeeper.SetParams(ctx, observertypes.Params{
		AdminPolicy: []*observertypes.Admin_Policy{
			{
				PolicyType: observertypes.Policy_Type_deploy_fungible_coin,
				Address:    admin,
			},
		},
	})
}

func TestKeeper_UpdateContractBytecode(t *testing.T) {
	t.Run("can update the bytecode from another contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		// set admin policy
		setAdminDeployFungibleCoin(ctx, zk, admin)

		// sample chainIDs and addresses
		chainList := zetacommon.DefaultChainsList()
		require.True(t, len(chainList) > 1)
		require.NotNil(t, chainList[0])
		require.NotNil(t, chainList[1])
		require.NotEqual(t, chainList[0].ChainId, chainList[1].ChainId)
		chainID1 := chainList[0].ChainId
		chainID2 := chainList[1].ChainId

		addr1 := sample.EthAddress()
		addr2 := sample.EthAddress()

		// deploy the system contract and a ZRC20 contract
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20 := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID1, "alpha", "alpha")

		// do some operation to populate the state
		_, err := k.DepositZRC20(ctx, zrc20, addr1, big.NewInt(100))
		require.NoError(t, err)
		_, err = k.DepositZRC20(ctx, zrc20, addr2, big.NewInt(200))
		require.NoError(t, err)

		// check the state
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

		// deploy new zrc20
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

		// update the bytecode
		res, err := k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			zrc20,
			newCodeAddress,
		))
		require.NoError(t, err)

		// check the returned new bytecode hash matches the one in the account
		acct := sdkk.EvmKeeper.GetAccount(ctx, zrc20)
		require.Equal(t, acct.CodeHash, res.NewBytecodeHash)

		// check the state
		// balances and total supply should remain
		// BYTECODE value is immutable and therefore part of the code, this value should change
		checkState()
		chainID, err = k.QueryChainIDFromContract(ctx, zrc20)
		require.NoError(t, err)
		require.Equal(t, chainID2, chainID.Int64())

		// can continue to interact with the contract
		_, err = k.DepositZRC20(ctx, zrc20, addr1, big.NewInt(1000))
		require.NoError(t, err)
		balance, err := k.BalanceOfZRC4(ctx, zrc20, addr1)
		require.NoError(t, err)
		require.Equal(t, int64(1100), balance.Int64())
		totalSupply, err := k.TotalSupplyZRC4(ctx, zrc20)
		require.NoError(t, err)
		require.Equal(t, int64(10001300), totalSupply.Int64())

		// can change again bytecode
		newCodeAddress, err = k.DeployZRC20Contract(
			ctx,
			"gamma",
			"GAMMA",
			18,
			chainID1,
			zetacommon.CoinType_ERC20,
			"gamma",
			big.NewInt(90_000),
		)
		require.NoError(t, err)
		_, err = k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			zrc20,
			newCodeAddress,
		))
		require.NoError(t, err)
		balance, err = k.BalanceOfZRC4(ctx, zrc20, addr1)
		require.NoError(t, err)
		require.Equal(t, int64(1100), balance.Int64())
		totalSupply, err = k.TotalSupplyZRC4(ctx, zrc20)
		require.NoError(t, err)
		require.Equal(t, int64(10001300), totalSupply.Int64())
		chainID, err = k.QueryChainIDFromContract(ctx, zrc20)
		require.NoError(t, err)
		require.Equal(t, chainID1, chainID.Int64())
	})

	t.Run("can update the bytecode of the wzeta connector contract", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		// deploy a connector
		setAdminDeployFungibleCoin(ctx, zk, admin)
		wzeta, _, _, oldConnector, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// deploy a new connector that will become official connector
		newConnector, err := k.DeployConnectorZEVM(ctx, wzeta)
		require.NoError(t, err)
		require.NotEmpty(t, newConnector)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, newConnector)

		// can update the bytecode of the new connector with the old connector contract
		_, err = k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			newConnector,
			oldConnector,
		))
		require.NoError(t, err)
	})

	t.Run("should fail if unauthorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeper(t)

		_, err := k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			sample.AccAddress(),
			sample.EthAddress(),
			sample.EthAddress(),
		))
		require.ErrorIs(t, err, sdkerrors.ErrUnauthorized)
	})

	t.Run("should fail invalid contract address", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeper(t)
		admin := sample.AccAddress()
		setAdminDeployFungibleCoin(ctx, zk, admin)

		_, err := k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			ethcommon.HexToAddress("invalid"),
			sample.EthAddress(),
		))
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should fail if can't get contract account", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		admin := sample.AccAddress()
		setAdminDeployFungibleCoin(ctx, zk, admin)
		contractAddr := sample.EthAddress()

		mockEVMKeeper.On(
			"GetAccount",
			mock.Anything,
			contractAddr,
		).Return(nil)

		_, err := k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			contractAddr,
			sample.EthAddress(),
		))
		require.ErrorIs(t, err, types.ErrContractNotFound)

		mockEVMKeeper.AssertExpectations(t)
	})

	t.Run("should fail neither a zrc20 nor wzeta connector", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		setAdminDeployFungibleCoin(ctx, zk, admin)
		wzeta, _, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// can't update the bytecode of the wzeta contract
		_, err := k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			wzeta,
			sample.EthAddress(),
		))
		require.ErrorIs(t, err, types.ErrInvalidContract)
	})

	t.Run("should fail if system contract not found", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.FungibleKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		setAdminDeployFungibleCoin(ctx, zk, admin)
		_, _, _, connector, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// remove system contract
		k.RemoveSystemContract(ctx)

		// can't update the bytecode of the wzeta contract
		_, err := k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			connector,
			sample.EthAddress(),
		))
		require.ErrorIs(t, err, types.ErrSystemContractNotFound)
	})

	t.Run("should fail if invalid bytecode address", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		admin := sample.AccAddress()
		setAdminDeployFungibleCoin(ctx, zk, admin)

		mockEVMKeeper.On(
			"GetAccount",
			mock.Anything,
			mock.Anything,
		).Return(&statedb.Account{})

		_, err := k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			sample.EthAddress(),
			ethcommon.HexToAddress("invalid"),
		))
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

		mockEVMKeeper.AssertExpectations(t)
	})

	t.Run("should fail if can't get new bytecode account", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		admin := sample.AccAddress()
		setAdminDeployFungibleCoin(ctx, zk, admin)
		contractAddr := sample.EthAddress()
		newBytecodeAddr := sample.EthAddress()

		mockEVMKeeper.On(
			"GetAccount",
			mock.Anything,
			contractAddr,
		).Return(&statedb.Account{})

		mockEVMKeeper.On(
			"GetAccount",
			mock.Anything,
			newBytecodeAddr,
		).Return(nil)

		_, err := k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			contractAddr,
			newBytecodeAddr,
		))
		require.ErrorIs(t, err, types.ErrContractNotFound)

		mockEVMKeeper.AssertExpectations(t)
	})

	t.Run("should fail if can't set account with new bytecode", func(t *testing.T) {
		k, ctx, _, zk := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock: true,
		})
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		admin := sample.AccAddress()
		setAdminDeployFungibleCoin(ctx, zk, admin)
		contractAddr := sample.EthAddress()
		newBytecodeAddr := sample.EthAddress()

		mockEVMKeeper.On(
			"GetAccount",
			mock.Anything,
			contractAddr,
		).Return(&statedb.Account{})

		mockEVMKeeper.On(
			"GetAccount",
			mock.Anything,
			newBytecodeAddr,
		).Return(&statedb.Account{})

		mockEVMKeeper.On(
			"SetAccount",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(errors.New("can't set account"))

		_, err := k.UpdateContractBytecode(ctx, types.NewMsgUpdateContractBytecode(
			admin,
			contractAddr,
			newBytecodeAddr,
		))
		require.ErrorIs(t, err, types.ErrSetBytecode)

		mockEVMKeeper.AssertExpectations(t)
	})
}
