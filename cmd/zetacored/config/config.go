package config

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/evm/cmd/evmd/config"
)

const (
	DisplayDenom = "zeta"
	BaseDenom    = "azeta"
	AppName      = "zetacored"
)

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(DisplayDenom, sdkmath.LegacyOneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(BaseDenom, sdkmath.LegacyNewDecWithPrec(1, config.BaseDenomUnit)); err != nil {
		panic(err)
	}
}
