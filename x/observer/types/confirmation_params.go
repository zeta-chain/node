package types

import (
	"errors"
)

// Validate checks that the ConfirmationParams is valid
func (cp ConfirmationParams) Validate() error {
	if cp.SafeInboundCount == 0 {
		return errors.New("SafeInboundCount must be greater than 0")
	}
	if cp.FastInboundCount > cp.SafeInboundCount {
		return errors.New("FastInboundCount must be less than or equal to SafeInboundCount")
	}
	if cp.SafeOutboundCount == 0 {
		return errors.New("SafeOutboundCount must be greater than 0")
	}
	if cp.FastOutboundCount > cp.SafeOutboundCount {
		return errors.New("FastOutboundCount must be less than or equal to SafeOutboundCount")
	}
	return nil
}
