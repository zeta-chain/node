package keeper

import (
	"github.com/zeta-chain/node/x/observer/types"
)

var _ types.QueryServer = Keeper{}
