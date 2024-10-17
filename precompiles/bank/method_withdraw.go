package bank

import (
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	ptypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// withdraw is used to withdraw cosmos coins minted using the bank's deposit function.
// The caller has to have enough cosmos coin on its cosmos account balance to withdraw the requested amount.
// After all check pass the bank will burn the cosmos coins and transfer the ZRC20 amount to the withdrawer.
// The cosmos coins have the denomination of "zrc20/0x12345" where 0x12345 is the ZRC20 address.
// Call this function using solidity with the following signature:
// From IBank.sol: function withdraw(address zrc20, uint256 amount) external returns (bool success);
// The address to be passed to the function is the ZRC20 address, like in 0x12345.
func (c *Contract) withdraw(
	ctx sdk.Context,
	evm *vm.EVM,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) (result []byte, err error) {
	// 1. Check everything is correct.
	if len(args) != 2 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// Unpack parameters for function withdraw.
	// function withdraw(address zrc20, uint256 amount) external returns (bool success);
	zrc20Addr, amount, err := unpackWithdrawArgs(args)
	if err != nil {
		return nil, err
	}

	// Get the correct caller address.
	caller, err := ptypes.GetEVMCallerAddress(evm, contract)
	if err != nil {
		return nil, err
	}

	// Get the cosmos address of the caller.
	// This address should have enough cosmos coin balance as the requested amount.
	fromAddr, err := ptypes.GetCosmosAddress(c.bankKeeper, caller)
	if err != nil {
		return nil, err
	}

	// Safety check: token has to be a non-paused whitelisted ZRC20.
	if err := c.fungibleKeeper.IsValidZRC20(ctx, zrc20Addr); err != nil {
		return nil, &ptypes.ErrInvalidToken{
			Got:    zrc20Addr.String(),
			Reason: err.Error(),
		}
	}

	// Caller has to have enough cosmos coin balance to withdraw the requested amount.
	coin := c.bankKeeper.GetBalance(ctx, fromAddr, ptypes.ZRC20ToCosmosDenom(zrc20Addr))
	if !coin.IsValid() {
		return nil, &ptypes.ErrInsufficientBalance{
			Requested: amount.String(),
			Got:       "invalid coin",
		}
	}

	if coin.Amount.LT(math.NewIntFromBigInt(amount)) {
		return nil, &ptypes.ErrInsufficientBalance{
			Requested: amount.String(),
			Got:       coin.Amount.String(),
		}
	}

	coinSet, err := ptypes.CreateCoinSet(ptypes.ZRC20ToCosmosDenom(zrc20Addr), amount)
	if err != nil {
		return nil, err
	}

	// Check if bank address has enough ZRC20 balance.
	if err := c.fungibleKeeper.CheckZRC20Balance(ctx, c.zrc20ABI, zrc20Addr, c.Address(), amount); err != nil {
		return nil, &ptypes.ErrInsufficientBalance{
			Requested: amount.String(),
			Got:       err.Error(),
		}
	}

	// 2. Effect: burn cosmos coin balance.
	if err := c.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAddr, types.ModuleName, coinSet); err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "SendCoinsFromAccountToModule",
			Got:  err.Error(),
		}
	}

	if err := c.bankKeeper.BurnCoins(ctx, types.ModuleName, coinSet); err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "BurnCoins",
			Got:  err.Error(),
		}
	}

	// 3. Interactions: send ZRC20.
	if err := c.fungibleKeeper.UnlockZRC20(ctx, c.zrc20ABI, zrc20Addr, caller, c.Address(), amount); err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "UnlockZRC20InBank",
			Got:  err.Error(),
		}
	}

	if err := c.addEventLog(ctx, evm.StateDB, WithdrawEventName, eventData{caller, zrc20Addr, fromAddr.String(), coinSet.Denoms()[0], amount}); err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "AddWithdrawLog",
			Got:  err.Error(),
		}
	}

	return method.Outputs.Pack(true)
}

func unpackWithdrawArgs(args []interface{}) (zrc20Addr common.Address, amount *big.Int, err error) {
	zrc20Addr, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, nil, &ptypes.ErrInvalidAddr{
			Got: zrc20Addr.String(),
		}
	}

	amount, ok = args[1].(*big.Int)
	if !ok || amount == nil || amount.Sign() <= 0 {
		return common.Address{}, nil, &ptypes.ErrInvalidAmount{
			Got: amount.String(),
		}
	}

	return zrc20Addr, amount, nil
}
