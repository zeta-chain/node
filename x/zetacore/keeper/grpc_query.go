package keeper

import (
	"github.com/zeta-chain/zetacore/x/zetacore/types"
)

var _ types.QueryServer = Keeper{}
