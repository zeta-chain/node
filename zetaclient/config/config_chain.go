package config

import (
	"sync"

	"github.com/zeta-chain/node/pkg/chains"
)

const (
	MaxBlocksPerPeriod = 100
)

// New constructs Config optionally with default values.
func New(setDefaults bool) Config {
	cfg := Config{
		EVMChainConfigs: make(map[int64]EVMConfig),
		BTCChainConfigs: make(map[int64]BTCConfig),

		mu: &sync.RWMutex{},
	}

	if setDefaults {
		cfg.EVMChainConfigs = evmChainsConfigs()
		cfg.BTCChainConfigs = btcChainsConfigs()
		cfg.SolanaConfig = solanaConfigLocalnet()
	}

	return cfg
}

// bitcoinConfigRegnet contains Bitcoin config for regnet
func bitcoinConfigRegnet() BTCConfig {
	return BTCConfig{
		RPCUsername: "e2etest",
		RPCPassword: "123",
		RPCHost:     "bitcoin:18443",
		RPCParams:   "regtest",
	}
}

// solanaConfigLocalnet contains config for Solana localnet
func solanaConfigLocalnet() SolanaConfig {
	return SolanaConfig{
		Endpoint: "http://solana:8899",
	}
}

// evmChainsConfigs contains EVM chain configs
// it contains list of EVM chains with empty endpoint except for localnet
func evmChainsConfigs() map[int64]EVMConfig {
	return map[int64]EVMConfig{
		chains.Ethereum.ChainId: {
			Chain: chains.Ethereum,
		},
		chains.BscMainnet.ChainId: {
			Chain: chains.BscMainnet,
		},
		chains.Goerli.ChainId: {
			Chain:    chains.Goerli,
			Endpoint: "",
		},
		chains.Sepolia.ChainId: {
			Chain:    chains.Sepolia,
			Endpoint: "",
		},
		chains.BscTestnet.ChainId: {
			Chain:    chains.BscTestnet,
			Endpoint: "",
		},
		chains.Mumbai.ChainId: {
			Chain:    chains.Mumbai,
			Endpoint: "",
		},
		chains.GoerliLocalnet.ChainId: {
			Chain:    chains.GoerliLocalnet,
			Endpoint: "http://eth:8545",
		},
	}
}

// btcChainsConfigs contains BTC chain configs
func btcChainsConfigs() map[int64]BTCConfig {
	return map[int64]BTCConfig{
		chains.BitcoinRegtest.ChainId: bitcoinConfigRegnet(),
	}
}
