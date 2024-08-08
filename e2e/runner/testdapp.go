package runner

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
)

// AssertTestDAppValues is a function that asserts the values of the test dapp
// this function uses TestDAppV2 for the assertions, in the future we should only use this contracts for all tests
// https://github.com/zeta-chain/node/issues/2655
func (r *E2ERunner) AssertTestDAppValues(equals bool, message string, amount *big.Int) {
	// check the payload was received on the contract
	actualMessage, err := r.TestDAppV2ZEVM.LastMessage(&bind.CallOpts{})
	require.NoError(r, err)

	// check the amount was received on the contract
	actualAmount, err := r.TestDAppV2ZEVM.LastAmount(&bind.CallOpts{})
	require.NoError(r, err)

	if equals {
		require.Equal(r, message, actualMessage)
		require.Equal(r, amount.Uint64(), actualAmount.Uint64())
	} else {
		require.NotEqual(r, message, actualMessage)
		require.NotEqual(r, amount.Uint64(), actualAmount.Uint64())
	}
}
