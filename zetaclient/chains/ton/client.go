package ton

import (
	"github.com/zeta-chain/node/zetaclient/chains/ton/repo"
	"github.com/zeta-chain/node/zetaclient/chains/ton/signer"
)

type TONClient interface { //nolint:revive
	repo.TONClient
	signer.TONClient
}
