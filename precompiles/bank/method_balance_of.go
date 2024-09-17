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
	zrc20Addr, addr, err := unpackBalanceOfArgs(args)
	if err != nil {
		return nil, err
	}

	// Get the counterpart cosmos address.
	toAddr, err := getCosmosAddress(c.bankKeeper, addr)
	if err != nil {
		return nil, err
	}

	// Bank Keeper GetBalance returns the specified Cosmos coin balance for a given address.
	// Check explicitly the balance is a non-negative non-nil value.
	coin := c.bankKeeper.GetBalance(ctx, toAddr, ZRC20ToCosmosDenom(zrc20Addr))
	if !coin.IsValid() {
		return nil, &ptypes.ErrInvalidCoin{
			Got:      coin.GetDenom(),
			Negative: coin.IsNegative(),
			Nil:      coin.IsNil(),
		}
	}

	return method.Outputs.Pack(coin.Amount.BigInt())
}

func unpackBalanceOfArgs(args []interface{}) (zrc20Addr common.Address, addr common.Address, err error) {
	zrc20Addr, ok := args[0].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, &ptypes.ErrInvalidAddr{
			Got: zrc20Addr.String(),
		}
	}

	addr, ok = args[1].(common.Address)
	if !ok {
		return common.Address{}, common.Address{}, &ptypes.ErrInvalidAddr{
			Got: addr.String(),
		}
	}

	return zrc20Addr, addr, nil
}
