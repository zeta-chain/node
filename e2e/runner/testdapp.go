package runner

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/pkg/contracts/testdappv2"
)

// AssertTestDAppZEVMValues is a function that asserts the values of the test dapp on the ZEVM
// this function uses TestDAppV2 for the assertions, in the future we should only use this contracts for all tests
// https://github.com/zeta-chain/node/issues/2655
func (r *E2ERunner) AssertTestDAppZEVMValues(equals bool, message string, amount *big.Int) {
	r.assertTestDAppValues(r.TestDAppV2ZEVM, equals, message, amount)
}

// AssertTestDAppEVMValues is a function that asserts the values of the test dapp on the external EVM
func (r *E2ERunner) AssertTestDAppEVMValues(equals bool, message string, amount *big.Int) {
	r.assertTestDAppValues(r.TestDAppV2EVM, equals, message, amount)
}

func (r *E2ERunner) assertTestDAppValues(
	testDApp *testdappv2.TestDAppV2,
	equals bool,
	message string,
	amount *big.Int,
) {
	// check the payload was received on the contract
	actualMessage, err := testDApp.LastMessage(&bind.CallOpts{})
	require.NoError(r, err)

	// check the amount was received on the contract
	actualAmount, err := testDApp.LastAmount(&bind.CallOpts{})
	require.NoError(r, err)

	if equals {
		require.Equal(r, message, actualMessage)
		require.Equal(r, amount.Uint64(), actualAmount.Uint64())
	} else {
		require.NotEqual(r, message, actualMessage)
		require.NotEqual(r, amount.Uint64(), actualAmount.Uint64())
	}
}

// EncodeGasCall encodes the payload for the gasCall function
func (r *E2ERunner) EncodeGasCall(message string) []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("gasCall", message)
	require.NoError(r, err)
	return encoded
}

// EncodeERC20Call encodes the payload for the erc20Call function
func (r *E2ERunner) EncodeERC20Call(erc20Addr ethcommon.Address, amount *big.Int, message string) []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("erc20Call", erc20Addr, amount, message)
	require.NoError(r, err)
	return encoded
}

// EncodeSimpleCall encodes the payload for the simpleCall function
func (r *E2ERunner) EncodeSimpleCall(message string) []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("simpleCall", message)
	require.NoError(r, err)
	return encoded
}
