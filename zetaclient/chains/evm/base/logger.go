package base

import "github.com/rs/zerolog"

// ObserverLogger is the base logger for chain observers
type ObserverLogger struct {
	// the parent logger for the chain observer
	Chain zerolog.Logger

	// the logger for inbound transactions
	Inbound zerolog.Logger

	// the logger for outbound transactions
	Outbound zerolog.Logger

	// the logger for the chain's gas price
	GasPrice zerolog.Logger

	// the logger for the compliance check
	Compliance zerolog.Logger
}
