package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgZetaConversionRateVoter{}, "zetacore/ZetaConversionRateVoter", nil)
	cdc.RegisterConcrete(&MsgAddToOutTxTracker{}, "zetacore/AddToOutTxTracker", nil)
	cdc.RegisterConcrete(&MsgRemoveFromOutTxTracker{}, "zetacore/RemoveFromOutTxTracker", nil)
	cdc.RegisterConcrete(&MsgCreateTSSVoter{}, "zetacore/CreateTSSVoter", nil)
	cdc.RegisterConcrete(&MsgGasBalanceVoter{}, "zetacore/GasBalanceVoter", nil)
	cdc.RegisterConcrete(&MsgGasPriceVoter{}, "zetacore/GasPriceVoter", nil)
	cdc.RegisterConcrete(&MsgNonceVoter{}, "zetacore/NonceVoter", nil)
	cdc.RegisterConcrete(&MsgVoteOnObservedOutboundTxResponse{}, "zetacore/ReceiveConfirmation", nil)
	cdc.RegisterConcrete(&MsgVoteOnObservedInboundTx{}, "zetacore/SendVoter", nil)
	cdc.RegisterConcrete(&MsgSetNodeKeys{}, "zetacore/SetNodeKeys", nil)

	// TODO : change RPC
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgZetaConversionRateVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddToOutTxTracker{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRemoveFromOutTxTracker{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateTSSVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgGasBalanceVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgGasPriceVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgNonceVoter{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVoteOnObservedOutboundTx{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVoteOnObservedInboundTx{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetNodeKeys{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
