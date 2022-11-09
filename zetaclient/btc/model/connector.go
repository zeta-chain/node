package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

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
	opReturn        = "OP_RETURN "
)

func NewConnectorEvent(amount *big.Float, raw []byte) (*ConnectorEvent, error) {
	if raw == nil || len(raw) < minLen {
		return nil, errors.New(invalidInput)
	}
	if raw[0] != ZetaMagicNumber {
		return nil, errors.New(noMagicNumber)
	}
	address := common.BytesToAddress(raw[1 : addrLen+1])
	return &ConnectorEvent{
		Amount:  amount,
		Address: address,
		Message: raw[addrLen+1:],
	}, nil
}

func (evt *ConnectorEvent) String() string {
	return fmt.Sprintf("Amount: %v, Address: %v, Message: %v", evt.Amount, evt.Address.Hex(), string(evt.Message))
}

func ParseRawEvents(rawEvents []*RawEvent) (*ConnectorEvent, error) {
	var amount *big.Float
	var evtStr string
	for _, raw := range rawEvents {
		// if raw event have an address is tss_address and value is tx amount
		if len(raw.Addresses) > 0 {
			amount = big.NewFloat(raw.Values[0])
			continue
		}
		// process memo field
		// check OP_RETURN
		if len(raw.ASM) >= (len(opReturn)+addrLen+1) && strings.HasPrefix(raw.ASM, opReturn) {
			evtStr = raw.ASM[len(opReturn):]
		}
	}
	evtBytes, err := hex.DecodeString(evtStr)
	if err != nil {
		return nil, err
	}
	return NewConnectorEvent(amount, evtBytes)
}
