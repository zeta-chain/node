package maintenance

import (
	"context"

	cometbft "github.com/cometbft/cometbft/types"

	observer "github.com/zeta-chain/node/x/observer/types"
)

type ZetacoreClient interface {
	NewBlockSubscriber(context.Context) (chan cometbft.EventDataNewBlock, error)
	GetKeyGen(context.Context) (observer.Keygen, error)
	GetTSS(context.Context) (observer.TSS, error)
	GetTSSHistory(context.Context) ([]observer.TSS, error)
	GetBlockHeight(context.Context) (int64, error)
	GetOperationalFlags(context.Context) (observer.OperationalFlags, error)
}
