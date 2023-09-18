package keeper_test

import (
	"testing"

	"cosmossdk.io/math"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

func TestKeeper_RefundAmountOnZetaChain(t *testing.T) {
	t.Run("should refund amount on zeta chain", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)
		asset := sample.EthAddress().String()
		sender := sample.EthAddress()
		chainID := getValidEthChainID(t)

		// deploy zrc20
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		zrc20Addr := deployZRC20(
			t,
			ctx,
			zk.FungibleKeeper,
			sdkk.EvmKeeper,
			chainID,
			"bar",
			asset,
			"bar",
		)

		err := k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
			}},
			math.NewUint(42),
		)
		require.NoError(t, err)

		// check amount deposited in balance
		balance, err := zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20Addr, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(42), balance.Uint64())

		// can refund again
		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: chainID,
				Sender:        sender.String(),
				Asset:         asset,
			}},
			math.NewUint(42),
		)
		require.NoError(t, err)
		balance, err = zk.FungibleKeeper.BalanceOfZRC4(ctx, zrc20Addr, sender)
		require.NoError(t, err)
		require.Equal(t, uint64(84), balance.Uint64())
	})

	t.Run("should fail with invalid cctx", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		err := k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: common.CoinType_Zeta,
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "unsupported coin type")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType: common.CoinType_Gas,
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "unsupported coin type")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: 999999,
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "only EVM chains are supported")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        "invalid",
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "invalid sender address")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        sample.EthAddress().String(),
			},
		},
			math.Uint{},
		)
		require.ErrorContains(t, err, "no amount to refund")

		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        sample.EthAddress().String(),
			}},
			math.ZeroUint(),
		)
		require.ErrorContains(t, err, "no amount to refund")

		// the foreign coin has not been set
		err = k.RefundAmountOnZetaChain(ctx, types.CrossChainTx{
			InboundTxParams: &types.InboundTxParams{
				CoinType:      common.CoinType_ERC20,
				SenderChainId: getValidEthChainID(t),
				Sender:        sample.EthAddress().String(),
				Asset:         sample.EthAddress().String(),
			}},
			math.NewUint(42),
		)
		require.ErrorContains(t, err, "zrc not found")
	})
}
