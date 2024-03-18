package bitcoin

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/require"
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

func TestBitcoinMerkleProof(t *testing.T) {
	blocks := LoadTestBlocks(t)

	t.Run("it should verify merkle proof", func(t *testing.T) {
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

			txns := getBlockTxs(t, blockVerbose)

			// Build a Merkle tree from the transaction hashes and verify each transaction
			mk := NewMerkle(txns)
			for i, tx := range txns {
				path, index, err := mk.BuildMerkleProof(i)
				require.NoError(t, err)

				// True proof should verify
				pass := Prove(*tx.Hash(), header.MerkleRoot, path, index)
				require.True(t, pass)

				// Fake proof should not verify
				fakeIndex := index ^ 0xffffffff // flip all bits
				pass = Prove(*tx.Hash(), header.MerkleRoot, path, fakeIndex)
				require.False(t, pass)
			}
		}
	})

	t.Run("it should fail if tree is empty", func(t *testing.T) {
		mt := Merkle{
			tree: []*chainhash.Hash{},
		}

		_, _, err := mt.BuildMerkleProof(0)
		require.Error(t, err)
	})

	t.Run("it should fail if tree len + 1 is not power of 2", func(t *testing.T) {
		mt := Merkle{
			tree: []*chainhash.Hash{{}, {}},
		}

		_, _, err := mt.BuildMerkleProof(0)
		require.Error(t, err)
	})

	t.Run("it should fail if txIndex out of range", func(t *testing.T) {
		mt := Merkle{
			tree: []*chainhash.Hash{{}},
		}

		_, _, err := mt.BuildMerkleProof(2)
		require.Error(t, err)
	})
}

// TODO: centralize test data
// https://github.com/zeta-chain/node/issues/1874
func LoadTestBlocks(t *testing.T) Blocks {
	file, err := os.Open("../testdata/test_blocks.json")
	require.NoError(t, err)
	defer file.Close()

	// Decode the JSON into the data struct
	var blocks Blocks
	err = json.NewDecoder(file).Decode(&blocks)
	require.NoError(t, err)

	return blocks
}

func unmarshalHeader(t *testing.T, headerBytes []byte) *wire.BlockHeader {
	var header wire.BlockHeader
	err := header.Deserialize(bytes.NewReader(headerBytes))
	require.NoError(t, err)
	return &header
}

func getBlockTxs(t *testing.T, blockVerbose *btcjson.GetBlockVerboseTxResult) []*btcutil.Tx {
	txns := []*btcutil.Tx{}
	for _, res := range blockVerbose.Tx {
		txBytes, err := hex.DecodeString(res.Hex)
		require.NoError(t, err)
		tx, err := btcutil.NewTxFromBytes(txBytes)
		require.NoError(t, err)
		txns = append(txns, tx)
	}
	return txns
}
