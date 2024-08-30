package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/proofs"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
)

func TestKeeper_VerifyProof(t *testing.T) {
	t.Run("should error if verification flags not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.Sepolia.ChainId, sample.Hash().String(), 1)
		require.ErrorContains(t, err, "proof verification is disabled for all chains")
	})

	t.Run("should error if verification not enabled for btc chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: false,
				},
			},
		})

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.BitcoinMainnet.ChainId, sample.Hash().String(), 1)
		require.ErrorIs(t, err, types.ErrBlockHeaderVerificationDisabled)
	})

	t.Run("should error if verification not enabled for evm chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.Ethereum.ChainId,
					Enabled: false,
				},
			},
		})
		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.Sepolia.ChainId, sample.Hash().String(), 1)
		require.ErrorIs(t, err, types.ErrBlockHeaderVerificationDisabled)
	})

	t.Run("should error if block header-based verification not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: false,
				},
				{
					ChainId: chains.Ethereum.ChainId,
					Enabled: false,
				},
			},
		})

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.ZetaChainPrivnet.ChainId, sample.Hash().String(), 1)
		require.ErrorContains(
			t,
			err,
			fmt.Sprintf("proof verification is disabled for chain %d", chains.ZetaChainPrivnet.ChainId),
		)
	})

	t.Run("should error if blockhash invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: true,
				},
				{
					ChainId: chains.Ethereum.ChainId,
					Enabled: true,
				},
			},
		})

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.BitcoinMainnet.ChainId, "invalid", 1)
		require.ErrorIs(t, err, types.ErrInvalidBlockHash)
	})

	t.Run("should error if block header not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: true,
				},
				{
					ChainId: chains.Ethereum.ChainId,
					Enabled: true,
				},
			},
		})

		_, err := k.VerifyProof(ctx, &proofs.Proof{}, chains.Sepolia.ChainId, sample.Hash().String(), 1)
		require.ErrorContains(t, err, "block header verification is disabled")
	})

	t.Run("should fail if proof can't be verified", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		proof, blockHeader, blockHash, txIndex, chainID, _ := sample.Proof(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: true,
				},
				{
					ChainId: chains.Ethereum.ChainId,
					Enabled: true,
				},
			},
		})

		k.SetBlockHeader(ctx, blockHeader)

		// providing wrong tx index
		_, err := k.VerifyProof(ctx, proof, chainID, blockHash, txIndex+1)
		require.ErrorContains(t, err, "block header verification is disabled")
	})

	t.Run("can verify a proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)

		proof, blockHeader, blockHash, txIndex, chainID, _ := sample.Proof(t)

		k.SetBlockHeaderVerification(ctx, types.BlockHeaderVerification{
			HeaderSupportedChains: []types.HeaderSupportedChain{
				{
					ChainId: chains.BitcoinMainnet.ChainId,
					Enabled: true,
				},
				{
					ChainId: chains.Sepolia.ChainId,
					Enabled: true,
				},
			},
		})

		k.SetBlockHeader(ctx, blockHeader)

		txBytes, err := k.VerifyProof(ctx, proof, chainID, blockHash, txIndex)
		require.NoError(t, err)
		require.NotNil(t, txBytes)
	})
}
