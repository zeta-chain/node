package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

func (m InboundParams) Validate() error {
	if m.Sender == "" {
		return fmt.Errorf("sender cannot be empty")
	}
	if chains.GetChainFromChainID(m.SenderChainId) == nil {
		return fmt.Errorf("invalid sender chain id %d", m.SenderChainId)
	}
	err := ValidateAddressForChain(m.Sender, m.SenderChainId)
	if err != nil {
		return err
	}

	if m.TxOrigin != "" {
		errTxOrigin := ValidateAddressForChain(m.TxOrigin, m.SenderChainId)
		if errTxOrigin != nil {
			return errTxOrigin
		}
	}
	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}
	err = ValidateHashForChain(m.ObservedHash, m.SenderChainId)
	if err != nil {
		return errors.Wrap(err, "invalid inbound tx observed hash")
	}
	if m.BallotIndex != "" {
		err = ValidateZetaIndex(m.BallotIndex)
		if err != nil {
			return errors.Wrap(err, "invalid inbound tx ballot index")
		}
	}
	return nil
}
