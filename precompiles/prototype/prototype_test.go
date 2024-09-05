package prototype

import (
	"encoding/json"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	"github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/testutil/keeper"
)

func Test_IPrototypeContract(t *testing.T) {
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, _, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	t.Run("should create contract and check address and ABI", func(t *testing.T) {
		contract := NewIPrototypeContract(k, appCodec, gasConfig)
		require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

		address := contract.Address()
		require.Equal(t, ContractAddress, address, "contract address should match the precompiled address")

		abi := contract.Abi()
		require.NotNil(t, abi, "contract ABI should not be nil")
	})

	t.Run("should check methods are present in ABI", func(t *testing.T) {
		contract := NewIPrototypeContract(k, appCodec, gasConfig)
		abi := contract.Abi()

		require.NotNil(t, abi.Methods[Bech32ToHexAddrMethodName], "bech32ToHexAddr method should be present in the ABI")
		require.NotNil(t, abi.Methods[Bech32ifyMethodName], "bech32ify method should be present in the ABI")
		require.NotNil(
			t,
			abi.Methods[GetGasStabilityPoolBalanceName],
			"getGasStabilityPoolBalance method should be present in the ABI",
		)
	})

	t.Run("should check gas requirements for methods", func(t *testing.T) {
		contract := NewIPrototypeContract(k, appCodec, gasConfig)
		abi := contract.Abi()
		var method [4]byte

		t.Run("bech32ToHexAddr", func(t *testing.T) {
			gasBech32ToHex := contract.RequiredGas(abi.Methods[Bech32ToHexAddrMethodName].ID)
			copy(method[:], abi.Methods[Bech32ToHexAddrMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				gasBech32ToHex,
				"bech32ToHexAddr method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				gasBech32ToHex,
			)
		})

		t.Run("bech32ify", func(t *testing.T) {
			gasBech32ify := contract.RequiredGas(abi.Methods[Bech32ifyMethodName].ID)
			copy(method[:], abi.Methods[Bech32ifyMethodName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				gasBech32ify,
				"bech32ify method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				gasBech32ify,
			)
		})

		t.Run("getGasStabilityPoolBalance", func(t *testing.T) {
			gasGetGasStabilityPoolBalance := contract.RequiredGas(abi.Methods[GetGasStabilityPoolBalanceName].ID)
			copy(method[:], abi.Methods[GetGasStabilityPoolBalanceName].ID[:4])
			baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
			require.Equal(
				t,
				GasRequiredByMethod[method]+baseCost,
				gasGetGasStabilityPoolBalance,
				"getGasStabilityPoolBalance method should require %d gas, got %d",
				GasRequiredByMethod[method]+baseCost,
				gasGetGasStabilityPoolBalance,
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

func Test_Bech32ToHexAddress(t *testing.T) {
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, _, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	contract := NewIPrototypeContract(k, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	methodID := abi.Methods[Bech32ToHexAddrMethodName]

	t.Run("should succeed with valid bech32 address", func(t *testing.T) {
		args := []interface{}{"zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u"}

		rawBytes, err := contract.Bech32ToHexAddr(&methodID, args)
		require.NoError(t, err, "Bech32ToHexAddr should not return an error")

		addr := common.BytesToAddress(rawBytes[12:])
		require.Equal(
			t,
			common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE"),
			addr,
			"Bech32ToHexAddr should return the correct address, got: %v",
			addr,
		)
	})

	t.Run("should fail if invalid argument type", func(t *testing.T) {
		args := []interface{}{1}

		_, err := contract.Bech32ToHexAddr(&methodID, args)
		require.Error(t, err, "expected invalid argument; wanted string; got: %T", args[0])
	})

	t.Run("should fail if invalid bech32 address", func(t *testing.T) {
		t.Run("invalid bech32 format", func(t *testing.T) {
			args := []interface{}{"foobar"}
			_, err := contract.Bech32ToHexAddr(&methodID, args)
			require.Error(t, err, "expected error; invalid bech32 address")
		})

		t.Run("invalid bech32 prefix", func(t *testing.T) {
			args := []interface{}{"foobar1"}
			_, err := contract.Bech32ToHexAddr(&methodID, args)
			require.Error(t, err, "expected error; invalid bech32 addresss")
		})

		t.Run("invalid bech32 decoding", func(t *testing.T) {
			args := []interface{}{"foobar1foobar"}
			_, err := contract.Bech32ToHexAddr(&methodID, args)
			require.Error(t, err, "expected error; decoding bech32 failed")
		})

		t.Run("invalid number of arguments", func(t *testing.T) {
			args := []interface{}{"zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u", "second argument"}
			_, err := contract.Bech32ToHexAddr(&methodID, args)
			require.Error(t, err, "expected invalid number of arguments; expected 1; got: 2")
			require.IsType(
				t,
				&types.ErrInvalidNumberOfArgs{},
				err,
				"expected error type: ErrInvalidNumberOfArgs, got: %T",
				err,
			)
		})
	})

	t.Run("should fail if empty address argument", func(t *testing.T) {
		args := []interface{}{""}
		_, err := contract.Bech32ToHexAddr(&methodID, args)
		require.Error(t, err, "expected error invalid bech32 address: %v", args[0])
	})
}

func Test_Bech32ify(t *testing.T) {
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, _, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	contract := NewIPrototypeContract(k, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	methodID := abi.Methods[Bech32ifyMethodName]

	t.Run("should succeed with zeta HRP", func(t *testing.T) {
		args := []interface{}{"zeta", common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE")}

		rawBytes, err := contract.Bech32ify(&methodID, args)
		require.NoError(t, err, "Bech32ify prefix zeta should not return an error")

		zetaAddr := string(rawBytes[64:107])
		require.Equal(
			t,
			"zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u",
			zetaAddr,
			"Bech32ify prefix zeta should return the correct address, got: %v",
			zetaAddr,
		)
	})

	t.Run("should succeed with cosmos HRP", func(t *testing.T) {
		args := []interface{}{"cosmos", common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE")}

		rawBytes, err := contract.Bech32ify(&methodID, args)
		require.NoError(t, err, "Bech32ify prefix cosmos should not return an error")

		zetaAddr := string(rawBytes[64:107])
		require.Equal(
			t,
			"cosmos1h8duy2dltz9xz0qqhm5wvcnj02upy887lqaq",
			zetaAddr,
			"Bech32ify prefix cosmos should return the correct address, got: %v",
			zetaAddr,
		)
	})

	t.Run("should fail with invalid arguments", func(t *testing.T) {
		t.Run("too many arguments", func(t *testing.T) {
			args := []interface{}{
				"zeta",
				common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE"),
				"third argument",
			}
			_, err := contract.Bech32ify(&methodID, args)
			require.Error(t, err, "expected invalid number of arguments; expected 2; got: 3")
			require.IsType(
				t,
				&types.ErrInvalidNumberOfArgs{},
				err,
				"expected error type: ErrInvalidNumberOfArgs, got: %T",
				err,
			)
		})

		t.Run("invalid HRP", func(t *testing.T) {
			args := []interface{}{1337, common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE")}
			_, err := contract.Bech32ify(&methodID, args)
			require.Error(t, err, "expected error invalid bech32 human readable prefix (HRP)")
		})

		t.Run("invalid hex address", func(t *testing.T) {
			args := []interface{}{"zeta", 1337}
			_, err := contract.Bech32ify(&methodID, args)
			require.Error(t, err, "expected error invalid hex address")
		})

		t.Run("empty HRP", func(t *testing.T) {
			args := []interface{}{"", common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE")}
			_, err := contract.Bech32ify(&methodID, args)
			require.Error(
				t,
				err,
				"expected error invalid bech32 human readable prefix (HRP). Please provide either an account, validator, or consensus address prefix (eg: cosmos, cosmosvaloper, cosmosvalcons)",
			)
		})
	})
}

func Test_GetGasStabilityPoolBalance(t *testing.T) {
	// Only check the function is called correctly inside the contract, and it returns the expected error.
	// Configuring a local environment for this contract would require deploying system contracts and gas pools.
	// This method is tested thoroughly in the e2e tests.
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, ctx, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	contract := NewIPrototypeContract(k, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	methodID := abi.Methods[GetGasStabilityPoolBalanceName]

	args := []interface{}{int64(1337)}

	_, err := contract.GetGasStabilityPoolBalance(ctx, &methodID, args)
	require.Error(
		t,
		err,
		"error calling fungible keeper: failed to get system contract variable: state variable not found",
	)

	t.Run("should fail with invalid arguments", func(t *testing.T) {
		t.Run("invalid number of arguments", func(t *testing.T) {
			args := []interface{}{int64(1337), "second argument"}
			_, err := contract.GetGasStabilityPoolBalance(ctx, &methodID, args)
			require.Error(t, err, "expected invalid number of arguments; expected 2; got: 3")
			require.IsType(
				t,
				&types.ErrInvalidNumberOfArgs{},
				err,
				"expected error type: ErrInvalidNumberOfArgs, got: %T",
				err,
			)
		})

		t.Run("invalid chainID", func(t *testing.T) {
			args := []interface{}{"foobar"}
			_, err := contract.GetGasStabilityPoolBalance(ctx, &methodID, args)
			require.Error(t, err, "expected int64, got: %T", args[0])
			require.IsType(t, types.ErrInvalidArgument{}, err, "expected error type: ErrInvalidArgument, got: %T", err)
		})
	})
}

func Test_InvalidMethod(t *testing.T) {
	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, _, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	contract := NewIPrototypeContract(k, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	_, doNotExist := abi.Methods["invalidMethod"]
	require.False(t, doNotExist, "invalidMethod should not be present in the ABI")
}

func Test_InvalidABI(t *testing.T) {
	IPrototypeMetaData.ABI = "invalid json"
	defer func() {
		if r := recover(); r != nil {
			require.IsType(t, &json.SyntaxError{}, r, "expected error type: json.SyntaxError, got: %T", r)
		}
	}()

	initABI()
}
