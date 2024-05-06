package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddOutboundTracker{}, "crosschain/AddOutboundTracker", nil)
	cdc.RegisterConcrete(&MsgAddInboundTracker{}, "crosschain/AddToInboundTracker", nil)
	cdc.RegisterConcrete(&MsgRemoveOutboundTracker{}, "crosschain/RemoveOutboundTracker", nil)
	cdc.RegisterConcrete(&MsgVoteGasPrice{}, "crosschain/VoteGasPrice", nil)
	cdc.RegisterConcrete(&MsgVoteOutbound{}, "crosschain/VoteOutbound", nil)
	cdc.RegisterConcrete(&MsgVoteInbound{}, "crosschain/VoteInbound", nil)
	cdc.RegisterConcrete(&MsgWhitelistERC20{}, "crosschain/WhitelistERC20", nil)
	cdc.RegisterConcrete(&MsgMigrateTssFunds{}, "crosschain/MigrateTssFunds", nil)
	cdc.RegisterConcrete(&MsgUpdateTssAddress{}, "crosschain/UpdateTssAddress", nil)
	cdc.RegisterConcrete(&MsgAbortStuckCCTX{}, "crosschain/AbortStuckCCTX", nil)
	cdc.RegisterConcrete(&MsgUpdateRateLimiterFlags{}, "crosschain/UpdateRateLimiterFlags", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddOutboundTracker{},
		&MsgAddInboundTracker{},
		&MsgRemoveOutboundTracker{},
		&MsgVoteGasPrice{},
		&MsgVoteOutbound{},
		&MsgVoteInbound{},
		&MsgWhitelistERC20{},
		&MsgMigrateTssFunds{},
		&MsgUpdateTssAddress{},
		&MsgAbortStuckCCTX{},
		&MsgUpdateRateLimiterFlags{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
