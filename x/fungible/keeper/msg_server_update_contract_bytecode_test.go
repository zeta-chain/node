package keeper_test

import (
	"errors"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/evm/x/vm/statedb"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/ptr"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func codeHashFromAddress(t *testing.T, ctx sdk.Context, k *keeper.Keeper, contractAddr string) string {
	res, err := k.CodeHash(ctx, &types.QueryCodeHashRequest{
		Address: contractAddr,
	})
	require.NoError(t, err)
	return res.CodeHash
}

func TestKeeper_UpdateContractBytecode(t *testing.T) {
	t.Run("can update the bytecode from another contract", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		// sample chainIDs and addresses
		chainList := chains.DefaultChainsList()
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
			coin.CoinType_ERC20,
			"beta",
			big.NewInt(90_000),
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		require.NoError(t, err)
		codeHash := codeHashFromAddress(t, ctx, k, newCodeAddress.Hex())

		// update the bytecode
		msg := types.NewMsgUpdateContractBytecode(
			admin,
			zrc20.Hex(),
			codeHash,
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateContractBytecode(ctx, msg)
		require.NoError(t, err)

		// check the returned new bytecode hash matches the one in the account
		acct := sdkk.EvmKeeper.GetAccount(ctx, zrc20)
		require.Equal(t, acct.CodeHash, ethcommon.HexToHash(codeHash).Bytes())

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
			coin.CoinType_ERC20,
			"gamma",
			big.NewInt(90_000),
			ptr.Ptr(sdkmath.NewUint(1000)),
		)
		codeHash = codeHashFromAddress(t, ctx, k, newCodeAddress.Hex())
		require.NoError(t, err)

		msg = types.NewMsgUpdateContractBytecode(
			admin,
			zrc20.Hex(),
			codeHash,
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)
		_, err = msgServer.UpdateContractBytecode(ctx, msg)
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
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		msgServer := keeper.NewMsgServerImpl(*k)

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		// deploy a connector
		wzeta, _, _, oldConnector, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		codeHash := codeHashFromAddress(t, ctx, k, oldConnector.Hex())

		// deploy a new connector that will become official connector
		newConnector, err := k.DeployConnectorZEVM(ctx, wzeta)
		require.NoError(t, err)
		require.NotEmpty(t, newConnector)
		assertContractDeployment(t, sdkk.EvmKeeper, ctx, newConnector)

		// can update the bytecode of the new connector with the old connector contract
		msg := types.NewMsgUpdateContractBytecode(
			admin,
			newConnector.Hex(),
			codeHash,
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateContractBytecode(ctx, msg)
		require.NoError(t, err)
	})

	t.Run("should fail if unauthorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateContractBytecode(
			admin,
			sample.EthAddress().Hex(),
			sample.Hash().Hex(),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.UpdateContractBytecode(ctx, msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail invalid contract address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateContractBytecode(
			admin,
			"invalid",
			sample.Hash().Hex(),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)

		_, err := msgServer.UpdateContractBytecode(ctx, msg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should fail if can't get contract account", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock:       true,
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		admin := sample.AccAddress()
		contractAddr := sample.EthAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		mockEVMKeeper.On(
			"GetAccount",
			mock.Anything,
			contractAddr,
		).Return(nil)

		msg := types.NewMsgUpdateContractBytecode(
			admin,
			contractAddr.Hex(),
			sample.Hash().Hex(),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateContractBytecode(ctx, msg)
		require.ErrorIs(t, err, types.ErrContractNotFound)

		mockEVMKeeper.AssertExpectations(t)
	})

	t.Run("should fail neither a zrc20 nor wzeta connector", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		wzeta, _, _, _, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// can't update the bytecode of the wzeta contract
		msg := types.NewMsgUpdateContractBytecode(
			admin,
			wzeta.Hex(),
			sample.Hash().Hex(),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateContractBytecode(ctx, msg)
		require.ErrorIs(t, err, types.ErrInvalidContract)
	})

	t.Run("should fail if system contract not found", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		_, _, _, connector, _ := deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)

		// remove system contract
		k.RemoveSystemContract(ctx)

		// can't update the bytecode of the wzeta contract
		msg := types.NewMsgUpdateContractBytecode(
			admin,
			connector.Hex(),
			sample.Hash().Hex(),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateContractBytecode(ctx, msg)
		require.ErrorIs(t, err, types.ErrSystemContractNotFound)
	})

	t.Run("should fail if can't set account with new bytecode", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
			UseEVMMock:       true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)
		admin := sample.AccAddress()

		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		contractAddr := sample.EthAddress()
		newCodeHash := sample.Hash().Hex()

		// set the contract as the connector
		k.SetSystemContract(ctx, types.SystemContract{
			ConnectorZevm: contractAddr.String(),
		})

		mockEVMKeeper.On(
			"GetAccount",
			mock.Anything,
			contractAddr,
		).Return(&statedb.Account{})

		mockEVMKeeper.On(
			"SetAccount",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(errors.New("can't set account"))

		msg := types.NewMsgUpdateContractBytecode(
			admin,
			contractAddr.Hex(),
			newCodeHash,
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateContractBytecode(ctx, msg)
		require.ErrorIs(t, err, types.ErrSetBytecode)

		mockEVMKeeper.AssertExpectations(t)
	})
}
