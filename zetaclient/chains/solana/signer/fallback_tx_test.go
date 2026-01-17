package signer

import (
	"errors"
	"testing"

	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"github.com/stretchr/testify/require"
)

func makeTestRpcError(message string, logs []string) *jsonrpc.RPCError {
	var rawLogs []any
	for _, l := range logs {
		rawLogs = append(rawLogs, l)
	}

	return &jsonrpc.RPCError{
		Code:    -32002,
		Message: message,
		Data: map[string]any{
			"logs": rawLogs,
		},
	}
}

func Test_ParseRPCErrorForFallback(t *testing.T) {
	gateway := "94U5AHQMKkV5txNJ17QPXWoh474PheGou6cNP2FEuL1d"

	t.Run("only gateway invoked, NonceMismatch", func(t *testing.T) {
		errorMsg := "Program " + gateway + " failed: custom program error: 0x1771"
		logs := []string{
			"Program ComputeBudget111111111111111111111111111111 invoke [1]",
			"Program ComputeBudget111111111111111111111111111111 success",
			"Program " + gateway + " invoke [1]",
			"Program log: Instruction: Execute",
			"Program log: AnchorError occurred. Error Code: NonceMismatch. Error Message: NonceMismatch.",
			errorMsg,
		}
		err := makeTestRpcError(
			"Transaction simulation failed: Error processing Instruction 0: custom program error: 0x1771",
			logs,
		)

		shouldUseFallbackTx, failureReason := parseRPCErrorForFallback(err, gateway)

		require.False(t, shouldUseFallbackTx)
		require.Empty(t, failureReason)
	})

	t.Run("another program invoked after gateway", func(t *testing.T) {
		errorMsg := "Program 4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc failed: custom program error: 0x1771"
		logs := []string{
			"Program " + gateway + " invoke [1]",
			"Program 4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc invoke [2]",
			"Program 4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc failed: custom program error: 0x1771",
		}
		err := makeTestRpcError(
			"Transaction simulation failed: Error processing Instruction 0: custom program error: 0x1771",
			logs,
		)

		shouldUseFallbackTx, failureReason := parseRPCErrorForFallback(err, gateway)

		require.True(t, shouldUseFallbackTx)
		require.Equal(t, errorMsg, failureReason)
	})

	t.Run("only gateway invoked but not NonceMismatch", func(t *testing.T) {
		errorMsg := "Program " + gateway + " failed: custom program error: 0x1771"
		logs := []string{
			"Program " + gateway + " invoke [1]",
			"Program log: Instruction: Execute",
			"Program log: AnchorError occurred. Error Code: SomeOtherError",
			"Program " + gateway + " failed: custom program error: 0x1771",
		}
		err := makeTestRpcError(
			"Transaction simulation failed: Error processing Instruction 0: custom program error: 0x1771",
			logs,
		)

		shouldUseFallbackTx, failureReason := parseRPCErrorForFallback(err, gateway)

		require.True(t, shouldUseFallbackTx)
		require.Equal(t, errorMsg, failureReason)
	})

	t.Run("gateway invoked from connected program, reentrancy", func(t *testing.T) {
		errorMsg := "Program " + gateway + " failed: custom program error: 0x1771"
		logs := []string{
			"Program " + gateway + " invoke [1]",
			"Program 4xEw862A2SEwMjofPkUyd4NEekmVJKJsdHkK3UkAtDrc invoke [2]",
			"Program " + gateway + " invoke [3]",
			"Program " + gateway + " failed: custom program error: 0x1771",
		}
		err := makeTestRpcError(
			"Transaction simulation failed: Error processing Instruction 0: custom program error: 0x1771",
			logs,
		)

		shouldUseFallbackTx, failureReason := parseRPCErrorForFallback(err, gateway)

		require.True(t, shouldUseFallbackTx)
		require.Equal(t, errorMsg, failureReason)
	})

	t.Run("invalid error type", func(t *testing.T) {
		err := errors.New("some generic error")
		shouldUseFallbackTx, failureReason := parseRPCErrorForFallback(err, gateway)

		require.False(t, shouldUseFallbackTx)
		require.Empty(t, failureReason)
	})

	t.Run("use fallback, transaction size error", func(t *testing.T) {
		err := &jsonrpc.RPCError{
			Code:    errorCodeJSONRPCInvalidParams,
			Message: "base64 encoded solana_transaction::versioned::VersionedTransaction too large: 2012 bytes (max: encoded/raw 1644/1232)",
			Data:    nil,
		}

		shouldUseFallbackTx, failureReason := parseRPCErrorForFallback(err, gateway)

		require.True(t, shouldUseFallbackTx)
		require.Equal(t, transactionSizeError, failureReason)
	})

	t.Run("don't use fallback, invalid params error but not transaction size", func(t *testing.T) {
		err := &jsonrpc.RPCError{
			Code:    errorCodeJSONRPCInvalidParams,
			Message: "some other invalid parameter error",
			Data:    nil,
		}

		shouldUseFallbackTx, failureReason := parseRPCErrorForFallback(err, gateway)

		require.False(t, shouldUseFallbackTx)
		require.Empty(t, failureReason)
	})
}
