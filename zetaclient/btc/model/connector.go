package model

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

type ConnectorEvent struct {
	DestChainID int32
	DestAddress common.Address
	Amount      uint64
}

func (evt *ConnectorEvent) ToBTCOP() string {
	return fmt.Sprintf("Zeta%6.6d%s%18.18d", evt.DestChainID, evt.DestAddress.Hex(), evt.Amount)
}
