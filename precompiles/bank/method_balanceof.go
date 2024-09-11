package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	ptypes "github.com/zeta-chain/node/precompiles/types"
)

func (c *Contract) balanceOf(
	ctx sdk.Context,
	method *abi.Method,
	args []interface{},
) (result []byte, err error) {
	if len(args) != 2 {
		return nil, &(ptypes.ErrInvalidNumberOfArgs{
			Got:    len(args),
			Expect: 2,
		})
	}

	// function balanceOf(address zrc20, address user) external view returns (uint256 balance);
	tokenAddr, addr := args[0].(common.Address), args[1].(common.Address)

	// common.Address has to be converted to AccAddress.
	accAddr := sdk.AccAddress(addr.Bytes())
	if accAddr.Empty() {
		return nil, &ptypes.ErrInvalidAddr{
			Got: accAddr.String(),
		}
	}

	// Convert ZRC20 address to a Cosmos denom formatted as "zevm/0x12345".
	tokenDenom := ZRC20ToCosmosDenom(tokenAddr)

	// Bank Keeper GetBalance returns the specified Cosmos coin balance for a given address.
	// Check explicitly the balance is a non-negative non-nil value.
	coin := c.bankKeeper.GetBalance(ctx, accAddr, tokenDenom)
	if !coin.IsValid() {
		return nil, &ptypes.ErrInvalidCoin{
			Got:      coin.GetDenom(),
			Negative: coin.IsNegative(),
			Nil:      coin.IsNil(),
		}
	}

	return method.Outputs.Pack(coin.Amount.BigInt())
}
