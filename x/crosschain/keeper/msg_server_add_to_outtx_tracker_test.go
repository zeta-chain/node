package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func getEthereumChainID() int64 {
	return 5 // Goerli

}
func TestMsgServer_AddToOutTxTracker(t *testing.T) {
	t.Run("add tracker admin", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		chainID := getEthereumChainID()
		txIndex, block, header, headerRLP, _, tx, err := sample.Proof()
		assert.NoError(t, err)
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
		assert.NoError(t, err)
		_, found := k.GetOutTxTracker(ctx, chainID, 0)
		assert.True(t, found)
	})
	t.Run("unable to add tracker admin exceeding maximum allowed length of hashlist without proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		chainID := getEthereumChainID()
		txIndex, block, header, headerRLP, _, tx, err := sample.Proof()
		assert.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		setupTssAndNonceToCctx(k, ctx, chainID, 0)
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: chainID,
			Nonce:   0,
			HashList: []*types.TxHashList{
				{
					TxHash:   "hash1",
					TxSigner: sample.AccAddress(),
					Proved:   false,
				},
				{
					TxHash:   "hash2",
					TxSigner: sample.AccAddress(),
					Proved:   false,
				},
			},
		})
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
		assert.NoError(t, err)
		tracker, found := k.GetOutTxTracker(ctx, chainID, 0)
		assert.True(t, found)
		assert.Equal(t, 2, len(tracker.HashList))
	})

	t.Run("fail add proof based tracker with wrong chainID", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getEthereumChainID()
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		assert.NoError(t, err)
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
		assert.ErrorIs(t, err, observertypes.ErrSupportedChains)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		assert.False(t, found)
	})

	t.Run("fail add proof based tracker with wrong nonce", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getEthereumChainID()
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		assert.NoError(t, err)
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
		assert.ErrorIs(t, err, types.ErrTxBodyVerificationFail)
		_, found := k.GetOutTxTracker(ctx, chainID, 1)
		assert.False(t, found)
	})

	t.Run("fail add proof based tracker with wrong tx_hash", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getEthereumChainID()
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		assert.NoError(t, err)
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
		assert.ErrorIs(t, err, types.ErrTxBodyVerificationFail)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		assert.False(t, found)
	})

	t.Run("fail proof based tracker with incorrect proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getEthereumChainID()
		txIndex, block, header, headerRLP, _, tx, err := sample.Proof()
		assert.NoError(t, err)
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
		assert.ErrorIs(t, err, types.ErrProofVerificationFail)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		assert.False(t, found)
	})
	t.Run("add proof based tracker with correct proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getEthereumChainID()
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		assert.NoError(t, err)
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
		assert.NoError(t, err)
		_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		assert.True(t, found)
	})
	t.Run("add proven txHash even if length of hashList is already 2", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getEthereumChainID()
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		assert.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: chainID,
			Nonce:   tx.Nonce(),
			HashList: []*types.TxHashList{
				{
					TxHash:   "hash1",
					TxSigner: sample.AccAddress(),
					Proved:   false,
				},
				{
					TxHash:   "hash2",
					TxSigner: sample.AccAddress(),
					Proved:   false,
				},
			},
		})
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
		assert.NoError(t, err)
		tracker, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		assert.True(t, found)
		assert.Equal(t, 3, len(tracker.HashList))
		// Proven tracker is prepended to the list
		assert.True(t, tracker.HashList[0].Proved)
		assert.False(t, tracker.HashList[1].Proved)
		assert.False(t, tracker.HashList[2].Proved)
	})
	t.Run("add proof for existing txHash", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := getEthereumChainID()
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		assert.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
		k.SetOutTxTracker(ctx, types.OutTxTracker{
			ChainId: chainID,
			Nonce:   tx.Nonce(),
			HashList: []*types.TxHashList{
				{
					TxHash:   tx.Hash().Hex(),
					TxSigner: sample.AccAddress(),
					Proved:   false,
				},
			},
		})
		tracker, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		assert.True(t, found)
		assert.False(t, tracker.HashList[0].Proved)
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
		assert.NoError(t, err)
		tracker, found = k.GetOutTxTracker(ctx, chainID, tx.Nonce())
		assert.True(t, found)
		assert.Equal(t, 1, len(tracker.HashList))
		assert.True(t, tracker.HashList[0].Proved)
	})
}

func setupTssAndNonceToCctx(k *keeper.Keeper, ctx sdk.Context, chainId, nonce int64) {

	tssPubKey := "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p"
	k.GetObserverKeeper().SetTSS(ctx, observertypes.TSS{
		TssPubkey: tssPubKey,
	})
	k.GetObserverKeeper().SetPendingNonces(ctx, observertypes.PendingNonces{
		Tss:       tssPubKey,
		NonceLow:  0,
		NonceHigh: 1,
		ChainId:   chainId,
	})
	cctx := types.CrossChainTx{
		Creator: "any",
		Index:   "0x123",
		CctxStatus: &types.Status{
			Status: types.CctxStatus_PendingOutbound,
		},
	}
	k.SetCrossChainTx(ctx, cctx)
	k.GetObserverKeeper().SetNonceToCctx(ctx, observertypes.NonceToCctx{
		ChainId:   chainId,
		Nonce:     nonce,
		CctxIndex: "0x123",
		Tss:       "zetapub1addwnpepq28c57cvcs0a2htsem5zxr6qnlvq9mzhmm76z3jncsnzz32rclangr2g35p",
	})
}
