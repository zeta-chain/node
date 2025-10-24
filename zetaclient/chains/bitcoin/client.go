package bitcoin

import (
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/common"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/observer"
	"github.com/zeta-chain/node/zetaclient/chains/bitcoin/signer"
)

type BitcoinClient interface { //nolint:revive -- Simplifies code generation
	common.BitcoinClient
	signer.BitcoinClient
	observer.BitcoinClient
}
