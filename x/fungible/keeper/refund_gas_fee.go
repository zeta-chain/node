package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/x/fungible/types"
)

// RefundRemainingGasFees refunds the remaining gas fees to the receiver
func (k Keeper) RefundRemainingGasFees(
	ctx sdk.Context,
	chainID int64,
	amount *big.Int,
	receiver ethcommon.Address,
) error {
	// get the gas zrc20 contract from the chain
	gasZRC20, err := k.QuerySystemContractGasCoinZRC20(ctx, big.NewInt(chainID))
	if err != nil {
		return err
	}

	// call deposit ZRC20 method
	return k.CallZRC20Deposit(
		ctx,
		types.ModuleAddressEVM,
		gasZRC20,
		receiver,
		amount,
		true,
	)
}
