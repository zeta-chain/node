package ethereum

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/testutil/testdata"
)

func TestProofGeneration(t *testing.T) {
	header, err := testdata.ReadEthHeader()
	require.NoError(t, err)

	var receipts types.Receipts
	for i := 0; i < testdata.TxsCount; i++ {
		receipt, err := testdata.ReadEthReceipt(i)
		require.NoError(t, err)
		receipts = append(receipts, &receipt)
	}

	// generate a trie from the receipts and compare the root hash with the one in the header
	receiptTree := NewTrie(receipts)
	require.EqualValues(t, header.ReceiptHash.Hex(), receiptTree.Trie.Hash().Hex())

	t.Run("generate proof for receipts", func(t *testing.T) {
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
	})

	t.Run("should error verify for negative key", func(t *testing.T) {
		proof, err := receiptTree.GenerateProof(0)
		require.NoError(t, err)

		verified, err := proof.Verify(header.ReceiptHash, -1)
		require.Error(t, err)
		require.Nil(t, verified)
	})

	t.Run("should not generate proof for negative tx index", func(t *testing.T) {
		proof, err := receiptTree.GenerateProof(-1)
		require.Error(t, err)
		require.Nil(t, proof)
	})

	t.Run("has key", func(t *testing.T) {
		proof, err := receiptTree.GenerateProof(0)
		require.NoError(t, err)
		require.Greater(t, len(proof.Keys), 0)

		proof2, err := receiptTree.GenerateProof(1)
		require.NoError(t, err)
		require.Equal(t, len(proof2.Keys), 3)

		key := proof.Keys[0]
		has, err := proof.Has(key)
		require.NoError(t, err)
		require.True(t, has)

		err = proof.Put(key, proof.Values[0])
		require.NoError(t, err)

		key2 := proof2.Keys[2]
		has, err = proof.Has(key2)
		require.NoError(t, err)
		require.False(t, has)
	})

	t.Run("get key", func(t *testing.T) {
		proof, err := receiptTree.GenerateProof(0)
		require.NoError(t, err)
		require.Greater(t, len(proof.Keys), 0)

		proof2, err := receiptTree.GenerateProof(1)
		require.NoError(t, err)
		require.Equal(t, len(proof2.Keys), 3)

		key := proof.Keys[0]
		_, err = proof.Get(key)
		require.NoError(t, err)

		key2 := proof2.Keys[2]
		_, err = proof.Get(key2)
		require.Error(t, err)
	})

	t.Run("delete key", func(t *testing.T) {
		proof, err := receiptTree.GenerateProof(0)
		require.NoError(t, err)
		require.Greater(t, len(proof.Keys), 0)

		proof2, err := receiptTree.GenerateProof(1)
		require.NoError(t, err)
		require.Equal(t, len(proof2.Keys), 3)

		key := proof.Keys[0]
		err = proof.Delete(key)
		require.NoError(t, err)

		err = proof.Delete(key)
		require.Error(t, err)

		key2 := proof2.Keys[2]
		err = proof.Delete(key2)
		require.Error(t, err)
	})
}
