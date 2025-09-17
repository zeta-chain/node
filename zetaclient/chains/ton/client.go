package ton

import (
	"github.com/zeta-chain/node/zetaclient/chains/ton/observer"
	"github.com/zeta-chain/node/zetaclient/chains/ton/signer"
)

type Client interface {
	observer.TONClient
	signer.TONClient
}
