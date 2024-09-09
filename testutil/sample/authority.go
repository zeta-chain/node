package sample

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/zeta-chain/node/pkg/chains"
	authoritytypes "github.com/zeta-chain/node/x/authority/types"
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
	return authoritytypes.ChainInfo{
		Chains: []chains.Chain{
			Chain(startChainID),
			Chain(startChainID + 1),
			Chain(startChainID + 2),
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

// MultipleSignerMessage is a sample message which has two signers instead of one. This is used to test cases when we have checks for number of signers such as authorized transactions.
type MultipleSignerMessage struct{}

var _ sdk.Msg = &MultipleSignerMessage{}

func (m *MultipleSignerMessage) Reset()               {}
func (m *MultipleSignerMessage) String() string       { return "MultipleSignerMessage" }
func (m *MultipleSignerMessage) ProtoMessage()        {}
func (m *MultipleSignerMessage) ValidateBasic() error { return nil }
func (m *MultipleSignerMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		sdk.MustAccAddressFromBech32(AccAddress()),
		sdk.MustAccAddressFromBech32(AccAddress()),
	}
}
