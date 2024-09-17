package bank

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	ptypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

// From IBank.sol: function deposit(address zrc20, uint256 amount) external returns (bool success);
func (c *Contract) deposit(
	ctx sdk.Context,
	evm *vm.EVM,
	contract *vm.Contract,
	method *abi.Method,
	args []interface{},
) (result []byte, err error) {
	// This function is developed using the Check - Effects - Interactions pattern:
	// 1. Check everything is correct.
	if len(args) != 2 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// Unpack parameters for function deposit.
	// function deposit(address zrc20, uint256 amount) external returns (bool success);
	zrc20Addr, amount, err := unpackDepositArgs(args)
	if err != nil {
		return nil, err
	}

	// Get the correct caller address.
	caller, err := getEVMCallerAddress(evm, contract)
	if err != nil {
		return nil, err
	}

	// Get the cosmos address of the caller.
	toAddr, err := getCosmosAddress(c.bankKeeper, caller)
	if err != nil {
		return nil, err
	}

	// Check for enough balance.
	// function balanceOf(address account) public view virtual override returns (uint256)
	resBalanceOf, err := c.CallContract(
		ctx,
		&c.fungibleKeeper,
		c.zrc20ABI,
		zrc20Addr,
		"balanceOf",
		[]interface{}{caller},
	)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "balanceOf",
			Got:  err.Error(),
		}
	}

	balance, ok := resBalanceOf[0].(*big.Int)
	if !ok || balance.Cmp(amount) < 0 {
		return nil, &ptypes.ErrInvalidAmount{
			Got: "not enough balance",
		}
	}

	// Check for enough bank's allowance.
	// function allowance(address owner, address spender) public view virtual override returns (uint256)
	resAllowance, err := c.CallContract(
		ctx,
		&c.fungibleKeeper,
		c.zrc20ABI,
		zrc20Addr,
		"allowance",
		[]interface{}{caller, ContractAddress},
	)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "allowance",
			Got:  err.Error(),
		}
	}

	allowance, ok := resAllowance[0].(*big.Int)
	if !ok || allowance.Cmp(amount) < 0 {
		return nil, &ptypes.ErrInvalidAmount{
			Got: "not enough allowance",
		}
	}

	// The process of creating a new cosmos coin is:
	// - Generate the new coin denom using ZRC20 address,
	//   this way we map ZRC20 addresses to cosmos denoms "zevm/0x12345".
	// - Mint coins.
	// - Send coins to the caller.
	coinSet, err := createCoinSet(ZRC20ToCosmosDenom(zrc20Addr), amount)
	if err != nil {
		return nil, err
	}

	// 2. Effect: subtract balance.
	// function transferFrom(address sender, address recipient, uint256 amount) public virtual override returns (bool)
	resTransferFrom, err := c.CallContract(
		ctx,
		&c.fungibleKeeper,
		c.zrc20ABI,
		zrc20Addr,
		"transferFrom",
		[]interface{}{caller, ContractAddress, amount},
	)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "transferFrom",
			Got:  err.Error(),
		}
	}

	transferred, ok := resTransferFrom[0].(bool)
	if !ok || !transferred {
		return nil, &ptypes.ErrUnexpected{
			When: "transferFrom",
			Got:  "transaction not successful",
		}
	}

	// 3. Interactions: create cosmos coin and send.
	err = c.bankKeeper.MintCoins(ctx, types.ModuleName, coinSet)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "MintCoins",
			Got:  err.Error(),
		}
	}

	err = c.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coinSet)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "SendCoinsFromModuleToAccount",
			Got:  err.Error(),
		}
	}

	if err := c.AddDepositLog(ctx, evm.StateDB, caller, zrc20Addr, ZRC20ToCosmosDenom(zrc20Addr), amount); err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "AddDepositLog",
			Got:  err.Error(),
		}
	}

	return method.Outputs.Pack(true)
}

func unpackDepositArgs(args []interface{}) (zrc20Addr common.Address, amount *big.Int, err error) {
	zrc20Addr, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, nil, &ptypes.ErrInvalidAddr{
			Got: zrc20Addr.String(),
		}
	}

	amount, ok = args[1].(*big.Int)
	if !ok || amount.Sign() < 0 || amount == nil || amount == new(big.Int) {
		return common.Address{}, nil, &ptypes.ErrInvalidAmount{
			Got: amount.String(),
		}
	}

	return zrc20Addr, amount, nil
}
