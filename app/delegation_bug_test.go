package app_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/node/testutil/sample"
)

// NOTE: File created in an arbitrary location

func TestDelegationBug(t *testing.T) {
	validator := sample.Validator(t, sample.Rand())
	validator.DelegatorShares = sdkmath.LegacyNewDec(3000060)
	validator.Tokens = sdkmath.NewInt(3000000)
	result := validator.TokensFromShares(sdkmath.LegacyNewDec(1000020)).TruncateInt()
	fmt.Println(result.String())
}
