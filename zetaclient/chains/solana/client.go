package solana

import (
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
	"github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	"github.com/zeta-chain/node/zetaclient/chains/solana/signer"
)

// TODO: Replace this interface for a repository interface.
// See: https://github.com/zeta-chain/node/issues/4224
//
//nolint:revive
type SolanaClient interface {
	observer.SolanaClient
	signer.SolanaClient
	repo.SolanaClient
}
