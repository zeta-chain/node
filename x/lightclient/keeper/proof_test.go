package keeper_test

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/proofs"
	"github.com/zeta-chain/zetacore/pkg/proofs/ethereum"
	"github.com/zeta-chain/zetacore/pkg/testdata"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/lightclient/types"
)

// generateProof generates a proof and block header
// returns the proof, block header, block hash, tx index, chain id, and tx hash
func generateProof(t *testing.T) (*proofs.Proof, proofs.BlockHeader, string, int64, int64, ethcommon.Hash) {
	header, err := testdata.ReadEthHeader()
	require.NoError(t, err)
	b, err := rlp.EncodeToBytes(&header)
	require.NoError(t, err)

	var txs ethtypes.Transactions
	for i := 0; i < testdata.TxsCount; i++ {
		tx, err := testdata.ReadEthTx(i)
		require.NoError(t, err)
		txs = append(txs, &tx)
	}
	txsTree := ethereum.NewTrie(txs)

	// choose 2 as the index of the tx to prove
	txIndex := 2
	proof, err := txsTree.GenerateProof(txIndex)
	require.NoError(t, err)

	chainID := chains.SepoliaChain().ChainId
	ethProof := proofs.NewEthereumProof(proof)
	ethHeader := proofs.NewEthereumHeader(b)
	blockHeader := proofs.BlockHeader{
		Height:     header.Number.Int64(),
		Hash:       header.Hash().Bytes(),
		ParentHash: header.ParentHash.Bytes(),
		ChainId:    chainID,
		Header:     ethHeader,
	}
	txHash := txs[txIndex].Hash()
	return ethProof, blockHeader, header.Hash().Hex(), int64(txIndex), chainID, txHash
}

func TestKeeper_VerifyProof(t *testing.T) {
	t.Run("should error if verification flags not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.SepoliaChain().ChainId, sample.Hash().String(), 1)
		require.ErrorIs(t, err, types.ErrVerificationFlagsNotFound)
	})

	t.Run("should error if verification not enabled for btc chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: true,
			BtcTypeChainEnabled: false,
		})

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.BtcMainnetChain().ChainId, sample.Hash().String(), 1)
		require.ErrorIs(t, err, types.ErrBlockHeaderVerificationDisabled)
	})

	t.Run("should error if verification not enabled for evm chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: false,
			BtcTypeChainEnabled: true,
		})

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.SepoliaChain().ChainId, sample.Hash().String(), 1)
		require.ErrorIs(t, err, types.ErrBlockHeaderVerificationDisabled)
	})

	t.Run("should error if block header-based verification not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: true,
			BtcTypeChainEnabled: true,
		})

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, 101, sample.Hash().String(), 1)
		require.ErrorIs(t, err, types.ErrChainNotSupported)
	})

	t.Run("should error if blockhash invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: true,
			BtcTypeChainEnabled: true,
		})

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.BtcMainnetChain().ChainId, "invalid", 1)
		require.ErrorIs(t, err, types.ErrInvalidBlockHash)
	})

	t.Run("should error if block header not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: true,
			BtcTypeChainEnabled: true,
		})

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.SepoliaChain().ChainId, sample.Hash().String(), 1)
		require.ErrorIs(t, err, types.ErrBlockHeaderNotFound)
	})

	t.Run("should fail if proof can't be verified", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		proof, blockHeader, blockHash, txIndex, chainID, _ := generateProof(t)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: true,
			BtcTypeChainEnabled: true,
		})

		k.SetBlockHeader(ctx, blockHeader)

		// providing wrong tx index
		_, err := k.VerifyProof(ctx, proof, chainID, blockHash, txIndex+1)
		require.ErrorIs(t, err, types.ErrProofVerificationFailed)
	})

	t.Run("can verify a proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		proof, blockHeader, blockHash, txIndex, chainID, _ := generateProof(t)

		k.SetVerificationFlags(ctx, types.VerificationFlags{
			EthTypeChainEnabled: true,
			BtcTypeChainEnabled: true,
		})

		k.SetBlockHeader(ctx, blockHeader)

		txBytes, err := k.VerifyProof(ctx, proof, chainID, blockHash, txIndex)
		require.NoError(t, err)
		require.NotNil(t, txBytes)
	})
}
