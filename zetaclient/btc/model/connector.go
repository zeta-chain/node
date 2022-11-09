package model

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ConnectorEvent struct {
	Amount  *big.Float
	Address common.Address
	Message []byte
}

const (
	ZetaMagicNumber = 0x5a // Z
	invalidInput    = "invalid input"
	invalidAmount   = "invalid amount"
	noMagicNumber   = "no magic number"
	addrLen         = 20 // len of expected address
	minLen          = addrLen + 1
)

func NewConnectorEvent(amt string, raw []byte) (*ConnectorEvent, error) {
	if raw == nil || len(raw) < minLen {
		return nil, errors.New(invalidInput)
	}
	if raw[0] != ZetaMagicNumber {
		return nil, errors.New(noMagicNumber)
	}
	address := common.BytesToAddress(raw[1 : addrLen+1])
	amount, ok := new(big.Float).SetString(amt)
	if !ok {
		return nil, errors.New(invalidAmount)
	}
	return &ConnectorEvent{
		Amount:  amount,
		Address: address,
		Message: raw[addrLen+1:],
	}, nil
}

func (evt *ConnectorEvent) String() string {
	return fmt.Sprintf("Amount: %v, Address: %v, Message: %v\n", evt.Amount, evt.Address.Hex(), string(evt.Message))
}
