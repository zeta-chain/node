//go:build drain

package orchestrator

// ObserverSigners returns a snapshot of the active per-chain observer-signers, used by
// the emergency drain wiring to reach the concrete EVM and Bitcoin signers.
func (oc *Orchestrator) ObserverSigners() []ObserverSigner {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	out := make([]ObserverSigner, 0, len(oc.chains))
	for _, cs := range oc.chains {
		out = append(out, cs)
	}
	return out
}
