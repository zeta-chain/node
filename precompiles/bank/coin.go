package bank

import (
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	ptypes "github.com/zeta-chain/node/precompiles/types"
)

// ZRC20ToCosmosDenom returns the cosmos coin address for a given ZRC20 address.
// This is converted to "zevm/{ZRC20Address}".
func ZRC20ToCosmosDenom(ZRC20Address common.Address) string {
	return ZEVMDenom + ZRC20Address.String()
}

func createCoinSet(tokenDenom string, amount *big.Int) (sdk.Coins, error) {
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

	return coinSet, nil
}
