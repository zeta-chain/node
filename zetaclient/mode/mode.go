// Package mode lists the execution modes for the ZetaClient.
package mode

import (
	"fmt"
)

type ClientMode uint8

const (
	// StandardMode represents the standard mode of execution for the ZetaClient.
	//
	// A standard observer-signer observes transactions from ZetaChain, participates in the TSS
	// signing rounds, and relays signed transactions to the appropriate connected chains.
	// Symmetrically, it observes transactions from the connected chains and broadcasts observation
	// votes to ZetaChain.
	StandardMode ClientMode = iota

	// DryMode represents the read-only execution mode for the ZetaClient.
	//
	// A dry observer-signer observes transactions from ZetaChain and the connected chains, but it
	// skips participating in TSS signing rounds and broadcasting transactions and observation
	// votes. In other words, it never mutates the state of the ZetaChain or the state of the
	// connected chains.
	DryMode

	// ChaosMode represents the chaos-testing execution mode for the ZetaClient.
	//
	// An observer-signer in chaos mode works as if in standard mode, but function calls that
	// interact with outside resources (e.g. ZetaChain, connected chains, TSS, and other nodes) may
	// intentionally fail.
	//
	// We use chaos mode to replicate unstable environments for testing.
	ChaosMode
)

// InvalidMode is not a valid client mode.
const InvalidMode ClientMode = 0b11111111

var ErrInvalidModeString = fmt.Errorf("invalid client mode string; should be %q, %q, or %q",
	stringStandard, stringDry, stringChaos)

const (
	stringStandard = "standard"
	stringDry      = "dry"
	stringChaos    = "chaos"
)

// New returns a new ClientMode given its string representation.
func New(s string) (ClientMode, error) {
	switch s {
	case stringStandard:
		return StandardMode, nil
	case stringDry:
		return DryMode, nil
	case stringChaos:
		return ChaosMode, nil
	default:
		return InvalidMode, ErrInvalidModeString
	}
}

func (mode ClientMode) String() string {
	switch mode {
	case StandardMode:
		return stringStandard
	case DryMode:
		return stringDry
	case ChaosMode:
		return stringChaos
	default:
		return fmt.Sprintf("invalid mode: %d", mode)
	}
}

func (mode ClientMode) IsDryMode() bool {
	return mode == DryMode
}

func (mode ClientMode) IsChaosMode() bool {
	return mode == ChaosMode
}
