package common

func EthChain() Chain {
	return Chain{
		ChainName: ChainName_eth_mainnet,
		ChainId:   1,
	}
}

func BscMainnetChain() Chain {
	return Chain{
		ChainName: ChainName_bsc_mainnet,
		ChainId:   56,
	}
}

func ZetaChain() Chain {
	return Chain{
		ChainName: ChainName_zeta_mainnet,
		ChainId:   7000,
	}
}

func ZetaTestnetChain() Chain {
	return Chain{
		ChainName: ChainName_zeta_testnet,
		ChainId:   7001,
	}
}

func ZetaMocknetChain() Chain {
	return Chain{
		ChainName: ChainName_zeta_mainnet,
		ChainId:   70000,
	}
}

func ZetaPrivnetChain() Chain {
	return Chain{
		ChainName: ChainName_zeta_mainnet,
		ChainId:   101,
	}
}

func BtcMainnetChain() Chain {
	return Chain{
		ChainName: ChainName_btc_mainnet,
		ChainId:   8332,
	}
}

func BtcChainID() int64 {
	return BtcMainnetChain().ChainId
}

func BtcDustOffset() int64 {
	return 2000
}

func PolygonChain() Chain {
	return Chain{
		ChainName: ChainName_polygon_mainnet,
		ChainId:   137,
	}
}

func BtcRegtestChain() Chain {
	return Chain{
		ChainName: ChainName_btc_regtest,
		ChainId:   18444,
	}
}

func GoerliChain() Chain {
	return Chain{
		ChainName: ChainName_goerli_localnet,
		ChainId:   1337,
	}
}

func DefaultChainsList() []*Chain {
	chains := []Chain{
		BtcMainnetChain(),
		BscMainnetChain(),
		EthChain(),
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
		BtcMainnetChain(),
		BscMainnetChain(),
		EthChain(),
	}
	var c []*Chain
	for i := 0; i < len(chains); i++ {
		c = append(c, &chains[i])
	}
	return c
}
