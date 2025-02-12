package signer

import "github.com/zeta-chain/node/zetaclient/chains/base"

// Signer Sui outbound transaction signer.
type Signer struct {
	*base.Signer
}

// New Signer constructor.
func New(baseSigner *base.Signer) *Signer {
	return &Signer{Signer: baseSigner}
}
