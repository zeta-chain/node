//go:build TESTNET
// +build TESTNET

package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_AddToInTxTracker(t *testing.T) {
	t.Run("add proof based tracker with correct proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx.Hash().Hex(),
			CoinType:  common.CoinType_Zeta,
			Proof:     proof,
			BlockHash: block.Hash().Hex(),
			TxIndex:   txIndex,
		})
		require.NoError(t, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx.Hash().Hex())
		require.True(t, found)
	})

	t.Run("fail to add proof based tracker with wrong tx hash", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    "fake_hash",
			CoinType:  common.CoinType_Zeta,
			Proof:     proof,
			BlockHash: block.Hash().Hex(),
			TxIndex:   txIndex,
		})
		require.ErrorIs(t, types.ErrTxBodyVerificationFail, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx.Hash().Hex())
		require.False(t, found)
	})

	t.Run("fail to add proof based tracker with wrong chain id", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, proof, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   97,
			TxHash:    tx.Hash().Hex(),
			CoinType:  common.CoinType_Zeta,
			Proof:     proof,
			BlockHash: block.Hash().Hex(),
			TxIndex:   txIndex,
		})
		require.ErrorIs(t, types.ErrTxBodyVerificationFail, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx.Hash().Hex())
		require.False(t, found)
	})

	t.Run("fail to add proof based tracker with wrong proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		chainID := int64(5)
		txIndex, block, header, headerRLP, _, tx, err := sample.Proof()
		require.NoError(t, err)
		setupVerificationParams(zk, ctx, txIndex, chainID, header, headerRLP, block)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err = msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx.Hash().Hex(),
			CoinType:  common.CoinType_Zeta,
			Proof:     common.NewEthereumProof(ethereum.NewProof()),
			BlockHash: block.Hash().Hex(),
			TxIndex:   txIndex,
		})
		require.ErrorIs(t, types.ErrProofVerificationFail, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx.Hash().Hex())
		require.False(t, found)
	})
	t.Run("normal user submit without proof", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)
		tx_hash := "string"
		chainID := int64(5)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx_hash,
			CoinType:  common.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
		require.ErrorIs(t, types.ErrTxBodyVerificationFail, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx_hash)
		require.False(t, found)
	})
	t.Run("admin add  tx tracker", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		tx_hash := "string"
		chainID := int64(5)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    tx_hash,
			CoinType:  common.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
		require.NoError(t, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx_hash)
		require.True(t, found)
	})
	t.Run("admin submit fake tracker", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		tx_hash := "string"
		chainID := int64(5)
		msgServer := keeper.NewMsgServerImpl(*k)
		_, err := msgServer.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   admin,
			ChainId:   chainID,
			TxHash:    "Malicious TX HASH",
			CoinType:  common.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
		require.NoError(t, err)
		_, found := k.GetInTxTracker(ctx, chainID, "Malicious TX HASH")
		require.True(t, found)
		_, found = k.GetInTxTracker(ctx, chainID, tx_hash)
		require.False(t, found)
	})
}

func setupVerificationParams(zk keepertest.ZetaKeepers, ctx sdk.Context, tx_index int64, chainID int64, header ethtypes.Header, headerRLP []byte, block *ethtypes.Block) {
	params := zk.ObserverKeeper.GetParams(ctx)
	params.ObserverParams = append(params.ObserverParams, &observerTypes.ObserverParams{
		Chain: &common.Chain{
			ChainId:   chainID,
			ChainName: common.ChainName_goerli_testnet,
		},
		BallotThreshold:       sdk.OneDec(),
		MinObserverDelegation: sdk.OneDec(),
		IsSupported:           true,
	})
	zk.ObserverKeeper.SetParams(ctx, params)
	zk.ObserverKeeper.SetBlockHeader(ctx, common.BlockHeader{
		Height:     block.Number().Int64(),
		Hash:       block.Hash().Bytes(),
		ParentHash: header.ParentHash.Bytes(),
		ChainId:    chainID,
		Header:     common.NewEthereumHeader(headerRLP),
	})
	zk.ObserverKeeper.SetCoreParams(ctx, observerTypes.CoreParamsList{CoreParams: []*observerTypes.CoreParams{
		{
			ChainId:                  chainID,
			ConnectorContractAddress: block.Transactions()[tx_index].To().Hex(),
		},
	}})
	zk.ObserverKeeper.SetCrosschainFlags(ctx, observerTypes.CrosschainFlags{
		BlockHeaderVerificationFlags: &observerTypes.BlockHeaderVerificationFlags{
			IsEthTypeChainEnabled: true,
			IsBtcTypeChainEnabled: false,
		},
	})
}
