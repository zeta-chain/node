package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/sdkconfig"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/authority/types"
)

func TestPolicies_Validate(t *testing.T) {
	sdkconfig.SetDefault(false)
	// use table driven tests to test the validation of policies
	tests := []struct {
		name        string
		policies    types.Policies
		errContains string
	}{
		{
			name:        "empty is valid",
			policies:    types.Policies{},
			errContains: "",
		},
		{
			name:        "default is valid",
			policies:    types.DefaultPolicies(),
			errContains: "",
		},
		{
			name: "policies with all group",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    sample.AccAddress(),
						PolicyType: types.PolicyType_groupEmergency,
					},
					{
						Address:    sample.AccAddress(),
						PolicyType: types.PolicyType_groupAdmin,
					},
					{
						Address:    sample.AccAddress(),
						PolicyType: types.PolicyType_groupOperational,
					},
				},
			},
			errContains: "",
		},
		{
			name: "valid if a policy type is not existing",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    sample.AccAddress(),
						PolicyType: types.PolicyType_groupEmergency,
					},
				},
			},
			errContains: "",
		},
		{
			name: "invalid if address is invalid",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    "invalid",
						PolicyType: types.PolicyType_groupEmergency,
					},
				},
			},
			errContains: "invalid address",
		},
		{
			name: "invalid if policy type is invalid",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    sample.AccAddress(),
						PolicyType: types.PolicyType(1000),
					},
				},
			},
			errContains: "invalid policy type",
		},
		{
			name: "invalid if duplicated policy type",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    sample.AccAddress(),
						PolicyType: types.PolicyType_groupEmergency,
					},
					{
						Address:    sample.AccAddress(),
						PolicyType: types.PolicyType_groupAdmin,
					},
					{
						Address:    sample.AccAddress(),
						PolicyType: types.PolicyType_groupEmergency,
					},
				},
			},
			errContains: "duplicate policy type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policies.Validate()
			if tt.errContains != "" {
				require.ErrorContains(t, err, tt.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPolicies_CheckSigner(t *testing.T) {
	signer := sample.AccAddress()
	tt := []struct {
		name           string
		policies       types.Policies
		signer         string
		policyRequired types.PolicyType
		expectedErr    error
	}{
		{
			name: "successfully check signer for policyType groupEmergency",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupEmergency,
					},
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupAdmin,
					},
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupOperational,
					},
				},
			},
			signer:         signer,
			policyRequired: types.PolicyType_groupEmergency,
			expectedErr:    nil,
		},
		{
			name: "successfully check signer for policyType groupOperational",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupEmergency,
					},
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupAdmin,
					},
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupOperational,
					},
				},
			},
			signer:         signer,
			policyRequired: types.PolicyType_groupOperational,
			expectedErr:    nil,
		},
		{
			name: "successfully check signer for policyType groupAdmin",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupEmergency,
					},
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupAdmin,
					},
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupOperational,
					},
				},
			},
			signer:         signer,
			policyRequired: types.PolicyType_groupAdmin,
			expectedErr:    nil,
		},
		{
			name: "signer not found",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupEmergency,
					},
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupAdmin,
					},
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupOperational,
					},
				},
			},
			signer:         sample.AccAddress(),
			policyRequired: types.PolicyType_groupEmergency,
			expectedErr:    types.ErrSignerDoesntMatch,
		},
		{
			name: "policy required not found",
			policies: types.Policies{
				Items: []*types.Policy{
					{
						Address:    signer,
						PolicyType: types.PolicyType_groupAdmin,
					},
					{
						Address:    sample.AccAddress(),
						PolicyType: types.PolicyType_groupOperational,
					},
				},
			},
			signer:         signer,
			policyRequired: types.PolicyType_groupEmergency,
			expectedErr:    types.ErrSignerDoesntMatch,
		},
		{
			name:           "empty policies",
			policies:       types.Policies{},
			signer:         signer,
			policyRequired: types.PolicyType_groupEmergency,
			expectedErr:    types.ErrSignerDoesntMatch,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.policies.CheckSigner(tc.signer, tc.policyRequired)
			require.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
