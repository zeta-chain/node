package types

// ZRC4Data represents the ZRC4 token details used to map
// the token to a Cosmos Coin
type ZRC4Data struct {
	Name     string
	Symbol   string
	Decimals uint8
}

// ZRC4StringResponse defines the string value from the call response
type ZRC4StringResponse struct {
	Value string
}

// ZRC4Uint8Response defines the uint8 value from the call response
type ZRC4Uint8Response struct {
	Value uint8
}

// ZRC4BoolResponse defines the bool value from the call response
type ZRC4BoolResponse struct {
	Value bool
}

// ZRC4StringResponse defines the string value from the call response
type UniswapV2FactoryByte32Response struct {
	Value [32]byte
}

// NewZRC4Data creates a new ZRC4Data instance
func NewZRC4Data(name, symbol string, decimals uint8) ZRC4Data {
	return ZRC4Data{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}
}
