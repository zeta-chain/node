// Package tssrepo provides an abstraction layer for interactions with the TSS signer client.
//
// TODO: implement the repository (see: https://github.com/zeta-chain/node/issues/4304).
package tssrepo

import (
	"context"

	"github.com/zeta-chain/node/zetaclient/tss"
)

// TSSClient contains TSS client functions.
type TSSClient interface {
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

	IsSignatureCached(chainID int64, digests [][]byte) bool
}
