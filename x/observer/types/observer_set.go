package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/zetacore/pkg/chains"
)

func (m *ObserverSet) Len() int {
	return len(m.ObserverList)
}

func (m *ObserverSet) LenUint() uint64 {
	return uint64(len(m.ObserverList))
}

// Validate observer mapper contains an existing chain
func (m *ObserverSet) Validate() error {
	for _, observerAddress := range m.ObserverList {
		_, err := sdk.AccAddressFromBech32(observerAddress)
		if err != nil {
			return err
		}
	}
	// Check for duplicates
	observers := make(map[string]bool)
	for _, observerAddress := range m.ObserverList {
		if _, ok := observers[observerAddress]; ok {
			return errors.Wrapf(ErrDuplicateObserver, "observer %s", observerAddress)
		}
		observers[observerAddress] = true
	}
	return nil
}

func CheckReceiveStatus(status chains.ReceiveStatus) error {
	switch status {
	case chains.ReceiveStatus_success:
		return nil
	case chains.ReceiveStatus_failed:
		return nil
	default:
		return ErrInvalidStatus
	}
}
