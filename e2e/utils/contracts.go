package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/testutil/contracts"
)

const (
	// ErrHashRevertFoo is the keccak256 hash of custom error "Foo()" on reverter contract
	ErrHashRevertFoo = "0xbfb4ebcf"
)

// MustHaveCalledExampleContract checks if the contract has been called correctly
func MustHaveCalledExampleContract(
	t require.TestingT,
	contract *testcontract.Example,
	amount *big.Int,
) {
	bar, err := contract.Bar(&bind.CallOpts{})
	require.NoError(t, err)
	require.Equal(
		t,
		0,
		bar.Cmp(amount),
		"cross-chain call failed bar value %s should be equal to amount %s",
		bar.String(),
		amount.String(),
	)
}
