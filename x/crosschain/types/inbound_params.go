package types

import (
	"fmt"

	"github.com/zeta-chain/zetacore/pkg/chains"
)

func (m InboundParams) Validate() error {
	if m.Sender == "" {
		return fmt.Errorf("sender cannot be empty")
	}

	if chains.GetChainFromChainID(m.SenderChainId) == nil {
		return fmt.Errorf("invalid sender chain id %d", m.SenderChainId)
	}

	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}

	// Disabled checks
	// TODO: Improve the checks, move the validation call to a new place and reenable
	// https://github.com/zeta-chain/node/issues/2234
	// https://github.com/zeta-chain/node/issues/2235
	//if err := ValidateAddressForChain(m.Sender, m.SenderChainId) err != nil {
	//	return err
	//}
	//if m.TxOrigin != "" {
	//	errTxOrigin := ValidateAddressForChain(m.TxOrigin, m.SenderChainId)
	//	if errTxOrigin != nil {
	//		return errTxOrigin
	//	}
	//}
	//if err := ValidateHashForChain(m.ObservedHash, m.SenderChainId); err != nil {
	//	return errors.Wrap(err, "invalid inbound tx observed hash")
	//}
	//if m.BallotIndex != "" {
	//	if err := ValidateCCTXIndex(m.BallotIndex); err != nil {
	//		return errors.Wrap(err, "invalid inbound tx ballot index")
	//	}
	//}
	return nil
}
