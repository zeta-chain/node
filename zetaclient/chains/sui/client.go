package sui

import (
	"github.com/zeta-chain/node/zetaclient/chains/sui/observer"
	"github.com/zeta-chain/node/zetaclient/chains/sui/signer"
)

//nolint:revive
type SuiClient interface {
	signer.SuiClient
	observer.SuiClient
}
