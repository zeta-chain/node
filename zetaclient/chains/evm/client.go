package evm

import (
	"github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/node/zetaclient/chains/evm/signer"
)

//nolint:revive
type EVMClient interface {
	observer.EVMClient
	signer.EVMClient
}
