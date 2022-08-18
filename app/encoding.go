package app

import (
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/ignite-hq/cli/ignite/pkg/cosmoscmd"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig(moduleBasics module.BasicManager) cosmoscmd.EncodingConfig {
	encodingConfig := cosmoscmd.MakeEncodingConfig(moduleBasics)
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
