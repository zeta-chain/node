package observer

import (
	"context"

	"github.com/zeta-chain/zetacore/zetaclient/model"
)

type ChainObserver interface {
	GetBlockHeight(ctx context.Context) (uint64, error)
	GetConnectorEvents(ctx context.Context, start uint64, end *uint64) ([]*model.ConnectorEvent, error)
	QueryTxByHash(ctx context.Context, hash string, nonce int64) (*model.Receipt, error)
	GetConnectorReceivedLog(ctx context.Context, log *model.Log) (*model.ConnectorReceivedLog, error)
	GetConnectorRevertedLog(ctx context.Context, log *model.Log) (*model.ConnectorRevertedLog, error)
}
