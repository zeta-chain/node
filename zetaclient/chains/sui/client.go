package sui

import (
	"github.com/zeta-chain/node/zetaclient/chains/sui/observer"
	"github.com/zeta-chain/node/zetaclient/chains/sui/signer"
)

type SuiClient interface { //nolint:revive
	signer.SuiClient
	observer.SuiClient
}
