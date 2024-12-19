package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/ptr"
	"github.com/zeta-chain/node/x/observer/types"
)

func TestOperationalFlags_Validate(t *testing.T) {
	tests := []struct {
		name        string
		of          types.OperationalFlags
		errContains string
	}{
		{
			name: "invalid restart height",
			of: types.OperationalFlags{
				RestartHeight: -1,
			},
			errContains: types.ErrOperationalFlagsRestartHeightNegative.Error(),
		},
		{
			name: "valid restart height",
			of: types.OperationalFlags{
				RestartHeight: 1,
			},
		},
		{
			name: "valid signer offset",
			of: types.OperationalFlags{
				SignerBlockTimeOffset: ptr.Ptr(time.Second),
			},
		},
		{
			name: "negative signer offset",
			of: types.OperationalFlags{
				SignerBlockTimeOffset: ptr.Ptr(-time.Second),
			},
			errContains: types.ErrOperationalFlagsSignerBlockTimeOffsetNegative.Error(),
		},
		{
			name: "signer offset limit exceeded",
			of: types.OperationalFlags{
				SignerBlockTimeOffset: ptr.Ptr(time.Minute),
			},
			errContains: types.ErrOperationalFlagsSignerBlockTimeOffsetLimit.Error(),
		},
		{
			name: "all flags valid",
			of: types.OperationalFlags{
				RestartHeight:         1,
				SignerBlockTimeOffset: ptr.Ptr(time.Second),
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
