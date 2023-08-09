package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddToOutTxTracker{}, "crosschain/AddToOutTxTracker", nil)
	cdc.RegisterConcrete(&MsgRemoveFromOutTxTracker{}, "crosschain/RemoveFromOutTxTracker", nil)
	cdc.RegisterConcrete(&MsgCreateTSSVoter{}, "crosschain/CreateTSSVoter", nil)
	cdc.RegisterConcrete(&MsgGasPriceVoter{}, "crosschain/GasPriceVoter", nil)
	cdc.RegisterConcrete(&MsgNonceVoter{}, "crosschain/NonceVoter", nil)
	cdc.RegisterConcrete(&MsgVoteOnObservedOutboundTxResponse{}, "crosschain/ReceiveConfirmation", nil)
	cdc.RegisterConcrete(&MsgVoteOnObservedInboundTx{}, "crosschain/SendVoter", nil)
	cdc.RegisterConcrete(&MsgSetNodeKeys{}, "crosschain/SetNodeKeys", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddToOutTxTracker{},
		&MsgRemoveFromOutTxTracker{},
		&MsgCreateTSSVoter{},
		&MsgGasPriceVoter{},
		&MsgNonceVoter{},
		&MsgVoteOnObservedOutboundTx{},
		&MsgVoteOnObservedInboundTx{},
		&MsgSetNodeKeys{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
