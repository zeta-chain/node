package repo

import (
	"context"
	"errors"
	"time"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/zetaclient/chains/ton/encoder"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
)

// GetTransactionsSince TODO.
func (repo *TONRepo) GetTransactionsSince(ctx context.Context,
	lastTx string,
) ([]ton.Transaction, error) {
	accountID := repo.Gateway.AccountID()

	lastLT, lastHash, err := encoder.DecodeTx(lastTx)
	if err != nil {
		return nil, errors.Join(ErrTransactionEncoding, err)
	}

	txs, err := repo.Client.GetTransactionsSince(ctx, accountID, lastLT, lastHash)
	if err != nil {
		return nil, errors.Join(ErrGetTransactionsSince, err)
	}

	return txs, nil
}

// func (repo *TONRepo) GetInboundTrackers(ctx context.Context) ([]types.InboundTracker, error) {
// 	chainID := repo.connectedChain.ChainId
// 	trackers, err := repo.ZetacoreClient.GetInboundTrackersForChain(ctx, chainID)
// 	if err != nil {
// 		return nil, errors.Join(ErrGetInboundTrackers, err)
// 	}
// 	return trackers, nil
// }

// ------------------------------------------------------------------------------------------------

// TODO: Duplicate; remove this before merge.
type TONClient interface {
	GetConfigParam(_ context.Context, index uint32) (*boc.Cell, error)

	GetBlockHeader(_ context.Context, blockID rpc.BlockIDExt) (rpc.BlockHeader, error)

	GetMasterchainInfo(context.Context) (rpc.MasterchainInfo, error)

	HealthCheck(context.Context) (time.Time, error)

	GetTransaction(_ context.Context,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) (ton.Transaction, error)

	GetTransactions(_ context.Context,
		count uint32,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) ([]ton.Transaction, error)

	GetTransactionsSince(_ context.Context,
		_ ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) ([]ton.Transaction, error)
}
