package observer

import "github.com/zeta-chain/node/zetaclient/chains/base"

// Observer SUI observer
type Observer struct {
	*base.Observer
}

// New Observer constructor.
func New(baseObserver *base.Observer) *Observer {
	return &Observer{
		Observer: baseObserver,
	}
}
