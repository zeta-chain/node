package app

import (
	evmenc "github.com/zeta-chain/ethermint/encoding"
	ethermint "github.com/zeta-chain/ethermint/types"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() ethermint.EncodingConfig {
	//encodingConfig := params.MakeEncodingConfig()
	encodingConfig := evmenc.MakeConfig(ModuleBasics)
	//std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
