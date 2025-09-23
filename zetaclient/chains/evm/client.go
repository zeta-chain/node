package evm

import (
	"github.com/zeta-chain/node/zetaclient/chains/evm/observer"
	"github.com/zeta-chain/node/zetaclient/chains/evm/signer"
)

type Client interface {
	observer.EVMClient
	signer.EVMClient
}
