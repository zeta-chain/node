package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	authoritytypes "github.com/zeta-chain/zetacore/x/authority/types"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func getEthereumChainID() int64 {
	return 5 // Goerli

}

// setEnabledChain sets the chain as enabled in chain params
func setEnabledChain(ctx sdk.Context, zk keepertest.ZetaKeepers, chainID int64) {
	zk.ObserverKeeper.SetChainParamsList(ctx, observertypes.ChainParamsList{ChainParams: []*observertypes.ChainParams{
		{
			ChainId:                  chainID,
			ConnectorContractAddress: sample.EthAddress().Hex(),
			BallotThreshold:          sdk.OneDec(),
			MinObserverDelegation:    sdk.OneDec(),
			IsSupported:              true,
		},
	}})
}

// setupTssAndNonceToCctx sets tss and nonce to cctx
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

func TestMsgServer_AddToOutTxTracker(t *testing.T) {
	t.Run("add tracker admin", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		chainID := getEthereumChainID()
		setupTssAndNonceToCctx(k, ctx, chainID, 0)
		setEnabledChain(ctx, zk, chainID)

		msgServer := keeper.NewMsgServerImpl(*k)

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    sample.Hash().Hex(),
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     0,
		})
		require.NoError(t, err)
		_, found := k.GetOutTxTracker(ctx, chainID, 0)
		require.True(t, found)
	})

	t.Run("unable to add tracker admin exceeding maximum allowed length of hashlist without proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseAuthorityMock: true,
		})

		admin := sample.AccAddress()
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		keepertest.MockIsAuthorized(&authorityMock.Mock, admin, authoritytypes.PolicyType_groupEmergency, true)

		chainID := getEthereumChainID()
		setupTssAndNonceToCctx(k, ctx, chainID, 0)
		setEnabledChain(ctx, zk, chainID)

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

		_, err := msgServer.AddToOutTxTracker(ctx, &types.MsgAddToOutTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    sample.Hash().Hex(),
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
			Nonce:     0,
		})
		require.NoError(t, err)
		tracker, found := k.GetOutTxTracker(ctx, chainID, 0)
		require.True(t, found)
		require.Equal(t, 2, len(tracker.HashList))
	})

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
}
