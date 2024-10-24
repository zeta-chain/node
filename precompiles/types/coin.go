package types

import (
	"math/big"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

const ZRC20DenomPrefix = "zrc20/"

// ZRC20ToCosmosDenom returns the cosmos coin address for a given ZRC20 address.
// This is converted to "zevm/{ZRC20Address}".
func ZRC20ToCosmosDenom(ZRC20Address common.Address) string {
	return ZRC20DenomPrefix + ZRC20Address.String()
}

func CreateCoinSet(tokenDenom string, amount *big.Int) (sdk.Coins, error) {
	coin := sdk.NewCoin(tokenDenom, math.NewIntFromBigInt(amount))
	if !coin.IsValid() {
		return nil, &ErrInvalidCoin{
			Got:      coin.GetDenom(),
			Negative: coin.IsNegative(),
			Nil:      coin.IsNil(),
		}
	}

	// A sdk.Coins (type []sdk.Coin) has to be created because it's the type expected by MintCoins
	// and SendCoinsFromModuleToAccount.
	// But sdk.Coins will only contain one coin, always.
	coinSet := sdk.NewCoins(coin)
	if !coinSet.IsValid() || coinSet.Empty() || coinSet.IsAnyNil() || coinSet == nil {
		return nil, &ErrInvalidCoin{
			Got:      coinSet.String(),
			Negative: coinSet.IsAnyNegative(),
			Nil:      coinSet.IsAnyNil(),
			Empty:    coinSet.Empty(),
		}
	}

	return coinSet, nil
}
