package types

import (
	"context"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

type EvmClient interface {
	TransactionReceipt(ctx context.Context, hash ethcommon.Hash)
}
