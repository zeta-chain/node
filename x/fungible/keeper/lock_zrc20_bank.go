package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

const (
	transferFrom = "transferFrom"
	allowance    = "allowance"
)

var (
	bankAddress common.Address = common.HexToAddress("0x0000000000000000000000000000000000000067")
	zeroAddress common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

// LockZRC20InBank locks ZRC20 tokens in the bank contract.
// The caller must have approved the bank contract to spend the amount of ZRC20 tokens.
func (k Keeper) LockZRC20InBank(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address, from common.Address,
	amount *big.Int,
) error {
	accAddress := sdk.AccAddress(bankAddress.Bytes())
	if k.GetAuthKeeper().GetAccount(ctx, accAddress) == nil {
		k.GetAuthKeeper().SetAccount(ctx, authtypes.NewBaseAccount(accAddress, nil, 0, 0))
	}

	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if from == zeroAddress {
		return fmt.Errorf("from address cannot be zero")
	}

	if zrc20Address == zeroAddress {
		return fmt.Errorf("zrc20 address cannot be zero")
	}

	if err := k.checkBankAllowance(ctx, zrc20ABI, from, zrc20Address, amount); err != nil {
		return err
	}

	args := []interface{}{from, bankAddress, amount}
	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		bankAddress,
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

func (k Keeper) checkBankAllowance(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	from, zrc20Address common.Address,
	amount *big.Int,
) error {
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	if from == zeroAddress {
		return fmt.Errorf("from address cannot be zero")
	}

	if zrc20Address == zeroAddress {
		return fmt.Errorf("zrc20 address cannot be zero")
	}

	args := []interface{}{from, bankAddress}
	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		bankAddress,
		zrc20Address,
		big.NewInt(0),
		nil,
		true,
		true,
		allowance,
		args...,
	)
	if err != nil {
		return err
	}

	if res.VmError != "" {
		return fmt.Errorf("%s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[allowance].Outputs.Unpack(res.Ret)
	if err != nil {
		return err
	}

	allowance, ok := ret[0].(*big.Int)
	if !ok {
		return fmt.Errorf("ZRC20 allowance returned an unexpected type")
	}

	if allowance.Cmp(amount) < 0 || allowance.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("invalid allowance, got: %s", allowance.String())
	}

	return nil
}
