package math

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPercentage(t *testing.T) {
	testCases := []struct {
		name        string
		numerator   *big.Int
		denominator *big.Int
		percentage  *big.Float
		fail        bool
	}{
		{
			name:        "positive percentage",
			numerator:   big.NewInt(165),
			denominator: big.NewInt(1000),
			percentage:  big.NewFloat(16.5),
			fail:        false,
		},
		{
			name:        "negative percentage",
			numerator:   big.NewInt(-165),
			denominator: big.NewInt(1000),
			percentage:  big.NewFloat(-16.5),
			fail:        false,
		},
		{
			name:        "zero denominator",
			numerator:   big.NewInt(1),
			denominator: big.NewInt(0),
			percentage:  nil,
			fail:        true,
		},
		{
			name:        "nil numerator",
			numerator:   nil,
			denominator: big.NewInt(1000),
			percentage:  nil,
			fail:        true,
		},
		{
			name:        "nil denominator",
			numerator:   big.NewInt(165),
			denominator: nil,
			percentage:  nil,
			fail:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			percentage := Percentage(tc.numerator, tc.denominator)
			fmt.Printf("percentage: %v\n", percentage)
			if tc.fail {
				require.Nil(t, percentage)
			} else {
				require.True(t, percentage.Cmp(tc.percentage) == 0)
			}
		})
	}
}
