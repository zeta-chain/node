package sui

// CoinType represents the coin type for the inbound
type CoinType string

const (
	// SUI is the coin type for SUI, native gas token
	SUI CoinType = "0000000000000000000000000000000000000000000000000000000000000002::sui::SUI"

	// SUIShort is the short coin type for SUI
	SUIShort CoinType = "0x2::sui::SUI"
)

// IsSUIType returns true if the given coin type is SUI
func IsSUIType(coinType CoinType) bool {
	return coinType == SUI || coinType == SUIShort
}
