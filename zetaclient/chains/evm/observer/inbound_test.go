package observer_test

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	zctx "github.com/zeta-chain/node/zetaclient/context"
	"github.com/zeta-chain/node/zetaclient/keys"

	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/zetaclient/chains/evm"
	"github.com/zeta-chain/node/zetaclient/chains/interfaces"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

func Test_CheckAndVoteInboundTokenZeta(t *testing.T) {
	// load archived ZetaSent inbound, receipt and cctx
	// https://etherscan.io/tx/0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76
	chain := chains.Ethereum
	confirmation := uint64(10)
	chainID := chain.ChainId
	chainParam := mocks.MockChainParams(chain.ChainId, confirmation)
	inboundHash := "0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76"

	ctx, _ := makeAppContext(t)

	t.Run("should pass for archived inbound, receipt and cctx", func(t *testing.T) {
		tx, receipt, cctx := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_Zeta,
		)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob, appContext := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		voteCtx := zctx.WithAppContext(context.Background(), appContext)

		ballot, err := ob.CheckAndVoteInboundTokenZeta(voteCtx, tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundParams.BallotIndex, ballot)
	})
	t.Run("should fail on unconfirmed inbound", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_Zeta,
		)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation - 1

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		_, err := ob.CheckAndVoteInboundTokenZeta(ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not been confirmed")
	})
	t.Run("should not act if no ZetaSent event", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_Zeta,
		)
		receipt.Logs = receipt.Logs[:2] // remove ZetaSent event
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		ballot, err := ob.CheckAndVoteInboundTokenZeta(ctx, tx, receipt, true)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
	})
	t.Run("should not act if emitter is not ZetaConnector", func(t *testing.T) {
		// Given tx from ETH
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t,
			TestDataDir,
			chains.Ethereum.ChainId,
			inboundHash,
			coin.CoinType_Zeta,
		)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		// Given BSC observer
		chain := chains.BscMainnet
		params := mocks.MockChainParams(chain.ChainId, confirmation)

		ob, _ := MockEVMObserver(
			t,
			chain,
			nil,
			nil,
			nil,
			nil,
			lastBlock,
			params,
		)

		// ACT
		_, err := ob.CheckAndVoteInboundTokenZeta(ctx, tx, receipt, true)

		// ASSERT
		require.ErrorContains(t, err, "emitter address mismatch")
	})
}

func Test_CheckAndVoteInboundTokenERC20(t *testing.T) {
	// load archived ERC20 inbound, receipt and cctx
	// https://etherscan.io/tx/0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da
	chain := chains.Ethereum
	confirmation := uint64(10)
	chainID := chain.ChainId
	chainParam := mocks.MockChainParams(chain.ChainId, confirmation)
	inboundHash := "0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da"

	ctx := context.Background()

	t.Run("should pass for archived inbound, receipt and cctx", func(t *testing.T) {
		tx, receipt, cctx := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_ERC20,
		)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		ballot, err := ob.CheckAndVoteInboundTokenERC20(ctx, tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundParams.BallotIndex, ballot)
	})
	t.Run("should fail on unconfirmed inbound", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_ERC20,
		)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation - 1

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		_, err := ob.CheckAndVoteInboundTokenERC20(ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not been confirmed")
	})
	t.Run("should not act if no Deposit event", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_ERC20,
		)
		receipt.Logs = receipt.Logs[:1] // remove Deposit event
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		ballot, err := ob.CheckAndVoteInboundTokenERC20(ctx, tx, receipt, true)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
	})
	t.Run("should not act if emitter is not ERC20 Custody", func(t *testing.T) {
		// ARRANGE
		// Given tx from ETH
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chains.Ethereum.ChainId,
			inboundHash,
			coin.CoinType_ERC20,
		)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		// Given BSC observer
		chain := chains.BscMainnet
		params := mocks.MockChainParams(chain.ChainId, confirmation)

		ob, _ := MockEVMObserver(
			t,
			chain,
			nil,
			nil,
			nil,
			nil,
			lastBlock,
			params,
		)

		// ACT
		_, err := ob.CheckAndVoteInboundTokenERC20(ctx, tx, receipt, true)

		// ASSERT
		require.ErrorContains(t, err, "emitter address mismatch")
	})
}

func Test_CheckAndVoteInboundTokenGas(t *testing.T) {
	// load archived Gas inbound, receipt and cctx
	// https://etherscan.io/tx/0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532
	chain := chains.Ethereum
	confirmation := uint64(10)
	chainID := chain.ChainId
	chainParam := mocks.MockChainParams(chain.ChainId, confirmation)
	inboundHash := "0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532"

	ctx := context.Background()

	t.Run("should pass for archived inbound, receipt and cctx", func(t *testing.T) {
		tx, receipt, cctx := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_Gas,
		)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		ballot, err := ob.CheckAndVoteInboundTokenGas(ctx, tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundParams.BallotIndex, ballot)
	})
	t.Run("should fail on unconfirmed inbound", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation - 1

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		_, err := ob.CheckAndVoteInboundTokenGas(ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not been confirmed")
	})
	t.Run("should not act if receiver is not TSS", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
		tx.To = testutils.OtherAddress1 // use other address
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		ballot, err := ob.CheckAndVoteInboundTokenGas(ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not TSS address")
		require.Equal(t, "", ballot)
	})
	t.Run("should not act if transaction failed", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
		receipt.Status = ethtypes.ReceiptStatusFailed
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		ballot, err := ob.CheckAndVoteInboundTokenGas(ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not a successful tx")
		require.Equal(t, "", ballot)
	})
	t.Run("should not act on nil message", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
		tx.Input = hex.EncodeToString([]byte(constant.DonationMessage)) // donation will result in nil message
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, lastBlock, chainParam)
		ballot, err := ob.CheckAndVoteInboundTokenGas(ctx, tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
	})
}

func Test_BuildInboundVoteMsgForZetaSentEvent(t *testing.T) {
	// load archived ZetaSent receipt
	// https://etherscan.io/tx/0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76
	chainID := int64(1)
	chain := chains.Ethereum
	inboundHash := "0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76"
	receipt := testutils.LoadEVMInboundReceipt(t, TestDataDir, chainID, inboundHash, coin.CoinType_Zeta)
	cctx := testutils.LoadCctxByInbound(t, chainID, coin.CoinType_Zeta, inboundHash)

	// parse ZetaSent event
	ob, app := MockEVMObserver(t, chain, nil, nil, nil, nil, 1, mocks.MockChainParams(1, 1))
	connector := mocks.MockConnectorNonEth(t, chainID)
	event := testutils.ParseReceiptZetaSent(receipt, connector)

	// create test compliance config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for archived ZetaSent event", func(t *testing.T) {
		msg := ob.BuildInboundVoteMsgForZetaSentEvent(app, event)
		require.NotNil(t, msg)
		require.Equal(t, cctx.InboundParams.BallotIndex, msg.Digest())
	})
	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		sender := event.ZetaTxSenderAddress.Hex()
		cfg.ComplianceConfig.RestrictedAddresses = []string{sender}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForZetaSentEvent(app, event)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		receiver := clienttypes.BytesToEthHex(event.DestinationAddress)
		cfg.ComplianceConfig.RestrictedAddresses = []string{receiver}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForZetaSentEvent(app, event)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if txOrigin is restricted", func(t *testing.T) {
		txOrigin := event.SourceTxOriginAddress.Hex()
		cfg.ComplianceConfig.RestrictedAddresses = []string{txOrigin}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForZetaSentEvent(app, event)
		require.Nil(t, msg)
	})
}

func Test_BuildInboundVoteMsgForDepositedEvent(t *testing.T) {
	// load archived Deposited receipt
	// https://etherscan.io/tx/0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da
	chain := chains.Ethereum
	chainID := chain.ChainId
	inboundHash := "0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da"
	tx, receipt := testutils.LoadEVMInboundNReceipt(t, TestDataDir, chainID, inboundHash, coin.CoinType_ERC20)
	cctx := testutils.LoadCctxByInbound(t, chainID, coin.CoinType_ERC20, inboundHash)

	// parse Deposited event
	ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, 1, mocks.MockChainParams(1, 1))
	custody := mocks.MockERC20Custody(t, chainID)
	event := testutils.ParseReceiptERC20Deposited(receipt, custody)
	sender := ethcommon.HexToAddress(tx.From)

	// create test compliance config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for archived Deposited event", func(t *testing.T) {
		msg := ob.BuildInboundVoteMsgForDepositedEvent(event, sender)
		require.NotNil(t, msg)
		require.Equal(t, cctx.InboundParams.BallotIndex, msg.Digest())
	})
	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{sender.Hex()}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForDepositedEvent(event, sender)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		receiver := clienttypes.BytesToEthHex(event.Recipient)
		cfg.ComplianceConfig.RestrictedAddresses = []string{receiver}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForDepositedEvent(event, sender)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg on donation transaction", func(t *testing.T) {
		event.Message = []byte(constant.DonationMessage)
		msg := ob.BuildInboundVoteMsgForDepositedEvent(event, sender)
		require.Nil(t, msg)
	})
}

func Test_BuildInboundVoteMsgForTokenSentToTSS(t *testing.T) {
	// load archived gas token transfer to TSS
	// https://etherscan.io/tx/0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532
	chain := chains.Ethereum
	chainID := chain.ChainId
	inboundHash := "0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532"
	tx, receipt := testutils.LoadEVMInboundNReceipt(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
	require.NoError(t, evm.ValidateEvmTransaction(tx))
	cctx := testutils.LoadCctxByInbound(t, chainID, coin.CoinType_Gas, inboundHash)

	// load archived gas token donation to TSS
	// https://etherscan.io/tx/0x52f214cf7b10be71f4d274193287d47bc9632b976e69b9d2cdeb527c2ba32155
	inboundHashDonation := "0x52f214cf7b10be71f4d274193287d47bc9632b976e69b9d2cdeb527c2ba32155"
	txDonation, receiptDonation := testutils.LoadEVMInboundNReceiptDonation(
		t,
		TestDataDir,
		chainID,
		inboundHashDonation,
		coin.CoinType_Gas,
	)
	require.NoError(t, evm.ValidateEvmTransaction(txDonation))

	// create test compliance config
	ob, _ := MockEVMObserver(t, chain, nil, nil, nil, nil, 1, mocks.MockChainParams(1, 1))
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for archived gas token transfer to TSS", func(t *testing.T) {
		msg := ob.BuildInboundVoteMsgForTokenSentToTSS(
			tx,
			ethcommon.HexToAddress(tx.From),
			receipt.BlockNumber.Uint64(),
		)
		require.NotNil(t, msg)
		require.Equal(t, cctx.InboundParams.BallotIndex, msg.Digest())
	})
	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{tx.From}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForTokenSentToTSS(
			tx,
			ethcommon.HexToAddress(tx.From),
			receipt.BlockNumber.Uint64(),
		)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		txCopy := &ethrpc.Transaction{}
		*txCopy = *tx
		message := hex.EncodeToString(ethcommon.HexToAddress(testutils.OtherAddress1).Bytes())
		txCopy.Input = message // use other address as receiver
		cfg.ComplianceConfig.RestrictedAddresses = []string{testutils.OtherAddress1}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForTokenSentToTSS(
			txCopy,
			ethcommon.HexToAddress(txCopy.From),
			receipt.BlockNumber.Uint64(),
		)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg on donation transaction", func(t *testing.T) {
		msg := ob.BuildInboundVoteMsgForTokenSentToTSS(txDonation,
			ethcommon.HexToAddress(txDonation.From), receiptDonation.BlockNumber.Uint64())
		require.Nil(t, msg)
	})
}

func Test_ObserveTSSReceiveInBlock(t *testing.T) {
	// https://etherscan.io/tx/0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532
	chain := chains.Ethereum
	chainID := chain.ChainId
	confirmation := uint64(1)
	chainParam := mocks.MockChainParams(chain.ChainId, confirmation)
	inboundHash := "0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532"

	// load archived tx and receipt
	tx, receipt := testutils.LoadEVMInboundNReceipt(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
	require.NoError(t, evm.ValidateEvmTransaction(tx))

	// load archived evm block
	// https://etherscan.io/block/19363323
	blockNumber := receipt.BlockNumber.Uint64()
	block := testutils.LoadEVMBlock(t, TestDataDir, chainID, blockNumber, true)

	// create mock zetacore client
	tss := mocks.NewTSSMainnet()
	lastBlock := receipt.BlockNumber.Uint64() + confirmation
	zetacoreClient := mocks.NewZetacoreClient(t).
		WithKeys(&keys.Keys{}).
		WithZetaChain().
		WithPostVoteInbound("", "").
		WithPostVoteInbound("", "")

	// test cases
	tests := []struct {
		name       string
		evmClient  interfaces.EVMRPCClient
		jsonClient interfaces.EVMJSONRPCClient
		errMsg     string
	}{
		{
			name: "should observe TSS receive in block",
			evmClient: func() interfaces.EVMRPCClient {
				// feed block number and receipt to mock client
				evmClient := mocks.NewEVMRPCClient(t)
				evmClient.On("BlockNumber", mock.Anything).Return(uint64(1000), nil)
				evmClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil)
				return evmClient
			}(),
			jsonClient: mocks.NewMockJSONRPCClient().WithBlock(block),
			errMsg:     "",
		},
		{
			name: "should not observe on error getting block",
			evmClient: func() interfaces.EVMRPCClient {
				// feed block number to allow construction of observer
				evmClient := mocks.NewEVMRPCClient(t)
				evmClient.On("BlockNumber", mock.Anything).Return(uint64(1000), nil)
				return evmClient
			}(),
			jsonClient: mocks.NewMockJSONRPCClient(), // no block
			errMsg:     "error getting block",
		},
		{
			name: "should not observe on error getting receipt",
			evmClient: func() interfaces.EVMRPCClient {
				// feed block number but RPC error on getting receipt
				evmClient := mocks.NewEVMRPCClient(t)
				evmClient.On("BlockNumber", mock.Anything).Return(uint64(1000), nil)
				evmClient.On("TransactionReceipt", mock.Anything, mock.Anything).Return(nil, errors.New("RPC error"))
				return evmClient
			}(),
			jsonClient: mocks.NewMockJSONRPCClient().WithBlock(block),
			errMsg:     "error getting receipt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob, _ := MockEVMObserver(t, chain, tt.evmClient, tt.jsonClient, zetacoreClient, tss, lastBlock, chainParam)
			err := ob.ObserveTSSReceiveInBlock(context.Background(), blockNumber)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func makeAppContext(t *testing.T) (context.Context, *zctx.AppContext) {
	var (
		app = zctx.New(config.New(false), nil, zerolog.New(zerolog.NewTestWriter(t)))
		ctx = context.Background()
	)

	return zctx.WithAppContext(ctx, app), app
}
