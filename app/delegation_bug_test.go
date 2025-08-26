package app_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/node/testutil/sample"
)

// NOTE: File created in an arbitrary location

func TestDelegationBug(t *testing.T) {
	// Values from zeta1j6c9cgutn9mg7dy3xh0qv5aksluuz0lwx4wncc to zetavaloper19pgg7htc67rhrk6jpvutpxwsrm8g3654jmrrye on mainnet
	{
		validator := sample.Validator(t, sample.Rand())
		validator.DelegatorShares = sdkmath.LegacyMustNewDecFromStr("303282111200117915174803.285525143883122398")

		tokens, ok := sdkmath.NewIntFromString("302373254591626134909587")
		require.True(t, ok)
		validator.Tokens = tokens
		result := validator.TokensFromShares(sdkmath.LegacyMustNewDecFromStr("1003005.744042141715136737"))
		resultTruncate := result.TruncateInt()
		resultRoundUp := validator.TokensFromSharesRoundUp(sdkmath.LegacyMustNewDecFromStr("1003005.744042141715136737"))

		fmt.Println(result.String())
		fmt.Println(resultTruncate.String()) // return 999999 while 1000000 tokens were delegated
		fmt.Println(resultRoundUp.String())
	}

	// Trying to reproduce with simpler values
	//{
	//	validator := sample.Validator(t, sample.Rand())
	//	validator.DelegatorShares = sdkmath.LegacyMustNewDecFromStr("30000600000000000.0000000000000003")
	//
	//	validator.Tokens = sdkmath.NewInt(30000000000000000)
	//	result := validator.TokensFromShares(sdkmath.LegacyMustNewDecFromStr("10000200000000000.0000000000000001")).TruncateInt()
	//	fmt.Println(result.String())
	//}
}
