// +build testnet

package cmd

const (
	Bech32PrefixAccAddr         = "tzeta"
	Bech32PrefixAccPub          = "tzetapub"
	Bech32PrefixValAddr         = "tzetav"
	Bech32PrefixValPub          = "tzetavpub"
	Bech32PrefixConsAddr        = "tzetac"
	Bech32PrefixConsPub         = "tzetacpub"
	DenomRegex                  = `[a-zA-Z][a-zA-Z0-9:\\/\\\-\\_\\.]{2,127}`
	METAChainCoinType    uint32 = 933
	METAChainHDPath      string = `m/44'/933'/0'/0/0`
	NET                         = "TESTNET"
)
