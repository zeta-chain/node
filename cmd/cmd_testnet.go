//go:build testnet
// +build testnet

package cmd

const (
	Bech32PrefixAccAddr         = "meta"
	Bech32PrefixAccPub          = "metapub"
	Bech32PrefixValAddr         = "metav"
	Bech32PrefixValPub          = "metavpub"
	Bech32PrefixConsAddr        = "metac"
	Bech32PrefixConsPub         = "metacpub"
	DenomRegex                  = `[a-zA-Z][a-zA-Z0-9:\\/\\\-\\_\\.]{2,127}`
	METAChainCoinType    uint32 = 933
	METAChainHDPath      string = `m/44'/933'/0'/0/0`
	NET                         = "MAINNET"
	CHAINID                     = "metacore"
)
