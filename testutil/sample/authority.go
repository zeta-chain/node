package sample

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
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

func PoliciesWithAdmin(address string) authoritytypes.Policies {
	return authoritytypes.Policies{
		Items: []*authoritytypes.Policy{
			{
				Address:    address,
				PolicyType: authoritytypes.PolicyType_groupEmergency,
			},
			{
				Address:    address,
				PolicyType: authoritytypes.PolicyType_groupAdmin,
			},
			{
				Address:    address,
				PolicyType: authoritytypes.PolicyType_groupOperational,
			},
		},
	}
}

func AdminMessage(signer string) sdk.Msg {
	return &crosschaintypes.MsgRefundAbortedCCTX{
		Creator: signer,
	}
}

func NonAdminMessage(signer string) sdk.Msg {
	return &crosschaintypes.MsgVoteOnObservedInboundTx{
		Creator: signer,
	}
}

func MultipleSignerMessage() sdk.Msg {
	return &TestMessage{}
}

type TestMessage struct{}

var _ sdk.Msg = &TestMessage{}

func (m *TestMessage) Reset()               {}
func (m *TestMessage) String() string       { return "TestMessage" }
func (m *TestMessage) ProtoMessage()        {}
func (m *TestMessage) ValidateBasic() error { return nil }
func (m *TestMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{
		sdk.MustAccAddressFromBech32(AccAddress()),
		sdk.MustAccAddressFromBech32(AccAddress()),
	}
}
func (m *TestMessage) Route() string { return "test" }
