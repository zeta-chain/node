package config

import (
	"encoding/json"
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

type CoreParams struct {
	ChainID                     int64
	BlockTimeExternalChain      uint64
	BlockTimeZetaChain          uint64
	GasPriceTicker              uint64
	InTxTicker                  uint64
	OutTxTicker                 uint64
	WatchUTXOTicker             uint64
	ConfCount                   uint64
	ConnectorContractAddress    string
	ZETATokenContractAddress    string
	ERC20CustodyContractAddress string
}

func NewCoreParams() *CoreParams {
	return &CoreParams{
		ChainID:                     0,
		BlockTimeExternalChain:      0,
		BlockTimeZetaChain:          0,
		GasPriceTicker:              0,
		InTxTicker:                  5,
		OutTxTicker:                 3,
		WatchUTXOTicker:             5,
		ConfCount:                   0,
		ConnectorContractAddress:    "",
		ZETATokenContractAddress:    "",
		ERC20CustodyContractAddress: "",
	}
}

func (c *CoreParams) UpdateFromCoreResponse(newConfig zetaObserverTypes.CoreParams) {
	c.BlockTimeZetaChain = newConfig.BlockTimeZeta
	c.BlockTimeExternalChain = newConfig.BlockTimeExternal
	c.GasPriceTicker = newConfig.GasPriceTicker
	c.ConfCount = newConfig.ConfirmationCount
	c.ConnectorContractAddress = newConfig.ConnectorContractAddress
	c.ZETATokenContractAddress = newConfig.ZetaTokenContractAddress
	c.ERC20CustodyContractAddress = newConfig.Erc20CustodyContractAddress
}

type EVMConfig struct {
	Client     *ethclient.Client
	Chain      common.Chain
	Endpoint   string
	CoreParams *CoreParams
}

type BTCConfig struct {
	// the following are rpcclient ConnConfig fields
	RPCUsername string
	RPCPassword string
	RPCEndpoint string
	RPCParams   string // "regtest", "mainnet", "testnet3"

	CoreParams *CoreParams
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
	EVMChainConfigs map[int64]*EVMConfig // TODO : chain to chain id
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

func (cp *CoreParams) UpdateCoreParams(params *zetaObserverTypes.CoreParams) {
	cp.ChainID = params.ChainId
	cp.BlockTimeExternalChain = params.BlockTimeExternal
	cp.BlockTimeZetaChain = params.BlockTimeZeta
	cp.GasPriceTicker = params.GasPriceTicker
	cp.InTxTicker = params.InTxTicker
	cp.OutTxTicker = params.OutTxTicker
	cp.WatchUTXOTicker = params.WatchUtxoTicker
	cp.ConfCount = params.ConfirmationCount
	cp.ConnectorContractAddress = params.ConnectorContractAddress
	cp.ZETATokenContractAddress = params.ZetaTokenContractAddress
	cp.ERC20CustodyContractAddress = params.Erc20CustodyContractAddress

}
