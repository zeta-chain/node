package sample

import authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"

func Policies() authoritytypes.Policies {
	return authoritytypes.Policies{
		PolicyAddresses: []*authoritytypes.PolicyAddress{
			{
				Address:    AccAddress(),
				PolicyType: authoritytypes.PolicyType_groupEmergency,
			},
			{
				Address:    AccAddress(),
				PolicyType: authoritytypes.PolicyType_groupAdmin,
			},
		},
	}
}
