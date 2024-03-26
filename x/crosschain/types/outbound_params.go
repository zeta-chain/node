package types

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/zeta-chain/zetacore/common"
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
	if common.GetChainFromChainID(m.ReceiverChainId) == nil {
		return fmt.Errorf("invalid receiver chain id %d", m.ReceiverChainId)
	}
	err := ValidateAddressForChain(m.Receiver, m.ReceiverChainId)
	if err != nil {
		return err
	}
	if m.Amount.IsNil() {
		return fmt.Errorf("amount cannot be nil")
	}
	if m.OutboundTxBallotIndex != "" {
		err = ValidateZetaIndex(m.OutboundTxBallotIndex)
		if err != nil {
			return errors.Wrap(err, "invalid outbound tx ballot index")
		}
	}
	if m.OutboundTxHash != "" {
		err = ValidateHashForChain(m.OutboundTxHash, m.ReceiverChainId)
		if err != nil {
			return errors.Wrap(err, "invalid outbound tx hash")
		}
	}
	return nil
}
