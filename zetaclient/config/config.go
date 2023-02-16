package config

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	common2 "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/common"
)

// ClientConfiguration
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
