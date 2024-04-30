package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgDeployFungibleCoinZRC20{}, "fungible/DeployFungibleCoinZRC20", nil)
	cdc.RegisterConcrete(&MsgDeploySystemContracts{}, "fungible/MsgDeploySystemContracts", nil)
	cdc.RegisterConcrete(&MsgRemoveForeignCoin{}, "fungible/RemoveForeignCoin", nil)
	cdc.RegisterConcrete(&MsgUpdateSystemContract{}, "fungible/UpdateSystemContract", nil)
	cdc.RegisterConcrete(&MsgUpdateZRC20WithdrawFee{}, "fungible/UpdateZRC20WithdrawFee", nil)
	cdc.RegisterConcrete(&MsgUpdateContractBytecode{}, "fungible/UpdateContractBytecode", nil)
	cdc.RegisterConcrete(&MsgPauseZRC20{}, "fungible/PauseZRC20", nil)
	cdc.RegisterConcrete(&MsgUnpauseZRC20{}, "fungible/UnpauseZRC20", nil)
	cdc.RegisterConcrete(&MsgUpdateZRC20LiquidityCap{}, "fungible/UpdateZRC20LiquidityCap", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDeployFungibleCoinZRC20{},
		&MsgDeploySystemContracts{},
		&MsgRemoveForeignCoin{},
		&MsgUpdateSystemContract{},
		&MsgUpdateZRC20WithdrawFee{},
		&MsgUpdateContractBytecode{},
		&MsgPauseZRC20{},
		&MsgUnpauseZRC20{},
		&MsgUpdateZRC20LiquidityCap{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
