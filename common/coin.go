package common

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetCoinType(coin string) (CoinType, error) {
	coinInt, err := strconv.ParseInt(coin, 10, 32)
	if err != nil {
		return CoinType_Cmd, err
	}
	if coinInt < 0 || coinInt > 3 {
		return CoinType_Cmd, fmt.Errorf("invalid coin type %d", coinInt)
	}
	// #nosec G701 always in range
	return CoinType(coinInt), nil
}

func GetAzetaDecFromAmountInZeta(zetaAmount string) (sdk.Dec, error) {
	zetaDec, err := sdk.NewDecFromStr(zetaAmount)
	if err != nil {
		return sdk.Dec{}, err
	}
	zetaToAzetaConvertionFactor, err := sdk.NewDecFromStr("1000000000000000000")
	if err != nil {
		return sdk.Dec{}, err
	}
	return zetaDec.Mul(zetaToAzetaConvertionFactor), nil
}
