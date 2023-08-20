package types

import (
	"errors"
	"fmt"

	"github.com/zeta-chain/zetacore/common"
)

// Validate observer mapper contains an existing chain
func (m *ObserverMapper) Validate() error {
	if m.ObserverChain == nil {
		return errors.New("observer chain is not set")
	}

	chains := common.DefaultChainsList()
	for _, chain := range chains {
		if *m.ObserverChain == *chain {
			return nil
		}
	}
	return fmt.Errorf("observer chain %d doesn't exist: ", m.ObserverChain.ChainName)
}

// VerifyObserverMapper verifies list of observer mappers
func VerifyObserverMapper(obs []*ObserverMapper) error {
	for _, mapper := range obs {
		if mapper != nil {
			err := mapper.Validate()
			if err != nil {
				return fmt.Errorf("observer mapper %s is invalid: %s", mapper.String(), err.Error())
			}
		}
	}
	return nil
}

func CheckReceiveStatus(status common.ReceiveStatus) error {
	switch status {
	case common.ReceiveStatus_Success:
		return nil
	case common.ReceiveStatus_Failed:
		return nil
	default:
		return ErrInvalidStatus
	}
}
