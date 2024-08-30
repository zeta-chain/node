package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/x/authority/types"
)

func TestPolicyType_Validate(t *testing.T) {
	tests := []struct {
		name       string
		policyType types.PolicyType
		wantErr    bool
	}{
		{
			name:       "valid groupEmergency",
			policyType: types.PolicyType_groupEmergency,
			wantErr:    false,
		},
		{
			name:       "valid groupOperational",
			policyType: types.PolicyType_groupOperational,
			wantErr:    false,
		},
		{
			name:       "valid groupAdmin",
			policyType: types.PolicyType_groupAdmin,
			wantErr:    false,
		},
		{
			name:       "invalid policy type",
			policyType: types.PolicyType(20),
			wantErr:    true,
		},
		{
			name:       "invalid policy type more than max length",
			policyType: types.PolicyType(len(types.PolicyType_name) + 1),
			wantErr:    true,
		},
		{
			name:       "empty policy type",
			policyType: types.PolicyType_groupEmpty,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policyType.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
