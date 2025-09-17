package sui

import (
	"github.com/zeta-chain/node/zetaclient/chains/sui/observer"
	"github.com/zeta-chain/node/zetaclient/chains/sui/signer"
)

type Client interface {
	signer.SuiClient
	observer.SuiClient
}
