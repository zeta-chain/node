package keeper_test

import (
	"context"
	"math/big"
	"testing"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/common/ethereum"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestMsgServer_AddToInTxTracker(t *testing.T) {
	t.Run("Add proof based tracker with correct proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		tx_hash := "string"
		chainID := int64(1)
		tx_index := int64(1)

		RPC_URL := "https://rpc.ankr.com/eth_goerli"
		client, err := ethclient.Dial(RPC_URL)
		require.NoError(t, err)
		bn := int64(9509129)
		block, err := client.BlockByNumber(context.Background(), big.NewInt(bn))
		headerRLP, _ := rlp.EncodeToBytes(block.Header())
		var header ethtypes.Header
		err = rlp.DecodeBytes(headerRLP, &header)
		require.NoError(t, err)

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
		k.SetTSS(ctx, types.TSS{})

		tr := ethereum.NewTrie(block.Transactions())
		var b []byte
		ib := rlp.AppendUint64(b, uint64(tx_index))
		proof := ethereum.NewProof()
		err = tr.Prove(ib, 0, proof)
		require.NoError(t, err)

		_, err = k.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx_hash,
			CoinType:  common.CoinType_Zeta,
			Proof:     common.NewEthereumProof(proof),
			BlockHash: block.Hash().Hex(),
			TxIndex:   tx_index,
		})
		require.NoError(t, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx_hash)
		require.True(t, found)
	})

	t.Run("Fail to add proof based tracker with wrong proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		tx_hash := "string"
		chainID := int64(1)
		tx_index := int64(1)

		RPC_URL := "https://rpc.ankr.com/eth_goerli"
		client, err := ethclient.Dial(RPC_URL)
		require.NoError(t, err)
		bn := int64(9509129)
		block, err := client.BlockByNumber(context.Background(), big.NewInt(bn))
		headerRLP, _ := rlp.EncodeToBytes(block.Header())
		var header ethtypes.Header
		err = rlp.DecodeBytes(headerRLP, &header)
		require.NoError(t, err)

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
		k.SetTSS(ctx, types.TSS{})

		_, err = k.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx_hash,
			CoinType:  common.CoinType_Zeta,
			Proof:     common.NewEthereumProof(ethereum.NewProof()),
			BlockHash: block.Hash().Hex(),
			TxIndex:   tx_index,
		})
		require.Error(t, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx_hash)
		require.False(t, found)
	})
	t.Run("normal user submit without proof", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		tx_hash := "string"
		chainID := int64(1)
		_, err := k.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
			Creator:   sample.AccAddress(),
			ChainId:   chainID,
			TxHash:    tx_hash,
			CoinType:  common.CoinType_Zeta,
			Proof:     nil,
			BlockHash: "",
			TxIndex:   0,
		})
		require.Error(t, err)
		_, found := k.GetInTxTracker(ctx, chainID, tx_hash)
		require.False(t, found)
	})
	t.Run("admin add  tx tracker with admin", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		admin := sample.AccAddress()
		setAdminPolicies(ctx, zk, admin)
		tx_hash := "string"
		chainID := int64(1)
		_, err := k.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
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
		chainID := int64(1)
		_, err := k.AddToInTxTracker(ctx, &types.MsgAddToInTxTracker{
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
