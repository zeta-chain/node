package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/pkg/crypto"
	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

const (
	transferFrom = "transferFrom"
	transfer     = "transfer"
	balanceOf    = "balanceOf"
	allowance    = "allowance"
)

var (
	ErrZRC20ZeroAddress    = fmt.Errorf("ZRC20 address cannot be zero")
	ErrZRC20NotWhiteListed = fmt.Errorf("ZRC20 is not whitelisted")
	ErrZRC20Paused         = fmt.Errorf("ZRC20 is paused")
	ErrZRC20NilABI         = fmt.Errorf("ZRC20 ABI is nil")
	ErrZeroAddress         = fmt.Errorf("address cannot be zero")
	ErrInvalidAmount       = fmt.Errorf("amount must be positive")
)

// LockZRC20 locks ZRC20 tokens in the bank contract.
// The caller must have approved the bank contract to spend the amount of ZRC20 tokens.
func (k Keeper) LockZRC20(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address, from common.Address,
	amount *big.Int,
) error {
	if zrc20ABI == nil {
		return ErrZRC20NilABI
	}

	if amount.Sign() <= 0 || amount == nil {
		return ErrInvalidAmount
	}

	if crypto.IsEmptyAddress(from) {
		return ErrZeroAddress
	}

	if crypto.IsEmptyAddress(zrc20Address) {
		return ErrZRC20ZeroAddress
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return err
	}

	if err := k.CheckFungibleZRC20Allowance(ctx, zrc20ABI, from, zrc20Address, amount); err != nil {
		return err
	}

	args := []interface{}{from, fungibletypes.ModuleAddressZEVM, amount}
	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		fungibletypes.ModuleAddressZEVM,
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
		return fmt.Errorf("EVM execution error in LockZRC20: %s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[transferFrom].Outputs.Unpack(res.Ret)
	if err != nil {
		return err
	}

	if len(ret) == 0 {
		return fmt.Errorf("no data returned from 'transferFrom' method")
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

// UnlockZRC20 unlocks ZRC20 tokens and sends them to the "to" address.
func (k Keeper) UnlockZRC20(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address, to common.Address,
	amount *big.Int,
) error {
	if zrc20ABI == nil {
		return ErrZRC20NilABI
	}

	if amount.Sign() <= 0 || amount == nil {
		return ErrInvalidAmount
	}

	if crypto.IsEmptyAddress(to) {
		return ErrZeroAddress
	}

	if crypto.IsEmptyAddress(zrc20Address) {
		return ErrZRC20ZeroAddress
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
		return fmt.Errorf("EVM execution error in UnlockZRC20: %s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[transfer].Outputs.Unpack(res.Ret)
	if err != nil {
		return err
	}

	if len(ret) == 0 {
		return fmt.Errorf("no data returned from 'transfer' method")
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

// CheckFungibleZRC20Allowance checks if the allowance of ZRC20 tokens,
// is equal or greater than the provided amount.
func (k Keeper) CheckFungibleZRC20Allowance(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	from, zrc20Address common.Address,
	amount *big.Int,
) error {
	if zrc20ABI == nil {
		return ErrZRC20NilABI
	}

	if amount.Sign() <= 0 || amount == nil {
		return ErrInvalidAmount
	}

	if crypto.IsEmptyAddress(from) {
		return ErrZeroAddress
	}

	if crypto.IsEmptyAddress(zrc20Address) {
		return ErrZRC20ZeroAddress
	}

	args := []interface{}{from, fungibletypes.ModuleAddressZEVM}
	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		fungibletypes.ModuleAddressZEVM,
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
		return fmt.Errorf("EVM execution error calling allowance: %s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[allowance].Outputs.Unpack(res.Ret)
	if err != nil {
		return err
	}

	if len(ret) == 0 {
		return fmt.Errorf("no data returned from 'allowance' method")
	}

	allowanceValue, ok := ret[0].(*big.Int)
	if !ok {
		return fmt.Errorf("ZRC20 allowance returned an unexpected type")
	}

	if allowanceValue.Cmp(amount) < 0 || allowanceValue.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("invalid allowance, got: %s", allowanceValue.String())
	}

	return nil
}

// CheckFungibleZRC20Balance checks if the balance of ZRC20 tokens,
// is equal or greater than the provided amount.
func (k Keeper) CheckFungibleZRC20Balance(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address common.Address,
	amount *big.Int,
) error {
	if zrc20ABI == nil {
		return ErrZRC20NilABI
	}

	if amount.Sign() <= 0 || amount == nil {
		return ErrInvalidAmount
	}

	if crypto.IsEmptyAddress(zrc20Address) {
		return ErrZRC20ZeroAddress
	}

	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		fungibletypes.ModuleAddressZEVM,
		zrc20Address,
		big.NewInt(0),
		nil,
		true,
		true,
		balanceOf,
		fungibletypes.ModuleAddressZEVM,
	)
	if err != nil {
		return err
	}

	if res.VmError != "" {
		return fmt.Errorf("EVM execution error calling balanceOf: %s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[balanceOf].Outputs.Unpack(res.Ret)
	if err != nil {
		return err
	}

	if len(ret) == 0 {
		return fmt.Errorf("no data returned from 'balanceOf' method")
	}

	balance, ok := ret[0].(*big.Int)
	if !ok {
		return fmt.Errorf("ZRC20 balanceOf returned an unexpected type")
	}

	if balance.Cmp(amount) == -1 {
		return fmt.Errorf("invalid balance, got: %s", balance.String())
	}

	return nil
}

// IsValidZRC20 returns an error whenever a ZRC20 is not whitelisted or paused.
func (k Keeper) IsValidZRC20(ctx sdk.Context, zrc20Address common.Address) error {
	if crypto.IsEmptyAddress(zrc20Address) {
		return ErrZRC20ZeroAddress
	}

	t, found := k.GetForeignCoins(ctx, zrc20Address.String())
	if !found {
		return ErrZRC20NotWhiteListed
	}

	if t.Paused {
		return ErrZRC20Paused
	}

	return nil
}
