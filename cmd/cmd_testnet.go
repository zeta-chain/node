// +build testnet

package cmd

const (
	Bech32PrefixAccAddr         = "tmeta"
	Bech32PrefixAccPub          = "tmetapub"
	Bech32PrefixValAddr         = "tmetav"
	Bech32PrefixValPub          = "tmetavpub"
	Bech32PrefixConsAddr        = "tmetac"
	Bech32PrefixConsPub         = "tmetacpub"
	DenomRegex                  = `[a-zA-Z][a-zA-Z0-9:\\/\\\-\\_\\.]{2,127}`
	METAChainCoinType    uint32 = 933
	METAChainHDPath      string = `m/44'/933'/0'/0/0`
	NET                         = "TESTNET"
)
