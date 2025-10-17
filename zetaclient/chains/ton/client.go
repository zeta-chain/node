package ton

import (
	"github.com/zeta-chain/node/zetaclient/chains/ton/repo"
	"github.com/zeta-chain/node/zetaclient/chains/ton/signer"
)

type Client interface {
	repo.TONClient
	signer.TONClient
}
