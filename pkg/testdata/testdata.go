package testdata

import (
	"embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

const (
	HeaderPath        = "ethereum/header.json"
	ReceiptPrefixPath = "ethereum/receipt_"
	TxPrefixPath      = "ethereum/tx_"
	TxsCount          = 81
)

//go:embed ethereum/*
var ethFiles embed.FS

//go:embed *
var testDataFiles embed.FS

// ReadEthHeader reads a header from a file.
// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func ReadEthHeader() (header types.Header, err error) {
	file, err := ethFiles.Open(HeaderPath)
	if err != nil {
		return header, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&header)
	return header, err
}

// ReadEthReceipt reads a receipt from a file.
// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func ReadEthReceipt(index int) (receipt types.Receipt, err error) {
	filePath := fmt.Sprintf("%s%d.json", ReceiptPrefixPath, index)

	file, err := ethFiles.Open(filePath)
	if err != nil {
		return receipt, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&receipt)
	return receipt, err
}

// ReadEthTx reads a tx from a file.
// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func ReadEthTx(index int) (tx types.Transaction, err error) {
	filePath := fmt.Sprintf("%s%d.json", TxPrefixPath, index)

	file, err := ethFiles.Open(filePath)
	if err != nil {
		return tx, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tx)
	return tx, err
}

type Block struct {
	TssAddress   string `json:"tssAddress"`
	Height       int    `json:"height"`
	Nonce        uint64 `json:"nonce"`
	Outboundid   string `json:"outTxid"`
	HeaderBase64 string `json:"headerBase64"`
	BlockBase64  string `json:"blockBase64"`
}

type Blocks struct {
	Blocks []Block `json:"blocks"`
}

// LoadTestBlocks loads test blocks from a file.
// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func LoadTestBlocks(t *testing.T) Blocks {
	file, err := testDataFiles.Open("test_blocks.json")
	require.NoError(t, err)
	defer file.Close()

	// Decode the JSON into the data struct
	var blocks Blocks
	err = json.NewDecoder(file).Decode(&blocks)
	require.NoError(t, err)

	return blocks
}
