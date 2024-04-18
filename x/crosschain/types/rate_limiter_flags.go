package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// Validate checks that the RateLimiterFlags is valid
func (r RateLimiterFlags) Validate() error {
	// window must not be negative
	if r.Window < 0 {
		return fmt.Errorf("window must be positive: %d", r.Window)
	}

	seen := make(map[string]bool)
	for _, conversion := range r.Conversions {
		// check no duplicated conversion
		if _, ok := seen[conversion.Zrc20]; ok {
			return fmt.Errorf("duplicated conversion: %s", conversion.Zrc20)
		}
		seen[conversion.Zrc20] = true

		// check conversion is valid
		if conversion.Rate.IsNil() {
			return fmt.Errorf("rate is nil for conversion: %s", conversion.Zrc20)
		}

		// check address is valid
		if !ethcommon.IsHexAddress(conversion.Zrc20) {
			return fmt.Errorf("invalid zrc20 address (%s)", conversion.Zrc20)
		}
	}

	return nil
}

// GetConversionRate returns the conversion rate for the given zrc20
func (r RateLimiterFlags) GetConversionRate(zrc20 string) (sdk.Dec, bool) {
	for _, conversion := range r.Conversions {
		if conversion.Zrc20 == zrc20 {
			return conversion.Rate, true
		}
	}
	return sdk.NewDec(0), false
}
