package proofs

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/testdata"
)

const numHeadersToTest = 100

func TestTrueEthereumHeader(t *testing.T) {
	var header ethtypes.Header
	// read file into a byte slice
	file, err := os.Open("../../testutil/testdata/eth_header_18495266.json")
	require.NoError(t, err)
	defer file.Close()
	headerBytes := make([]byte, 4096)
	n, err := file.Read(headerBytes)
	require.NoError(t, err)
	err = header.UnmarshalJSON(headerBytes[:n])
	require.NoError(t, err)
	var buffer bytes.Buffer
	err = header.EncodeRLP(&buffer)
	require.NoError(t, err)

	headerData := NewEthereumHeader(buffer.Bytes())
	err = headerData.Validate(header.Hash().Bytes(), 1, 18495266)
	require.NoError(t, err)

	parentHash, err := headerData.ParentHash()
	require.NoError(t, err)
	require.Equal(t, header.ParentHash.Bytes(), parentHash)

	err = headerData.ValidateTimestamp(time.Now())
	require.NoError(t, err)
}

func TestFalseEthereumHeader(t *testing.T) {
	var header ethtypes.Header
	// read file into a byte slice
	file, err := os.Open("../../testutil/testdata/eth_header_18495266.json")
	require.NoError(t, err)
	defer file.Close()
	headerBytes := make([]byte, 4096)
	n, err := file.Read(headerBytes)
	require.NoError(t, err)
	err = header.UnmarshalJSON(headerBytes[:n])
	require.NoError(t, err)
	hash := header.Hash()
	header.Number = big.NewInt(18495267)
	var buffer bytes.Buffer
	err = header.EncodeRLP(&buffer)
	require.NoError(t, err)

	headerData := NewEthereumHeader(buffer.Bytes())
	err = headerData.Validate(hash.Bytes(), 1, 18495267)
	require.Error(t, err)
}

func TestTrueBitcoinHeader(t *testing.T) {
	blocks := testdata.LoadTestBlocks(t)

	for _, b := range blocks.Blocks {
		// Deserialize the header bytes from base64
		headerBytes, err := base64.StdEncoding.DecodeString(b.HeaderBase64)
		require.NoError(t, err)
		header := unmarshalHeader(t, headerBytes)

		// Validate
		validateTrueBitcoinHeader(t, header, headerBytes)
	}
}

func TestFakeBitcoinHeader(t *testing.T) {
	blocks := testdata.LoadTestBlocks(t)

	for _, b := range blocks.Blocks {
		// Deserialize the header bytes from base64
		headerBytes, err := base64.StdEncoding.DecodeString(b.HeaderBase64)
		require.NoError(t, err)
		header := unmarshalHeader(t, headerBytes)

		// Validate
		validateFakeBitcoinHeader(t, header, headerBytes)
	}
}

func TestNonExistentHeaderType(t *testing.T) {
	headerData := HeaderData{}

	err := headerData.Validate([]byte{}, 18332, 0)
	require.EqualError(t, err, "unrecognized header type")

	parentHash, err := headerData.ParentHash()
	require.EqualError(t, err, "unrecognized header type")
	require.Nil(t, parentHash)

	err = headerData.ValidateTimestamp(time.Now())
	require.ErrorContains(t, err, "unrecognized header type")
}

func BitcoinHeaderValidationLiveTest(t *testing.T) {
	client := createBTCClient(t)
	bn, err := client.GetBlockCount()
	require.NoError(t, err)
	fmt.Printf("Verifying block headers in block range [%d, %d]\n", bn-numHeadersToTest+1, bn)

	for height := bn - numHeadersToTest + 1; height <= bn; height++ {
		blockHash, err := client.GetBlockHash(height)
		require.NoError(t, err)

		// Get the block header
		header, err := client.GetBlockHeader(blockHash)
		require.NoError(t, err)
		headerBytes := marshalHeader(t, header)

		// Validate true header
		validateTrueBitcoinHeader(t, header, headerBytes)

		// Validate fake header
		validateFakeBitcoinHeader(t, header, headerBytes)

		fmt.Printf("Block header verified for block: %d hash: %s\n", height, blockHash)
	}
}

func createBTCClient(t *testing.T) *rpcclient.Client {
	connCfg := &rpcclient.ConnConfig{
		Host:         "127.0.0.1:18332",
		User:         "user",
		Pass:         "pass",
		HTTPPostMode: true,
		DisableTLS:   true,
		Params:       "testnet3",
	}
	client, err := rpcclient.New(connCfg, nil)
	require.NoError(t, err)
	return client
}

func copyHeader(header *wire.BlockHeader) *wire.BlockHeader {
	copyHeader := &wire.BlockHeader{
		Version:    header.Version,
		PrevBlock:  chainhash.Hash{},
		MerkleRoot: chainhash.Hash{},
		Timestamp:  header.Timestamp,
		Bits:       header.Bits,
		Nonce:      header.Nonce,
	}
	copy(copyHeader.PrevBlock[:], header.PrevBlock[:])
	copy(copyHeader.MerkleRoot[:], header.MerkleRoot[:])

	return copyHeader
}

func marshalHeader(t *testing.T, header *wire.BlockHeader) []byte {
	var headerBuf bytes.Buffer
	err := header.Serialize(&headerBuf)
	require.NoError(t, err)
	return headerBuf.Bytes()
}

func unmarshalHeader(t *testing.T, headerBytes []byte) *wire.BlockHeader {
	var header wire.BlockHeader
	err := header.Deserialize(bytes.NewReader(headerBytes))
	require.NoError(t, err)
	return &header
}

func validateTrueBitcoinHeader(t *testing.T, header *wire.BlockHeader, headerBytes []byte) {
	blockHash := header.BlockHash()

	headerData := NewBitcoinHeader(headerBytes)
	// True Bitcoin header should pass validation
	err := headerData.Validate(blockHash[:], 18332, 0)
	require.NoError(t, err)

	// True Bitcoin header should pass timestamp validation
	err = headerData.ValidateTimestamp(time.Now())
	require.NoError(t, err)

	parentHash, err := headerData.ParentHash()
	require.NoError(t, err)
	require.Equal(t, header.PrevBlock.CloneBytes(), parentHash)
}

func validateFakeBitcoinHeader(t *testing.T, header *wire.BlockHeader, headerBytes []byte) {
	blockHash := header.BlockHash()

	// Incorrect header length should fail validation
	err := ValidateBitcoinHeader(headerBytes[:79], blockHash[:], 18332)
	require.Error(t, err)

	// Incorrect version should fail validation
	fakeHeader := copyHeader(header)
	fakeHeader.Version = 0
	fakeBytes := marshalHeader(t, fakeHeader)
	fakeHash := fakeHeader.BlockHash()
	err = ValidateBitcoinHeader(fakeBytes, fakeHash[:], 18332)
	require.Error(t, err)

	// Incorrect timestamp should fail validation
	// Case1: timestamp is before genesis block
	fakeHeader = copyHeader(header)
	fakeHeader.Timestamp = chaincfg.TestNet3Params.GenesisBlock.Header.Timestamp.Add(-time.Second)
	fakeBytes = marshalHeader(t, fakeHeader)
	fakeHash = fakeHeader.BlockHash()
	err = ValidateBitcoinHeader(fakeBytes, fakeHash[:], 18332)
	require.Error(t, err)
	// Case2: timestamp is after 2 hours in the future
	fakeHeader = copyHeader(header)
	fakeHeader.Timestamp = header.Timestamp.Add(time.Second * (blockchain.MaxTimeOffsetSeconds + 1))
	fakeBytes = marshalHeader(t, fakeHeader)
	err = NewBitcoinHeader(fakeBytes).ValidateTimestamp(header.Timestamp)
	require.Error(t, err)

	// Incorrect block hash should fail validation
	fakeHeader = copyHeader(header)
	header.Nonce = 0
	fakeBytes = marshalHeader(t, header)
	err = ValidateBitcoinHeader(fakeBytes, blockHash[:], 18332)
	require.Error(t, err)

	// PoW not satisfied should fail validation
	fakeHash = fakeHeader.BlockHash()
	err = ValidateBitcoinHeader(fakeBytes, fakeHash[:], 18332)
	require.Error(t, err)
}
