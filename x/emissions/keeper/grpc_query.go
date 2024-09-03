package keeper

import (
	"github.com/zeta-chain/node/x/emissions/types"
)

var _ types.QueryServer = Keeper{}
