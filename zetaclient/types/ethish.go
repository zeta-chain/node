package types

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeta-chain/zetacore/common"
)

type PoolContract int64

const (
	UniswapV2 PoolContract = iota
	UniswapV3
	Fixed
)

type PoolTokenOrder int64

const (
	ZETAETH PoolTokenOrder = iota
	ETHZETA
)

type ChainETHish struct {
	ConnectorABI             abi.ABI
	ChainID                  *big.Int
	ConnectorContractAddress string
	PoolContractAddress      string
	PoolContract             PoolContract
	PoolTokenOrder           PoolTokenOrder
	ZETATokenContractAddress string
	Client                   *ethclient.Client
	Name                     common.Chain
	Topics                   [][]ethcommon.Hash
	BlockTime                uint64
	Endpoint                 string
}

func BytesToEthHex(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

func HashToAddress(h ethcommon.Hash) ethcommon.Address {
	return ethcommon.BytesToAddress(h.Bytes()[12:32])
}
