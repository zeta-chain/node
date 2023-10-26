package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddToOutTxTracker{}, "crosschain/AddToOutTxTracker", nil)
	cdc.RegisterConcrete(&MsgAddToInTxTracker{}, "crosschain/AddToInTxTracker", nil)
	cdc.RegisterConcrete(&MsgRemoveFromOutTxTracker{}, "crosschain/RemoveFromOutTxTracker", nil)
	cdc.RegisterConcrete(&MsgCreateTSSVoter{}, "crosschain/CreateTSSVoter", nil)
	cdc.RegisterConcrete(&MsgGasPriceVoter{}, "crosschain/GasPriceVoter", nil)
	cdc.RegisterConcrete(&MsgNonceVoter{}, "crosschain/NonceVoter", nil)
	cdc.RegisterConcrete(&MsgVoteOnObservedOutboundTx{}, "crosschain/VoteOnObservedOutboundTx", nil)
	cdc.RegisterConcrete(&MsgVoteOnObservedInboundTx{}, "crosschain/VoteOnObservedInboundTx", nil)
	cdc.RegisterConcrete(&MsgWhitelistERC20{}, "crosschain/WhitelistERC20", nil)
	cdc.RegisterConcrete(&MsgMigrateTssFunds{}, "crosschain/MigrateTssFunds", nil)
	cdc.RegisterConcrete(&MsgUpdateTssAddress{}, "crosschain/UpdateTssAddress", nil)
	cdc.RegisterConcrete(&MsgSetNodeKeys{}, "crosschain/SetNodeKeys", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddToOutTxTracker{},
		&MsgAddToInTxTracker{},
		&MsgRemoveFromOutTxTracker{},
		&MsgCreateTSSVoter{},
		&MsgGasPriceVoter{},
		&MsgNonceVoter{},
		&MsgVoteOnObservedOutboundTx{},
		&MsgVoteOnObservedInboundTx{},
		&MsgWhitelistERC20{},
		&MsgMigrateTssFunds{},
		&MsgUpdateTssAddress{},
		&MsgSetNodeKeys{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
