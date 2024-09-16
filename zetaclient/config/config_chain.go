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
		BitcoinConfig:   BTCConfig{},

		mu: &sync.RWMutex{},
	}

	if setDefaults {
		cfg.BitcoinConfig = bitcoinConfigRegnet()
		cfg.SolanaConfig = solanaConfigLocalnet()
		cfg.EVMChainConfigs = evmChainsConfigs()
	}

	return cfg
}

// bitcoinConfigRegnet contains Bitcoin config for regnet
func bitcoinConfigRegnet() BTCConfig {
	return BTCConfig{
		RPCUsername:     "smoketest", // smoketest is the previous name for E2E test, we keep this name for compatibility between client versions in upgrade test
		RPCPassword:     "123",
		RPCHost:         "bitcoin:18443",
		RPCParams:       "regtest",
		RPCAlertLatency: 60,
	}
}

// solanaConfigLocalnet contains config for Solana localnet
func solanaConfigLocalnet() SolanaConfig {
	return SolanaConfig{
		Endpoint:        "http://solana:8899",
		RPCAlertLatency: 60,
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
			Chain:           chains.GoerliLocalnet,
			Endpoint:        "http://eth:8545",
			RPCAlertLatency: 60,
		},
	}
}
