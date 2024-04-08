package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
)

// TestKeeper_GetChainState tests get, and set chain state
func TestKeeper_GetChainState(t *testing.T) {
	k, ctx, _, _ := keepertest.LightclientKeeper(t)
	_, found := k.GetChainState(ctx, 42)
	require.False(t, found)

	k.SetChainState(ctx, sample.ChainState(42))
	_, found = k.GetChainState(ctx, 42)
	require.True(t, found)
}

func TestKeeper_GetAllChainStates(t *testing.T) {
	k, ctx, _, _ := keepertest.LightclientKeeper(t)
	c1 := sample.ChainState(42)
	c2 := sample.ChainState(43)
	c3 := sample.ChainState(44)

	k.SetChainState(ctx, c1)
	k.SetChainState(ctx, c2)
	k.SetChainState(ctx, c3)

	list := k.GetAllChainStates(ctx)
	require.Len(t, list, 3)
	require.Contains(t, list, c1)
	require.Contains(t, list, c2)
	require.Contains(t, list, c3)
}

//t.Run("unable to add tracker admin exceeding maximum allowed length of hashlist without proof", func(t *testing.T) {
//	k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
//		UseAuthorityMock: true,
//	})
//
//	admin := sample.AccAddress()
//	authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
//	keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)
//
//	chainID := getEthereumChainID()
//	setupTssAndNonceToCctx(k, ctx, chainID, 0, types.CctxStatus_PendingOutbound)
//	setEnabledChain(ctx, zk, chainID)
//
//	k.SetOutTxTracker(ctx, types.OutTxTracker{
//		ChainId: chainID,
//		Nonce:   0,
//		HashList: []*types.TxHashList{
//			{
//				TxHash:   "hash1",
//				TxSigner: sample.AccAddress(),
//				Proved:   false,
//			},
//			{
//				TxHash:   "hash2",
//				TxSigner: sample.AccAddress(),
//				Proved:   false,
//			},
//		},
//	})
//
//	msgServer := keeper.NewMsgServerImpl(*k)
//
//	_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
//		Creator:   admin,
//		ChainId:   chainID,
//		TxHash:    sample.Hash().Hex(),
//		Proof:     nil,
//		BlockHash: "",
//		TxIndex:   0,
//		Nonce:     0,
//	})
//	require.NoError(t, err)
//	tracker, found := k.GetOutTxTracker(ctx, chainID, 0)
//	require.True(t, found)
//	require.Equal(t, 2, len(tracker.HashList))
//})

// Commented out as these tests don't work without using RPC
// TODO: Reenable these tests
// https://github.com/zeta-chain/node/issues/1875
//t.Run("fail add proof based tracker with wrong chainID", func(t *testing.T) {
//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
//
//	chainID := getEthereumChainID()
//
//	txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
//	require.NoError(t, err)
//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
//	setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
//
//	msgServer := keeper.NewMsgServerImpl(*k)
//
//	_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
//		Creator:   sample.AccAddress(),
//		ChainId:   97,
//		TxHash:    tx.Hash().Hex(),
//		Proof:     proof,
//		BlockHash: block.Hash().Hex(),
//		TxIndex:   txIndex,
//		Nonce:     tx.Nonce(),
//	})
//	require.ErrorIs(t, err, observertypes.ErrSupportedChains)
//	_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
//	require.False(t, found)
//})
//
//t.Run("fail add proof based tracker with wrong nonce", func(t *testing.T) {
//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
//
//	chainID := getEthereumChainID()
//
//	txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
//	require.NoError(t, err)
//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
//	setupTssAndNonceToCctx(k, ctx, chainID, 1)
//
//	msgServer := keeper.NewMsgServerImpl(*k)
//
//	_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
//		Creator:   sample.AccAddress(),
//		ChainId:   chainID,
//		TxHash:    tx.Hash().Hex(),
//		Proof:     proof,
//		BlockHash: block.Hash().Hex(),
//		TxIndex:   txIndex,
//		Nonce:     1,
//	})
//	require.ErrorIs(t, err, types.ErrTxBodyVerificationFail)
//	_, found := k.GetOutTxTracker(ctx, chainID, 1)
//	require.False(t, found)
//})
//
//t.Run("fail add proof based tracker with wrong tx_hash", func(t *testing.T) {
//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
//
//	chainID := getEthereumChainID()
//
//	txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
//	require.NoError(t, err)
//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
//	setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
//
//	msgServer := keeper.NewMsgServerImpl(*k)
//	_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
//		Creator:   sample.AccAddress(),
//		ChainId:   chainID,
//		TxHash:    "wrong_hash",
//		Proof:     proof,
//		BlockHash: block.Hash().Hex(),
//		TxIndex:   txIndex,
//		Nonce:     tx.Nonce(),
//	})
//	require.ErrorIs(t, err, types.ErrTxBodyVerificationFail)
//	_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
//	require.False(t, found)
//})
//
//t.Run("fail proof based tracker with incorrect proof", func(t *testing.T) {
//
//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
//	chainID := getEthereumChainID()
//
//	txIndex, block, header, headerRLP, _, tx, err := sample.Proof()
//	require.NoError(t, err)
//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
//	setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
//
//	msgServer := keeper.NewMsgServerImpl(*k)
//
//	_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
//		Creator:   sample.AccAddress(),
//		ChainId:   chainID,
//		TxHash:    tx.Hash().Hex(),
//		Proof:     common.NewEthereumProof(ethereum.NewProof()),
//		BlockHash: block.Hash().Hex(),
//		TxIndex:   txIndex,
//		Nonce:     tx.Nonce(),
//	})
//	require.ErrorIs(t, err, types.ErrProofVerificationFail)
//	_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
//	require.False(t, found)
//})
//
//t.Run("add proof based tracker with correct proof", func(t *testing.T) {
//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
//
//	chainID := getEthereumChainID()
//	txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
//	require.NoError(t, err)
//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
//	setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
//
//	msgServer := keeper.NewMsgServerImpl(*k)
//
//	_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
//		Creator:   sample.AccAddress(),
//		ChainId:   chainID,
//		TxHash:    tx.Hash().Hex(),
//		Proof:     proof,
//		BlockHash: block.Hash().Hex(),
//		TxIndex:   txIndex,
//		Nonce:     tx.Nonce(),
//	})
//	require.NoError(t, err)
//	_, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
//	require.True(t, found)
//})
//
//t.Run("add proven txHash even if length of hashList is already 2", func(t *testing.T) {
//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
//
//	chainID := getEthereumChainID()
//
//	txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
//	require.NoError(t, err)
//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
//	setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
//	k.SetOutTxTracker(ctx, types.OutTxTracker{
//		ChainId: chainID,
//		Nonce:   tx.Nonce(),
//		HashList: []*types.TxHashList{
//			{
//				TxHash:   "hash1",
//				TxSigner: sample.AccAddress(),
//				Proved:   false,
//			},
//			{
//				TxHash:   "hash2",
//				TxSigner: sample.AccAddress(),
//				Proved:   false,
//			},
//		},
//	})
//
//	msgServer := keeper.NewMsgServerImpl(*k)
//
//	_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
//		Creator:   sample.AccAddress(),
//		ChainId:   chainID,
//		TxHash:    tx.Hash().Hex(),
//		Proof:     proof,
//		BlockHash: block.Hash().Hex(),
//		TxIndex:   txIndex,
//		Nonce:     tx.Nonce(),
//	})
//	require.NoError(t, err)
//	tracker, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
//	require.True(t, found)
//	require.Equal(t, 3, len(tracker.HashList))
//	// Proven tracker is prepended to the list
//	require.True(t, tracker.HashList[0].Proved)
//	require.False(t, tracker.HashList[1].Proved)
//	require.False(t, tracker.HashList[2].Proved)
//})
//
//t.Run("add proof for existing txHash", func(t *testing.T) {
//	k, ctx, _, zk := keepertest.CrosschainKeeper(t)
//
//	chainID := getEthereumChainID()
//
//	txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
//	require.NoError(t, err)
//	setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
//	setupTssAndNonceToCctx(k, ctx, chainID, int64(tx.Nonce()))
//	k.SetOutTxTracker(ctx, types.OutTxTracker{
//		ChainId: chainID,
//		Nonce:   tx.Nonce(),
//		HashList: []*types.TxHashList{
//			{
//				TxHash:   tx.Hash().Hex(),
//				TxSigner: sample.AccAddress(),
//				Proved:   false,
//			},
//		},
//	})
//	tracker, found := k.GetOutTxTracker(ctx, chainID, tx.Nonce())
//	require.True(t, found)
//	require.False(t, tracker.HashList[0].Proved)
//
//	msgServer := keeper.NewMsgServerImpl(*k)
//
//	_, err = msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
//		Creator:   sample.AccAddress(),
//		ChainId:   chainID,
//		TxHash:    tx.Hash().Hex(),
//		Proof:     proof,
//		BlockHash: block.Hash().Hex(),
//		TxIndex:   txIndex,
//		Nonce:     tx.Nonce(),
//	})
//	require.NoError(t, err)
//	tracker, found = k.GetOutTxTracker(ctx, chainID, tx.Nonce())
//	require.True(t, found)
//	require.Equal(t, 1, len(tracker.HashList))
//	require.True(t, tracker.HashList[0].Proved)
//})
