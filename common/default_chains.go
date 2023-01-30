package common

func EthChain() Chain {
	return Chain{
		ChainName: ChainName_eth_mainnet,
		ChainId:   1,
	}
}

func GoerliChain() Chain {
	return Chain{
		ChainName: ChainName_goerli_testnet,
		ChainId:   5,
	}
}

func GoerliLocalNetChain() Chain {
	return Chain{
		ChainName: ChainName_goerli_localnet,
		ChainId:   1337,
	}
}

func BscMainnetChain() Chain {
	return Chain{
		ChainName: ChainName_bsc_mainnet,
		ChainId:   56,
	}
}

func BscTestnetChain() Chain {
	return Chain{
		ChainName: ChainName_bsc_testnet,
		ChainId:   97,
	}
}

func BaobabChain() Chain {
	return Chain{
		ChainName: ChainName_baobab_testnet,
		ChainId:   1001,
	}
}
func ZetaChain() Chain {
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

func BtcTestNetChain() Chain {
	return Chain{
		ChainName: ChainName_btc_testnet,
		ChainId:   18332,
	}
}

func PolygonChain() Chain {
	return Chain{
		ChainName: ChainName_polygon_mainnet,
		ChainId:   137,
	}
}

func MumbaiChain() Chain {
	return Chain{
		ChainName: ChainName_mumbai_testnet,
		ChainId:   80001,
	}
}

func DefaultChainsList() []*Chain {
	chains := []Chain{
		BtcTestNetChain(),
		BtcMainnetChain(),
		PolygonChain(),
		MumbaiChain(),
		BaobabChain(),
		BscTestnetChain(),
		BscMainnetChain(),
		EthChain(),
		GoerliChain(),
		GoerliLocalNetChain(),
		ZetaChain(),
	}
	var c []*Chain
	for i := 0; i < len(chains); i++ {
		c = append(c, &chains[i])
	}
	return c
}
