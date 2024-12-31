package config

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethermint "github.com/zeta-chain/ethermint/types"
)

const (
	DisplayDenom = "zeta"
	BaseDenom    = "azeta"
	AppName      = "zetacored"
)

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(DisplayDenom, sdk.OneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(BaseDenom, sdk.NewDecWithPrec(1, ethermint.BaseDenomUnit)); err != nil {
		panic(err)
	}
}
