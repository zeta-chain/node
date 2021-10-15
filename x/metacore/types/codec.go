package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
	cdc.RegisterConcrete(&MsgCreateSendVoter{}, "metacore/CreateSendVoter", nil)

	cdc.RegisterConcrete(&MsgTxoutConfirmationVoter{}, "metacore/TxoutConfirmationVoter", nil)

	cdc.RegisterConcrete(&MsgSetNodeKeys{}, "metacore/SetNodeKeys", nil)

	cdc.RegisterConcrete(&MsgCreateTxinVoter{}, "metacore/CreateTxinVoter", nil)

}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateSendVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTxoutConfirmationVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetNodeKeys{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateTxinVoter{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	//amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
