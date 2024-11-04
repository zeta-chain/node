package types

import (
	"math/big"
	"strings"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/zeta-chain/node/cmd/zetacored/config"
)

// ZRC20ToCosmosDenom returns the cosmos coin address for a given ZRC20 address.
// This is converted to "zrc20/{ZRC20Address}".
func ZRC20ToCosmosDenom(ZRC20Address common.Address) string {
	return config.ZRC20DenomPrefix + ZRC20Address.String()
}

func CreateCoinSet(zrc20address common.Address, amount *big.Int) (sdk.Coins, error) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	denom := ZRC20ToCosmosDenom(zrc20address)

	coin := sdk.NewCoin(denom, math.NewIntFromBigInt(amount))
	if !coin.IsValid() {
		return nil, &ErrInvalidCoin{
			Got:      coin.GetDenom(),
			Negative: coin.IsNegative(),
			Nil:      coin.IsNil(),
		}
	}

	// A sdk.Coins (type []sdk.Coin) has to be created because it's the type expected by MintCoins
	// and SendCoinsFromModuleToAccount.
	// But coinSet will only contain one coin, always.
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

// CoinIsZRC20 checks if a given coin is a ZRC20 coin based on its denomination.
func CoinIsZRC20(denom string) bool {
	return strings.HasPrefix(denom, config.ZRC20DenomPrefix)
}
