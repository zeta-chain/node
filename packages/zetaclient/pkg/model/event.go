package model

import "math/big"

type ConnectorEvent struct {
	SourceTxOriginAddress string
	ZetaTxSenderAddress   string
	DestinationChainId    *big.Int
	DestinationAddress    []byte
	ZetaValueAndGas       *big.Int
	DestinationGasLimit   *big.Int
	Message               []byte
	ZetaParams            []byte
	BlockNumber           uint64
	TxHash                string
	//	Raw                   []byte
}

type EventFilter string

const (
	FilterAllEvents      = EventFilter("all")
	FilterReceivedEvents = EventFilter("received")
	FilterRevertedEvents = EventFilter("reverted")
)
