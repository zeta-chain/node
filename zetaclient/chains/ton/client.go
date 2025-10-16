package ton

import (
	"github.com/zeta-chain/node/zetaclient/chains/ton/repo"
	"github.com/zeta-chain/node/zetaclient/chains/ton/signer"
)

//nolint:revive
type TONClient interface {
	repo.TONClient
	signer.TONClient
}
