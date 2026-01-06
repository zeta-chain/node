package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	evmtypes "github.com/cosmos/evm/x/vm/types"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/zrc20.sol"

	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
	"github.com/zeta-chain/node/x/fungible/keeper"
	"github.com/zeta-chain/node/x/fungible/types"
)

func TestKeeper_UpdateZRC20WithdrawFee(t *testing.T) {
	t.Run("can update the withdraw fee", func(t *testing.T) {
		k, ctx, sdkk, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		chainID := getValidChainID(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// set coin admin
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		keepertest.MockGetChainListEmpty(&authorityMock.Mock)

		// deploy the system contract and a ZRC20 contract
		deploySystemContracts(t, ctx, k, sdkk.EvmKeeper)
		zrc20Addr := setupGasCoin(t, ctx, k, sdkk.EvmKeeper, chainID, "alpha", "alpha")

		// initial protocol fee is zero
		protocolFee, err := k.QueryProtocolFlatFee(ctx, zrc20Addr)
		require.NoError(t, err)
		require.Zero(t, protocolFee.Uint64())

		// can update the protocol fee and gas limit
		msg := types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			zrc20Addr.String(),
			math.NewUint(42),
			math.NewUint(42),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateZRC20WithdrawFee(ctx, msg)
		require.NoError(t, err)

		// can query the updated fee
		protocolFee, err = k.QueryProtocolFlatFee(ctx, zrc20Addr)
		require.NoError(t, err)
		require.Equal(t, uint64(42), protocolFee.Uint64())
		gasLimit, err := k.QueryGasLimit(ctx, zrc20Addr)
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit.Uint64())

		// can update protocol fee only
		msg = types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			zrc20Addr.String(),
			math.NewUint(43),
			math.Uint{},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateZRC20WithdrawFee(ctx, msg)
		require.NoError(t, err)
		protocolFee, err = k.QueryProtocolFlatFee(ctx, zrc20Addr)
		require.NoError(t, err)
		require.Equal(t, uint64(43), protocolFee.Uint64())
		gasLimit, err = k.QueryGasLimit(ctx, zrc20Addr)
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit.Uint64())

		// can update gas limit only
		msg = types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			zrc20Addr.String(),
			math.Uint{},
			math.NewUint(44),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateZRC20WithdrawFee(ctx, msg)
		require.NoError(t, err)
		protocolFee, err = k.QueryProtocolFlatFee(ctx, zrc20Addr)
		require.NoError(t, err)
		require.Equal(t, uint64(43), protocolFee.Uint64())
		gasLimit, err = k.QueryGasLimit(ctx, zrc20Addr)
		require.NoError(t, err)
		require.Equal(t, uint64(44), gasLimit.Uint64())
	})

	t.Run("should fail if not authorized", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)

		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			sample.EthAddress().String(),
			math.NewUint(42),
			math.Uint{},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, authoritytypes.ErrUnauthorized)
		_, err := msgServer.UpdateZRC20WithdrawFee(ctx, msg)
		require.ErrorIs(t, err, authoritytypes.ErrUnauthorized)
	})

	t.Run("should fail if invalid zrc20 address", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			"invalid_address",
			math.NewUint(42),
			math.Uint{},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateZRC20WithdrawFee(ctx, msg)
		require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)
	})

	t.Run("should fail if can't retrieve the foreign coin", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		msg := types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			sample.EthAddress().String(),
			math.NewUint(42),
			math.Uint{},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateZRC20WithdrawFee(ctx, msg)
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if can't query old fee", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseAuthorityMock: true,
		})

		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)

		// setup
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)
		zrc20 := sample.EthAddress()

		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20.String()))

		// the method shall fail since we only set the foreign coin manually in the store but didn't deploy the contract
		msg := types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			zrc20.String(),
			math.NewUint(42),
			math.Uint{},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err := msgServer.UpdateZRC20WithdrawFee(ctx, msg)
		require.ErrorIs(t, err, types.ErrContractCall)
	})

	t.Run("should fail if contract call for setting new protocol fee fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock:       true,
			UseAuthorityMock: true,
		})
		msgServer := keeper.NewMsgServerImpl(*k)
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)

		// setup
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		zrc20Addr := sample.EthAddress()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20Addr.String()))

		// evm mocks
		mockEVMKeeper.On("EstimateGas", mock.Anything, mock.Anything).Maybe().Return(
			&evmtypes.EstimateGasResponse{Gas: 1000},
			nil,
		)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		// this is the query (commit == false)
		zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		protocolFlatFee, err := zrc20ABI.Methods["PROTOCOL_FLAT_FEE"].Outputs.Pack(big.NewInt(42))
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: protocolFlatFee})

		gasLimit, err := zrc20ABI.Methods["GAS_LIMIT"].Outputs.Pack(big.NewInt(42))
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: gasLimit})

		// this is the update call (commit == true)
		mockEVMKeeper.MockEVMFailCallOnce()

		msg := types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			zrc20Addr.String(),
			math.NewUint(42),
			math.Uint{},
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateZRC20WithdrawFee(ctx, msg)
		require.ErrorIs(t, err, types.ErrContractCall)

		mockEVMKeeper.AssertExpectations(t)
	})

	t.Run("should fail if query gas limit fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.FungibleKeeperWithMocks(t, keepertest.FungibleMockOptions{
			UseEVMMock:       true,
			UseAuthorityMock: true,
		})
		k.GetAuthKeeper().GetModuleAccount(ctx, types.ModuleName)
		msgServer := keeper.NewMsgServerImpl(*k)
		mockEVMKeeper := keepertest.GetFungibleEVMMock(t, k)

		// setup
		admin := sample.AccAddress()
		authorityMock := keepertest.GetFungibleAuthorityMock(t, k)

		zrc20Addr := sample.EthAddress()
		k.SetForeignCoins(ctx, sample.ForeignCoins(t, zrc20Addr.String()))

		// evm mocks
		mockEVMKeeper.On("EstimateGas", mock.Anything, mock.Anything).Maybe().Return(
			&evmtypes.EstimateGasResponse{Gas: 1000},
			nil,
		)
		mockEVMKeeper.On("WithChainID", mock.Anything).Maybe().Return(ctx)
		mockEVMKeeper.On("ChainID").Maybe().Return(big.NewInt(1))

		// this is the query (commit == false)
		zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
		require.NoError(t, err)
		protocolFlatFee, err := zrc20ABI.Methods["PROTOCOL_FLAT_FEE"].Outputs.Pack(big.NewInt(42))
		require.NoError(t, err)
		mockEVMKeeper.MockEVMSuccessCallOnceWithReturn(&evmtypes.MsgEthereumTxResponse{Ret: protocolFlatFee})

		_, err = zrc20ABI.Methods["GAS_LIMIT"].Outputs.Pack(big.NewInt(42))
		require.NoError(t, err)
		mockEVMKeeper.MockEVMFailCallOnce()

		msg := types.NewMsgUpdateZRC20WithdrawFee(
			admin,
			zrc20Addr.String(),
			math.Uint{},
			math.NewUint(42),
		)
		keepertest.MockCheckAuthorization(&authorityMock.Mock, msg, nil)
		_, err = msgServer.UpdateZRC20WithdrawFee(ctx, msg)
		require.ErrorIs(t, err, types.ErrContractCall)

		mockEVMKeeper.AssertExpectations(t)
	})
}
