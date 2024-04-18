package keeper

import (
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

var _ types.QueryServer = Keeper{}
