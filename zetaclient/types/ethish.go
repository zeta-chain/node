package types

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/common"
)

type ChainETHish struct {
	ConnectorABI                abi.ABI
	ConnectorContractAddress    string
	ZETATokenContractAddress    string
	ERC20CustodyContractAddress string
	Client                      *ethclient.Client
	Chain                       common.Chain
	Topics                      [][]ethcommon.Hash
	BlockTime                   uint64
	Endpoint                    string
	OutTxObservePeriod          uint64
}

func BytesToEthHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}
