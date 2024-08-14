// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)

package types

const (
	// ErrDifferentOrigin is raised when an approval is set but the origin address is not the same as the spender.
	ErrDifferentOrigin = "tx origin address %s does not match the delegator address %s"
	// ErrInvalidAddr is raised when the origin address is invalid
	ErrInvalidAddr = "invalid addr: %s"
	// ErrInvalidDelegator is raised when the delegator address is not valid.
	ErrInvalidDelegator = "invalid delegator address: %s"
	// ErrInvalidNumberOfArgs is raised when the number of arguments is not what is expected.
	ErrInvalidNumberOfArgs = "invalid number of arguments; expected %d; got: %d"
	// ErrInvalidMethod is raised when the method is invalid.
	ErrInvalidMethod = "unknown method: %s"
)