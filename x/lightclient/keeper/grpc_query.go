package keeper

import (
	"github.com/zeta-chain/node/x/lightclient/types"
)

var _ types.QueryServer = Keeper{}
