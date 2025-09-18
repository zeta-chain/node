package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgAddOutboundTracker{}, "crosschain/AddOutboundTracker", nil)
	cdc.RegisterConcrete(&MsgAddInboundTracker{}, "crosschain/AddInboundTracker", nil)
	cdc.RegisterConcrete(&MsgRemoveOutboundTracker{}, "crosschain/RemoveOutboundTracker", nil)
	cdc.RegisterConcrete(&MsgVoteGasPrice{}, "crosschain/VoteGasPrice", nil)
	cdc.RegisterConcrete(&MsgVoteOutbound{}, "crosschain/VoteOutbound", nil)
	cdc.RegisterConcrete(&MsgVoteInbound{}, "crosschain/VoteInbound", nil)
	cdc.RegisterConcrete(&MsgWhitelistERC20{}, "crosschain/WhitelistERC20", nil)
	cdc.RegisterConcrete(&MsgMigrateTssFunds{}, "crosschain/MigrateTssFunds", nil)
	cdc.RegisterConcrete(&MsgUpdateTssAddress{}, "crosschain/UpdateTssAddress", nil)
	cdc.RegisterConcrete(&MsgAbortStuckCCTX{}, "crosschain/AbortStuckCCTX", nil)
	cdc.RegisterConcrete(&MsgUpdateRateLimiterFlags{}, "crosschain/UpdateRateLimiterFlags", nil)
	cdc.RegisterConcrete(&MsgRemoveInboundTracker{}, "crosschain/RemoveInboundTracker", nil)

	// legacy messages defined for backward compatibility
	cdc.RegisterConcrete(&MsgAddToInTxTracker{}, "crosschain/AddToInTxTracker", nil)
	cdc.RegisterConcrete(&MsgAddToOutTxTracker{}, "crosschain/AddToOutTxTracker", nil)
	cdc.RegisterConcrete(&MsgRemoveFromOutTxTracker{}, "crosschain/RemoveFromOutTxTracker", nil)
	cdc.RegisterConcrete(&MsgVoteOnObservedOutboundTx{}, "crosschain/VoteOnObservedOutboundTx", nil)
	cdc.RegisterConcrete(&MsgVoteOnObservedInboundTx{}, "crosschain/VoteOnObservedInboundTx", nil)
	cdc.RegisterConcrete(&MsgGasPriceVoter{}, "crosschain/GasPriceVoter", nil)
	cdc.RegisterConcrete(&MsgMigrateERC20CustodyFunds{}, "crosschain/MigrateERC20CustodyFunds", nil)
	cdc.RegisterConcrete(&MsgUpdateERC20CustodyPauseStatus{}, "crosschain/UpdateERC20CustodyPauseStatus", nil)
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
		&MsgRemoveInboundTracker{},

		// legacy messages defined for backward compatibility
		&MsgAddToInTxTracker{},
		&MsgAddToOutTxTracker{},
		&MsgRemoveFromOutTxTracker{},
		&MsgVoteOnObservedOutboundTx{},
		&MsgVoteOnObservedInboundTx{},
		&MsgGasPriceVoter{},
		&MsgMigrateERC20CustodyFunds{},
		&MsgUpdateERC20CustodyPauseStatus{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
