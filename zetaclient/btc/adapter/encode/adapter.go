package encode

import (
	"context"

	"github.com/zeta-chain/zetacore/zetaclient/btc/model"
)

type Encoder interface {
	Encode(context.Context, *model.Event) ([][]byte, error)
	Decode(context.Context, [][]byte) (*model.Event, error)
}
