package ethereum

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func TestProofGeneration(t *testing.T) {
	RPC_URL := "https://rpc.ankr.com/eth_sepolia"
	client, err := ethclient.Dial(RPC_URL)
	if err != nil {
		t.Fatal(err)
	}
	bn := int64(5000000)
	block, err := client.BlockByNumber(context.Background(), big.NewInt(bn))
	if err != nil {
		t.Fatal(err)
	}

	// Write the block to a JSON file
	if err := writeJSONToFile(block.Header(), "/Users/lucas/Desktop/header.json"); err != nil {
		t.Fatal(err)
	}

	headerRLP, _ := rlp.EncodeToBytes(block.Header())
	t.Logf("block header size %d\n", len(headerRLP))

	var header types.Header
	rlp.DecodeBytes(headerRLP, &header)

	t.Logf("block %d\n", block.Number())
	t.Logf("  tx root %x\n", header.TxHash)

	//ttt := new(trie.Trie)
	tr := NewTrie(block.Transactions())
	t.Logf("  sha2    %x\n", tr.Hash())
	if tr.Hash() != header.TxHash {
		t.Fatal("tx root mismatch")
	} else {
		t.Logf("  tx root hash & block tx root match\n")
	}

	var indexBuf []byte
	for i := 0; i < len(block.Transactions()); i++ {

		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))

		proof := NewProof()
		tr.Prove(indexBuf, 0, proof)
		t.Logf("proof len %d\n", len(proof.Keys))
		value, err := proof.Verify(block.Header().TxHash, i)
		//value, err := trie.VerifyProof(tr.trie.Hash(), indexBuf, proof)
		t.Logf("pass? %v\n", err == nil)
		//t.Logf("value %x\n", value)

		var txx types.Transaction
		txx.UnmarshalBinary(value)
		t.Logf("  tx       %+v\n", txx.To().Hex())
		t.Logf("  tx  hash %+v\n", txx.Hash().Hex())
		if txx.Hash() != block.Transactions()[i].Hash() {
			t.Fatal("tx hash mismatch")
		} else {
			t.Logf("  tx hash & block tx hash match\n")
		}
		signer := types.NewLondonSigner(txx.ChainId())
		sender, err := types.Sender(signer, &txx)
		t.Logf("  tx from %s\n", sender.Hex())
	}

	//for k, v := range proof.Proof {
	//	key, _ := base64.StdEncoding.DecodeString(k)
	//	t.Logf("k: %x, v: %x\n", key, v)
	//}

	{
		var receipts types.Receipts
		for i, tx := range block.Transactions() {
			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				t.Fatal(err)
			}
			receipts = append(receipts, receipt)
			time.Sleep(200 * time.Millisecond)

			// Write the receipt to a JSON file
			if i == 0 {
				if err := writeJSONToFile(receipt, "/Users/lucas/Desktop/receipt_0.json"); err != nil {
					t.Fatal(err)
				}
			}
			if i == 1 {
				if err := writeJSONToFile(receipt, "/Users/lucas/Desktop/receipt_1.json"); err != nil {
					t.Fatal(err)
				}
			}
		}

		receiptTree := NewTrie(receipts)
		t.Logf("  block receipt root %x\n", block.Header().ReceiptHash)
		t.Logf("  receipt tree root  %x\n", receiptTree.Hash())
		if receiptTree.Hash() != block.Header().ReceiptHash {
			t.Fatal("receipt root mismatch")
		} else {
			t.Logf("  receipt root hash & block receipt root match\n")
		}

		i := 1
		proof := NewProof()
		indexBuf = rlp.AppendUint64(indexBuf[:0], uint64(i))
		err = receiptTree.Prove(indexBuf, 0, proof)
		if err != nil {
			t.Fatal(err)
		}

		// NOTE: eth receipts only hashes the following fields
		// 	data := &receiptRLP{r.statusEncoding(), r.CumulativeGasUsed, r.Bloom, r.Logs}
		value, err := trie.VerifyProof(block.Header().ReceiptHash, indexBuf, proof)
		t.Logf("pass? %v\n", err == nil)
		t.Logf("value %x\n", value)
		value, err = proof.Verify(block.Header().ReceiptHash, i)
		if err != nil {
			t.Fatal(err)
		}

		var receipt types.Receipt
		receipt.UnmarshalBinary(value)

		t.Logf("  receipt %+v\n", receipt)
		t.Logf("  receipt tx hash %+v\n", receipt.TxHash.Hex())

		for _, log := range receipt.Logs {
			t.Logf("  log %+v\n", log)
		}
	}
}

// Function to write a struct to a JSON file
func writeJSONToFile(data interface{}, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return err
	}
	return nil
}

// Function to read a JSON file into a struct
func readJSONFromFile(data interface{}, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	return nil
}
