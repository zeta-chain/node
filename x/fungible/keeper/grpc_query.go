package keeper

import (
	"github.com/zeta-chain/node/x/fungible/types"
)

var _ types.QueryServer = Keeper{}
