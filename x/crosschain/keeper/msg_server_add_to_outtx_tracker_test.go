//go:build TESTNET
// +build TESTNET

package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestMsgServer_AddToOutTxTracker(t *testing.T) {
	t.Run("Add tracker admin", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		chainID := int64(5)
		txIndex, block, header, headerRLP, _, tx, err := sample.Proof()
		require.NoError(t, err)
		SetupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		SetupTss(k, ctx)
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

	t.Run("Fail add proof based tracker with wrong chainID", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		SetupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		SetupTss(k, ctx)
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
		require.Error(t, err)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		require.False(t, found)
	})

	t.Run("Fail add proof based tracker with wrong nonce", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		SetupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		SetupTss(k, ctx)
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
		require.Error(t, err)
		_, found := k.GetOutTxTracker(ctx, chainID, 1)
		require.False(t, found)
	})

	t.Run("Fail add proof based tracker with wrong tx_hash", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		SetupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		SetupTss(k, ctx)
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
		require.Error(t, err)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		require.False(t, found)
	})
	t.Run("Add proof based tracker with correct proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		SetupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		SetupTss(k, ctx)
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

func SetupTss(k *keeper.Keeper, ctx sdk.Context) {
	k.SetTSS(ctx, types.TSS{
		TssPubkey: "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
	})
}
