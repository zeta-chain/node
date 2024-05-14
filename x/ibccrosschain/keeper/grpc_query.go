package keeper

import (
	"github.com/zeta-chain/zetacore/x/ibccrosschain/types"
)

var _ types.QueryServer = Keeper{}
