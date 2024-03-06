package keeper

import (
	"github.com/zeta-chain/zetacore/x/authority/types"
)

var _ types.QueryServer = Keeper{}
