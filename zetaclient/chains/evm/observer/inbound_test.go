package observer

import (
	"encoding/hex"
	"errors"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/node/pkg/chains"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/constant"
	"github.com/zeta-chain/node/zetaclient/chains/evm/client"
	"github.com/zeta-chain/node/zetaclient/chains/evm/common"
	"github.com/zeta-chain/node/zetaclient/config"
	"github.com/zeta-chain/node/zetaclient/testutils"
	"github.com/zeta-chain/node/zetaclient/testutils/mocks"
	clienttypes "github.com/zeta-chain/node/zetaclient/types"
)

func Test_CheckAndVoteInboundTokenZeta(t *testing.T) {
	// load archived ZetaSent inbound, receipt and cctx
	// https://etherscan.io/tx/0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76
	chain := chains.Ethereum
	chainID := chain.ChainId
	inboundHash := "0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76"

	t.Run("should pass for archived inbound, receipt and cctx", func(t *testing.T) {
		ob := newTestSuite(t)

		tx, receipt, cctx := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_Zeta,
		)
		require.NoError(t, common.ValidateEvmTransaction(tx))

		ob.WithLastBlock(receipt.BlockNumber.Uint64() + ob.chainParams.InboundConfirmationSafe())

		ballot, err := ob.checkAndVoteInboundTokenZeta(ob.ctx, tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundParams.BallotIndex, ballot)
	})
	t.Run("should fail on unconfirmed inbound", func(t *testing.T) {
		ob := newTestSuite(t)

		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_Zeta,
		)
		require.NoError(t, common.ValidateEvmTransaction(tx))

		ob.WithLastBlock(receipt.BlockNumber.Uint64() + ob.chainParams.InboundConfirmationSafe() - 2)

		_, err := ob.checkAndVoteInboundTokenZeta(ob.ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not been confirmed")
	})
	t.Run("should not act if no ZetaSent event", func(t *testing.T) {
		ob := newTestSuite(t)

		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_Zeta,
		)
		receipt.Logs = receipt.Logs[:2] // remove ZetaSent event
		require.NoError(t, common.ValidateEvmTransaction(tx))

		ob.WithLastBlock(receipt.BlockNumber.Uint64() + ob.chainParams.InboundConfirmationSafe())

		ballot, err := ob.checkAndVoteInboundTokenZeta(ob.ctx, tx, receipt, true)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
	})
	t.Run("should not act if emitter is not ZetaConnector", func(t *testing.T) {
		// Given observer with another chain to trigger logic for
		// different evm address (based on mocked chain params)
		ob := newTestSuite(t, func(cfg *testSuiteConfig) { cfg.chain = &chains.BscMainnet })

		// Given tx from ETH
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t,
			TestDataDir,
			chains.Ethereum.ChainId,
			inboundHash,
			coin.CoinType_Zeta,
		)
		require.NoError(t, common.ValidateEvmTransaction(tx))

		ob.WithLastBlock(receipt.BlockNumber.Uint64() + ob.chainParams.InboundConfirmationSafe())

		// ACT
		_, err := ob.checkAndVoteInboundTokenZeta(ob.ctx, tx, receipt, true)

		// ASSERT
		require.ErrorContains(t, err, "emitter address mismatch")
	})
}

func Test_CheckAndVoteInboundTokenERC20(t *testing.T) {
	// load archived ERC20 inbound, receipt and cctx
	// https://etherscan.io/tx/0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da
	chain := chains.Ethereum
	chainID := chain.ChainId
	inboundHash := "0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da"

	t.Run("should pass for archived inbound, receipt and cctx", func(t *testing.T) {
		ob := newTestSuite(t)

		tx, receipt, cctx := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_ERC20,
		)
		require.NoError(t, common.ValidateEvmTransaction(tx))

		ob.WithLastBlock(receipt.BlockNumber.Uint64() + ob.chainParams.InboundConfirmationSafe())

		ballot, err := ob.checkAndVoteInboundTokenERC20(ob.ctx, tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundParams.BallotIndex, ballot)
	})
	t.Run("should fail on unconfirmed inbound", func(t *testing.T) {
		ob := newTestSuite(t)

		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_ERC20,
		)
		require.NoError(t, common.ValidateEvmTransaction(tx))

		ob.WithLastBlock(receipt.BlockNumber.Uint64() + ob.chainParams.InboundConfirmationSafe() - 2)

		_, err := ob.checkAndVoteInboundTokenERC20(ob.ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not been confirmed")
	})
	t.Run("should not act if no Deposit event", func(t *testing.T) {
		ob := newTestSuite(t)

		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_ERC20,
		)
		receipt.Logs = receipt.Logs[:1] // remove Deposit event
		require.NoError(t, common.ValidateEvmTransaction(tx))

		ob.WithLastBlock(receipt.BlockNumber.Uint64() + ob.chainParams.InboundConfirmationSafe())

		ballot, err := ob.checkAndVoteInboundTokenERC20(ob.ctx, tx, receipt, true)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
	})
	t.Run("should not act if emitter is not ERC20 Custody", func(t *testing.T) {
		// ARRANGE
		// Given observer with different chain (thus chain params) to have different evm addresses
		ob := newTestSuite(t, func(cfg *testSuiteConfig) { cfg.chain = &chains.BscMainnet })

		// Given tx from ETH
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chains.Ethereum.ChainId,
			inboundHash,
			coin.CoinType_ERC20,
		)
		require.NoError(t, common.ValidateEvmTransaction(tx))

		ob.WithLastBlock(receipt.BlockNumber.Uint64() + ob.chainParams.InboundConfirmationSafe())

		// ACT
		_, err := ob.checkAndVoteInboundTokenERC20(ob.ctx, tx, receipt, true)

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
	inboundHash := "0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532"

	t.Run("should pass for archived inbound, receipt and cctx", func(t *testing.T) {
		tx, receipt, cctx := testutils.LoadEVMInboundNReceiptNCctx(
			t,
			TestDataDir,
			chainID,
			inboundHash,
			coin.CoinType_Gas,
		)
		require.NoError(t, common.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := newTestSuite(t)
		ob.WithLastBlock(lastBlock)

		ballot, err := ob.checkAndVoteInboundTokenGas(ob.ctx, tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundParams.BallotIndex, ballot)
	})
	t.Run("should fail on unconfirmed inbound", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
		require.NoError(t, common.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation - 2

		ob := newTestSuite(t)
		ob.WithLastBlock(lastBlock)

		_, err := ob.checkAndVoteInboundTokenGas(ob.ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not been confirmed")
	})
	t.Run("should not act if receiver is not TSS", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
		tx.To = testutils.OtherAddress1 // use other address
		require.NoError(t, common.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := newTestSuite(t)
		ob.WithLastBlock(lastBlock)

		ballot, err := ob.checkAndVoteInboundTokenGas(ob.ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not TSS address")
		require.Equal(t, "", ballot)
	})
	t.Run("should not act if transaction failed", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
		receipt.Status = ethtypes.ReceiptStatusFailed
		require.NoError(t, common.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := newTestSuite(t)
		ob.WithLastBlock(lastBlock)

		ballot, err := ob.checkAndVoteInboundTokenGas(ob.ctx, tx, receipt, false)
		require.ErrorContains(t, err, "not a successful tx")
		require.Equal(t, "", ballot)
	})
	t.Run("should not act on nil message", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMInboundNReceiptNCctx(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
		tx.Input = hex.EncodeToString([]byte(constant.DonationMessage)) // donation will result in nil message
		require.NoError(t, common.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := newTestSuite(t)
		ob.WithLastBlock(lastBlock)

		ballot, err := ob.checkAndVoteInboundTokenGas(ob.ctx, tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
	})
}

func Test_BuildInboundVoteMsgForZetaSentEvent(t *testing.T) {
	// load archived ZetaSent receipt
	// https://etherscan.io/tx/0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76
	chainID := int64(1)
	inboundHash := "0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76"
	receipt := testutils.LoadEVMInboundReceipt(t, TestDataDir, chainID, inboundHash, coin.CoinType_Zeta)
	cctx := testutils.LoadCctxByInbound(t, chainID, coin.CoinType_Zeta, inboundHash)

	// parse ZetaSent event
	ob := newTestSuite(t)

	connector := mocks.MockConnectorNonEth(t, chainID)
	event := testutils.ParseReceiptZetaSent(receipt, connector)

	// create test compliance config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for archived ZetaSent event", func(t *testing.T) {
		msg := ob.buildInboundVoteMsgForZetaSentEvent(ob.appContext, event)
		require.NotNil(t, msg)
		require.Equal(t, cctx.InboundParams.BallotIndex, msg.Digest())
	})
	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		sender := event.ZetaTxSenderAddress.Hex()
		cfg.ComplianceConfig.RestrictedAddresses = []string{sender}
		config.SetRestrictedAddressesFromConfig(cfg)
		msg := ob.buildInboundVoteMsgForZetaSentEvent(ob.appContext, event)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		receiver := clienttypes.BytesToEthHex(event.DestinationAddress)
		cfg.ComplianceConfig.RestrictedAddresses = []string{receiver}
		config.SetRestrictedAddressesFromConfig(cfg)
		msg := ob.buildInboundVoteMsgForZetaSentEvent(ob.appContext, event)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if txOrigin is restricted", func(t *testing.T) {
		txOrigin := event.SourceTxOriginAddress.Hex()
		cfg.ComplianceConfig.RestrictedAddresses = []string{txOrigin}
		config.SetRestrictedAddressesFromConfig(cfg)
		msg := ob.buildInboundVoteMsgForZetaSentEvent(ob.appContext, event)
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
	ob := newTestSuite(t)
	custody := mocks.MockERC20Custody(t, chainID)
	event := testutils.ParseReceiptERC20Deposited(receipt, custody)
	sender := ethcommon.HexToAddress(tx.From)

	// create test compliance config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for archived Deposited event", func(t *testing.T) {
		msg := ob.buildInboundVoteMsgForDepositedEvent(event, sender)
		require.NotNil(t, msg)
		require.Equal(t, cctx.InboundParams.BallotIndex, msg.Digest())
	})
	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{sender.Hex()}
		config.SetRestrictedAddressesFromConfig(cfg)
		msg := ob.buildInboundVoteMsgForDepositedEvent(event, sender)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		receiver := clienttypes.BytesToEthHex(event.Recipient)
		cfg.ComplianceConfig.RestrictedAddresses = []string{receiver}
		config.SetRestrictedAddressesFromConfig(cfg)
		msg := ob.buildInboundVoteMsgForDepositedEvent(event, sender)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg on donation transaction", func(t *testing.T) {
		event.Message = []byte(constant.DonationMessage)
		msg := ob.buildInboundVoteMsgForDepositedEvent(event, sender)
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
	require.NoError(t, common.ValidateEvmTransaction(tx))
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
	require.NoError(t, common.ValidateEvmTransaction(txDonation))

	// create test compliance config
	ob := newTestSuite(t)
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for archived gas token transfer to TSS", func(t *testing.T) {
		msg := ob.buildInboundVoteMsgForTokenSentToTSS(
			tx,
			ethcommon.HexToAddress(tx.From),
			receipt.BlockNumber.Uint64(),
		)
		require.NotNil(t, msg)
		require.Equal(t, cctx.InboundParams.BallotIndex, msg.Digest())
	})
	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{tx.From}
		config.SetRestrictedAddressesFromConfig(cfg)
		msg := ob.buildInboundVoteMsgForTokenSentToTSS(
			tx,
			ethcommon.HexToAddress(tx.From),
			receipt.BlockNumber.Uint64(),
		)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		txCopy := &client.Transaction{}
		*txCopy = *tx
		message := hex.EncodeToString(ethcommon.HexToAddress(testutils.OtherAddress1).Bytes())
		txCopy.Input = message // use other address as receiver
		cfg.ComplianceConfig.RestrictedAddresses = []string{testutils.OtherAddress1}
		config.SetRestrictedAddressesFromConfig(cfg)
		msg := ob.buildInboundVoteMsgForTokenSentToTSS(
			txCopy,
			ethcommon.HexToAddress(txCopy.From),
			receipt.BlockNumber.Uint64(),
		)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg on donation transaction", func(t *testing.T) {
		msg := ob.buildInboundVoteMsgForTokenSentToTSS(txDonation,
			ethcommon.HexToAddress(txDonation.From), receiptDonation.BlockNumber.Uint64())
		require.Nil(t, msg)
	})
}

func Test_ObserveTSSReceiveInBlock(t *testing.T) {
	// https://etherscan.io/tx/0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532
	chain := chains.Ethereum
	chainID := chain.ChainId
	inboundHash := "0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532"

	// load archived tx and receipt
	tx, receipt := testutils.LoadEVMInboundNReceipt(t, TestDataDir, chainID, inboundHash, coin.CoinType_Gas)
	require.NoError(t, common.ValidateEvmTransaction(tx))

	// load archived evm block
	// https://etherscan.io/block/19363323
	blockNumber := receipt.BlockNumber.Uint64()
	block := testutils.LoadEVMBlock(t, TestDataDir, chainID, blockNumber, true)

	// test cases
	tests := []struct {
		name               string
		mockEVMClient      func(m *mocks.EVMClient)
		mockZetacoreClient func(m *mocks.ZetacoreClient)
		errMsg             string
	}{
		{
			name: "should observe TSS receive in block",
			mockEVMClient: func(m *mocks.EVMClient) {
				// feed block number and receipt to mock client
				m.On("BlockNumber", mock.Anything).Return(uint64(1000), nil)
				m.On("TransactionReceipt", mock.Anything, mock.Anything).Return(receipt, nil)
				m.On("BlockByNumberCustom", mock.Anything, mock.Anything).Return(block, nil)
			},
			mockZetacoreClient: func(m *mocks.ZetacoreClient) {
				m.On("GetCctxByHash", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
			},
			errMsg: "",
		},
		{
			name: "should not observe on error getting block",
			mockEVMClient: func(m *mocks.EVMClient) {
				// feed block number to allow construction of observer
				m.On("BlockNumber", mock.Anything).Unset()
				m.On("BlockByNumberCustom", mock.Anything, mock.Anything).Unset()
				m.On("BlockNumber", mock.Anything).Return(uint64(0), errors.New("RPC error"))
				m.On("BlockByNumberCustom", mock.Anything, mock.Anything).Return(nil, errors.New("RPC error"))
			},
			mockZetacoreClient: nil,
			errMsg:             "error getting block",
		},
		{
			name: "should not observe on error getting receipt",
			mockEVMClient: func(m *mocks.EVMClient) {
				// feed block number but RPC error on getting receipt
				m.On("BlockNumber", mock.Anything).Return(uint64(1000), nil)
				m.On("TransactionReceipt", mock.Anything, mock.Anything).Return(nil, errors.New("RPC error"))
				m.On("BlockByNumberCustom", mock.Anything, mock.Anything).Return(block, nil)
			},
			mockZetacoreClient: nil,
			errMsg:             "error getting receipt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ob := newTestSuite(t)
			ob.WithLastBlock(receipt.BlockNumber.Uint64() + ob.chainParams.InboundConfirmationSafe())

			if tt.mockEVMClient != nil {
				tt.mockEVMClient(ob.evmMock)
			}

			if tt.mockZetacoreClient != nil {
				tt.mockZetacoreClient(ob.zetacore)
			}

			err := ob.observeTSSReceiveInBlock(ob.ctx, blockNumber)
			if tt.errMsg != "" {
				require.ErrorContains(t, err, tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}
