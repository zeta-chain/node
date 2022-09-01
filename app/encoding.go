package app

import (
	evmenc "github.com/evmos/ethermint/encoding"
	//"github.com/zeta-chain/zetacore/app/params"
	"github.com/cosmos/cosmos-sdk/simapp/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	//encodingConfig := params.MakeEncodingConfig()
	encodingConfig := evmenc.MakeConfig(ModuleBasics)
	//std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
