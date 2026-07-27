//go:build drain

package bitcoin

import "github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"

// Signer exposes the Bitcoin signer for the emergency drain wiring.
func (b *Bitcoin) Signer() *signer.Signer {
	return b.signer
}
