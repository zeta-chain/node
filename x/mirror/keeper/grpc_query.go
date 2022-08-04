package keeper

import (
	"github.com/zeta-chain/zetacore/x/mirror/types"
)

var _ types.QueryServer = Keeper{}
