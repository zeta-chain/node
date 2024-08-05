package app

import (
	"cosmossdk.io/simapp/params"
	evmenc "github.com/zeta-chain/ethermint/encoding"
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
