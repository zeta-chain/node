package infra

import (
	"github.com/zeta-chain/zetacore/zetaclient/adapters/observer"
)

var _ observer.ChainObserver = (*BitcoinChainObserver)(nil)

type BitcoinChainObserver struct {
}

func NewBitcoinChainObserver() *BitcoinChainObserver {
	return &BitcoinChainObserver{}
}
