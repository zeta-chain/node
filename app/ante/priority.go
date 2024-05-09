package ante

import (
	"math"

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
	newCtx := ctx.WithPriority(math.MaxInt64)
	return next(newCtx, tx, simulate)
}
