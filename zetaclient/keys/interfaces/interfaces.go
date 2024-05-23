package interfaces

import (
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ObserverKeys is the interface for observer's keys
type ObserverKeys interface {
	GetSignerInfo() *ckeys.Record
	GetOperatorAddress() sdk.AccAddress
	GetAddress() (sdk.AccAddress, error)
	GetPrivateKey(password string) (cryptotypes.PrivKey, error)
	GetKeybase() ckeys.Keyring
	GetHotkeyPassword() string
}
