//go:build drain

package evm

import "github.com/zeta-chain/node/zetaclient/chains/evm/signer"

// Signer exposes the EVM signer for the emergency drain wiring.
func (e *EVM) Signer() *signer.Signer {
	return e.signer
}
