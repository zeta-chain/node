package repo

import (
	"context"
	"time"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
)

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
