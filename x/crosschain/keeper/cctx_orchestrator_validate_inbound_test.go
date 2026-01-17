package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/crypto"
	keepertest "github.com/zeta-chain/node/testutil/keeper"
	"github.com/zeta-chain/node/testutil/sample"
	"github.com/zeta-chain/node/x/crosschain/types"
	observerTypes "github.com/zeta-chain/node/x/observer/types"
)

func TestKeeper_ValidateInbound(t *testing.T) {
	t.Run("successfully validate inbound", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseObserverMock:  true,
				UseFungibleMock:  true,
				UseAuthorityMock: true,
			})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		inboundBlockHeight := uint64(420)
		inboundHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		receiverChain := chains.Goerli
		senderChain := chains.Goerli
		sender := sample.EthAddress()
		tssList := sample.TssList(3)

		// Set up mocks for CheckIfTSSMigrationTransfer
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).Return(senderChain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})
		// setup Mocks for GetTSS
		observerMock.On("GetTSS", mock.Anything).Return(tss, true)
		// setup Mocks for IsInboundEnabled
		observerMock.On("IsInboundEnabled", ctx).Return(true)
		// setup mocks for Initiate Outbound
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).
			Return(observerTypes.ChainNonces{Nonce: 1}, true)
		observerMock.On("GetPendingNonces", mock.Anything, mock.Anything, mock.Anything).
			Return(observerTypes.PendingNonces{NonceHigh: 1}, true)
		observerMock.On("SetChainNonces", mock.Anything, mock.Anything).Return(nil)
		observerMock.On("SetPendingNonces", mock.Anything, mock.Anything).Return(nil)
		// setup Mocks for SaveCCTXUpdate
		observerMock.On("SetNonceToCctx", mock.Anything, mock.Anything).Return(nil)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     senderChain.ChainId,
			MedianIndex: 0,
			Prices:      []uint64{100},
		})

		// call InitiateOutbound
		msg := types.MsgVoteInbound{
			Creator:            creator,
			Sender:             sender.String(),
			SenderChainId:      senderChain.ChainId,
			Receiver:           receiver.String(),
			ReceiverChain:      receiverChain.ChainId,
			Amount:             amount,
			Message:            message,
			InboundHash:        inboundHash.String(),
			InboundBlockHeight: inboundBlockHeight,
			CallOptions: &types.CallOptions{
				GasLimit: gasLimit,
			},
			CoinType:   cointType,
			TxOrigin:   sender.String(),
			Asset:      asset,
			EventIndex: eventIndex,
		}

		_, err := k.ValidateInbound(ctx, &msg, false)
		require.NoError(t, err)
		require.Len(t, k.GetAllCrossChainTx(ctx), 1)
	})

	t.Run("fail if tss not found", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseObserverMock:  true,
				UseFungibleMock:  true,
				UseAuthorityMock: true,
			})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		inboundBlockHeight := uint64(420)
		inboundHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		receiverChain := chains.Goerli
		senderChain := chains.Goerli
		sender := sample.EthAddress()

		// setup Mocks for GetTSS
		observerMock.On("GetTSS", mock.Anything).Return(tss, false)
		// setup Mocks for IsInboundEnabled

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     senderChain.ChainId,
			MedianIndex: 0,
			Prices:      []uint64{100},
		})

		// call InitiateOutbound
		msg := types.MsgVoteInbound{
			Creator:            creator,
			Sender:             sender.String(),
			SenderChainId:      senderChain.ChainId,
			Receiver:           receiver.String(),
			ReceiverChain:      receiverChain.ChainId,
			Amount:             amount,
			Message:            message,
			InboundHash:        inboundHash.String(),
			InboundBlockHeight: inboundBlockHeight,
			CallOptions: &types.CallOptions{
				GasLimit: gasLimit,
			},
			CoinType:   cointType,
			TxOrigin:   sender.String(),
			Asset:      asset,
			EventIndex: eventIndex,
		}

		_, err := k.ValidateInbound(ctx, &msg, false)
		require.ErrorIs(t, err, types.ErrCannotFindTSSKeys)
	})

	t.Run("fail if InitiateOutbound fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseObserverMock:  true,
				UseFungibleMock:  true,
				UseAuthorityMock: true,
			})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		inboundBlockHeight := uint64(420)
		inboundHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		receiverChain := chains.Goerli
		senderChain := chains.Goerli
		sender := sample.EthAddress()
		tssList := sample.TssList(3)

		// Set up mocks for CheckIfTSSMigrationTransfer
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).Return(senderChain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})
		// setup Mocks for GetTSS
		observerMock.On("GetTSS", mock.Anything).Return(tss, true)
		// setup Mocks for IsInboundEnabled
		observerMock.On("IsInboundEnabled", ctx).Return(true)
		// setup mocks for Initiate Outbound
		observerMock.On("GetChainNonces", mock.Anything, mock.Anything).
			Return(observerTypes.ChainNonces{Nonce: 1}, false)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     senderChain.ChainId,
			MedianIndex: 0,
			Prices:      []uint64{100},
		})

		// call InitiateOutbound
		msg := types.MsgVoteInbound{
			Creator:            creator,
			Sender:             sender.String(),
			SenderChainId:      senderChain.ChainId,
			Receiver:           receiver.String(),
			ReceiverChain:      receiverChain.ChainId,
			Amount:             amount,
			Message:            message,
			InboundHash:        inboundHash.String(),
			InboundBlockHeight: inboundBlockHeight,
			CallOptions: &types.CallOptions{
				GasLimit: gasLimit,
			},
			CoinType:   cointType,
			TxOrigin:   sender.String(),
			Asset:      asset,
			EventIndex: eventIndex,
		}

		_, err := k.ValidateInbound(ctx, &msg, false)
		require.ErrorIs(t, err, types.ErrCannotFindReceiverNonce)
	})

	t.Run("fail if inbound is disabled", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseObserverMock:  true,
				UseFungibleMock:  true,
				UseAuthorityMock: true,
			})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		inboundBlockHeight := uint64(420)
		inboundHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		tss := sample.Tss()
		receiverChain := chains.Goerli
		senderChain := chains.Goerli
		sender := sample.EthAddress()
		tssList := sample.TssList(3)

		// Set up mocks for CheckIfTSSMigrationTransfer
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).Return(senderChain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})
		// setup Mocks for GetTSS
		observerMock.On("GetTSS", mock.Anything).Return(tss, true)
		// setup Mocks for IsInboundEnabled
		observerMock.On("IsInboundEnabled", ctx).Return(false)

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     senderChain.ChainId,
			MedianIndex: 0,
			Prices:      []uint64{100},
		})

		// call InitiateOutbound
		msg := types.MsgVoteInbound{
			Creator:            creator,
			Sender:             sender.String(),
			SenderChainId:      senderChain.ChainId,
			Receiver:           receiver.String(),
			ReceiverChain:      receiverChain.ChainId,
			Amount:             amount,
			Message:            message,
			InboundHash:        inboundHash.String(),
			InboundBlockHeight: inboundBlockHeight,
			CallOptions: &types.CallOptions{
				GasLimit: gasLimit,
			},
			CoinType:   cointType,
			TxOrigin:   sender.String(),
			Asset:      asset,
			EventIndex: eventIndex,
		}

		_, err := k.ValidateInbound(ctx, &msg, false)
		require.ErrorIs(t, err, observerTypes.ErrInboundDisabled)
	})

	t.Run("fails when CheckIfTSSMigrationTransfer fails", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseObserverMock:  true,
				UseFungibleMock:  true,
				UseAuthorityMock: true,
			})

		// Setup mock data
		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		receiver := sample.EthAddress()
		creator := sample.AccAddress()
		amount := sdkmath.NewUint(42)
		message := "test"
		inboundBlockHeight := uint64(420)
		inboundHash := sample.Hash()
		gasLimit := uint64(100)
		asset := "test-asset"
		eventIndex := uint64(1)
		cointType := coin.CoinType_ERC20
		receiverChain := chains.Goerli
		senderChain := chains.Goerli
		sender := sample.EthAddress()
		tssList := sample.TssList(3)

		// setup Mocks for GetTSS
		observerMock.On("GetTSS", mock.Anything).Return(tssList[0], true)

		// Set up mocks for CheckIfTSSMigrationTransfer
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", mock.Anything, senderChain.ChainId).Return(senderChain, false)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		k.SetGasPrice(ctx, types.GasPrice{
			ChainId:     senderChain.ChainId,
			MedianIndex: 0,
			Prices:      []uint64{100},
		})

		// call InitiateOutbound
		msg := types.MsgVoteInbound{
			Creator:            creator,
			Sender:             sender.String(),
			SenderChainId:      senderChain.ChainId,
			Receiver:           receiver.String(),
			ReceiverChain:      receiverChain.ChainId,
			Amount:             amount,
			Message:            message,
			InboundHash:        inboundHash.String(),
			InboundBlockHeight: inboundBlockHeight,
			CallOptions: &types.CallOptions{
				GasLimit: gasLimit,
			},
			CoinType:   cointType,
			TxOrigin:   sender.String(),
			Asset:      asset,
			EventIndex: eventIndex,
		}

		_, err := k.ValidateInbound(ctx, &msg, false)
		require.ErrorIs(t, err, observerTypes.ErrSupportedChains)
	})
}
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

		err := k.CheckIfTSSMigrationTransfer(ctx, &msg)
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

		err := k.CheckIfTSSMigrationTransfer(ctx, &msg)
		require.NoError(t, err)
	})

	t.Run("fails when chain is not supported", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
				UseObserverMock:  true,
			})

		observerMock := keepertest.GetCrosschainObserverMock(t, k)
		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		chain := chains.Chain{
			ChainId: 999,
		}
		tssList := sample.TssList(3)
		sender := sample.AccAddress()

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, false)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err := k.CheckIfTSSMigrationTransfer(ctx, &msg)
		require.ErrorIs(t, err, observerTypes.ErrSupportedChains)
	})

	t.Run("skips check when an older tss address is invalid for bitcoin chain", func(t *testing.T) {
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

		err := k.CheckIfTSSMigrationTransfer(ctx, &msg)
		require.NoError(t, err)
	})

	t.Run("skips check when an older tss address is invalid for evm chain", func(t *testing.T) {
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

		err := k.CheckIfTSSMigrationTransfer(ctx, &msg)
		require.NoError(t, err)
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
		sender, err := crypto.GetTSSAddrEVM(tssList[0].TssPubkey)
		require.NoError(t, err)

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender.String(),
		}

		err = k.CheckIfTSSMigrationTransfer(ctx, &msg)
		require.ErrorIs(t, err, types.ErrMigrationFromOldTss)
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
		sender, err := crypto.GetTSSAddrBTC(tssList[0].TssPubkey, bitcoinParams)
		require.NoError(t, err)

		// Set up mocks
		observerMock.On("GetAllTSS", ctx).Return(tssList)
		observerMock.On("GetSupportedChainFromChainID", ctx, chain.ChainId).Return(chain, true)
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err = k.CheckIfTSSMigrationTransfer(ctx, &msg)
		require.ErrorIs(t, err, types.ErrMigrationFromOldTss)
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
			ChainId:     999,
			Network:     chains.Network_btc,
			Consensus:   chains.Consensus_bitcoin,
			CctxGateway: chains.CCTXGateway_observers,
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

		err := k.CheckIfTSSMigrationTransfer(ctx, &msg)
		require.ErrorContains(t, err, "no Bitcoin network params for chain ID: 999")
	})

	t.Run("fails if gateway is not observer ", func(t *testing.T) {
		k, ctx, _, _ := keepertest.CrosschainKeeperWithMocks(t,
			keepertest.CrosschainMockOptions{
				UseAuthorityMock: true,
			})

		authorityMock := keepertest.GetCrosschainAuthorityMock(t, k)
		chain := chains.GoerliLocalnet
		chain.CctxGateway = chains.CCTXGateway_zevm
		sender := sample.AccAddress()

		// Set up mocks
		authorityMock.On("GetAdditionalChainList", ctx).Return([]chains.Chain{})

		msg := types.MsgVoteInbound{
			SenderChainId: chain.ChainId,
			Sender:        sender,
		}

		err := k.CheckIfTSSMigrationTransfer(ctx, &msg)
		require.NoError(t, err)
	})
}
