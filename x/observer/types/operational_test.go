package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestOperationalFlags_Validate(t *testing.T) {
	tests := []struct {
		name        string
		of          types.OperationalFlags
		errContains string
	}{
		{
			name: "invalid operational flags",
			of: types.OperationalFlags{
				RestartHeight: -1,
			},
			errContains: types.ErrOperationalFlagsRestartHeightNegative.Error(),
		},
		{
			name: "valid",
			of: types.OperationalFlags{
				RestartHeight: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.of.Validate()
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}
