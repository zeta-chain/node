package sui

import (
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Gateway error codes
// see: https://github.com/zeta-chain/protocol-contracts-sui/blob/e5a756e473da884dcbc59b574b387a7a365ac823/sources/gateway.move#L14-L21
const (
	ErrCodeAlreadyWhitelisted     uint64 = 0
	ErrCodeInvalidReceiverAddress uint64 = 1
	ErrCodeNotWhitelisted         uint64 = 2
	ErrCodeNonceMismatch          uint64 = 3
	ErrCodePayloadTooLong         uint64 = 4
	ErrCodeInactiveWithdrawCap    uint64 = 5
	ErrCodeInactiveWhitelistCap   uint64 = 6
	ErrCodeDepositPaused          uint64 = 7
)

var (
	// ErrParseEvent event parse error
	ErrParseEvent = errors.New("event parse error")

	// ErrObjectOwnership is the error returned when a wrong object ownership is used in withdraw_and_call
	ErrObjectOwnership = errors.New("wrong object ownership")

	// ErrInvalidPayload is the error returned when a invalid payload format is used in withdraw_and_call
	ErrInvalidPayload = errors.New("invalid payload")

	// retryableOutboundErrCodes are the outbound execution (if failed) error codes that are retryable.
	// The list is used to determine if a withdraw_and_call should fallback if rejected by the network.
	// Note: keep this list in sync with the actual implementation in `gateway.move`
	retryableOutboundErrCodes = []uint64{
		ErrCodeNotWhitelisted,
		ErrCodeNonceMismatch,
		ErrCodeInactiveWithdrawCap,
	}

	// moveAbortRegex is the MoveAbort error regex pattern: "MoveAbort(..., <code>) ..."
	moveAbortRegex = regexp.MustCompile(`MoveAbort\(.+?,\s*(\d+)\)`)

	// commandIndexRegex extracts the command index from any execution error message.
	// For example: "MoveAbort(...) in command 2"
	commandIndexRegex = regexp.MustCompile(`in\s+command\s+(\d+)$`)
)

// MoveAbort represents a MoveAbort execution error.
// see: https://github.com/MystenLabs/sui-rust-sdk/blob/65eb9f3ad63b98f5b04465963d340e53b301a149/crates/sui-sdk-types/src/execution_status.rs#L173
type MoveAbort struct {
	Message string
	Code    uint64
}

// NewMoveAbortFromExecutionError creates a MoveAbort struct from Sui 'ExecutionError::MoveAbort' execution error message.
// Example: "MoveAbort(MoveLocation { module: ModuleId { address: a5f027339b7e04e5d55c2ac90ea71d616870aa21d9f16fd0237a2a42e67c9f3e, name: Identifier("gateway") }, function: 11, instruction: 37, function_name: Some("withdraw_impl") }, 3) in command 0"
func NewMoveAbortFromExecutionError(errorMsg string) (abort MoveAbort, err error) {
	matches := moveAbortRegex.FindStringSubmatch(errorMsg)
	if len(matches) != 2 {
		return abort, errors.Errorf("unable to extract code from error string: %s", errorMsg)
	}

	codeStr := matches[1]
	code, err := strconv.ParseUint(codeStr, 10, 64)
	if err != nil {
		return abort, errors.Wrapf(err, "unable to convert code to uint64: %s", codeStr)
	}
	return MoveAbort{
		Message: errorMsg,
		Code:    code,
	}, nil
}

// IsRetryable returns true if the MoveAbort error code is in the retryable error list.
func (m MoveAbort) IsRetryable() bool {
	return slices.Contains(retryableOutboundErrCodes, m.Code)
}

// IsRetryableExecutionError checks if the error message is a retryable error.
// Sui withdraw and withdrawAndCall may fail with unknown execution errors, zetaclient
// retries the outbound on known errors and cancels the outbound on any unknown errors.
func IsRetryableExecutionError(errorMsg string) (bool, error) {
	// cmdIndex is optional in Sui error message
	cmdIndex, err := extractCommandIndex(errorMsg)

	switch {
	case err != nil:
		// command index not found, it's fine
		// cancel this tx with unknown command index
		return false, nil
	case cmdIndex == 0:
		// 'withdraw_impl' error
		// The gateway 'withdraw_impl' may fail with three types of known MoveAbort errors.
		// If it does, the scheduler should retry the outbound until it succeeds:
		// 	- MoveAbort (ErrCodeNotWhitelisted)
		// 	- MoveAbort (ErrCodeNonceMismatch)
		// 	- MoveAbort (ErrCodeInactiveWithdrawCap)
		if strings.HasPrefix(errorMsg, "MoveAbort") {
			moveAbort, err := NewMoveAbortFromExecutionError(errorMsg)
			if err != nil {
				return false, errors.Wrap(err, "unable to create MoveAbort from execution error")
			}
			return moveAbort.IsRetryable(), nil
		}
		return false, nil
	case slices.Contains([]uint16{1, 2, 3, 4}, cmdIndex):
		// cancel tx if any one of the remaining commands failed
		// command 1: gas budget transfer error
		// command 2: 'set_message_context' error
		// command 3: 'on_call' error
		// command 4: 'reset_message_context' error
		return false, nil
	default: // never happen
		return false, errors.Errorf("invalid command index: %d", cmdIndex)
	}
}

// extractCommandIndex extracts the command index from an execution error message if present.
// The command index is optional in the Sui execution error message.
// see: https://github.com/MystenLabs/sui/blob/8a8b5e54c59762f2da57c8ff1e76d571e7015492/crates/sui-types/src/execution_status.rs#L25
func extractCommandIndex(errorMsg string) (uint16, error) {
	matches := commandIndexRegex.FindStringSubmatch(errorMsg)
	if len(matches) != 2 {
		return 0, errors.New("no command index found")
	}

	cmdIndex, err := strconv.ParseUint(matches[1], 10, 16)
	if err != nil {
		return 0, errors.Wrap(err, "unable to convert command index to uint16")
	}

	// #nosec G103 always in range
	return uint16(cmdIndex), nil
}
