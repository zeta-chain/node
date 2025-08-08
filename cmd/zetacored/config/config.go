package config

import (
	"path/filepath"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/config"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/spf13/viper"

	"github.com/zeta-chain/node/app/eips"
)

const (
	DisplayDenom  = "zeta"
	BaseDenom     = "azeta"
	AppName       = "zetacored"
	BaseDenomUnit = 18
)

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(DisplayDenom, sdkmath.LegacyOneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(BaseDenom, sdkmath.LegacyNewDecWithPrec(1, BaseDenomUnit)); err != nil {
		panic(err)
	}
}

// GetChainIDFromHome returns the chain ID from the client configuration
// in the given home directory.
func GetChainIDFromHome(home string) (string, error) {
	v := viper.New()
	v.AddConfigPath(filepath.Join(home, "config"))
	v.SetConfigName("client")
	v.SetConfigType("toml")

	if err := v.ReadInConfig(); err != nil {
		return "", err
	}
	conf := new(config.ClientConfig)

	if err := v.Unmarshal(conf); err != nil {
		return "", err
	}

	return conf.ChainID, nil
}

// CosmosEVMActivators defines a map of opcode modifiers associated
// with a key defining the corresponding EIP.
var CosmosEVMActivators = map[int]func(*vm.JumpTable){
	0o000: eips.Enable0000,
	0o001: eips.Enable0001,
	0o002: eips.Enable0002,
}
