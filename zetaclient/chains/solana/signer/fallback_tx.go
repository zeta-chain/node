package signer

import (
	"strings"

	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
)

const (
	transactionSizeError = "transaction is too large"

	// Solana uses JSON-RPC error code -32602 for invalid params
	// see: https://github.com/paritytech/jsonrpc/blob/dc9550b4b0d8bf409d025eba7e9b229b67af9401/core/src/types/error.rs#L32
	errorCodeJSONRPCInvalidParams = -32602
)

// parseRPCErrorForFallback parse error as RPCError and verifies if fallback tx should be used
// and which failure reason to attach to fallback tx
func parseRPCErrorForFallback(err error, program string) (useFallback bool, failureReason string) {
	rpcErr, ok := err.(*jsonrpc.RPCError)
	if !ok {
		return false, ""
	}

	// Special handling for transaction size errors
	if isTransactionSizeError(rpcErr) {
		return true, transactionSizeError
	}

	if !strings.Contains(rpcErr.Message, "Error processing Instruction") {
		return false, ""
	}

	dataMap, ok := rpcErr.Data.(map[string]interface{})
	if !ok {
		return false, ""
	}

	rawLogs, ok := dataMap["logs"].([]interface{})
	if !ok {
		return false, ""
	}

	logs := parseLogs(rawLogs)

	// if any other program invoked after gateway OR nonce mismatch not present, fallback
	shouldUseFallbackTx := programInvokedAfterTargetInLogs(logs, program) || !containsNonceMismatch(logs)
	if !shouldUseFallbackTx {
		return false, ""
	}

	// get failure reason from logs
	return true, getFailureReason(logs)
}

// getFailureReason returns first log that is in format "Program <P_ID> <error> failed"
func getFailureReason(logs []string) string {
	var failures []string
	for _, line := range logs {
		if strings.HasPrefix(line, "Program ") && strings.Contains(line, " failed") {
			failures = append(failures, line)
		}
	}

	if len(failures) == 0 {
		return ""
	}

	// returning first one is enough, since that is original program where program failed
	return failures[0]
}

// programInvokedAfterTargetInLogs checks if there is Program <P_ID> invoke after target program invoke log
func programInvokedAfterTargetInLogs(logs []string, targetProgram string) bool {
	foundTarget := false

	for _, line := range logs {
		if strings.HasPrefix(line, "Program ") && strings.Contains(line, " invoke") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				programID := fields[1]

				if programID == targetProgram {
					foundTarget = true
					continue
				}

				if foundTarget {
					// Found another program after target
					return true
				}
			}
		}
	}
	return false
}

// parseLogs parses raw logs to string
func parseLogs(rawLogs []interface{}) []string {
	var logs []string
	for _, l := range rawLogs {
		if str, ok := l.(string); ok {
			logs = append(logs, str)
		}
	}
	return logs
}

// containsNonceMismatch checks if some log contains NonceMismatch
func containsNonceMismatch(logs []string) bool {
	for _, log := range logs {
		if strings.Contains(log, "NonceMismatch") {
			return true
		}
	}
	return false
}

// isTransactionSizeError checks if the error matches transaction size error patterns
//
// the invalid transaction size error message will be like:
// "base64 encoded solana_transaction::versioned::VersionedTransaction too large: 2012 bytes (max: encoded/raw 1644/1232)"
//
// see:
//
//	https://github.com/solana-labs/solana/blob/bfacaf616fa4a1c57e2a337fcc864c92c25815a0/rpc/src/rpc.rs#L4506
//	https://github.com/solana-labs/solana/blob/bfacaf616fa4a1c57e2a337fcc864c92c25815a0/rpc/src/rpc.rs#L4521
func isTransactionSizeError(rpcErr *jsonrpc.RPCError) bool {
	if rpcErr.Code == errorCodeJSONRPCInvalidParams && strings.Contains(rpcErr.Message, "too large") &&
		strings.Contains(rpcErr.Message, "bytes (max: encoded/raw") {
		return true
	}

	return false
}
