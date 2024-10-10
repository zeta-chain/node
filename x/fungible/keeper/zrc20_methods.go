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
	allowance    = "allowance"
	balanceOf    = "balanceOf"
	totalSupply  = "totalSupply"
	transfer     = "transfer"
	transferFrom = "transferFrom"
)

// ZRC20Allowance returns the ZRC20 allowance for a given spender.
func (k Keeper) ZRC20Allowance(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address, owner, spender common.Address,
) (*big.Int, error) {
	if zrc20ABI == nil {
		return nil, fungibletypes.ErrZRC20NilABI
	}

	if crypto.IsEmptyAddress(owner) || crypto.IsEmptyAddress(spender) {
		return nil, fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return nil, err
	}

	// function allowance(address owner, address spender)
	args := []interface{}{owner, spender}
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
		return nil, err
	}

	if res.VmError != "" {
		return nil, fmt.Errorf("EVM execution error calling allowance: %s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[allowance].Outputs.Unpack(res.Ret)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("no data returned from 'allowance' method")
	}

	allowanceValue, ok := ret[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("ZRC20 allowance returned an unexpected type")
	}

	return allowanceValue, nil
}

// ZRC20BalanceOf checks the ZRC20 balance of a given EOA.
func (k Keeper) ZRC20BalanceOf(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address, owner common.Address,
) (*big.Int, error) {
	if zrc20ABI == nil {
		return nil, fungibletypes.ErrZRC20NilABI
	}

	if crypto.IsEmptyAddress(owner) {
		return nil, fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return nil, err
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
		owner,
	)
	if err != nil {
		return nil, err
	}

	if res.VmError != "" {
		return nil, fmt.Errorf("EVM execution error calling balanceOf: %s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[balanceOf].Outputs.Unpack(res.Ret)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("no data returned from 'balanceOf' method")
	}

	balance, ok := ret[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("ZRC20 balanceOf returned an unexpected type")
	}

	return balance, nil
}

// ZRC20TotalSupply returns the total supply of a ZRC20 token.
func (k Keeper) ZRC20TotalSupply(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address common.Address,
) (*big.Int, error) {
	if zrc20ABI == nil {
		return nil, fungibletypes.ErrZRC20NilABI
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return nil, err
	}

	// function totalSupply() public view virtual override returns (uint256)
	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		fungibletypes.ModuleAddressZEVM,
		zrc20Address,
		big.NewInt(0),
		nil,
		true,
		true,
		totalSupply,
	)
	if err != nil {
		return nil, err
	}

	if res.VmError != "" {
		return nil, fmt.Errorf("EVM execution error calling totalSupply: %s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[totalSupply].Outputs.Unpack(res.Ret)
	if err != nil {
		return nil, err
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("no data returned from 'totalSupply' method")
	}

	totalSupply, ok := ret[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("ZRC20 totalSupply returned an unexpected type")
	}

	return totalSupply, nil
}

// ZRC20Transfer transfers ZRC20 tokens from the sender to the recipient.
func (k Keeper) ZRC20Transfer(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address, from, to common.Address,
	amount *big.Int,
) (bool, error) {
	if zrc20ABI == nil {
		return false, fungibletypes.ErrZRC20NilABI
	}

	if crypto.IsEmptyAddress(from) || crypto.IsEmptyAddress(to) {
		return false, fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return false, err
	}

	// transfer from the EOA locking the assets to the owner.
	args := []interface{}{to, amount}
	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		from,
		zrc20Address,
		big.NewInt(0),
		nil,
		true,
		true,
		transfer,
		args...,
	)
	if err != nil {
		return false, err
	}

	if res.VmError != "" {
		return false, fmt.Errorf("EVM execution error in transfer: %s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[transfer].Outputs.Unpack(res.Ret)
	if err != nil {
		return false, err
	}

	if len(ret) == 0 {
		return false, fmt.Errorf("no data returned from 'transfer' method")
	}

	transferred, ok := ret[0].(bool)
	if !ok {
		return false, fmt.Errorf("transfer returned an unexpected value")
	}

	return transferred, nil
}

// ZRC20TransferFrom transfers ZRC20 tokens from the owner to the spender.
// The transaction is started by the spender.
// This requires the spender to have been approved by the owner.
func (k Keeper) ZRC20TransferFrom(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address, from, to common.Address,
	amount *big.Int,
) (bool, error) {
	if zrc20ABI == nil {
		return false, fungibletypes.ErrZRC20NilABI
	}

	if crypto.IsEmptyAddress(from) || crypto.IsEmptyAddress(to) {
		return false, fungibletypes.ErrZeroAddress
	}

	if err := k.IsValidZRC20(ctx, zrc20Address); err != nil {
		return false, err
	}

	args := []interface{}{from, to, amount}
	res, err := k.CallEVM(
		ctx,
		*zrc20ABI,
		to,
		zrc20Address,
		big.NewInt(0),
		nil,
		true,
		true,
		transferFrom,
		args...,
	)
	if err != nil {
		return false, err
	}

	if res.VmError != "" {
		return false, fmt.Errorf("EVM execution error in transferFrom: %s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[transferFrom].Outputs.Unpack(res.Ret)
	if err != nil {
		return false, err
	}

	if len(ret) == 0 {
		return false, fmt.Errorf("no data returned from 'transferFrom' method")
	}

	transferred, ok := ret[0].(bool)
	if !ok {
		return false, fmt.Errorf("transferFrom returned an unexpected value")
	}

	return transferred, nil
}
