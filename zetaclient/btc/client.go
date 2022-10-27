package btc

import (
	"github.com/zeta-chain/zetacore/zetaclient/btc/model"
)

type Client interface {
	GetBlockHeight() (int64, error)
	GetBlockHash(int64) (string, error)
	GetEventsByHash(string) ([]*model.Event, error)
}
