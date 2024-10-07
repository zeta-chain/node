package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

const transferFrom = "transferFrom"

var zeroAddress common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")

// LockZRC20 locks ZRC20 tokens in the bank contract.
// The caller must have approved the bank contract to spend the amount of ZRC20 tokens.
func (k Keeper) LockZRC20(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address, from common.Address,
	amount *big.Int,
) error {
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive, got: %s", amount.String())
	}

	if from == zeroAddress {
		return fmt.Errorf("from address cannot be zero")
	}

	if zrc20Address == zeroAddress {
		return fmt.Errorf("zrc20 address cannot be zero")
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return err
	}

	if err := k.CheckFungibleZRC20Allowance(ctx, zrc20ABI, from, zrc20Address, amount); err != nil {
		return err
	}

	args := []interface{}{from, fungibletypes.ModuleAddressEVM, amount}
	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		fungibletypes.ModuleAddressEVM,
		zrc20Address,
		big.NewInt(0),
		nil,
		true,
		true,
		transferFrom,
		args...,
	)
	if err != nil {
		return err
	}

	if res.VmError != "" {
		return fmt.Errorf("%s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[transferFrom].Outputs.Unpack(res.Ret)
	if err != nil {
		return err
	}

	transferred, ok := ret[0].(bool)
	if !ok {
		return fmt.Errorf("transferFrom returned an unexpected value")
	}

	if !transferred {
		return fmt.Errorf("transferFrom not successful")
	}

	return nil
}
