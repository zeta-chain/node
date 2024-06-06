package types

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// DefaultPolicyAddress is the default value for policy address
	DefaultPolicyAddress = "zeta1afk9zr2hn2jsac63h4hm60vl9z3e5u69gndzf7c99cqge3vzwjzsxn0x73"
)

// DefaultPolicies returns the default value for policies
func DefaultPolicies() Policies {
	return Policies{
		Items: []*Policy{
			{
				Address:    DefaultPolicyAddress,
				PolicyType: PolicyType_groupEmergency,
			},
			{
				Address:    DefaultPolicyAddress,
				PolicyType: PolicyType_groupOperational,
			},
			{
				Address:    DefaultPolicyAddress,
				PolicyType: PolicyType_groupAdmin,
			},
		},
	}
}

// Validate performs basic validation of policies
func (p Policies) Validate() error {
	policyTypeMap := make(map[PolicyType]bool)

	// for each policy, check address, policy type, and ensure no duplicate policy types
	for _, policy := range p.Items {
		_, err := sdk.AccAddressFromBech32(policy.Address)
		if err != nil {
			return fmt.Errorf("invalid address: %s", err)
		}

		if err := policy.PolicyType.Validate(); err != nil {
			return err
		}

		if policyTypeMap[policy.PolicyType] {
			return fmt.Errorf("duplicate policy type: %s", policy.PolicyType)
		}
		policyTypeMap[policy.PolicyType] = true
	}

	return nil
}

// CheckSigner checks if the signer is authorized for the given policy type
func (p Policies) CheckSigner(signer string, policyRequired PolicyType) error {
	for _, policy := range p.Items {
		if policy.Address == signer && policy.PolicyType == policyRequired {
			return nil
		}
	}
	return errors.Wrap(ErrSignerDoesntMatch, fmt.Sprintf("signer: %s, policy required for message: %s ",
		signer, policyRequired.String()))
}
