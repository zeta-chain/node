package common

import "strconv"

func GetCoinType(coin string) CoinType {
	coinInt, err := strconv.ParseInt(coin, 10, 64)
	if err != nil {
		panic(err)
	}
	return CoinType(coinInt)
}
