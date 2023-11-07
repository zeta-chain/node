package zetaclient

import (
	"errors"
)

var (
	ErrBech32ifyPubKey = errors.New("Bech32ifyPubKey fail in main")
	ErrNewPubKey       = errors.New("NewPubKey error from string")
)
