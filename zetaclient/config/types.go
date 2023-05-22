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
	GasPriceTicker              uint64
	InTxTicker                  uint64
	OutTxTicker                 uint64
	WatchUTXOTicker             uint64
	ConfCount                   uint64
	ConnectorContractAddress    string
	ZETATokenContractAddress    string
	ERC20CustodyContractAddress string
	OutboundTxScheduleInterval  int64
	OutboundTxScheduleLookahead int64
}

func NewCoreParams() *CoreParams {
	return &CoreParams{
		ChainID:                     0,
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
	RPCHost     string
	RPCParams   string // "regtest", "mainnet", "testnet3"

	CoreParams *CoreParams
}

type Config struct {
	Peer          string
	PublicIP      string
	LogFormat     string
	LogLevel      zerolog.Level
	LogSampler    bool
	PreParamsPath string
	KeygenBlock   int64
	KeyGenPubKeys []string
	ChainID       string
	ZetaCoreURL   string
	AuthzGranter  string
	AuthzHotkey   string

	ChainsEnabled       []common.Chain
	EVMChainConfigs     map[int64]*EVMConfig // TODO : chain to chain id
	BitcoinConfig       *BTCConfig
	P2PDiagnostic       bool
	ConfigUpdateTicker  uint64
	P2PDiagnosticTicker uint64
	TssPath             string
}

func (c Config) GetAuthzHotkey() string {
	return c.AuthzHotkey
}

func (c Config) String() string {
	s, _ := json.MarshalIndent(c, "", "\t")
	return string(s)
}

func (cp *CoreParams) UpdateCoreParams(params *zetaObserverTypes.CoreParams) {
	cp.ChainID = params.ChainId
	cp.GasPriceTicker = params.GasPriceTicker
	cp.InTxTicker = params.InTxTicker
	cp.OutTxTicker = params.OutTxTicker
	cp.WatchUTXOTicker = params.WatchUtxoTicker
	cp.ConfCount = params.ConfirmationCount
	cp.ConnectorContractAddress = params.ConnectorContractAddress
	cp.ZETATokenContractAddress = params.ZetaTokenContractAddress
	cp.ERC20CustodyContractAddress = params.Erc20CustodyContractAddress
	cp.OutboundTxScheduleInterval = params.OutboundTxScheduleInterval
	cp.OutboundTxScheduleLookahead = params.OutboundTxScheduleLookahead
}
