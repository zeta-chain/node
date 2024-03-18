package common_test

import (
	"errors"
	"os"
	"testing"

	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/bitcoin"
	"github.com/zeta-chain/zetacore/common/ethereum"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

const (
	headerPath        = "./ethereum/testdata/header.json"
	txPrefixPath      = "./ethereum/testdata/tx_"
	receiptPrefixPath = "./ethereum/testdata/receipt_"
	txCount           = 81
	numBlocksToTest   = 100
)

type Block struct {
	TssAddress   string `json:"tssAddress"`
	Height       int    `json:"height"`
	Nonce        uint64 `json:"nonce"`
	OutTxid      string `json:"outTxid"`
	HeaderBase64 string `json:"headerBase64"`
	BlockBase64  string `json:"blockBase64"`
}

type Blocks struct {
	Blocks []Block `json:"blocks"`
}

// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func LoadTestBlocks(t *testing.T) Blocks {
	file, err := os.Open("./testdata/test_blocks.json")
	require.NoError(t, err)
	defer file.Close()

	// Decode the JSON into the data struct
	var blocks Blocks
	err = json.NewDecoder(file).Decode(&blocks)
	require.NoError(t, err)

	return blocks
}

func Test_IsErrorInvalidProof(t *testing.T) {
	require.False(t, common.IsErrorInvalidProof(nil))
	require.False(t, common.IsErrorInvalidProof(errors.New("foo")))
	invalidProofErr := errors.New("foo")
	invalidProof := common.NewErrInvalidProof(invalidProofErr)
	require.True(t, common.IsErrorInvalidProof(invalidProof))
	require.Equal(t, invalidProofErr.Error(), invalidProof.Error())
}

func TestBitcoinMerkleProof(t *testing.T) {
	blocks := LoadTestBlocks(t)

	for _, b := range blocks.Blocks {
		// Deserialize the header bytes from base64
		headerBytes, err := base64.StdEncoding.DecodeString(b.HeaderBase64)
		require.NoError(t, err)
		header := unmarshalHeader(t, headerBytes)

		// Deserialize the block bytes from base64
		blockBytes, err := base64.StdEncoding.DecodeString(b.BlockBase64)
		require.NoError(t, err)
		blockVerbose := &btcjson.GetBlockVerboseTxResult{}
		err = json.Unmarshal(blockBytes, blockVerbose)
		require.NoError(t, err)

		// Validate block
		validateBitcoinBlock(t, header, headerBytes, blockVerbose, b.OutTxid, b.TssAddress, b.Nonce)
	}
}

func TestEthereumMerkleProof(t *testing.T) {
	header, err := readHeader()
	require.NoError(t, err)
	b, err := rlp.EncodeToBytes(&header)
	require.NoError(t, err)

	headerData := common.NewEthereumHeader(b)
	t.Run("should verify tx proof", func(t *testing.T) {
		var txs types.Transactions
		for i := 0; i < txCount; i++ {
			tx, err := readTx(i)
			require.NoError(t, err)
			txs = append(txs, &tx)
		}

		// generate a trie from the txs and compare the root hash with the one in the header
		txsTree := ethereum.NewTrie(txs)
		require.EqualValues(t, header.TxHash.Hex(), txsTree.Trie.Hash().Hex())

		for i := range txs {
			// generate a proof for each tx and verify it
			proof, err := txsTree.GenerateProof(i)
			require.NoError(t, err)

			ethProof := common.NewEthereumProof(proof)

			_, err = ethProof.Verify(headerData, i)
			require.NoError(t, err)
		}
	})

	t.Run("should fail to verify receipts proof", func(t *testing.T) {
		var receipts types.Receipts
		for i := 0; i < txCount; i++ {
			receipt, err := readReceipt(i)
			require.NoError(t, err)
			receipts = append(receipts, &receipt)
		}

		// generate a trie from the receipts and compare the root hash with the one in the header
		txsTree := ethereum.NewTrie(receipts)
		require.EqualValues(t, header.ReceiptHash.Hex(), txsTree.Trie.Hash().Hex())

		for i := range receipts {
			// generate a proof for each receipt and verify it
			proof, err := txsTree.GenerateProof(i)
			require.NoError(t, err)

			ethProof := common.NewEthereumProof(proof)

			_, err = ethProof.Verify(headerData, i)
			require.Error(t, err)
		}
	})
}

func BitcoinMerkleProofLiveTest(t *testing.T) {
	client := createBTCClient(t)
	bn, err := client.GetBlockCount()
	require.NoError(t, err)
	fmt.Printf("Verifying transactions in block range [%d, %d]\n", bn-numBlocksToTest+1, bn)

	// Verify all transactions in the past 'numBlocksToTest' blocks
	for height := bn - numBlocksToTest + 1; height <= bn; height++ {
		blockHash, err := client.GetBlockHash(height)
		require.NoError(t, err)

		// Get the block header
		header, err := client.GetBlockHeader(blockHash)
		require.NoError(t, err)
		headerBytes := marshalHeader(t, header)
		target := blockchain.CompactToBig(header.Bits)

		// Get the block with verbose transactions
		blockVerbose, err := client.GetBlockVerboseTx(blockHash)
		require.NoError(t, err)

		// Validate block
		validateBitcoinBlock(t, header, headerBytes, blockVerbose, "", "", 0)

		fmt.Printf("Verification succeeded for block: %d hash: %s root: %s target: %064x transactions: %d\n", height, blockHash, header.MerkleRoot, target, len(blockVerbose.Tx))
	}
}

func validateBitcoinBlock(t *testing.T, _ *wire.BlockHeader, headerBytes []byte, blockVerbose *btcjson.GetBlockVerboseTxResult, outTxid string, tssAddress string, nonce uint64) {
	// Deserialization should work for each transaction in the block
	txns := []*btcutil.Tx{}
	txBodies := [][]byte{}
	for _, res := range blockVerbose.Tx {
		txBytes, err := hex.DecodeString(res.Hex)
		require.NoError(t, err)
		tx, err := btcutil.NewTxFromBytes(txBytes)
		require.NoError(t, err)

		// Validate Tss SegWit transaction if it's an outTx
		if res.Txid == outTxid {
			msg := &crosschaintypes.MsgAddToOutTxTracker{
				ChainId: common.BtcTestNetChain().ChainId,
				Nonce:   nonce,
				TxHash:  outTxid,
			}
			err = keeper.VerifyBTCOutTxBody(msg, txBytes, tssAddress)
			require.NoError(t, err)
		}
		txns = append(txns, tx)
		txBodies = append(txBodies, txBytes)
	}

	// Build a Merkle tree from the transaction hashes and verify each transaction
	mk := bitcoin.NewMerkle(txns)
	for i := range txns {
		path, index, err := mk.BuildMerkleProof(i)
		require.NoError(t, err)

		// True proof should verify
		proof := common.NewBitcoinProof(txBodies[i], path, index)
		txBytes, err := proof.Verify(common.NewBitcoinHeader(headerBytes), 0)
		require.NoError(t, err)
		require.Equal(t, txBytes, txBodies[i])

		// Fake proof should not verify
		fakeIndex := index ^ 0xffffffff // flip all bits
		fakeProof := common.NewBitcoinProof(txBodies[i], path, fakeIndex)
		txBytes, err = fakeProof.Verify(common.NewBitcoinHeader(headerBytes), 0)
		require.Error(t, err)
		require.Nil(t, txBytes)
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

// readTx reads a tx from a file.
// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func readTx(index int) (tx types.Transaction, err error) {
	filePath := fmt.Sprintf("%s%d.json", txPrefixPath, index)

	file, err := os.Open(filePath)
	if err != nil {
		return tx, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tx)
	return tx, err
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
