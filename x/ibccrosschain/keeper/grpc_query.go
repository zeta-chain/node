package keeper

import (
	"github.com/zeta-chain/node/x/ibccrosschain/types"
)

var _ types.QueryServer = Keeper{}
