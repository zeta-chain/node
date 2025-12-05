package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddObserver{}, "observer/AddObserver", nil)
	cdc.RegisterConcrete(&MsgRemoveObserver{}, "observer/RemoveObserver", nil)
	cdc.RegisterConcrete(&MsgUpdateChainParams{}, "observer/UpdateChainParams", nil)
	cdc.RegisterConcrete(&MsgRemoveChainParams{}, "observer/RemoveChainParams", nil)
	cdc.RegisterConcrete(&MsgVoteBlockHeader{}, "observer/VoteBlockHeader", nil)
	cdc.RegisterConcrete(&MsgVoteBlame{}, "observer/VoteBlame", nil)
	cdc.RegisterConcrete(&MsgUpdateKeygen{}, "observer/UpdateKeygen", nil)
	cdc.RegisterConcrete(&MsgUpdateObserver{}, "observer/UpdateObserver", nil)
	cdc.RegisterConcrete(&MsgResetChainNonces{}, "observer/ResetChainNonces", nil)
	cdc.RegisterConcrete(&MsgVoteTSS{}, "observer/VoteTSS", nil)
	cdc.RegisterConcrete(&MsgEnableCCTX{}, "observer/EnableCCTX", nil)
	cdc.RegisterConcrete(&MsgDisableCCTX{}, "observer/DisableCCTX", nil)
	cdc.RegisterConcrete(&MsgDisableFastConfirmation{}, "observer/DisableFastConfirmation", nil)
	cdc.RegisterConcrete(&MsgUpdateGasPriceIncreaseFlags{}, "observer/UpdateGasPriceIncreaseFlags", nil)
	cdc.RegisterConcrete(&MsgUpdateOperationalFlags{}, "observer/UpdateOperationalFlags", nil)
	cdc.RegisterConcrete(&MsgUpdateOperationalChainParams{}, "observer/UpdateOperationalChainParams", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddObserver{},
		&MsgRemoveObserver{},
		&MsgUpdateChainParams{},
		&MsgRemoveChainParams{},
		&MsgVoteBlame{},
		&MsgUpdateKeygen{},
		&MsgVoteBlockHeader{},
		&MsgUpdateObserver{},
		&MsgResetChainNonces{},
		&MsgVoteTSS{},
		&MsgEnableCCTX{},
		&MsgDisableCCTX{},
		&MsgDisableFastConfirmation{},
		&MsgUpdateGasPriceIncreaseFlags{},
		&MsgUpdateOperationalFlags{},
		&MsgUpdateOperationalChainParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
