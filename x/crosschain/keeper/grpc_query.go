package keeper

import (
	"github.com/zeta-chain/node/x/crosschain/types"
)

var _ types.QueryServer = Keeper{}
