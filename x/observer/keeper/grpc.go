package keeper

import (
	"github.com/zeta-chain/zetacore/x/observer/types"
)

var _ types.QueryServer = Keeper{}
