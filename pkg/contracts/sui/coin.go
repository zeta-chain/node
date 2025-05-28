package sui

// CoinType represents the coin type for the SUI token
type CoinType string

const (
	// SUI is the coin type for SUI, native gas token
	SUI CoinType = "0000000000000000000000000000000000000000000000000000000000000002::sui::SUI"

	// SUIShort is the short coin type for SUI
	SUIShort CoinType = "0x2::sui::SUI"

	// MISTPerSUI is the number of mist in one SUI
	MistPerSUI = uint64(1e9)
)

// IsSUICoinType returns true if the given coin type is SUI
func IsSUICoinType(coinType CoinType) bool {
	return coinType == SUI || coinType == SUIShort
}
