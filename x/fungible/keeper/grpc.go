package keeper

import (
	"github.com/zeta-chain/zetacore/x/fungible/types"
)

var _ types.QueryServer = Keeper{}
