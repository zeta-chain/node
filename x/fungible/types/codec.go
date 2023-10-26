package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDeployFungibleCoinZRC20{}, "fungible/DeployFungibleCoinZRC20", nil)
	cdc.RegisterConcrete(&MsgRemoveForeignCoin{}, "fungible/RemoveForeignCoin", nil)
	cdc.RegisterConcrete(&MsgUpdateSystemContract{}, "fungible/UpdateSystemContract", nil)
	cdc.RegisterConcrete(&MsgUpdateZRC20WithdrawFee{}, "fungible/UpdateZRC20WithdrawFee", nil)
	cdc.RegisterConcrete(&MsgUpdateContractBytecode{}, "fungible/UpdateContractBytecode", nil)
	cdc.RegisterConcrete(&MsgUpdateZRC20PausedStatus{}, "fungible/UpdateZRC20PausedStatus", nil)
	cdc.RegisterConcrete(&MsgUpdateZRC20LiquidityCap{}, "fungible/UpdateZRC20LiquidityCap", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeployFungibleCoinZRC20{},
		&MsgRemoveForeignCoin{},
		&MsgUpdateSystemContract{},
		&MsgUpdateZRC20WithdrawFee{},
		&MsgUpdateContractBytecode{},
		&MsgUpdateZRC20PausedStatus{},
		&MsgUpdateZRC20LiquidityCap{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
