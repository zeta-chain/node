package common_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
)

func Test_GetAzetaDecFromAmountInZeta(t *testing.T) {
	tt := []struct {
		name       string
		zetaAmount string
	}{
		{
			name:       "valid zeta amount",
			zetaAmount: "210000000",
		},
		{
			name:       "very high zeta amount",
			zetaAmount: "21000000000000000000",
		},
		{
			name:       "very low zeta amount",
			zetaAmount: "1",
		},
		{
			name:       "zero zeta amount",
			zetaAmount: "0",
		},
		{
			name:       "decimal zeta amount",
			zetaAmount: "0.1",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := common.GetAzetaDecFromAmountInZeta(tc.zetaAmount)
			assert.NoError(t, err)
		})
	}

}
