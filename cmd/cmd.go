//go:build !testnet && !mocknet
// +build !testnet,!mocknet

package cmd

const (
	Bech32PrefixAccAddr         = "zeta"
	Bech32PrefixAccPub          = "zetapub"
	Bech32PrefixValAddr         = "zetav"
	Bech32PrefixValPub          = "zetavpub"
	Bech32PrefixConsAddr        = "zetac"
	Bech32PrefixConsPub         = "zetacpub"
	DenomRegex                  = `[a-zA-Z][a-zA-Z0-9:\\/\\\-\\_\\.]{2,127}`
	ZETAChainCoinType    uint32 = 933
	ZetaChainHDPath      string = `m/44'/933'/0'/0/0`
	NET                         = "MAINNET"
	CHAINID                     = "testing"
)
