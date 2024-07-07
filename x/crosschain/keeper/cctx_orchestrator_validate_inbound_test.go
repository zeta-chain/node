package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/pkg/chains"
	"github.com/zeta-chain/zetacore/pkg/crypto"
	keepertest "github.com/zeta-chain/zetacore/testutil/keeper"
	"github.com/zeta-chain/zetacore/testutil/sample"
	"github.com/zeta-chain/zetacore/x/crosschain/types"
	observerTypes "github.com/zeta-chain/zetacore/x/observer/types"
)

func TestKeeper_CheckMigration(t *testing.T) {
	t.Run("Do not return error if sender is not a TSS address for evm chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
				UseObserverMock:  true,
			})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		chain := chains.Goerli
		tssList := sample.TssList(3)
		sender := sample.AccAddress()

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err := k.CheckMigration(ctx, &msg)
		require.NoError(t, err)
	})

	t.Run("Do not return error if sender is not a TSS address for btc chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
				UseObserverMock:  true,
			})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		chain := chains.BitcoinTestnet
		tssList := sample.TssList(3)
		sender := sample.AccAddress()

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err := k.CheckMigration(ctx, &msg)
		require.NoError(t, err)
	})

	t.Run("fails when chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
				UseObserverMock:  true,
			})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		chain := chains.Chain{
			ChainId: 999,
		}
		tssList := sample.TssList(3)
		sender := sample.AccAddress()

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, false)

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err := k.CheckMigration(ctx, &msg)
		require.ErrorIs(t, err, observerTypes.ErrSupportedChains)
	})

	t.Run("fails when tss address is invalid for bitcoin chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
				UseObserverMock:  true,
			})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		chain := chains.BitcoinTestnet
		tssList := sample.TssList(3)
		tssList[0].TssPubkey = "invalid"
		sender := sample.AccAddress()

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err := k.CheckMigration(ctx, &msg)
		require.ErrorIs(t, err, types.ErrInvalidAddress)
	})

	t.Run("fails when tss address is invalid for evm chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
				UseObserverMock:  true,
			})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		chain := chains.Goerli
		tssList := sample.TssList(3)
		tssList[0].TssPubkey = "invalid"
		sender := sample.AccAddress()

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err := k.CheckMigration(ctx, &msg)
		require.ErrorIs(t, err, types.ErrInvalidAddress)
	})

	t.Run("fails when sender is a TSS address for evm chain for evm chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
				UseObserverMock:  true,
			})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		chain := chains.Goerli
		tssList := sample.TssList(3)
		sender, err := crypto.GetTssAddrEVM(tssList[0].TssPubkey)
		require.NoError(t, err)

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender.String(),
		}

		err = k.CheckMigration(ctx, &msg)
		require.ErrorIs(t, err, types.ErrTssAddress)
	})

	t.Run("fails when sender is a TSS address for btc chain for btc chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
				UseObserverMock:  true,
			})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		chain := chains.BitcoinTestnet
		tssList := sample.TssList(3)
		bitcoinParams, err := chains.BitcoinNetParamsFromChainID(chain.ChainId)
		require.NoError(t, err)
		sender, err := crypto.GetTssAddrBTC(tssList[0].TssPubkey, bitcoinParams)
		require.NoError(t, err)

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err = k.CheckMigration(ctx, &msg)
		require.ErrorIs(t, err, types.ErrTssAddress)
	})

	t.Run("fails if bitcoin network params not found for BTC chain", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
				UseObserverMock:  true,
			})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		chain := chains.Chain{
			ChainId:   999,
			Consensus: chains.Consensus_bitcoin,
		}
		tssList := sample.TssList(3)
		sender := sample.AccAddress()

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{chain})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err := k.CheckMigration(ctx, &msg)
		require.ErrorContains(t, err, "no Bitcoin net params for chain ID: 999")
	})
}
