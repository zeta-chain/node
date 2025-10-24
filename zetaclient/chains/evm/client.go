package evm

import (
	"github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/node/zetaclient/chains/evm/signer"
)

type EVMClient interface { //nolint:revive -- Simplifies code generation
	observer.EVMClient
	signer.EVMClient
}
