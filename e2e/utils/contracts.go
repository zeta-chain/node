package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/testutil/contracts/example"
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
		amount.Uint64(),
		bar.Uint64(),
	)
}

// MustHaveCalledExampleContractWithMsg checks if the contract has been called correctly with correct amount and msg
func MustHaveCalledExampleContractWithMsg(
	t require.TestingT,
	contract *testcontract.Example,
	amount *big.Int,
	msg []byte,
) {
	bar, err := contract.Bar(&bind.CallOpts{})
	require.NoError(t, err)
	require.Equal(
		t,
		amount.Uint64(),
		bar.Uint64(),
	)

	lastMsg, err := contract.LastMessage(&bind.CallOpts{})
	require.NoError(t, err)
	require.Equal(t, string(msg), string(lastMsg))
}
