package utils

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/example"
)

const (
	// ErrHashRevertFoo is the keccak256 hash of custom error "Foo()" on reverter contract
	ErrHashRevertFoo = "0xbfb4ebcf"
)

// WaitAndVerifyExampleContractCall waits for the example contract to be called with the expected amount and sender
// This function is to tolerate the fact that the contract state may not be synced across all nodes behind a RPC.
func WaitAndVerifyExampleContractCall(
	t require.TestingT,
	contract *testcontract.Example,
	amount *big.Int,
	sender []byte,
) {
	// wait until the contract gets called with the expected amount and sender or timeout
	startTime := time.Now()
	checkInterval := 2 * time.Second

	for {
		time.Sleep(checkInterval)
		require.False(t, time.Since(startTime) > nodeSyncTolerance, "timeout waiting for contract state to update")

		bar, err := contract.Bar(&bind.CallOpts{})
		require.NoError(t, err)

		actualSender, err := contract.LastSender(&bind.CallOpts{})
		require.NoError(t, err)

		// stop only if both amount and sender are matching the expected values
		if bar.Cmp(amount) == 0 && bytes.Equal(actualSender, sender) {
			return
		}
	}
}

// WaitAndVerifyExampleContractCallWithMsg waits for the example contract to be called with the expected amount, msg and sender
// This function is to tolerate the fact that the contract state may not be synced across all nodes behind a RPC.
func WaitAndVerifyExampleContractCallWithMsg(
	t require.TestingT,
	contract *testcontract.Example,
	amount *big.Int,
	msg []byte,
	sender []byte,
) {
	// wait until the contract gets called with the expected amount, msg and sender or timeout
	startTime := time.Now()
	checkInterval := 2 * time.Second

	for {
		time.Sleep(checkInterval)
		require.False(t, time.Since(startTime) > nodeSyncTolerance, "timeout waiting for contract state to update")

		bar, err := contract.Bar(&bind.CallOpts{})
		require.NoError(t, err)

		lastMsg, err := contract.LastMessage(&bind.CallOpts{})
		require.NoError(t, err)

		actualSender, err := contract.LastSender(&bind.CallOpts{})
		require.NoError(t, err)

		// stop only if amount, msg and sender are all matching the expected values
		if bar.Cmp(amount) == 0 &&
			bytes.Equal(lastMsg, msg) &&
			bytes.Equal(actualSender, sender) {
			return
		}
	}
}
