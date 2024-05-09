package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.AnteDecorator = SystemPriorityDecorator{}

// SystemPriorityDecorator adds bigger priority for system messages
type SystemPriorityDecorator struct {
}

// NewSystemPriorityDecorator creates a decorator to add bigger priority for system messages
func NewSystemPriorityDecorator() SystemPriorityDecorator {
	return SystemPriorityDecorator{}
}

// AnteHandle implements AnteDecorator
func (vad SystemPriorityDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	newCtx := ctx.WithPriority(500000000) // arbirtrary value, to be revisited, maybe relative to current context.Priority (eg. double it)
	return next(newCtx, tx, simulate)
}
