package types

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

func (m OutboundParams) GetGasPriceUInt64() (uint64, error) {
	gasPrice, err := strconv.ParseUint(m.GasPrice, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse cctx gas price %s: %s", m.GasPrice, err.Error())
	}

	return gasPrice, nil
}

func (m OutboundParams) GetGasPriorityFeeUInt64() (uint64, error) {
	// noop
	if m.GasPriorityFee == "" {
		return 0, nil
	}

	fee, err := strconv.ParseUint(m.GasPriorityFee, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to parse priority fee from %q", m.GasPriorityFee)
	}

	return fee, nil
}

func (m OutboundParams) Validate() error {
	if m.Receiver == "" {
		return fmt.Errorf("receiver cannot be empty")
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
