package mocks

import (
	"github.com/zeta-chain/node/zetaclient/chains/solana/observer"
)

//go:generate mockery --name solanaRepo --structname SolanaRepo --filename solana.go --output ../

//nolint:unused // used for code gen
type solanaRepo interface {
	observer.SolanaRepo
}
