package app_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/zeta-chain/node/testutil/sample"
)

// NOTE: File created in an arbitrary location

func TestDelegationBug(t *testing.T) {
	// Values from zeta1j6c9cgutn9mg7dy3xh0qv5aksluuz0lwx4wncc to zetavaloper19pgg7htc67rhrk6jpvutpxwsrm8g3654jmrrye on mainnet
	{
		totalShares := sdkmath.LegacyMustNewDecFromStr("303282111200117915174803.285525143883122398")
		delShares := sdkmath.LegacyMustNewDecFromStr("1003005.744042141715136737")

		tokens, ok := sdkmath.NewIntFromString("302373254591626134909587")
		require.True(t, ok)

		validator := sample.Validator(t, sample.Rand())
		validator.DelegatorShares = totalShares
		validator.Tokens = tokens

		// https://github.com/cosmos/cosmos-sdk/blob/79fcc30f7eb642b028ebe6bf1a6fa29334c7de75/x/staking/keeper/grpc_query.go#L597C26-L597C72
		// the calculation in the query returns 999999 while 1000000 tokens were delegated
		result := validator.TokensFromShares(delShares).TruncateInt()
		require.EqualValues(t, int64(999999), result.Int64())

		// using TokensFromSharesRoundUp fixes the issue
		resultRoundUp := validator.TokensFromSharesRoundUp(delShares).TruncateInt()
		require.EqualValues(t, int64(1000000), resultRoundUp.Int64())

		// https://github.com/cosmos/cosmos-sdk/blob/79fcc30f7eb642b028ebe6bf1a6fa29334c7de75/x/staking/keeper/delegation.go#L1379
		// this is how the max number of token that can be undelegated is calculated
		sharesTruncated, err := validator.SharesFromTokensTruncated(sdkmath.NewInt(1000000))
		require.NoError(t, err)

		// if sharesTruncated.GT(delShares) undelegate message fail
		// if false it means 1000000 tokens can be undelegated
		require.False(t, sharesTruncated.GT(delShares))
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
