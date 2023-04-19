package config

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
)

type ClientConfiguration struct {
	ChainHost       string `json:"chain_host" mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc" mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name" mapstructure:"signer_name"`
	SignerPasswd    string
}

type EVMCommonConfig struct {
	ChainID                     int64
	BlockTimeExternalChain      uint64
	BlockTimeZetaChain          uint64
	GasPriceTicker              uint64
	ConfCount                   uint64
	ConnectorContractAddress    string
	ZETATokenContractAddress    string
	ERC20CustodyContractAddress string
}
type EVMConfig struct {
	ConnectorABI abi.ABI
	Client       *ethclient.Client
	Chain        common.Chain
	Endpoint     string
	CommonConfig *EVMCommonConfig
}

type BTCConfigConfig struct {
	WatchInTxPeriod     uint64
	WatchGasPricePeriod uint64
	WatchUTXOSPeriod    uint64
}

type BTCConfig struct {
	// the following are rpcclient ConnConfig fields
	RPCUsername string
	RPCPassword string
	RPCEndpoint string
	RPCParams   string // "regtest", "mainnet", "testnet3"

	WatchInTxPeriod     uint64
	WatchGasPricePeriod uint64
	WatchUTXOSPeriod    uint64
	CommonConfig        *BTCConfigConfig
}

type Config struct {
	ValidatorName string
	Peer          string
	LogConsole    bool
	LogLevel      zerolog.Level
	PreParamsPath string
	KeygenBlock   int64
	ChainID       string
	ZetaCoreURL   string
	AuthzGranter  string
	AuthzHotkey   string

	ChainsEnabled   []common.Chain
	EVMChainConfigs map[string]*EVMConfig // TODO : chain to chain id
	BitcoinConfig   *BTCConfig
}

func (c Config) GetAuthzHotkey() string {
	return c.AuthzHotkey
}
