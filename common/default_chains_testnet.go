//go:build TESTNET
// +build TESTNET

package common

func GoerliChain() Chain {
	return Chain{
		ChainName: ChainName_goerli_testnet,
		ChainId:   5,
	}
}

func BscTestnetChain() Chain {
	return Chain{
		ChainName: ChainName_bsc_testnet,
		ChainId:   97,
	}
}

func ZetaChain() Chain {
	return Chain{
		ChainName: ChainName_zeta_testnet,
		ChainId:   7001,
	}
}

func BtcTestNetChain() Chain {
	return Chain{
		ChainName: ChainName_btc_testnet,
		ChainId:   18332,
	}
}

func BtcChainID() int64 {
	return BtcTestNetChain().ChainId
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
		MumbaiChain(),
		BscTestnetChain(),
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
		BtcTestNetChain(),
		MumbaiChain(),
		BscTestnetChain(),
		GoerliChain(),
	}
	var c []*Chain
	for i := 0; i < len(chains); i++ {
		c = append(c, &chains[i])
	}
	return c
}
