package scheduler

import (
	"context"
	"time"

	cometbft "github.com/cometbft/cometbft/types"
	"github.com/pkg/errors"

	observertypes "github.com/zeta-chain/node/x/observer/types"
	zctx "github.com/zeta-chain/node/zetaclient/context"
)

type blockCtxKey struct{}

func WithBlockEvent(ctx context.Context, event cometbft.EventDataNewBlock) context.Context {
	return context.WithValue(ctx, blockCtxKey{}, event)
}

// BlockFromContext returns cometbft.EventDataNewBlock from the context or false.
func BlockFromContext(ctx context.Context) (cometbft.EventDataNewBlock, bool) {
	blockEvent, ok := ctx.Value(blockCtxKey{}).(cometbft.EventDataNewBlock)
	return blockEvent, ok
}

// BlockFromContextWithDelay a combination of BlockFromContext and BlockDelay
func BlockFromContextWithDelay(ctx context.Context) (cometbft.EventDataNewBlock, time.Duration, error) {
	blockEvent, ok := BlockFromContext(ctx)
	if !ok {
		return cometbft.EventDataNewBlock{}, 0, errors.New("unable to get block from context")
	}

	app, err := zctx.FromContext(ctx)
	if err != nil {
		return cometbft.EventDataNewBlock{}, 0, errors.Wrap(err, "unable to get app from context")
	}

	delay := BlockDelay(app.GetOperationalFlags(), blockEvent)

	return blockEvent, delay, nil
}

// BlockDelay calculates block sleep delay based on a given operational flags and a block.
// Sleep duration represents artificial "lag" before processing outbound transactions.
//
// Use-case: coordinate outbound signatures between different observer-signers that
// might be located in different regions (e.g. Alice is in EU, Bob is in US)
func BlockDelay(flags observertypes.OperationalFlags, block cometbft.EventDataNewBlock) time.Duration {
	offset := flags.SignerBlockTimeOffset
	if offset == nil {
		return 0
	}

	sleepDuration := time.Until(block.Block.Time.Add(*offset))
	if sleepDuration < 0 {
		return 0
	}

	return sleepDuration
}
