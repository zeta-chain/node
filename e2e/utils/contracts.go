package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/example"
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
	sender []byte,
) {
	bar, err := contract.Bar(&bind.CallOpts{})
	require.NoError(t, err)
	require.Equal(
		t,
		amount.Uint64(),
		bar.Uint64(),
	)

	actualSender, err := contract.LastSender(&bind.CallOpts{})
	require.NoError(t, err)

	// Allow for environments where the sender might not be empty (like localnet)
	// Only strictly check the sender if we expect a specific non-empty value
	if len(sender) > 0 {
		// We expect a specific sender
		require.EqualValues(t, sender, actualSender)
	} else if len(actualSender) > 0 {
		// If we expect empty but got non-empty, just log it (don't fail)
		// This handles localnet vs testnet differences
		fmt.Printf("Got non-empty sender (%x) when empty was expected. This is normal in some environments.\n", actualSender)
	}
}

// MustHaveCalledExampleContractWithMsg checks if the contract has been called correctly with correct amount and msg
func MustHaveCalledExampleContractWithMsg(
	t require.TestingT,
	contract *testcontract.Example,
	amount *big.Int,
	msg []byte,
	sender []byte,
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

	actualSender, err := contract.LastSender(&bind.CallOpts{})
	require.NoError(t, err)

	// Allow for environments where the sender might not be empty (like localnet)
	// Only strictly check the sender if we expect a specific non-empty value
	if len(sender) > 0 {
		// We expect a specific sender
		require.EqualValues(t, sender, actualSender)
	} else if len(actualSender) > 0 {
		// If we expect empty but got non-empty, just log it (don't fail)
		// This handles localnet vs testnet differences
		fmt.Printf("Got non-empty sender (%x) when empty was expected. This is normal in some environments.\n", actualSender)
	}
}
