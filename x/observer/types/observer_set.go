package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
)

func (m *ObserverSet) Len() int {
	return len(m.ObserverList)
}

func (m *ObserverSet) LenUint() uint64 {
	return uint64(len(m.ObserverList))
}

// Validate observer set verifies that the observer set is valid
// - All observer addresses are valid
// - No duplicate observer addresses
func (m *ObserverSet) Validate() error {
	observers := make(map[string]struct{})
	for _, observerAddress := range m.ObserverList {
		// Check for valid observer addresses
		_, err := sdk.AccAddressFromBech32(observerAddress)
		if err != nil {
			return errors.Wrapf(ErrInvalidObserverAddress, "observer %s err %s", observerAddress, err.Error())
		}
		// Check for duplicates
		if _, ok := observers[observerAddress]; ok {
			return errors.Wrapf(ErrDuplicateObserver, "observer %s", observerAddress)
		}

		observers[observerAddress] = struct{}{}
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
