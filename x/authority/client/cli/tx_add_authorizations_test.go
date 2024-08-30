package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/x/authority/client/cli"
	"github.com/zeta-chain/node/x/authority/types"
)

func Test_GetPolicyType(t *testing.T) {
	tt := []struct {
		name                string
		policyTypeString    string
		expectedPolicyType  types.PolicyType
		expecterErrorString string
	}{
		{
			name:                "groupEmergency",
			policyTypeString:    "0",
			expectedPolicyType:  types.PolicyType_groupEmergency,
			expecterErrorString: "",
		},
		{
			name:                "groupOperational",
			policyTypeString:    "1",
			expectedPolicyType:  types.PolicyType_groupOperational,
			expecterErrorString: "",
		},
		{
			name:                "groupAdmin",
			policyTypeString:    "2",
			expectedPolicyType:  types.PolicyType_groupAdmin,
			expecterErrorString: "",
		},
		{
			name:                "groupEmpty",
			policyTypeString:    "3",
			expectedPolicyType:  types.PolicyType_groupEmpty,
			expecterErrorString: "invalid policy type value",
		},
		{
			name:                "string literal for policy type not accepted",
			policyTypeString:    "groupEmergency",
			expectedPolicyType:  types.PolicyType_groupEmpty,
			expecterErrorString: "failed to parse policy type",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			policyType, err := cli.GetPolicyType(tc.policyTypeString)
			require.Equal(t, tc.expectedPolicyType, policyType)
			if tc.expectedPolicyType == types.PolicyType_groupEmpty {
				require.ErrorContains(t, err, tc.expecterErrorString)
			}
		})
	}

}
