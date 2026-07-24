package snapshot

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/require"
)

func TestDBRecords(t *testing.T) {
	// ARRANGE: one claimable account (2.5 ZETA) + remainder (1.5 ZETA)
	res := &Result{
		Supply:    sdkmath.NewIntWithDecimal(4, 18),                 // 4 ZETA
		Remainder: sdkmath.NewIntWithDecimal(15, 17),                // 1.5 ZETA
		Accounts: []Account{{
			Canonical: "0x" + "01",
			Class:     classEOA,
			Liquid:    sdkmath.NewIntWithDecimal(25, 17), // 2.5 ZETA
			Staked:    sdkmath.ZeroInt(),
			Unbonding: sdkmath.ZeroInt(),
		}},
	}

	// ACT
	records := res.DBRecords()

	// ASSERT: header + one claimable + one remainder row
	require.Equal(t, DBHeader, records[0])
	require.Len(t, records, 3)

	claimable := records[1]
	require.Equal(t, "2500000000000000000", claimable[5], "total_azeta")
	require.Equal(t, "2500000000", claimable[6], "total_9dec floors 18->9 decimals")
	require.Equal(t, claimStatusClaimable, claimable[7])

	remainder := records[2]
	require.Equal(t, classRemainder, remainder[1])
	require.Equal(t, "1500000000000000000", remainder[5])
	require.Equal(t, claimStatusMultisig, remainder[7])

	// re-summing total_azeta across all rows equals supply
	sum := sdkmath.ZeroInt()
	for _, row := range records[1:] {
		v, ok := sdkmath.NewIntFromString(row[5])
		require.True(t, ok)
		sum = sum.Add(v)
	}
	require.Equal(t, res.Supply.String(), sum.String())

	// SPL cap = floored per-row totals
	require.Equal(t, "4000000000", res.SPLCap().String())
}
