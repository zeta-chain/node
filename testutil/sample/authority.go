package sample

import (
	"fmt"

	"github.com/zeta-chain/zetacore/pkg/chains"
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

func ChainInfo(startChainID int64) authoritytypes.ChainInfo {
	chain1 := Chain(startChainID)
	chain2 := Chain(startChainID + 1)
	chain3 := Chain(startChainID + 2)

	return authoritytypes.ChainInfo{
		Chains: []chains.Chain{
			*chain1,
			*chain2,
			*chain3,
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

func Authorization() authoritytypes.Authorization {
	return authoritytypes.Authorization{
		MsgUrl:           "ABC",
		AuthorizedPolicy: authoritytypes.PolicyType_groupOperational,
	}
}
