package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgServer_AddToOutTxTracker(t *testing.T) {
	t.Run("add tracker admin", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		chainID := int64(5)
		txIndex, block, header, headerRLP, _, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		setupTssAndNonceToCctx(k, ctx, chainID, 0)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    tx.Hash().Hex(),
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     0,
		})
		require.NoError(t, err)
		_, found := k.GetOutTxTracker(ctx, chainID, 0)
		require.True(t, found)
	})

	t.Run("fail add proof based tracker with wrong chainID", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   97,
			TxHash:    tx.Hash().Hex(),
			Proof:     proof,
			BlockHash: block.Hash().Hex(),
			TxIndex:   txIndex,
			Nonce:     tx.Nonce(),
		})
		require.ErrorIs(t, types.ErrTxBodyVerificationFail, err)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		require.False(t, found)
	})

	t.Run("fail add proof based tracker with wrong nonce", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		setupTssAndNonceToCctx(k, ctx, chainID, 1)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx.Hash().Hex(),
			Proof:     proof,
			BlockHash: block.Hash().Hex(),
			TxIndex:   txIndex,
			Nonce:     1,
		})
		require.ErrorIs(t, types.ErrTxBodyVerificationFail, err)
		_, found := k.GetOutTxTracker(ctx, chainID, 1)
		require.False(t, found)
	})

	t.Run("fail add proof based tracker with wrong tx_hash", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    "wrong_hash",
			Proof:     proof,
			BlockHash: block.Hash().Hex(),
			TxIndex:   txIndex,
			Nonce:     tx.Nonce(),
		})
		require.ErrorIs(t, types.ErrTxBodyVerificationFail, err)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		require.False(t, found)
	})

	t.Run("fail proof based tracker with incorrect proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, _, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx.Hash().Hex(),
			Proof:     common.NewEthereumProof(ethereum.NewProof()),
			BlockHash: block.Hash().Hex(),
			TxIndex:   txIndex,
			Nonce:     tx.Nonce(),
		})
		require.ErrorIs(t, types.ErrProofVerificationFail, err)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		require.False(t, found)
	})
	t.Run("add proof based tracker with correct proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx.Hash().Hex(),
			Proof:     proof,
			BlockHash: block.Hash().Hex(),
			TxIndex:   txIndex,
			Nonce:     tx.Nonce(),
		})
		require.NoError(t, err)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		require.True(t, found)
	})
}

func setupTssAndNonceToCctx(k *keeper.Keeper, ctx sdk.Context, chainId, nonce int64) {
	k.SetTSS(ctx, types.TSS{
		TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
	})
	cctx := types.CrossChainTx{
		Creator: "any",
		Index:   "0x123",
		CctxStatus: &types.Status{
			Status: types.CctxStatus_PendingOutbound,
		},
	}
	k.SetCrossChainTx(ctx, cctx)
	k.SetNonceToCctx(ctx, types.NonceToCctx{
		ChainId:   chainId,
		Nonce:     nonce,
		CctxIndex: "0x123",
		Tss:       "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
	})
}
