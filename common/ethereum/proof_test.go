package ethereum

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

const (
	headerPath        = "./testdata/header.json"
	receiptPrefixPath = "./testdata/receipt_"
	receiptCount      = 81
)

func TestProofGeneration(t *testing.T) {
	header, err := readHeader()
	require.NoError(t, err)

	var receipts types.Receipts
	for i := 0; i < receiptCount; i++ {
		receipt, err := readReceipt(i)
		require.NoError(t, err)
		receipts = append(receipts, &receipt)
	}

	// generate a trie from the receipts and compare the root hash with the one in the header
	receiptTree := NewTrie(receipts)
	require.EqualValues(t, header.ReceiptHash.Hex(), header.ReceiptHash.Hex())

	for i, receipt := range receipts {
		// generate a proof for each receipt and verify it
		proof, err := receiptTree.GenerateProof(i)
		require.NoError(t, err)

		verified, err := proof.Verify(header.ReceiptHash, i)
		require.NoError(t, err)

		// recover the receipt from the proof and compare it with the original receipt
		// NOTE: eth receipts only hashes the following fields
		// data := &receiptRLP{r.statusEncoding(), r.CumulativeGasUsed, r.Bloom, r.Logs}
		var verifiedReceipt types.Receipt
		err = verifiedReceipt.UnmarshalBinary(verified)
		require.NoError(t, err)
		require.EqualValues(t, receipt.Status, verifiedReceipt.Status)
		require.EqualValues(t, receipt.CumulativeGasUsed, verifiedReceipt.CumulativeGasUsed)
		require.EqualValues(t, receipt.Bloom.Bytes(), verifiedReceipt.Bloom.Bytes())
		require.EqualValues(t, len(receipt.Logs), len(verifiedReceipt.Logs))
		for i, log := range receipt.Logs {
			require.EqualValues(t, log.Address, verifiedReceipt.Logs[i].Address)
			require.EqualValues(t, log.Topics, verifiedReceipt.Logs[i].Topics)
			require.EqualValues(t, log.Data, verifiedReceipt.Logs[i].Data)
		}
	}

}

// readHeader reads a header from a file.
// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func readHeader() (header types.Header, err error) {
	file, err := os.Open(headerPath)
	if err != nil {
		return header, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&header)
	return header, err
}

// readReceipt reads a receipt from a file.
// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func readReceipt(index int) (receipt types.Receipt, err error) {
	filePath := fmt.Sprintf("%s%d.json", receiptPrefixPath, index)

	file, err := os.Open(filePath)
	if err != nil {
		return receipt, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&receipt)
	return receipt, err
}
