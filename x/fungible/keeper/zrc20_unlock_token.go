package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

const transfer = "transfer"

// UnlockZRC20 unlocks ZRC20 tokens and sends them to the "to" address.
func (k Keeper) UnlockZRC20(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address, to common.Address,
	amount *big.Int,
) error {
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive, got: %s", amount.String())
	}

	if to == zeroAddress {
		return fmt.Errorf("from address cannot be zero")
	}

	if zrc20Address == zeroAddress {
		return fmt.Errorf("zrc20 address cannot be zero")
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return err
	}

	if err := k.CheckFungibleZRC20Balance(ctx, zrc20ABI, zrc20Address, amount); err != nil {
		return err
	}

	args := []interface{}{to, amount}
	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		fungibletypes.ModuleAddressZEVM,
		zrc20Address,
		big.NewInt(0),
		nil,
		true,
		true,
		transfer,
		args...,
	)
	if err != nil {
		return err
	}

	if res.VmError != "" {
		return fmt.Errorf("%s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[transfer].Outputs.Unpack(res.Ret)
	if err != nil {
		return err
	}

	transferred, ok := ret[0].(bool)
	if !ok {
		return fmt.Errorf("transfer returned an unexpected value")
	}

	if !transferred {
		return fmt.Errorf("transfer not successful")
	}

	return nil
}
