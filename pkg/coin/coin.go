package coin

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
)

func AzetaPerZeta() sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(1e18)
}

func GetCoinType(coin string) (CoinType, error) {
	coinInt, err := strconv.ParseInt(coin, 10, 32)
	if err != nil {
		return CoinType_Cmd, err
	}

	// check boundaries of the enum
	if coinInt < 0 || coinInt > int64(len(CoinType_name)) {
		return CoinType_Cmd, fmt.Errorf("invalid coin type %d", coinInt)
	}

	// #nosec G115 always in range
	return CoinType(coinInt), nil
}

func GetAzetaDecFromAmountInZeta(zetaAmount string) (sdkmath.LegacyDec, error) {
	zetaDec, err := sdkmath.LegacyNewDecFromStr(zetaAmount)
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	zetaToAzetaConvertionFactor := sdkmath.LegacyNewDecFromInt(sdkmath.NewInt(1000000000000000000))
	return zetaDec.Mul(zetaToAzetaConvertionFactor), nil
}

func (c CoinType) SupportsRefund() bool {
	return c == CoinType_ERC20 || c == CoinType_Gas || c == CoinType_Zeta
}
