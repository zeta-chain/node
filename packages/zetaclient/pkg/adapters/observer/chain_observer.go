package observer

import (
	"context"
	"math/big"

	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/store"
	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/model"
)

type ChainObserver interface {
	GetBlockHeight(ctx context.Context) (uint64, error)
	GetZetaPrice(ctx context.Context) (*big.Int, uint64, error)
	GetGasPrice(ctx context.Context) (uint64, error)
	// received, reverted
	GetConnectorEvents(ctx context.Context, start uint64, end *uint64, filter model.EventFilter) ([]*model.ConnectorEvent, error)
	GetTxByHash(ctx context.Context, hash string, nonce int64) (*model.Receipt, error)
	// out tx
	PrepareTx(ctx context.Context, outTx *model.OutTx) (*model.OutTx, error)
	SendTx(ctx context.Context, outTx *model.OutTx) (*model.OutTxReceipt, error)
	DB() store.Repository
}
