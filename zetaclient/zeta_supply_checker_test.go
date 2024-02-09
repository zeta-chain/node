package zetaclient_test

import (
	"os"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/zetaclient"
)

func MustNewIntFromString(val string) sdkmath.Int {
	v, ok := sdkmath.NewIntFromString(val)
	if !ok {
		panic("invalid int")
	}
	return v
}
func TestZetaSupplyChecker_ValidateZetaSupply(t *testing.T) {
	tt := []struct {
		name                     string
		abortedTxAmount          sdkmath.Int
		zetaInTransit            sdkmath.Int
		genesisAmounts           sdkmath.Int
		externalChainTotalSupply sdkmath.Int
		zetaTokenSupplyOnNode    sdkmath.Int
		ethLockedAmount          sdkmath.Int
		validate                 require.BoolAssertionFunc
	}{
		{
			name:                     "1 zeta cctx in progress",
			abortedTxAmount:          MustNewIntFromString("0"),
			zetaInTransit:            MustNewIntFromString("1000000000000000000"),
			externalChainTotalSupply: MustNewIntFromString("9000000000000000000"),
			genesisAmounts:           MustNewIntFromString("1000000000000000000"),
			zetaTokenSupplyOnNode:    MustNewIntFromString("1000000000000000000"),
			ethLockedAmount:          MustNewIntFromString("10000000000000000000"),
			validate: func(t require.TestingT, b bool, i ...interface{}) {
				require.True(t, b, i...)
			},
		},
		// Todo add more scenarios
		//https://github.com/zeta-chain/node/issues/1375
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
			tc.validate(t, zetaclient.ValidateZetaSupply(logger, tc.abortedTxAmount, tc.zetaInTransit, tc.genesisAmounts, tc.externalChainTotalSupply, tc.zetaTokenSupplyOnNode, tc.ethLockedAmount))
		})
	}
}
