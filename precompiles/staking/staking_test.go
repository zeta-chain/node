package staking

import (
	"encoding/json"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	"github.com/zeta-chain/zetacore/testutil/keeper"
)

func Test_IStakingContract(t *testing.T) {
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	keys, memKeys, tkeys, allKeys := keeper.StoreKeys()
	cdc := keeper.NewCodec()
	sdkKeepers := keeper.NewSDKKeepersWithKeys(cdc, keys, memKeys, tkeys, allKeys)

	gasConfig := storetypes.TransientGasConfig()

	t.Run("should create contract and check address and ABI", func(t *testing.T) {
		contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
		require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

		address := contract.Address()
		require.Equal(t, ContractAddress, address, "contract address should match the precompiled address")

		abi := contract.Abi()
		require.NotNil(t, abi, "contract ABI should not be nil")
	})

	t.Run("should check methods are present in ABI", func(t *testing.T) {
		contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
		abi := contract.Abi()

		require.NotNil(t, abi.Methods[DelegateMethodName], "delegate method should be present in the ABI")
		require.NotNil(t, abi.Methods[UndelegateMethodName], "undelegate method should be present in the ABI")
		require.NotNil(
			t,
			abi.Methods[RedelegateMethodName],
			"redelegate method should be present in the ABI",
		)
	})

	t.Run("should check gas requirements for methods", func(t *testing.T) {
		contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
		abi := contract.Abi()
		var method [4]byte

		t.Run("delegate", func(t *testing.T) {
			delegate := contract.RequiredGas(abi.Methods[DelegateMethodName].ID)
			copy(method[:], abi.Methods[DelegateMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				delegate,
				"delegate method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				delegate,
			)
		})

		t.Run("undelegate", func(t *testing.T) {
			undelegate := contract.RequiredGas(abi.Methods[UndelegateMethodName].ID)
			copy(method[:], abi.Methods[UndelegateMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				undelegate,
				"undelegate method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				undelegate,
			)
		})

		t.Run("redelegate", func(t *testing.T) {
			redelegate := contract.RequiredGas(abi.Methods[RedelegateMethodName].ID)
			copy(method[:], abi.Methods[RedelegateMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				redelegate,
				"redelegate method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				redelegate,
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
	keys, memKeys, tkeys, allKeys := keeper.StoreKeys()
	cdc := keeper.NewCodec()
	sdkKeepers := keeper.NewSDKKeepersWithKeys(cdc, keys, memKeys, tkeys, allKeys)
	gasConfig := storetypes.TransientGasConfig()

	contract := NewIStakingContract(&sdkKeepers.StakingKeeper, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIStakingContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	_, doNotExist := abi.Methods["invalidMethod"]
	require.False(t, doNotExist, "invalidMethod should not be present in the ABI")
}

func Test_InvalidABI(t *testing.T) {
	IStakingMetaData.ABI = "invalid json"
	defer func() {
		if r := recover(); r != nil {
			require.IsType(t, &json.SyntaxError{}, r, "expected error type: json.SyntaxError, got: %T", r)
		}
	}()

	initABI()
}
