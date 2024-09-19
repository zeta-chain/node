package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/proofs"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/lightclient/types"
)

// TODO: Add test for Bitcoin proof verification
// https://github.com/zeta-chain/node/issues/1994

func TestKeeper_Prove(t *testing.T) {
	t.Run("should error if req is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		_, err := k.Prove(wctx, nil)
		require.Error(t, err)
	})

	t.Run("should error if block hash is invalid", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		proof, _, _, txIndex, _, hash := sample.Proof(t)

		_, err := k.Prove(wctx, &types.QueryProveRequest{
			ChainId:   1000,
			TxHash:    hash.Hex(),
			Proof:     proof,
			BlockHash: "invalid",
			TxIndex:   txIndex,
		})
		require.ErrorContains(t, err, "cannot convert hash to bytes for chain")
	})

	t.Run("should error if block header not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		proof, _, blockHash, txIndex, chainID, hash := sample.Proof(t)

		_, err := k.Prove(wctx, &types.QueryProveRequest{
			ChainId:   chainID,
			TxHash:    hash.Hex(),
			Proof:     proof,
			BlockHash: blockHash,
			TxIndex:   txIndex,
		})
		require.ErrorContains(t, err, "block header not found")
	})

	t.Run("should returns response with proven false if invalid proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		proof, blockHeader, blockHash, txIndex, chainID, hash := sample.Proof(t)

		k.SetBlockHeader(ctx, blockHeader)

		res, err := k.Prove(wctx, &types.QueryProveRequest{
			ChainId:   chainID,
			TxHash:    hash.Hex(),
			Proof:     proof,
			BlockHash: blockHash,
			TxIndex:   txIndex + 1, // change txIndex to make it invalid
		})
		require.NoError(t, err)
		require.False(t, res.Valid)
	})

	t.Run("should returns response with proven true if valid proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		proof, blockHeader, blockHash, txIndex, chainID, hash := sample.Proof(t)

		k.SetBlockHeader(ctx, blockHeader)

		res, err := k.Prove(wctx, &types.QueryProveRequest{
			ChainId:   chainID,
			TxHash:    hash.Hex(),
			Proof:     proof,
			BlockHash: blockHash,
			TxIndex:   txIndex,
		})
		require.NoError(t, err)
		require.True(t, res.Valid)
	})

	t.Run("should error if error during proof verification", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		proof, blockHeader, blockHash, txIndex, chainID, hash := sample.Proof(t)

		// corrupt the block header
		blockHeader.Header = proofs.HeaderData{}

		k.SetBlockHeader(ctx, blockHeader)

		_, err := k.Prove(wctx, &types.QueryProveRequest{
			ChainId:   chainID,
			TxHash:    hash.Hex(),
			Proof:     proof,
			BlockHash: blockHash,
			TxIndex:   txIndex,
		})
		require.Error(t, err)
	})

	t.Run("should error if tx hash mismatch", func(t *testing.T) {
		k, ctx, _, _ := keepertest.LightclientKeeper(t)
		wctx := sdk.WrapSDKContext(ctx)

		proof, blockHeader, blockHash, txIndex, chainID, _ := sample.Proof(t)

		k.SetBlockHeader(ctx, blockHeader)

		_, err := k.Prove(wctx, &types.QueryProveRequest{
			ChainId:   chainID,
			TxHash:    sample.Hash().Hex(), // change tx hash to make it invalid
			Proof:     proof,
			BlockHash: blockHash,
			TxIndex:   txIndex,
		})
		require.ErrorContains(t, err, "tx hash mismatch")
	})
}
