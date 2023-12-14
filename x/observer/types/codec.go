package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddObserver{}, "observer/AddObserver", nil)
	cdc.RegisterConcrete(&MsgUpdateCoreParams{}, "observer/UpdateCoreParams", nil)
	cdc.RegisterConcrete(&MsgRemoveCoreParams{}, "observer/RemoveCoreParams", nil)
	cdc.RegisterConcrete(&MsgAddBlameVote{}, "crosschain/AddBlameVote", nil)
	cdc.RegisterConcrete(&MsgUpdateCrosschainFlags{}, "crosschain/UpdateCrosschainFlags", nil)
	cdc.RegisterConcrete(&MsgUpdateKeygen{}, "crosschain/UpdateKeygen", nil)
	cdc.RegisterConcrete(&MsgAddBlockHeader{}, "crosschain/AddBlockHeader", nil)
	cdc.RegisterConcrete(&MsgUpdateObserver{}, "observer/UpdateObserver", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddObserver{},
		&MsgUpdateCoreParams{},
		&MsgRemoveCoreParams{},
		&MsgAddBlameVote{},
		&MsgUpdateCrosschainFlags{},
		&MsgUpdateKeygen{},
		&MsgAddBlockHeader{},
		&MsgUpdateObserver{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
