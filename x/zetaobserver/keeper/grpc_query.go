package keeper

import (
	"github.com/zeta-chain/zetacore/x/zetaobserver/types"
)

var _ types.QueryServer = Keeper{}
