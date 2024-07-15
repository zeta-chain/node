package coin

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func AzetaPerZeta() sdk.Dec {
	return sdk.NewDec(1e18)
}

func GetCoinType(coin string) (CoinType, error) {
	coinInt, err := strconv.ParseInt(coin, 10, 32)
	if err != nil {
		return CoinType_Cmd, err
	}
	if coinInt < 0 || coinInt > 3 {
		return CoinType_Cmd, fmt.Errorf("invalid coin type %d", coinInt)
	}
	// #nosec G115 always in range
	return CoinType(coinInt), nil
}

func GetAzetaDecFromAmountInZeta(zetaAmount string) (sdk.Dec, error) {
	zetaDec, err := sdk.NewDecFromStr(zetaAmount)
	if err != nil {
		return sdk.Dec{}, err
	}
	zetaToAzetaConvertionFactor := sdk.NewDecFromInt(sdk.NewInt(1000000000000000000))
	return zetaDec.Mul(zetaToAzetaConvertionFactor), nil
}
