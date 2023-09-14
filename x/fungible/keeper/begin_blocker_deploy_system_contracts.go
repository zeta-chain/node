//go:build !PRIVNET && !TESTNET
// +build !PRIVNET,!TESTNET

package keeper

import (
	"context"
)

func (k Keeper) BlockOneDeploySystemContracts(_ context.Context) error {
	return nil
}
func (k Keeper) TestUpdateSystemContractAddress(_ context.Context) error {
	return nil
}
func (k Keeper) TestUpdateZRC20WithdrawFee(_ context.Context) error {
	return nil
}
