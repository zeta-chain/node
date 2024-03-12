package v7

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/observer/types"
)

// observerKeeper prevents circular dependency
type observerKeeper interface {
	GetParams(ctx sdk.Context) (params types.Params)
	GetAuthorityKeeper() types.AuthorityKeeper
}

// MigrateStore performs in-place store migrations from v6 to v7
func MigrateStore(ctx sdk.Context, observerKeeper observerKeeper) error {
	ctx.Logger().Info("Migrating observer store from v6 to v7")
	return MigratePolicies(ctx, observerKeeper)
}

// MigratePolicies migrates policies from observer to authority
func MigratePolicies(ctx sdk.Context, observerKeeper observerKeeper) error {
	params := observerKeeper.GetParams(ctx)
	authorityKeeper := observerKeeper.GetAuthorityKeeper()

	var policies authoritytypes.Policies

	// convert observer policies to authority policies
	for _, adminPolicy := range params.AdminPolicy {
		if adminPolicy != nil {

			// Convert group1 -> emergency and group2 -> admin
			policyType := authoritytypes.PolicyType_groupAdmin
			if adminPolicy.PolicyType == types.Policy_Type_group1 {
				policyType = authoritytypes.PolicyType_groupEmergency
			}

			policies.Items = append(policies.Items, &authoritytypes.Policy{
				Address:    adminPolicy.Address,
				PolicyType: policyType,
			})
		}
	}

	// ensure policies are valid
	if err := policies.Validate(); err != nil {
		return err
	}

	// set policies in authority
	authorityKeeper.SetPolicies(ctx, policies)
	return nil
}
