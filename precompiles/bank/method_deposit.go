package bank

import (
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
	// This function is developed using the
	// Check - Effects - Interactions pattern:
	// 1. Check everything is correct.
	if len(args) != 2 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// function deposit(address zrc20, uint256 amount) external returns (bool success);
	ZRC20Addr, amount := args[0].(common.Address), args[1].(*big.Int)
	if amount.Sign() < 0 || amount == nil || amount == new(big.Int) {
		return nil, &ptypes.ErrInvalidAmount{
			Got: amount.String(),
		}
	}

	// Initialize the ZRC20 ABI, as we need to call the balanceOf and allowance methods.
	ZRC20ABI, err := zrc20.ZRC20MetaData.GetAbi()
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "ZRC20MetaData.GetAbi",
			Got:  err.Error(),
		}
	}

	// Check for enough balance.
	// function balanceOf(address account) public view virtual override returns (uint256)
	argsBalanceOf := []interface{}{caller}

	resBalanceOf, err := c.CallContract(ctx, ZRC20ABI, ZRC20Addr, "balanceOf", argsBalanceOf)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "balanceOf",
			Got:  err.Error(),
		}
	}

	balance := resBalanceOf[0].(*big.Int)
	if balance.Cmp(amount) < 0 {
		return nil, &ptypes.ErrUnexpected{
			When: "balance0f",
			Got:  "not enough balance",
		}
	}

	// Check for enough allowance.
	// function allowance(address owner, address spender) public view virtual override returns (uint256)
	argsAllowance := []interface{}{caller, ContractAddress}

	resAllowance, err := c.CallContract(ctx, ZRC20ABI, ZRC20Addr, "allowance", argsAllowance)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "allowance",
			Got:  err.Error(),
		}
	}

	allowance := resAllowance[0].(*big.Int)
	if allowance.Cmp(amount) < 0 {
		return nil, &ptypes.ErrUnexpected{
			When: "allowance",
			Got:  "not enough allowance",
		}
	}

	// Handle the toAddr:
	// check it's valid and not blocked.
	toAddr := sdk.AccAddress(caller.Bytes())
	if toAddr.Empty() {
		return nil, &ptypes.ErrInvalidAddr{
			Got:    toAddr.String(),
			Reason: "empty address",
		}
	}

	if c.bankKeeper.BlockedAddr(toAddr) {
		return nil, &ptypes.ErrInvalidAddr{
			Got:    toAddr.String(),
			Reason: "blocked by bank keeper",
		}
	}

	// The process of creating a new cosmos coin is:
	// - Generate the new coin denom using ZRC20 address,
	//   this way we map ZRC20 addresses to cosmos denoms "zevm/0x12345".
	// - Mint coins.
	// - Send coins to the caller.
	tokenDenom := ZRC20ToCosmosDenom(ZRC20Addr)
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

	if !c.bankKeeper.IsSendEnabledCoin(ctx, coin) {
		return nil, &ptypes.ErrUnexpected{
			When: "IsSendEnabledCoins",
			Got:  "coin not enabled to be sent",
		}
	}

	// 2. Effect: subtract balance.
	// function transferFrom(address sender, address recipient, uint256 amount) public virtual override returns (bool)
	argsTransferFrom := []interface{}{caller, ContractAddress, amount}

	resTransferFrom, err := c.CallContract(ctx, ZRC20ABI, ZRC20Addr, "transferFrom", argsTransferFrom)
	if err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "transferFrom",
			Got:  err.Error(),
		}
	}

	transferred := resTransferFrom[0].(bool)
	if !transferred {
		return nil, &ptypes.ErrUnexpected{
			When: "transferFrom",
			Got:  "transaction not successful",
		}
	}

	// 3. Interactions: create cosmos coin and send.
	if err := c.bankKeeper.MintCoins(ctx, types.ModuleName, coinSet); err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "MintCoins",
			Got:  err.Error(),
		}
	}

	if err := c.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAddr, coinSet); err != nil {
		return nil, &ptypes.ErrUnexpected{
			When: "SendCoinsFromModuleToAccount",
			Got:  err.Error(),
		}
	}

	return method.Outputs.Pack(true)
}
