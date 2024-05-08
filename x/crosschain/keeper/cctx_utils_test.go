package keeper_test

import (
	"fmt"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/coin"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/keeper"
	crosschainkeeper "github.com/zeta-chain/zetacore/x/crosschain/keeper"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestGetRevertGasLimit(t *testing.T) {
	t.Run("should return 0 if no inbound tx params", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{})
		require.NoError(t, err)
		require.Equal(t, uint64(0), gasLimit)
	})

	t.Run("should return 0 if coin type is not gas or erc20", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType: coin.CoinType_Zeta,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(0), gasLimit)
	})

	t.Run("should return the gas limit of the gas token", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		gas := setupGasCoin(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper, chainID, "foo", "FOO")

		_, err := zk.FungibleKeeper.UpdateZRC20GasLimit(ctx, gas, big.NewInt(42))
		require.NoError(t, err)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{

			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit)
	})

	t.Run("should return the gas limit of the associated asset", func(t *testing.T) {
		k, ctx, sdkk, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID()
		deploySystemContracts(t, ctx, zk.FungibleKeeper, sdkk.EvmKeeper)
		asset := sample.EthAddress().String()
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

		_, err := zk.FungibleKeeper.UpdateZRC20GasLimit(ctx, zrc20Addr, big.NewInt(42))
		require.NoError(t, err)

		gasLimit, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{

			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: chainID,
				Asset:         asset,
			}})
		require.NoError(t, err)
		require.Equal(t, uint64(42), gasLimit)
	})

	t.Run("should fail if no gas coin found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{

			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: 999999,
			}})
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if query gas limit for gas coin fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID()

		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().String(),
			ForeignChainId:       chainID,
			CoinType:             coin.CoinType_Gas,
		})

		// no contract deployed therefore will fail
		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{

			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_Gas,
				SenderChainId: chainID,
			}})
		require.ErrorIs(t, err, fungibletypes.ErrContractCall)
	})

	t.Run("should fail if no asset found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeper(t)

		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{

			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: 999999,
			}})
		require.ErrorIs(t, err, types.ErrForeignCoinNotFound)
	})

	t.Run("should fail if query gas limit for asset fails", func(t *testing.T) {
		k, ctx, _, zk := keepertest.CrosschainKeeper(t)
		k.GetAuthKeeper().GetModuleAccount(ctx, fungibletypes.ModuleName)

		chainID := getValidEthChainID()
		asset := sample.EthAddress().String()

		zk.FungibleKeeper.SetForeignCoins(ctx, fungibletypes.ForeignCoins{
			Zrc20ContractAddress: sample.EthAddress().String(),
			ForeignChainId:       chainID,
			CoinType:             coin.CoinType_ERC20,
			Asset:                asset,
		})

		// no contract deployed therefore will fail
		_, err := k.GetRevertGasLimit(ctx, types.CrossChainTx{
			InboundParams: &types.InboundParams{
				CoinType:      coin.CoinType_ERC20,
				SenderChainId: chainID,
				Asset:         asset,
			}})
		require.ErrorIs(t, err, fungibletypes.ErrContractCall)
	})
}

func TestGetAbortedAmount(t *testing.T) {
	amount := sdkmath.NewUint(100)
	t.Run("should return the inbound amount if outbound not present", func(t *testing.T) {
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				Amount: amount,
			},
		}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, amount, a)
	})
	t.Run("should return the amount outbound amount", func(t *testing.T) {
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				Amount: sdkmath.ZeroUint(),
			},
			OutboundParams: []*types.OutboundParams{
				{Amount: amount},
			},
		}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, amount, a)
	})
	t.Run("should return the zero if outbound amount is not present and inbound is 0", func(t *testing.T) {
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				Amount: sdkmath.ZeroUint(),
			},
		}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, sdkmath.ZeroUint(), a)
	})
	t.Run("should return the zero if no amounts are present", func(t *testing.T) {
		cctx := types.CrossChainTx{}
		a := crosschainkeeper.GetAbortedAmount(cctx)
		require.Equal(t, sdkmath.ZeroUint(), a)
	})
}

func Test_IsPending(t *testing.T) {
	tt := []struct {
		status   types.CctxStatus
		expected bool
	}{
		{types.CctxStatus_PendingInbound, false},
		{types.CctxStatus_PendingOutbound, true},
		{types.CctxStatus_PendingRevert, true},
		{types.CctxStatus_Reverted, false},
		{types.CctxStatus_Aborted, false},
		{types.CctxStatus_OutboundMined, false},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("status %s", tc.status), func(t *testing.T) {
			require.Equal(t, tc.expected, crosschainkeeper.IsPending(&types.CrossChainTx{CctxStatus: &types.Status{Status: tc.status}}))
		})
	}
}

func TestKeeper_UpdateNonce(t *testing.T) {
	t.Run("should error if supported chain is nil", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(nil)

		err := k.UpdateNonce(ctx, 5, nil)
		require.Error(t, err)
	})

	t.Run("should error if chain nonces not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{
			ChainName: 5,
			ChainId:   5,
		})
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).Return(observertypes.ChainNonces{}, false)
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				Amount: sdkmath.ZeroUint(),
			},
			OutboundParams: []*types.OutboundParams{
				{Amount: sdkmath.NewUint(1)},
			},
		}
		err := k.UpdateNonce(ctx, 5, &cctx)
		require.Error(t, err)
	})

	t.Run("should error if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{
			ChainName: 5,
			ChainId:   5,
		})
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).Return(observertypes.ChainNonces{
			Nonce: 100,
		}, true)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, false)
		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				Amount: sdkmath.ZeroUint(),
			},
			OutboundParams: []*types.OutboundParams{
				{Amount: sdkmath.NewUint(1)},
			},
		}
		err := k.UpdateNonce(ctx, 5, &cctx)
		require.Error(t, err)
		require.Equal(t, uint64(100), cctx.GetCurrentOutboundParam().TssNonce)
	})

	t.Run("should error if pending nonces not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{
			ChainName: 5,
			ChainId:   5,
		})
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).Return(observertypes.ChainNonces{
			Nonce: 100,
		}, true)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, true)
		observerMock.On("GetPendingNonces", mock.Anything, mock.Anything, mock.Anything).Return(observertypes.PendingNonces{}, false)

		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				Amount: sdkmath.ZeroUint(),
			},
			OutboundParams: []*types.OutboundParams{
				{Amount: sdkmath.NewUint(1)},
			},
		}
		err := k.UpdateNonce(ctx, 5, &cctx)
		require.Error(t, err)
	})

	t.Run("should error if nonces not equal", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{
			ChainName: 5,
			ChainId:   5,
		})
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).Return(observertypes.ChainNonces{
			Nonce: 100,
		}, true)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, true)
		observerMock.On("GetPendingNonces", mock.Anything, mock.Anything, mock.Anything).Return(observertypes.PendingNonces{
			NonceHigh: 99,
		}, true)

		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				Amount: sdkmath.ZeroUint(),
			},
			OutboundParams: []*types.OutboundParams{
				{Amount: sdkmath.NewUint(1)},
			},
		}
		err := k.UpdateNonce(ctx, 5, &cctx)
		require.Error(t, err)
	})

	t.Run("should update nonces", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t, keepertest.CrosschainMockOptions{
			UseObserverMock: true,
		})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, mock.Anything).Return(&chains.Chain{
			ChainName: 5,
			ChainId:   5,
		})
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).Return(observertypes.ChainNonces{
			Nonce: 100,
		}, true)
		observerMock.On("GetTSS", mock.Anything).Return(observertypes.TSS{}, true)
		observerMock.On("GetPendingNonces", mock.Anything, mock.Anything, mock.Anything).Return(observertypes.PendingNonces{
			NonceHigh: 100,
		}, true)

		observerMock.On("SetChainNonces", mock.Anything, mock.Anything).Once()
		observerMock.On("SetPendingNonces", mock.Anything, mock.Anything).Once()

		cctx := types.CrossChainTx{
			InboundParams: &types.InboundParams{
				Amount: sdkmath.ZeroUint(),
			},
			OutboundParams: []*types.OutboundParams{
				{Amount: sdkmath.NewUint(1)},
			},
		}
		err := k.UpdateNonce(ctx, 5, &cctx)
		require.NoError(t, err)
	})
}

func TestKeeper_SortCctxsByHeightAndChainId(t *testing.T) {
	// cctx1
	cctx1 := sample.CrossChainTx(t, "1-1")
	cctx1.GetCurrentOutboundParam().ReceiverChainId = 1
	cctx1.InboundParams.ObservedExternalHeight = 10

	// cctx2
	cctx2 := sample.CrossChainTx(t, "1-2")
	cctx2.GetCurrentOutboundParam().ReceiverChainId = 1
	cctx2.InboundParams.ObservedExternalHeight = 13

	// cctx3
	cctx3 := sample.CrossChainTx(t, "56-1")
	cctx3.GetCurrentOutboundParam().ReceiverChainId = 56
	cctx3.InboundParams.ObservedExternalHeight = 13

	// cctx4
	cctx4 := sample.CrossChainTx(t, "56-2")
	cctx4.GetCurrentOutboundParam().ReceiverChainId = 56
	cctx4.InboundParams.ObservedExternalHeight = 16

	// sort by height
	cctxs := []*types.CrossChainTx{cctx1, cctx2, cctx3, cctx4}
	keeper.SortCctxsByHeightAndChainID(cctxs)

	// check order
	require.Len(t, cctxs, 4)
	require.Equal(t, cctx1, cctxs[0])
	require.Equal(t, cctx2, cctxs[1])
	require.Equal(t, cctx3, cctxs[2])
	require.Equal(t, cctx4, cctxs[3])
}
