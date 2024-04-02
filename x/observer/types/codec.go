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
	cdc.RegisterConcrete(&MsgAddBlameVote{}, "crosschain/AddBlameVote", nil)
	cdc.RegisterConcrete(&MsgUpdateCrosschainFlags{}, "crosschain/UpdateCrosschainFlags", nil)
	cdc.RegisterConcrete(&MsgUpdateKeygen{}, "crosschain/UpdateKeygen", nil)
	cdc.RegisterConcrete(&MsgVoteBlockHeader{}, "crosschain/VoteBlockHeader", nil)
	cdc.RegisterConcrete(&MsgUpdateObserver{}, "observer/UpdateObserver", nil)
	cdc.RegisterConcrete(&MsgResetChainNonces{}, "observer/ResetChainNonces", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddObserver{},
		&MsgUpdateChainParams{},
		&MsgRemoveChainParams{},
		&MsgAddBlameVote{},
		&MsgUpdateCrosschainFlags{},
		&MsgUpdateKeygen{},
		&MsgVoteBlockHeader{},
		&MsgUpdateObserver{},
		&MsgResetChainNonces{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
