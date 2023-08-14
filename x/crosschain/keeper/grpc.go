package keeper

import (
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

var _ types.QueryServer = Keeper{}
