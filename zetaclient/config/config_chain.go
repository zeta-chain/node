package config

import (
	"sync"

	"github.com/zeta-chain/zetacore/pkg/chains"
)

const (
	MaxBlocksPerPeriod = 100
)

const (
	// connectorAbiString is the ABI of the connector contract
	// TODO(revamp): we should be able to use info from Go binding
	connectorAbiString = `
[{"inputs":[{"internalType":"address","name":"_zetaTokenAddress","type":"address"},{"internalType":"address","name":"_tssAddress","type":"address"},{"internalType":"address","name":"_tssAddressUpdater","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"account","type":"address"}],"name":"Unpaused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"originSenderAddress","type":"bytes"},{"indexed":true,"internalType":"uint256","name":"originChainId","type":"uint256"},{"indexed":true,"internalType":"address","name":"destinationAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":true,"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"ZetaReceived","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"originSenderAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"originChainId","type":"uint256"},{"indexed":true,"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"indexed":true,"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":true,"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"ZetaReverted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"originSenderAddress","type":"address"},{"indexed":false,"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"gasLimit","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"zetaParams","type":"bytes"}],"name":"ZetaSent","type":"event"},{"inputs":[],"name":"getLockedAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"originSenderAddress","type":"bytes"},{"internalType":"uint256","name":"originChainId","type":"uint256"},{"internalType":"address","name":"destinationAddress","type":"address"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"onReceive","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"originSenderAddress","type":"address"},{"internalType":"uint256","name":"originChainId","type":"uint256"},{"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"bytes32","name":"internalSendHash","type":"bytes32"}],"name":"onRevert","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"pause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceTssAddressUpdater","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"uint256","name":"destinationChainId","type":"uint256"},{"internalType":"bytes","name":"destinationAddress","type":"bytes"},{"internalType":"uint256","name":"gasLimit","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"},{"internalType":"uint256","name":"zetaAmount","type":"uint256"},{"internalType":"bytes","name":"zetaParams","type":"bytes"}],"internalType":"struct ZetaInterfaces.SendInput","name":"input","type":"tuple"}],"name":"send","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"tssAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"tssAddressUpdater","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"unpause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_tssAddress","type":"address"}],"name":"updateTssAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"zetaToken","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"}]`

	// erc20CustodyAbiString is the ABI of the erc20 custodu contract
	// TODO(revamp): we should be able to use info from Go binding
	erc20CustodyAbiString = `
[{"inputs":[{"internalType":"address","name":"_TSSAddress","type":"address"},{"internalType":"address","name":"_TSSAddressUpdater","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"InvalidSender","type":"error"},{"inputs":[],"name":"InvalidTSSUpdater","type":"error"},{"inputs":[],"name":"IsPaused","type":"error"},{"inputs":[],"name":"NotPaused","type":"error"},{"inputs":[],"name":"NotWhitelisted","type":"error"},{"inputs":[],"name":"ZeroAddress","type":"error"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes","name":"recipient","type":"bytes"},{"indexed":false,"internalType":"address","name":"asset","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":false,"internalType":"bytes","name":"message","type":"bytes"}],"name":"Deposited","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"sender","type":"address"}],"name":"Paused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"sender","type":"address"}],"name":"Unpaused","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"asset","type":"address"}],"name":"Unwhitelisted","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"asset","type":"address"}],"name":"Whitelisted","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"recipient","type":"address"},{"indexed":false,"internalType":"address","name":"asset","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"}],"name":"Withdrawn","type":"event"},{"inputs":[],"name":"TSSAddress","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"TSSAddressUpdater","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"recipient","type":"bytes"},{"internalType":"address","name":"asset","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"bytes","name":"message","type":"bytes"}],"name":"deposit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"pause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"paused","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"renounceTSSAddressUpdater","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"unpause","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"asset","type":"address"}],"name":"unwhitelist","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"_address","type":"address"}],"name":"updateTSSAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"asset","type":"address"}],"name":"whitelist","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"}],"name":"whitelisted","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"address","name":"asset","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

// GetConnectorABI returns the ABI of the connector contract
func GetConnectorABI() string {
	return connectorAbiString
}

// GetERC20CustodyABI returns the ABI of the erc20 custody contract
func GetERC20CustodyABI() string {
	return erc20CustodyAbiString
}

// New constructs Config optionally with default values.
func New(setDefaults bool) Config {
	cfg := Config{
		EVMChainConfigs: make(map[int64]EVMConfig),
		BitcoinConfig:   BTCConfig{},

		mu: &sync.RWMutex{},
	}

	if setDefaults {
		cfg.BitcoinConfig = bitcoinConfigRegnet()
		cfg.EVMChainConfigs = evmChainsConfigs()
	}

	return cfg
}

// bitcoinConfigRegnet contains Bitcoin config for regnet
func bitcoinConfigRegnet() BTCConfig {
	return BTCConfig{
		RPCUsername: "smoketest", // smoketest is the previous name for E2E test, we keep this name for compatibility between client versions in upgrade test
		RPCPassword: "123",
		RPCHost:     "bitcoin:18443",
		RPCParams:   "regtest",
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
