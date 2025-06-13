package mocks

import (
	"context"
	"time"

	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
)

//go:generate mockery --name tonRPC --structname TONRPC --filename ton_rpc.go --output ./
//nolint:unused // used for code gen
type tonRPC interface {
	GetAccountState(ctx context.Context, acc ton.AccountID) (rpc.Account, error)
	GetBlockHeader(ctx context.Context, blockID rpc.BlockIDExt) (rpc.BlockHeader, error)
	GetConfigParam(ctx context.Context, index uint32) (*boc.Cell, error)
	GetMasterchainInfo(ctx context.Context) (rpc.MasterchainInfo, error)
	GetTransaction(ctx context.Context, acc ton.AccountID, lt uint64, hash ton.Bits256) (ton.Transaction, error)
	GetTransactions(
		ctx context.Context,
		count uint32,
		accountID ton.AccountID,
		lt uint64,
		hash ton.Bits256,
	) ([]ton.Transaction, error)
	GetTransactionsSince(
		ctx context.Context,
		acc ton.AccountID,
		oldestLT uint64,
		oldestHash ton.Bits256,
	) (txs []ton.Transaction, err error)
	HealthCheck(ctx context.Context) (time.Time, error)
	SendMessage(ctx context.Context, payload []byte) (uint32, error)
}

var _ tonRPC = &rpc.Client{}
