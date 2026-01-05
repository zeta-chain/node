// Package interfaces provides interfaces for clients and signers for the chain to interact with
package interfaces

import (
	"context"

	"github.com/zeta-chain/node/zetaclient/tss"
)

// TSSSigner is the interface for the TSS signer.
type TSSSigner interface {
	PubKey() tss.PubKey

	Sign(_ context.Context,
		data []byte,
		height uint64,
		nonce uint64,
		chainID int64,
	) ([65]byte, error)

	SignBatch(_ context.Context,
		digests [][]byte,
		height uint64,
		nonce uint64,
		chainID int64,
	) ([][65]byte, error)
}
