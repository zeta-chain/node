package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgFungibleTestMsg{}, "fungible/FungibleTestMsg", nil)
	cdc.RegisterConcrete(&MsgDeployFungibleCoinZRC4{}, "fungible/DeployFungibleCoinZRC4", nil)
	cdc.RegisterConcrete(&MsgDeployGasPriceOracle{}, "fungible/DeployGasPriceOracle", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFungibleTestMsg{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeployFungibleCoinZRC4{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeployGasPriceOracle{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
