package sample

import (
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
