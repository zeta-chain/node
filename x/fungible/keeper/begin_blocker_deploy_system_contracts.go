//go:build !PRIVNET && !TESTNET && !MOCK_MAINNET
// +build !PRIVNET,!TESTNET,!MOCK_MAINNET

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
