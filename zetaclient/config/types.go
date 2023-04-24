package config

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/zeta-chain/zetacore/common"
	zetaObserverTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

type ClientConfiguration struct {
	ChainHost       string `json:"chain_host" mapstructure:"chain_host"`
	ChainRPC        string `json:"chain_rpc" mapstructure:"chain_rpc"`
	ChainHomeFolder string `json:"chain_home_folder" mapstructure:"chain_home_folder"`
	SignerName      string `json:"signer_name" mapstructure:"signer_name"`
	SignerPasswd    string
}

type EVMChainsCoreParams struct {
	ChainID                     int64
	BlockTimeExternalChain      uint64
	BlockTimeZetaChain          uint64
	GasPriceTicker              uint64
	ConfCount                   uint64
	ConnectorContractAddress    string
	ZETATokenContractAddress    string
	ERC20CustodyContractAddress string
}

func NewCoreParams() *EVMChainsCoreParams {
	return &EVMChainsCoreParams{
		ChainID:                     0,
		BlockTimeExternalChain:      0,
		BlockTimeZetaChain:          0,
		GasPriceTicker:              0,
		ConfCount:                   0,
		ConnectorContractAddress:    "",
		ZETATokenContractAddress:    "",
		ERC20CustodyContractAddress: "",
	}
}

func (c *EVMChainsCoreParams) UpdateFromCoreResponse(newConfig zetaObserverTypes.ClientParams) {
	c.BlockTimeZetaChain = newConfig.BlockTimeZeta
	c.BlockTimeExternalChain = newConfig.BlockTimeExternal
	c.GasPriceTicker = newConfig.GasPriceTicker
	c.ConfCount = newConfig.ConfirmationCount
	c.ConnectorContractAddress = newConfig.ConnectorContractAddress
	c.ZETATokenContractAddress = newConfig.ZetaTokenContractAddress
	c.ERC20CustodyContractAddress = newConfig.Erc20CustodyContractAddress
}

type EVMConfig struct {
	ConnectorABI abi.ABI
	Client       *ethclient.Client
	Chain        common.Chain
	Endpoint     string
	CoreParams   *EVMChainsCoreParams
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

func (c Config) String() string {
	s, _ := json.MarshalIndent(c, "", "\t")
	return string(s)
}

func (c Config) PrintEVMConfigs() string {
	s, _ := json.MarshalIndent(c.EVMChainConfigs, "", "\t")
	return string(s)
}

func (c Config) PrintBTCConfigs() string {
	s, _ := json.MarshalIndent(c.BitcoinConfig, "", "\t")
	return string(s)
}

func (c Config) PrintSupportedChains() string {
	s, _ := json.MarshalIndent(c.ChainsEnabled, "", "\t")
	return string(s)
}
