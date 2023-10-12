//go:build PRIVNET
// +build PRIVNET

package common

func GoerliChain() Chain {
	return Chain{
		ChainName: ChainName_goerli_localnet,
		ChainId:   1337,
	}
}

func ZetaChain() Chain {
	return Chain{
		ChainName: ChainName_zeta_mainnet,
		ChainId:   101,
	}
}

func BtcRegtestChain() Chain {
	return Chain{
		ChainName: ChainName_btc_regtest,
		ChainId:   18444,
	}
}

func BtcChainID() int64 {
	return BtcRegtestChain().ChainId
}

func DefaultChainsList() []*Chain {
	chains := []Chain{
		BtcRegtestChain(),
		GoerliChain(),
		ZetaChain(),
	}
	var c []*Chain
	for i := 0; i < len(chains); i++ {
		c = append(c, &chains[i])
	}
	return c
}

func ExternalChainList() []*Chain {
	chains := []Chain{
		BtcRegtestChain(),
		GoerliChain(),
	}
	var c []*Chain
	for i := 0; i < len(chains); i++ {
		c = append(c, &chains[i])
	}
	return c
}
