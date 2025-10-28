// Package mode lists the execution modes for the zetaclient.
package mode

import "fmt"

type ClientMode uint8

const (
	// StandardMode represents the standard execution mode for the zetaclient.
	//
	// An observer-signer in standard mode observes transactions from ZetaChain, signs them, and
	// relays them to the appropriate connected chains. Symmetrically, it observes transactions from
	// the connected chains and relays them to ZetaChain.
	StandardMode ClientMode = iota

	// DryMode represents the read-only execution mode for the zetaclient.
	//
	// An observer-signer in dry-mode only observes the transactions from ZetaChain and the
	// connected chains, without signing them or otherwise mutating the state of the ZetaChain or
	// the state of the connected chains.
	DryMode

	// ChaosMode represents the chaos-testing execution mode for the zetaclient.
	//
	// An observer-signer in chaos-mode works as if in standard mode, but function calls that
	// interact with outside resources (mainly ZetaChain, connected chains, and other nodes) may
	// intentionally fail.
	//
	// We use ChaosMode to replicate unstable environments for testing.
	ChaosMode
)

func (mode ClientMode) String() string {
	switch mode {
	case StandardMode:
		return "standard-mode"
	case DryMode:
		return "dry-mode"
	case ChaosMode:
		return "chaos-mode"
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
