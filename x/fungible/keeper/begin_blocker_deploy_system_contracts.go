//go:build !PRIVNET && !TESTNET
// +build !PRIVNET,!TESTNET

package keeper

import (
	"context"
)

func (k Keeper) BlockOneDeploySystemContracts(goCtx context.Context) error {
	return nil
}
