package keeper

import (
	"github.com/zeta-chain/node/x/authority/types"
)

var _ types.QueryServer = Keeper{}
