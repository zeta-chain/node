// Package cmd provides cosmos constants for ZetaClient.
package cmd

import (
	"sync"

	cosmos "github.com/cosmos/cosmos-sdk/types"
)

const (
	Bech32PrefixAccAddr         = "zeta"
	Bech32PrefixAccPub          = "zetapub"
	Bech32PrefixValAddr         = "zetav"
	Bech32PrefixValPub          = "zetavpub"
	Bech32PrefixConsAddr        = "zetac"
	Bech32PrefixConsPub         = "zetacpub"
	DenomRegex                  = `[a-zA-Z][a-zA-Z0-9:\\/\\\-\\_\\.]{2,127}`
	ZetaChainHDPath      string = `m/44'/60'/0'/0/0`
)

var setupConfig sync.Once

// SetupCosmosConfig configures basic Cosmos parameters.
// This function is required because some parts of ZetaClient rely on these constants.
func SetupCosmosConfig() {
	setupConfig.Do(setupCosmosConfig)
}

func setupCosmosConfig() {
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.SetFullFundraiserPath(ZetaChainHDPath)
	cosmos.SetCoinDenomRegex(func() string { return DenomRegex })
}
