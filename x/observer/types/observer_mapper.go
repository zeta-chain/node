package types

import (
	"github.com/zeta-chain/zetacore/common"
)

func (m *ObserverMapper) Validate() bool {
	chains := DefaultChainsList()
	for _, chain := range chains {
		if m.ObserverChain == chain {
			return true
		}
	}
	return false
}

func VerifyObserverMapper(obs []*ObserverMapper) bool {
	for _, mapper := range obs {
		ok := mapper.Validate()
		if !ok {
			return ok
		}
	}
	return true
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
