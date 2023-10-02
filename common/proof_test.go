package common_test

import (
	"errors"
	"os"
	"testing"

	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/bitcoin"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
)

const numBlocksToTest = 100

type Block struct {
	TssAddress   string `json:"tssAddress"`
	Height       int    `json:"height"`
	OutTxid      string `json:"outTxid"`
	HeaderBase64 string `json:"headerBase64"`
	BlockBase64  string `json:"blockBase64"`
}

type Blocks struct {
	Blocks []Block `json:"blocks"`
}

func LoadTestBlocks() Blocks {
	file, err := os.Open("./bitcoin/test_blocks.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode the JSON into the data struct
	var blocks Blocks
	err = json.NewDecoder(file).Decode(&blocks)
	if err != nil {
		log.Fatal(err)
	}

	return blocks
}

func Test_IsErrorInvalidProof(t *testing.T) {
	require.False(t, common.IsErrorInvalidProof(nil))
	require.False(t, common.IsErrorInvalidProof(errors.New("foo")))
	require.True(t, common.IsErrorInvalidProof(common.NewErrInvalidProof(errors.New("foo"))))
}

func TestBitcoinMerkleProof(t *testing.T) {
	blocks := LoadTestBlocks()

	for _, b := range blocks.Blocks {
		// Deserialize the header bytes from base64
		headerBytes, err := base64.StdEncoding.DecodeString(b.HeaderBase64)
		if err != nil {
			t.Error(err)
		}
		header := unmarshalHeader(headerBytes)

		// Deserialize the block bytes from base64
		blockBytes, err := base64.StdEncoding.DecodeString(b.BlockBase64)
		if err != nil {
			t.Error(err)
		}
		blockVerbose := &btcjson.GetBlockVerboseTxResult{}
		err = json.Unmarshal(blockBytes, blockVerbose)
		if err != nil {
			t.Error(err)
		}

		// Validate block
		validateBitcoinBlock(t, header, headerBytes, blockVerbose, b.OutTxid, b.TssAddress)
	}
}

func TestBitcoinMerkleProofLiveTest(t *testing.T) {
	client := createBTCClient(t)
	bn, err := client.GetBlockCount()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Verifying transactions in block range [%d, %d]\n", bn-numBlocksToTest+1, bn)

	// Verify all transactions in the past 'numBlocksToTest' blocks
	for height := bn - numBlocksToTest + 1; height <= bn; height++ {
		blockHash, err := client.GetBlockHash(height)
		if err != nil {
			t.Error(err)
		}

		// Get the block header
		header, err := client.GetBlockHeader(blockHash)
		if err != nil {
			t.Error(err)
		}
		headerBytes := marshalHeader(header)
		target := blockchain.CompactToBig(header.Bits)

		// Get the block with verbose transactions
		blockVerbose, err := client.GetBlockVerboseTx(blockHash)
		if err != nil {
			t.Error(err)
		}

		// Validate block
		validateBitcoinBlock(t, header, headerBytes, blockVerbose, "", "")

		fmt.Printf("Verification succeeded for block: %d hash: %s root: %s target: %064x transactions: %d\n", height, blockHash, header.MerkleRoot, target, len(blockVerbose.Tx))
	}
}

func validateBitcoinBlock(t *testing.T, header *wire.BlockHeader, headerBytes []byte, blockVerbose *btcjson.GetBlockVerboseTxResult, outTxid string, tssAddress string) {
	// Deserialization should work for each transaction in the block
	txns := []*btcutil.Tx{}
	txBodies := [][]byte{}
	for _, res := range blockVerbose.Tx {
		txBytes, err := hex.DecodeString(res.Hex)
		if err != nil {
			log.Fatalf("error decoding transaction hex: %v", err)
		}
		tx, err := btcutil.NewTxFromBytes(txBytes)
		if err != nil {
			log.Fatalf("error deserializing transaction: %v", err)
		}

		// Validate Tss SegWit transaction if it's an outTx
		if res.Txid == outTxid {
			keeper.ValidateBTCOutTxBody(nil, txBytes, tssAddress)
		}
		txns = append(txns, tx)
		txBodies = append(txBodies, txBytes)
	}

	// Build a Merkle tree from the transaction hashes and verify each transaction
	mk := bitcoin.NewMerkle(txns)
	for i := range txns {
		path, index, err := mk.BuildMerkleProof(i)
		if err != nil {
			log.Fatalf("Error building merkle proof: %v", err)
		}

		// True proof should verify
		proof := common.NewBitcoinProof(txBodies[i], path, index)
		txBytes, err := proof.Verify(common.NewBitcoinHeader(headerBytes), 0)
		if err != nil {
			log.Fatal("Merkle proof verification failed")
		}
		if !bytes.Equal(txBytes, txBodies[i]) {
			log.Fatalf("Transaction body mismatch")
		}

		// Fake proof should not verify
		fakeIndex := index ^ 0xffffffff // flip all bits
		fakeProof := common.NewBitcoinProof(txBodies[i], path, fakeIndex)
		txBytes, err = fakeProof.Verify(common.NewBitcoinHeader(headerBytes), 0)
		if err == nil || txBytes != nil {
			log.Fatalf("Merkle proof should not verify")
		}
	}
}
