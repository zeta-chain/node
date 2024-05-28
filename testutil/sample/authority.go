package sample

import (
	"fmt"

	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
)

func Policies() authoritytypes.Policies {
	return authoritytypes.Policies{
		Items: []*authoritytypes.Policy{
			{
				Address:    AccAddress(),
				PolicyType: authoritytypes.PolicyType_groupEmergency,
			},
			{
				Address:    AccAddress(),
				PolicyType: authoritytypes.PolicyType_groupAdmin,
			},
			{
				Address:    AccAddress(),
				PolicyType: authoritytypes.PolicyType_groupOperational,
			},
		},
	}
}

func AuthorizationList(val string) authoritytypes.AuthorizationList {
	return authoritytypes.AuthorizationList{
		Authorizations: []authoritytypes.Authorization{
			{
				MsgUrl:           fmt.Sprintf("/zetachain/%d%s", 0, val),
				AuthorizedPolicy: authoritytypes.PolicyType_groupEmergency,
			},
			{
				MsgUrl:           fmt.Sprintf("/zetachain/%d%s", 1, val),
				AuthorizedPolicy: authoritytypes.PolicyType_groupAdmin,
			},
			{
				MsgUrl:           fmt.Sprintf("/zetachain/%d%s", 2, val),
				AuthorizedPolicy: authoritytypes.PolicyType_groupOperational,
			},
		},
	}
}
