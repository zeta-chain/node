package types

import (
	"fmt"

	"github.com/pkg/errors"
)

func (m InboundParams) Validate() error {
	if m.Sender == "" {
		return fmt.Errorf("sender cannot be empty")
	}

	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}

	if err := ValidateAddressForChain(m.Sender, m.SenderChainId); err != nil {
		return err
	}
	if m.TxOrigin != "" {
		if err := ValidateAddressForChain(m.TxOrigin, m.SenderChainId); err != nil {
			return err
		}
	}
	if err := ValidateHashForChain(m.ObservedHash, m.SenderChainId); err != nil {
		return errors.Wrap(err, "invalid inbound tx observed hash")
	}
	if m.BallotIndex != "" {
		if err := ValidateCCTXIndex(m.BallotIndex); err != nil {
			return errors.Wrap(err, "invalid inbound tx ballot index")
		}
	}
	return nil
}
