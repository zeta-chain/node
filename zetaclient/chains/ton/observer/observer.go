package observer

import (
	"context"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/base"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
)

type Observer struct {
	base.Observer

	client    *liteapi.Client
	gatewayID ton.AccountID
}

var _ interfaces.ChainObserver = (*Observer)(nil)

func New(bo *base.Observer, client *liteapi.Client, gatewayID ton.AccountID) (*Observer, error) {
	bo.LoadLastTxScanned()

	return &Observer{
		Observer:  *bo,
		gatewayID: gatewayID,
		client:    client,
	}, nil
}

func (ob *Observer) Start(ctx context.Context) {
	// todo
}

func (ob *Observer) Stop() {
	// todo
}

func (ob *Observer) VoteOutboundIfConfirmed(ctx context.Context, cctx *types.CrossChainTx) (bool, error) {
	// todo
	return false, nil
}
