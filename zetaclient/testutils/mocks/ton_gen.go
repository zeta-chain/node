package mocks

import (
	"context"
	"time"

	"github.com/tonkeeper/tongo/liteclient"
	"github.com/tonkeeper/tongo/tlb"
	"github.com/tonkeeper/tongo/ton"

	"github.com/zeta-chain/node/zetaclient/chains/ton/config"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
)

//go:generate mockery --name tonClient --structname TONLiteClient --filename ton_liteclient.go --output ./
//nolint:unused // used for code gen
type tonClient interface {
	config.Getter
	GetMasterchainInfo(ctx context.Context) (liteclient.LiteServerMasterchainInfoC, error)
	GetBlockHeader(ctx context.Context, blockID ton.BlockIDExt, mode uint32) (tlb.BlockInfo, error)
	GetTransactionsSince(ctx context.Context, acc ton.AccountID, lt uint64, hash ton.Bits256) ([]ton.Transaction, error)
	GetFirstTransaction(ctx context.Context, acc ton.AccountID) (*ton.Transaction, int, error)
	GetTransaction(ctx context.Context, acc ton.AccountID, lt uint64, hash ton.Bits256) (ton.Transaction, error)
	HealthCheck(ctx context.Context) (time.Time, error)
	GetAccountState(ctx context.Context, accountID ton.AccountID) (tlb.ShardAccount, error)
	SendMessage(ctx context.Context, payload []byte) (uint32, error)
}

var _ tonClient = &liteapi.Client{}
