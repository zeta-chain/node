package config

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
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

type ChainETHish struct {
	ConnectorABI                abi.ABI
	ConnectorContractAddress    string
	ZETATokenContractAddress    string
	ERC20CustodyContractAddress string
	Client                      *ethclient.Client
	Chain                       common.Chain
	Topics                      [][]common2.Hash
	BlockTime                   uint64
	ConfCount                   uint64
	Endpoint                    string
	OutTxObservePeriod          uint64
}

type ChainBitcoinish struct {
	// the following are rpcclient ConnConfig fields
	RPCUsername string
	RPCPassword string
	RPCEndpoint string
	RPCParams   string // "regtest", "mainnet", "testnet3"

	WatchInTxPeriod     uint64
	WatchGasPricePeriod uint64
	WatchUTXOSPeriod    uint64
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

	ChainsEnabled []common.Chain
	ChainConfigs  map[string]*ChainETHish
	BitcoinConfig *ChainBitcoinish
}

func (c Config) GetAuthzHotkey() string {
	return c.AuthzHotkey
}
