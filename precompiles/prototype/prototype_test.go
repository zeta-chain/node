package prototype

import (
	"encoding/json"
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	"github.com/zeta-chain/zetacore/precompiles/types"
	"github.com/zeta-chain/zetacore/testutil/keeper"
)

func Test_IPrototypeContract(t *testing.T) {
	/*
		Configuration
	*/

	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, _, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	/*
		Contract and ABI tests
	*/

	// Create a new IPrototypeContract instance and get Address and Abi.
	contract := NewIPrototypeContract(k, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

	address := contract.Address()
	require.Equal(t, ContractAddress, address, "contract address should match the precompiled address")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	// Check all methods are present in the ABI.
	bech32ToHex := abi.Methods[Bech32ToHexAddrMethodName]
	require.NotNil(t, bech32ToHex, "bech32ToHexAddr method should be present in the ABI")

	bech32ify := abi.Methods[Bech32ifyMethodName]
	require.NotNil(t, bech32ify, "bech32ify method should be present in the ABI")

	getGasStabilityPoolBalance := abi.Methods[GetGasStabilityPoolBalanceName]
	require.NotNil(t, getGasStabilityPoolBalance, "getGasStabilityPoolBalance method should be present in the ABI")

	/*
		Gas tests
	*/

	// Check all methods use the correct gas amount.
	var method [4]byte

	gasBech32ToHex := contract.RequiredGas(bech32ToHex.ID)
	copy(method[:], bech32ToHex.ID[:4])
	baseCost := uint64(len(method)) * gasConfig.WriteCostPerByte
	require.Equal(
		t,
		GasRequiredByMethod[method]+baseCost,
		gasBech32ToHex,
		"bech32ToHexAddr method should require %d gas, got %d",
		GasRequiredByMethod[method]+baseCost,
		gasBech32ToHex)

	gasBech32ify := contract.RequiredGas(bech32ify.ID)
	copy(method[:], bech32ify.ID[:4])
	baseCost = uint64(len(method)) * gasConfig.WriteCostPerByte
	require.Equal(
		t,
		GasRequiredByMethod[method]+baseCost,
		gasBech32ify,
		"bech32ify method should require %d gas, got %d",
		GasRequiredByMethod[method]+baseCost,
		gasBech32ify)

	gasGetGasStabilityPoolBalance := contract.RequiredGas(getGasStabilityPoolBalance.ID)
	copy(method[:], getGasStabilityPoolBalance.ID[:4])
	baseCost = uint64(len(method)) * gasConfig.WriteCostPerByte
	require.Equal(
		t,
		GasRequiredByMethod[method]+baseCost,
		gasGetGasStabilityPoolBalance,
		"getGasStabilityPoolBalance method should require %d gas, got %d",
		GasRequiredByMethod[method]+baseCost,
		gasGetGasStabilityPoolBalance)

	// Can not happen, but check the gas for an invalid method.
	// At runtime if the method does not exist in the ABI, it returns an error.
	invalidMethodBytes := []byte("invalidMethod")
	gasInvalidMethod := contract.RequiredGas(invalidMethodBytes)
	baseCost = uint64(len(method)) * gasConfig.WriteCostPerByte
	require.Equal(
		t,
		uint64(0),
		gasInvalidMethod,
		"invalid method should require %d gas, got %d",
		uint64(0),
		gasInvalidMethod)
}

func Test_Bech32ToHexAddress(t *testing.T) {
	/*
		Configuration
	*/

	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, _, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	// Create contract and get ABI.
	contract := NewIPrototypeContract(k, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	// Test Bech32HexAddr method. Should succeed.
	methodID := abi.Methods[Bech32ToHexAddrMethodName]
	args := make([]interface{}, 0)
	args = append(args, "zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u")

	rawBytes, err := contract.Bech32ToHexAddr(&methodID, args)
	require.NoError(t, err, "Bech32ToHexAddr should not return an error")

	// Discard the first 12 bytes, the address is the last 20 bytes.
	addr := common.BytesToAddress(rawBytes[12:])
	require.Equal(
		t,
		common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE"),
		addr,
		"Bech32ToHexAddr should return the correct address, got: %v",
		addr,
	)

	// Test Bech32HexAddr method. Should fail with invalid number of arguments.
	args = append(args, "second argument")
	_, err = contract.Bech32ToHexAddr(&methodID, args)
	require.Error(t, err, "expected invalid number of arguments; expected 1; got: 2")
	require.IsType(t, &types.ErrInvalidNumberOfArgs{}, err, "expected error type: ErrInvalidNumberOfArgs, got: %T", err)

	// Test Bech32HexAddr method. Should fail with invalid address.
	argsInvalid := make([]interface{}, 0)
	argsInvalid = append(argsInvalid, "")
	_, errInvalid := contract.Bech32ToHexAddr(&methodID, argsInvalid)
	require.Error(t, errInvalid, "expected error invalid bech32 address: %v", argsInvalid[0])
}

func Test_Bech32ify(t *testing.T) {
	/*
		Configuration
	*/

	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, _, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	// Create contract and get ABI.
	contract := NewIPrototypeContract(k, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	// Test Bech32ify method.
	methodID := abi.Methods[Bech32ifyMethodName]
	args := make([]interface{}, 0)
	args = append(args, "zeta")
	args = append(args, common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE"))

	rawBytes, err := contract.Bech32ify(&methodID, args)
	require.NoError(t, err, "Bech32ify should not return an error")

	// Manually extract the address from the raw bytes.
	zetaAddr := string(rawBytes[64:107])
	require.Equal(
		t,
		"zeta1h8duy2dltz9xz0qqhm5wvcnj02upy887fyn43u",
		string(zetaAddr),
		"Bech32ify should return the correct address, got: %v",
		zetaAddr,
	)

	// Test for invalid number of arguments.
	args = append(args, "third argument")
	_, err = contract.Bech32ify(&methodID, args)
	require.Error(t, err, "expected invalid number of arguments; expected 2; got: 3")
	require.IsType(t, &types.ErrInvalidNumberOfArgs{}, err, "expected error type: ErrInvalidNumberOfArgs, got: %T", err)

	// Test for invalid bech32 human readable prefix.
	argsInvalidBech32 := make([]interface{}, 0)
	argsInvalidBech32 = append(argsInvalidBech32, 1337)
	argsInvalidBech32 = append(argsInvalidBech32, common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE"))
	_, errInvalidBech32 := contract.Bech32ify(&methodID, argsInvalidBech32)
	require.Error(t, errInvalidBech32, "expected error invalid bech32 human readable prefix (HRP)")

	// Test for invalid hex address.
	argsInvalidHexAddress := make([]interface{}, 0)
	argsInvalidHexAddress = append(argsInvalidHexAddress, "zeta")
	argsInvalidHexAddress = append(argsInvalidHexAddress, 1337)
	_, errInvalidHexAddress := contract.Bech32ify(&methodID, argsInvalidHexAddress)
	require.Error(t, errInvalidHexAddress, "expected error invalid hex address")

	// Test for invalid bech32 human readable prefix.
	argsInvalidEmptyPrefix := make([]interface{}, 0)
	argsInvalidEmptyPrefix = append(argsInvalidEmptyPrefix, "")
	argsInvalidEmptyPrefix = append(argsInvalidEmptyPrefix, common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE"))
	_, errInvalidEmptyPrefix := contract.Bech32ify(&methodID, argsInvalidEmptyPrefix)
	require.Error(t, errInvalidEmptyPrefix, "expected error invalid bech32 human readable prefix (HRP). Please provide a either an account, validator or consensus address prefix (eg: cosmos, cosmosvaloper, cosmosvalcons)")
}

func Test_GetGasStabilityPoolBalance(t *testing.T) {
	/*
		Configuration
	*/

	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, ctx, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	// Create contract and get ABI.
	contract := NewIPrototypeContract(k, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	// Test GetGasStabilityPoolBalance method.
	// Only check the function is called correctly inside the contract, and it returns the expected error.
	// Configuring a local environment for this contract would require deploying system contracts and gas pools.
	// This method is tested thoroughly in the e2e tests.
	methodID := abi.Methods[GetGasStabilityPoolBalanceName]
	args := make([]interface{}, 0)
	args = append(args, int64(1337))

	_, err := contract.GetGasStabilityPoolBalance(ctx, &methodID, args)
	require.Error(
		t,
		err,
		"error calling fungible keeper: failed to get system contract variable: state variable not found",
	)

	// Test for invalid number of arguments.
	args = append(args, "second argument")
	_, err = contract.GetGasStabilityPoolBalance(ctx, &methodID, args)
	require.Error(t, err, "expected invalid number of arguments; expected 2; got: 3")
	require.IsType(t, &types.ErrInvalidNumberOfArgs{}, err, "expected error type: ErrInvalidNumberOfArgs, got: %T", err)

	// Test for invalid chainID.
	argsInvalid := make([]interface{}, 0)
	argsInvalid = append(argsInvalid, "foobar")
	_, errInvalid := contract.GetGasStabilityPoolBalance(ctx, &methodID, argsInvalid)
	require.Error(t, errInvalid, "expected int64, got: %T", argsInvalid[0])
	require.IsType(t, types.ErrInvalidArgument{}, errInvalid, "expected error type: ErrInvalidArgument, got: %T", errInvalid)
}

func Test_InvalidMethod(t *testing.T) {
	/*
		Configuration
	*/

	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, _, _, _ := keeper.FungibleKeeper(t)
	gasConfig := storetypes.TransientGasConfig()

	// Create contract and get ABI.
	contract := NewIPrototypeContract(k, appCodec, gasConfig)
	require.NotNil(t, contract, "NewIPrototypeContract() should not return a nil contract")

	abi := contract.Abi()
	require.NotNil(t, abi, "contract ABI should not be nil")

	// Test for non existent method.
	_, doNotExist := abi.Methods["invalidMethod"]
	require.False(t, doNotExist, "invalidMethod should not be present in the ABI")
}

func Test_MissingABI(t *testing.T) {
	prototypeABI = ""
	defer func() {
		if r := recover(); r != nil {
			require.Equal(t, "missing prototype ABI", r, "expected error: missing ABI, got: %v", r)
		}
	}()

	initABI()
}

func Test_InvalidABI(t *testing.T) {
	prototypeABI = "invalid json"
	defer func() {
		if r := recover(); r != nil {
			require.IsType(t, &json.SyntaxError{}, r, "expected error type: json.SyntaxError, got: %T", r)
		}
	}()

	initABI()
}