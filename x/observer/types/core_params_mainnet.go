//go:build !PRIVNET && !TESTNET
// +build !PRIVNET,!TESTNET

package types

func GetCoreParams() CoreParamsList {
	params := CoreParamsList{
		CoreParams: []*CoreParams{},
	}
	return params
}
