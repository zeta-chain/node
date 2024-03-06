package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/app"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/authority/types"
)

// setConfig sets the global config to use zeta chain's bech32 prefixes
func setConfig() {
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cfg.Seal()
}

func TestPolicies_Validate(t *testing.T) {
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
			name:        "regular valid policies",
			policies:    sample.Policies(),
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
