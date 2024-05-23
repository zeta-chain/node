package types

import (
	"fmt"
	"strconv"

	"github.com/zeta-chain/zetacore/pkg/chains"
)

func (m OutboundTxParams) GetGasPrice() (uint64, error) {
	gasPrice, err := strconv.ParseUint(m.OutboundTxGasPrice, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse cctx gas price %s: %s", m.OutboundTxGasPrice, err.Error())
	}

	return gasPrice, nil
}

func (m OutboundTxParams) Validate() error {
	if m.Receiver == "" {
		return fmt.Errorf("receiver cannot be empty")
	}
	if chains.GetChainFromChainID(m.ReceiverChainId) == nil {
		return fmt.Errorf("invalid receiver chain id %d", m.ReceiverChainId)
	}

	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}

	// Disabled checks
	// TODO: Improve the checks, move the validation call to a new place and reenable
	// https://github.com/zeta-chain/node/issues/2234
	// https://github.com/zeta-chain/node/issues/2235
	//if err := ValidateAddressForChain(m.Receiver, m.ReceiverChainId); err != nil {
	//	return err
	//}
	//if m.BallotIndex != "" {
	//
	//	if err := ValidateCCTXIndex(m.BallotIndex); err != nil {
	//		return errors.Wrap(err, "invalid outbound tx ballot index")
	//	}
	//}
	//if m.Hash != "" {
	//	if err := ValidateHashForChain(m.Hash, m.ReceiverChainId); err != nil {
	//		return errors.Wrap(err, "invalid outbound tx hash")
	//	}
	//}

	return nil
}
