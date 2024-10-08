package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	fungibletypes "github.com/zeta-chain/node/x/fungible/types"
)

const balanceOf = "balanceOf"

// CheckFungibleZRC20Balance checks if the balance of ZRC20 tokens,
// is equal or greater than the provided amount.
func (k Keeper) CheckFungibleZRC20Balance(
	ctx sdk.Context,
	zrc20ABI *abi.ABI,
	zrc20Address common.Address,
	amount *big.Int,
) error {
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive, got: %s", amount.String())
	}

	if zrc20Address == zeroAddress {
		return fmt.Errorf("zrc20 address cannot be zero")
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
		return fmt.Errorf("%s", res.VmError)
	}

	ret, err := zrc20ABI.Methods[balanceOf].Outputs.Unpack(res.Ret)
	if err != nil {
		return err
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
