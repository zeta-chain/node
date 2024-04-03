package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddObserver{}, "observer/AddObserver", nil)
	cdc.RegisterConcrete(&MsgUpdateChainParams{}, "observer/UpdateChainParams", nil)
	cdc.RegisterConcrete(&MsgRemoveChainParams{}, "observer/RemoveChainParams", nil)
	cdc.RegisterConcrete(&MsgAddBlameVote{}, "observer/AddBlameVote", nil)
	cdc.RegisterConcrete(&MsgUpdateCrosschainFlags{}, "observer/UpdateCrosschainFlags", nil)
	cdc.RegisterConcrete(&MsgUpdateKeygen{}, "observer/UpdateKeygen", nil)
	cdc.RegisterConcrete(&MsgAddBlockHeader{}, "observer/AddBlockHeader", nil)
	cdc.RegisterConcrete(&MsgUpdateObserver{}, "observer/UpdateObserver", nil)
	cdc.RegisterConcrete(&MsgResetChainNonces{}, "observer/ResetChainNonces", nil)
	cdc.RegisterConcrete(&MsgVoteTSS{}, "observer/VoteTSS", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddObserver{},
		&MsgUpdateChainParams{},
		&MsgRemoveChainParams{},
		&MsgAddBlameVote{},
		&MsgUpdateCrosschainFlags{},
		&MsgUpdateKeygen{},
		&MsgAddBlockHeader{},
		&MsgUpdateObserver{},
		&MsgResetChainNonces{},
		&MsgVoteTSS{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
