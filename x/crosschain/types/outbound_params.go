package types

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/pkg/chains"
)

func (m OutboundParams) GetGasPriceUInt64() (uint64, error) {
	gasPrice, err := strconv.ParseUint(m.GasPrice, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("unable to parse cctx gas price %s: %s", m.GasPrice, err.Error())
	}

	return gasPrice, nil
}

func (m OutboundParams) Validate() error {
	if m.Receiver == "" {
		return fmt.Errorf("receiver cannot be empty")
	}
	if chains.GetChainFromChainID(m.ReceiverChainId) == nil {
		return fmt.Errorf("invalid receiver chain id %d", m.ReceiverChainId)
	}
	err := ValidateAddressForChain(m.Receiver, m.ReceiverChainId)
	if err != nil {
		return err
	}
	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}
	if m.BallotIndex != "" {
		err = ValidateZetaIndex(m.BallotIndex)
		if err != nil {
			return errors.Wrap(err, "invalid outbound tx ballot index")
		}
	}
	if m.Hash != "" {
		err = ValidateHashForChain(m.Hash, m.ReceiverChainId)
		if err != nil {
			return errors.Wrap(err, "invalid outbound tx hash")
		}
	}
	return nil
}
