package runner

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testdappv2"
)

const (
	// The nodes of external chain take to to sync the state, for example:
	// The BSC chain outbound takes 15 seconds to confirm and still not enough for the node to sync,
	// so we need proper tolerance to ensure the nodes behind RPC are in sync
	externalNodeSyncTolerance = 30 * time.Second
)

// AssertTestDAppZEVMCalled is a function that asserts the values of the test dapp on the ZEVM
// this function uses TestDAppV2 for the assertions, in the future we should only use this contracts for all tests
// https://github.com/zeta-chain/node/issues/2655
func (r *E2ERunner) AssertTestDAppZEVMCalled(expectedCalled bool, message string, sender []byte, amount *big.Int) {
	r.waitAndVerifyTestDAppCall(r.TestDAppV2ZEVM, expectedCalled, message, sender, amount)
}

// AssertTestDAppEVMCalled is a function that asserts the values of the test dapp on the external EVM
func (r *E2ERunner) AssertTestDAppEVMCalled(expectedCalled bool, message string, amount *big.Int) {
	r.waitAndVerifyTestDAppCall(r.TestDAppV2EVM, expectedCalled, message, nil, amount)
}

// waitAndVerifyTestDAppCall waits for the test dapp to be called with the expected message, sender and amount
// This function is to tolerate the fact that the dApp state may not be synced across all nodes behind a RPC.
func (r *E2ERunner) waitAndVerifyTestDAppCall(
	testDApp *testdappv2.TestDAppV2,
	expectedCalled bool,
	message string,
	expectedSender []byte,
	expectedAmount *big.Int,
) {
	// do simply assert if the dApp is NOT expected to be called, no need to wait
	if !expectedCalled {
		called, err := testDApp.GetCalledWithMessage(&bind.CallOpts{}, message)
		require.NoError(r, err)
		require.EqualValues(r, false, called)
		return
	}

	// wait until the dApp gets called with the expected message and amount or timeout
	startTime := time.Now()
	checkInterval := 2 * time.Second

	for {
		time.Sleep(checkInterval)
		require.False(r, time.Since(startTime) > externalNodeSyncTolerance, "timeout waiting for dApp state to update")

		called, err := testDApp.GetCalledWithMessage(&bind.CallOpts{}, message)
		require.NoError(r, err)

		sender, err := testDApp.GetSenderWithMessage(&bind.CallOpts{}, message)
		require.NoError(r, err)

		amount, err := testDApp.GetAmountWithMessage(&bind.CallOpts{}, message)
		require.NoError(r, err)

		// if sender is provided, check if it matches the actual sender, otherwise skip the check
		senderMatched := len(expectedSender) == 0 || bytes.Equal(sender, expectedSender)

		// stop only if sender, message and amount are matching the expected values
		if called == expectedCalled && amount.Cmp(expectedAmount) == 0 && senderMatched {
			return
		}
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

// EncodeGasCallRevert encodes the payload for the gasCall function that reverts
func (r *E2ERunner) EncodeGasCallRevert() []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("gasCall", "revert")
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

// EncodeERC20CallRevert encodes the payload for the erc20Call function that reverts
func (r *E2ERunner) EncodeERC20CallRevert(erc20Addr ethcommon.Address, amount *big.Int) []byte {
	abi, err := testdappv2.TestDAppV2MetaData.GetAbi()
	require.NoError(r, err)

	// encode the message
	encoded, err := abi.Pack("erc20Call", erc20Addr, amount, "revert")
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
