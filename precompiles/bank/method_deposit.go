package bank

import (
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/protocol-contracts/v2/pkg/zrc20.sol"

	ptypes "github.com/zeta-chain/node/precompiles/types"
	"github.com/zeta-chain/node/x/fungible/types"
)

func (c *Contract) deposit(
	ctx sdk.Context,
	method *abi.Method,
	caller common.Address,
	args []interface{},
) (result []byte, err error) {
	fmt.Printf("DEBUG: deposit()\n")
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
	zrc20Addr, ok := args[0].(common.Address)
	if !ok {
		return nil, &ptypes.ErrInvalidAddr{
			Got: zrc20Addr.String(),
		}
	}

	amount, ok := args[1].(*big.Int)
	if !ok || amount.Sign() < 0 || amount == nil || amount == new(big.Int) {
		return nil, &ptypes.ErrInvalidAmount{
			Got: amount.String(),
		}
	}
	fmt.Printf("DEBUG: deposit(): zrc20Addr (ERC20ZRC20) %s\n", zrc20Addr.String())
	fmt.Printf("DEBUG: deposit(): caller %s\n", caller.String())

	// Initialize the ZRC20 ABI, as we need to call the balanceOf and allowance methods.
	zrc20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		fmt.Printf("DEBUG: deposit(): zrc20.ZRC20MetaData.GetAbi() error %s\n", err.Error())
		return nil, &ptypes.ErrUnexpected{
			When: "ZRC20MetaData.GetAbi()",
			Got:  err.Error(),
		}
	}
	fmt.Printf("DEBUG: deposit(): zrc20ABI %v\n", zrc20ABI)

	// Check for enough balance.
	// function balanceOf(address account) public view virtual override returns (uint256)
	resBalanceOf, err := c.CallContract(
		ctx,
		&c.fungibleKeeper,
		zrc20ABI,
		ContractAddress,
		zrc20Addr,
		"balanceOf",
		true,
		[]interface{}{caller},
	)
	if err != nil {
		fmt.Printf("DEBUG: deposit(): balanceOf c.CallContract error %s\n", err.Error())
		return nil, &ptypes.ErrUnexpected{
			When: "balanceOf",
			Got:  err.Error(),
		}
	}
	fmt.Printf("DEBUG: deposit(): resBalanceOf %v\n", resBalanceOf)

	balance, ok := resBalanceOf[0].(*big.Int)
	if !ok || balance.Cmp(amount) < 0 {
		return nil, &ptypes.ErrInvalidAmount{
			Got: "not enough balance",
		}
	}
	fmt.Printf("DEBUG: deposit(): balanceOf caller %v\n", balance.Uint64())

	// Check for enough allowance.
	// function allowance(address owner, address spender) public view virtual override returns (uint256)
	resAllowance, err := c.CallContract(
		ctx,
		&c.fungibleKeeper,
		zrc20ABI,
		ContractAddress,
		zrc20Addr,
		"allowance",
		true,
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
	fmt.Printf("DEBUG: deposit(): allowance caller %v\n", allowance.Uint64())

	// Handle the toAddr:
	// check it's valid and not blocked.
	toAddr := sdk.AccAddress(caller.Bytes())
	if toAddr.Empty() {
		return nil, &ptypes.ErrInvalidAddr{
			Got:    toAddr.String(),
			Reason: "empty address",
		}
	}
	fmt.Printf("DEBUG: deposit(): caller toAddr %s\n", toAddr.String())

	if c.bankKeeper.BlockedAddr(toAddr) {
		return nil, &ptypes.ErrInvalidAddr{
			Got:    toAddr.String(),
			Reason: "destination address blocked by bank keeper",
		}
	}

	// The process of creating a new cosmos coin is:
	// - Generate the new coin denom using ZRC20 address,
	//   this way we map ZRC20 addresses to cosmos denoms "zevm/0x12345".
	// - Mint coins.
	// - Send coins to the caller.
	tokenDenom := ZRC20ToCosmosDenom(zrc20Addr)

	coin := sdk.NewCoin(tokenDenom, math.NewIntFromBigInt(amount))
	if !coin.IsValid() {
		return nil, &ptypes.ErrInvalidCoin{
			Got:      coin.GetDenom(),
			Negative: coin.IsNegative(),
			Nil:      coin.IsNil(),
		}
	}

	// A sdk.Coins (type []sdk.Coin) has to be created because it's the type expected by MintCoins
	// and SendCoinsFromModuleToAccount.
	// But sdk.Coins will only contain one coin, always.
	coinSet := sdk.NewCoins(coin)
	if !coinSet.IsValid() {
		return nil, &ptypes.ErrInvalidCoin{
			Got:      coinSet.Sort().GetDenomByIndex(0),
			Negative: coinSet.IsAnyNegative(),
			Nil:      coinSet.IsAnyNil(),
		}
	}

	// 2. Effect: subtract balance.
	// function transferFrom(address sender, address recipient, uint256 amount) public virtual override returns (bool)
	resTransferFrom, err := c.CallContract(
		ctx,
		&c.fungibleKeeper,
		zrc20ABI,
		ContractAddress,
		zrc20Addr,
		"transferFrom",
		true,
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
	fmt.Printf("DEBUG: deposit(): transferred %v\n", transferred)

	// 3. Interactions: create cosmos coin and send.
	err = c.bankKeeper.MintCoins(ctx, types.ModuleName, coinSet)
	if err != nil {
		fmt.Printf("DEBUG: deposit(): MintCoins error %s\n", err.Error())
		return nil, &ptypes.ErrUnexpected{
			When: "MintCoins",
			Got:  err.Error(),
		}
	}
	fmt.Printf("DEBUG: deposit(): MintCoins finished\n")

	err = c.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coinSet)
	if err != nil {
		fmt.Printf("DEBUG: deposit(): SendCoinsFromModuleToAccount error %s\n", err.Error())
		return nil, &ptypes.ErrUnexpected{
			When: "SendCoinsFromModuleToAccount",
			Got:  err.Error(),
		}
	}
	fmt.Printf("DEBUG: deposit(): SendCoinsFromModuleToAccount finished\n")

	return method.Outputs.Pack(true)
}
