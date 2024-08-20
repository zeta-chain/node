package prototype

import (
	"testing"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	ethermint "github.com/zeta-chain/ethermint/types"
	"github.com/zeta-chain/zetacore/testutil/keeper"
)

func Test_IPrototypeContract(t *testing.T) {
	/*
		Configuration
	*/

	var encoding ethermint.EncodingConfig
	appCodec := encoding.Codec
	k, ctx, _, _ := keeper.FungibleKeeper(t)
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

	/*
		Methods tests
	*/

	// Test Bech32HexAddr method.
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

	// Test Bech32ify method.
	methodID = abi.Methods[Bech32ifyMethodName]
	args = make([]interface{}, 0)
	args = append(args, "zeta")
	args = append(args, common.HexToAddress("0xB9Dbc229Bf588A613C00BEE8e662727AB8121cfE"))

	rawBytes, err = contract.Bech32ify(&methodID, args)
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

	// Test GetGasStabilityPoolBalance method.
	// Only check the function is called correctly inside the contract, and it returns the expected error.
	// Configuring a local environment for this contract would require deploying system contracts and gas pools.
	// This method is tested thoroughly in the e2e tests.
	methodID = abi.Methods[GetGasStabilityPoolBalanceName]
	args = make([]interface{}, 0)
	args = append(args, int64(1337))

	rawBytes, err = contract.GetGasStabilityPoolBalance(ctx, &methodID, args)
	require.Error(
		t,
		err,
		"error calling fungible keeper: failed to get system contract variable: state variable not found",
	)
}
