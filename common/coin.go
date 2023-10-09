package common

import "strconv"

func GetCoinType(coin string) (CoinType, error) {
	coinInt, err := strconv.ParseInt(coin, 10, 32)
	if err != nil {
		return CoinType_Cmd, err
	}
	return CoinType(coinInt), nil
}
