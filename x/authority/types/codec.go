package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdatePolicies{}, "authority/UpdatePolicies", nil)
	cdc.RegisterConcrete(&MsgUpdateChainInfo{}, "authority/UpdateChainInfo", nil)
	cdc.RegisterConcrete(&MsgRemoveChainInfo{}, "authority/RemoveChainInfo", nil)
	cdc.RegisterConcrete(&MsgAddAuthorization{}, "authority/AddAuthorization", nil)
	cdc.RegisterConcrete(&MsgRemoveAuthorization{}, "authority/RemoveAuthorization", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdatePolicies{},
		&MsgUpdateChainInfo{},
		&MsgRemoveChainInfo{},
		&MsgAddAuthorization{},
		&MsgRemoveAuthorization{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
