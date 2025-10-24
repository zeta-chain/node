package sui

import (
	"github.com/zeta-chain/node/zetaclient/chains/sui/observer"
	"github.com/zeta-chain/node/zetaclient/chains/sui/signer"
)

type SuiClient interface { //nolint:revive -- Simplifies code generation
	signer.SuiClient
	observer.SuiClient
}
