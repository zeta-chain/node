package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

const allowance = "allowance"

// CheckFungibleZRC20Allowance checks if the allowance of ZRC20 tokens,
// is equal or greater than the provided amount.
func (k Keeper) CheckFungibleZRC20Allowance(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	from, zrc20Address common.Address,
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
