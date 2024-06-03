package cli_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/x/authority/client/cli"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

func Test_GetPolicyType(t *testing.T) {
	tt := []struct {
		name                 string
		policyTypeString     string
		expecectedPolicyType types.PolicyType
		expecterErrorString  string
	}{
		{
			name:                 "groupEmergency",
			policyTypeString:     "0",
			expecectedPolicyType: types.PolicyType_groupEmergency,
			expecterErrorString:  "",
		},
		{
			name:                 "groupOperational",
			policyTypeString:     "1",
			expecectedPolicyType: types.PolicyType_groupOperational,
			expecterErrorString:  "",
		},
		{
			name:                 "groupAdmin",
			policyTypeString:     "2",
			expecectedPolicyType: types.PolicyType_groupAdmin,
			expecterErrorString:  "",
		},
		{
			name:                 "groupEmpty",
			policyTypeString:     "3",
			expecectedPolicyType: types.PolicyType_groupEmpty,
			expecterErrorString:  "invalid policy type value",
		},
		{
			name:                 "string literal for policy type not accepted",
			policyTypeString:     "groupEmergency",
			expecectedPolicyType: types.PolicyType_groupEmpty,
			expecterErrorString:  "failed to parse policy type",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			policyType, err := cli.GetPolicyType(tc.policyTypeString)
			require.Equal(t, tc.expecectedPolicyType, policyType)
			if tc.expecectedPolicyType == types.PolicyType_groupEmpty {
				require.ErrorContains(t, err, tc.expecterErrorString)
			}
		})
	}

}
