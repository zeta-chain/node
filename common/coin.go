package common

import (
	"fmt"
	"strconv"
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
