package signer

import (
	"strings"

	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
)

// parseRPCErrorForFallback parse error as RPCError and verifies if fallback tx should be used
// and which failure reason to attach to fallback tx
func parseRPCErrorForFallback(err error, program string) (bool, string) {
	rpcErr, ok := err.(*jsonrpc.RPCError)
	if !ok {
		return false, ""
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
	failureReason := getFailureReason(logs)
	return true, failureReason
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
