package bank

import (
	"encoding/json"
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	"github.com/zeta-chain/node/testutil/keeper"
)

func Test_IBankContract(t *testing.T) {
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	fungibleKeeper, ctx, keepers, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	t.Run("should create contract and check address and ABI", func(t *testing.T) {
		contract := NewIBankContract(ctx, keepers.BankKeeper, *fungibleKeeper, appCodec, gasConfig)
		require.NotNil(t, contract, "NewIBankContract() should not return a nil contract")

		address := contract.Address()
		require.Equal(t, ContractAddress, address, "contract address should match the precompiled address")

		abi := contract.Abi()
		require.NotNil(t, abi, "contract ABI should not be nil")
	})

	t.Run("should check methods are present in ABI", func(t *testing.T) {
		contract := NewIBankContract(ctx, keepers.BankKeeper, *fungibleKeeper, appCodec, gasConfig)
		abi := contract.Abi()

		require.NotNil(t, abi.Methods[DepositMethodName], "deposit method should be present in the ABI")
		require.NotNil(t, abi.Methods[WithdrawMethodName], "withdraw method should be present in the ABI")
		require.NotNil(t, abi.Methods[BalanceOfMethodName], "balanceOf method should be present in the ABI")
	})

	t.Run("should check gas requirements for methods", func(t *testing.T) {
		contract := NewIBankContract(ctx, keepers.BankKeeper, *fungibleKeeper, appCodec, gasConfig)
		abi := contract.Abi()
		var method [4]byte

		t.Run("deposit", func(t *testing.T) {
			gasDeposit := contract.RequiredGas(abi.Methods[DepositMethodName].ID)
			copy(method[:], abi.Methods[DepositMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				gasDeposit,
				"deposit method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				gasDeposit,
			)
		})

		t.Run("withdraw", func(t *testing.T) {
			gasWithdraw := contract.RequiredGas(abi.Methods[WithdrawMethodName].ID)
			copy(method[:], abi.Methods[WithdrawMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				gasWithdraw,
				"withdraw method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				gasWithdraw,
			)
		})

		t.Run("balanceOf", func(t *testing.T) {
			gasBalanceOf := contract.RequiredGas(abi.Methods[BalanceOfMethodName].ID)
			copy(method[:], abi.Methods[BalanceOfMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				gasBalanceOf,
				"balanceOf method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				gasBalanceOf,
			)
		})

		t.Run("invalid method", func(t *testing.T) {
			invalidMethodBytes := []byte("invalidMethod")
			gasInvalidMethod := contract.RequiredGas(invalidMethodBytes)
			require.Equal(
				t,
				uint64(0),
				gasInvalidMethod,
				"invalid method should require %d gas, got %d",
				uint64(0),
				gasInvalidMethod,
			)
		})
	})
}

func Test_InvalidMethod(t *testing.T) {
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	fungibleKeeper, ctx, keepers, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	contract := NewIBankContract(ctx, keepers.BankKeeper, *fungibleKeeper, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIBankContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	_, doNotExist := abi.Methods["invalidMethod"]
	require.False(t, doNotExist, "invalidMethod should not be present in the ABI")
}

func Test_InvalidABI(t *testing.T) {
	IBankMetaData.ABI = "invalid json"
	defer func() {
		if r := recover(); r != nil {
			require.IsType(t, &json.SyntaxError{}, r, "expected error type: json.SyntaxError, got: %T", r)
		}
	}()

	initABI()
}
