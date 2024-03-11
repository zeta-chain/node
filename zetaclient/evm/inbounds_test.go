package evm_test

import (
	"encoding/hex"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/onrik/ethrpc"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/zetacore/common"
	observertypes "github.com/zeta-chain/zetacore/x/observer/types"
	"github.com/zeta-chain/zetacore/zetaclient/config"
	"github.com/zeta-chain/zetacore/zetaclient/evm"
	"github.com/zeta-chain/zetacore/zetaclient/interfaces"
	"github.com/zeta-chain/zetacore/zetaclient/testutils"
	"github.com/zeta-chain/zetacore/zetaclient/testutils/stub"
	clienttypes "github.com/zeta-chain/zetacore/zetaclient/types"
)

// MockEVMClient creates a mock ChainClient with custom chain, TSS, params etc
func MockEVMClient(
	chain common.Chain,
	tss interfaces.TSSSigner,
	lastBlock uint64,
	params observertypes.ChainParams) *evm.ChainClient {
	client := &evm.ChainClient{
		Tss: tss,
		Mu:  &sync.Mutex{},
	}
	client.WithChain(chain)
	client.WithZetaClient(stub.NewZetaCoreBridge())
	client.SetLastBlockHeight(lastBlock)
	client.SetChainParams(params)
	return client
}

func TestEVM_CheckAndVoteInboundTokenZeta(t *testing.T) {
	// load archived ZetaSent intx, receipt and cctx
	// https://etherscan.io/tx/0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76
	chainID := int64(1)
	intxHash := "0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76"
	confirmation := uint64(10)

	t.Run("should pass for archived intx, receipt and cctx", func(t *testing.T) {
		tx, receipt, cctx := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_Zeta)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		ballot, err := ob.CheckAndVoteInboundTokenZeta(tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundTxParams.InboundTxBallotIndex, ballot)
	})
	t.Run("should fail on unconfirmed intx", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_Zeta)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation - 1

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		_, err := ob.CheckAndVoteInboundTokenZeta(tx, receipt, false)
		require.ErrorContains(t, err, "not been confirmed")
	})
	t.Run("should not act if no ZetaSent event", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_Zeta)
		receipt.Logs = receipt.Logs[:2] // remove ZetaSent event
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		ballot, err := ob.CheckAndVoteInboundTokenZeta(tx, receipt, true)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
	})
	t.Run("should not act if emitter is not ZetaConnector", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_Zeta)
		chainID = 56 // use BSC chain connector
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		_, err := ob.CheckAndVoteInboundTokenZeta(tx, receipt, true)
		require.ErrorContains(t, err, "emitter address mismatch")
	})
}

func TestEVM_CheckAndVoteInboundTokenERC20(t *testing.T) {
	// load archived ERC20 intx, receipt and cctx
	// https://etherscan.io/tx/0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da
	chainID := int64(1)
	intxHash := "0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da"
	confirmation := uint64(10)

	t.Run("should pass for archived intx, receipt and cctx", func(t *testing.T) {
		tx, receipt, cctx := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_ERC20)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		ballot, err := ob.CheckAndVoteInboundTokenERC20(tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundTxParams.InboundTxBallotIndex, ballot)
	})
	t.Run("should fail on unconfirmed intx", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_ERC20)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation - 1

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		_, err := ob.CheckAndVoteInboundTokenERC20(tx, receipt, false)
		require.ErrorContains(t, err, "not been confirmed")
	})
	t.Run("should not act if no Deposit event", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_ERC20)
		receipt.Logs = receipt.Logs[:1] // remove Deposit event
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		ballot, err := ob.CheckAndVoteInboundTokenERC20(tx, receipt, true)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
	})
	t.Run("should not act if emitter is not ERC20 Custody", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_ERC20)
		chainID = 56 // use BSC chain ERC20 custody
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		_, err := ob.CheckAndVoteInboundTokenERC20(tx, receipt, true)
		require.ErrorContains(t, err, "emitter address mismatch")
	})
}

func TestEVM_CheckAndVoteInboundTokenGas(t *testing.T) {
	// load archived Gas intx, receipt and cctx
	// https://etherscan.io/tx/0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532
	chainID := int64(1)
	intxHash := "0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532"
	confirmation := uint64(10)

	t.Run("should pass for archived intx, receipt and cctx", func(t *testing.T) {
		tx, receipt, cctx := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_Gas)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		ballot, err := ob.CheckAndVoteInboundTokenGas(tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, cctx.InboundTxParams.InboundTxBallotIndex, ballot)
	})
	t.Run("should fail on unconfirmed intx", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_Gas)
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation - 1

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		_, err := ob.CheckAndVoteInboundTokenGas(tx, receipt, false)
		require.ErrorContains(t, err, "not been confirmed")
	})
	t.Run("should not act if receiver is not TSS", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_Gas)
		tx.To = testutils.OtherAddress // use other address
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		ballot, err := ob.CheckAndVoteInboundTokenGas(tx, receipt, false)
		require.ErrorContains(t, err, "not TSS address")
		require.Equal(t, "", ballot)
	})
	t.Run("should not act if transaction failed", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_Gas)
		receipt.Status = ethtypes.ReceiptStatusFailed
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		ballot, err := ob.CheckAndVoteInboundTokenGas(tx, receipt, false)
		require.ErrorContains(t, err, "not a successful tx")
		require.Equal(t, "", ballot)
	})
	t.Run("should not act on nil message", func(t *testing.T) {
		tx, receipt, _ := testutils.LoadEVMIntxNReceiptNCctx(t, chainID, intxHash, common.CoinType_Gas)
		tx.Input = hex.EncodeToString([]byte(common.DonationMessage)) // donation will result in nil message
		require.NoError(t, evm.ValidateEvmTransaction(tx))
		lastBlock := receipt.BlockNumber.Uint64() + confirmation

		ob := MockEVMClient(common.EthChain(), stub.NewTSSMainnet(), lastBlock, stub.MockChainParams(chainID, confirmation))
		ballot, err := ob.CheckAndVoteInboundTokenGas(tx, receipt, false)
		require.NoError(t, err)
		require.Equal(t, "", ballot)
	})
}

func TestEVM_BuildInboundVoteMsgForZetaSentEvent(t *testing.T) {
	// load archived ZetaSent receipt
	// https://etherscan.io/tx/0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76
	chainID := int64(1)
	intxHash := "0xf3935200c80f98502d5edc7e871ffc40ca898e134525c42c2ae3cbc5725f9d76"
	receipt := testutils.LoadEVMIntxReceipt(t, chainID, intxHash, common.CoinType_Zeta)
	cctx := testutils.LoadEVMIntxCctx(t, chainID, intxHash, common.CoinType_Zeta)

	// parse ZetaSent event
	ob := MockEVMClient(common.EthChain(), nil, 1, stub.MockChainParams(1, 1))
	connector := stub.MockConnectorNonEth(chainID)
	event := testutils.ParseReceiptZetaSent(receipt, connector)

	// create test compliance config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for archived ZetaSent event", func(t *testing.T) {
		msg := ob.BuildInboundVoteMsgForZetaSentEvent(event)
		require.NotNil(t, msg)
		require.Equal(t, cctx.InboundTxParams.InboundTxBallotIndex, msg.Digest())
	})
	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		sender := event.ZetaTxSenderAddress.Hex()
		cfg.ComplianceConfig.RestrictedAddresses = []string{sender}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForZetaSentEvent(event)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		receiver := clienttypes.BytesToEthHex(event.DestinationAddress)
		cfg.ComplianceConfig.RestrictedAddresses = []string{receiver}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForZetaSentEvent(event)
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if txOrigin is restricted", func(t *testing.T) {
		txOrigin := event.SourceTxOriginAddress.Hex()
		cfg.ComplianceConfig.RestrictedAddresses = []string{txOrigin}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForZetaSentEvent(event)
		require.Nil(t, msg)
	})
}

func TestEVM_BuildInboundVoteMsgForDepositedEvent(t *testing.T) {
	// load archived Deposited receipt
	// https://etherscan.io/tx/0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da
	chainID := int64(1)
	intxHash := "0x4ea69a0e2ff36f7548ab75791c3b990e076e2a4bffeb616035b239b7d33843da"
	tx, receipt := testutils.LoadEVMIntxNReceipt(t, chainID, intxHash, common.CoinType_ERC20)
	cctx := testutils.LoadEVMIntxCctx(t, chainID, intxHash, common.CoinType_ERC20)

	// parse Deposited event
	ob := MockEVMClient(common.EthChain(), nil, 1, stub.MockChainParams(1, 1))
	custody := stub.MockERC20Custody(chainID)
	event := testutils.ParseReceiptERC20Deposited(receipt, custody)
	sender := ethcommon.HexToAddress(tx.From)

	// create test compliance config
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for archived Deposited event", func(t *testing.T) {
		msg := ob.BuildInboundVoteMsgForDepositedEvent(event, sender)
		require.NotNil(t, msg)
		require.Equal(t, cctx.InboundTxParams.InboundTxBallotIndex, msg.Digest())
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
		event.Message = []byte(common.DonationMessage)
		msg := ob.BuildInboundVoteMsgForDepositedEvent(event, sender)
		require.Nil(t, msg)
	})
}

func TestEVM_BuildInboundVoteMsgForTokenSentToTSS(t *testing.T) {
	// load archived gas token transfer to TSS
	// https://etherscan.io/tx/0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532
	chainID := int64(1)
	intxHash := "0xeaec67d5dd5d85f27b21bef83e01cbdf59154fd793ea7a22c297f7c3a722c532"
	tx, receipt := testutils.LoadEVMIntxNReceipt(t, chainID, intxHash, common.CoinType_Gas)
	require.NoError(t, evm.ValidateEvmTransaction(tx))
	cctx := testutils.LoadEVMIntxCctx(t, chainID, intxHash, common.CoinType_Gas)

	// load archived gas token donation to TSS
	// https://etherscan.io/tx/0x52f214cf7b10be71f4d274193287d47bc9632b976e69b9d2cdeb527c2ba32155
	inTxHashDonation := "0x52f214cf7b10be71f4d274193287d47bc9632b976e69b9d2cdeb527c2ba32155"
	txDonation, receiptDonation := testutils.LoadEVMIntxNReceiptDonation(t, chainID, inTxHashDonation, common.CoinType_Gas)
	require.NoError(t, evm.ValidateEvmTransaction(txDonation))

	// create test compliance config
	ob := MockEVMClient(common.EthChain(), nil, 1, stub.MockChainParams(1, 1))
	cfg := config.Config{
		ComplianceConfig: config.ComplianceConfig{},
	}

	t.Run("should return vote msg for archived gas token transfer to TSS", func(t *testing.T) {
		msg := ob.BuildInboundVoteMsgForTokenSentToTSS(tx, ethcommon.HexToAddress(tx.From), receipt.BlockNumber.Uint64())
		require.NotNil(t, msg)
		require.Equal(t, cctx.InboundTxParams.InboundTxBallotIndex, msg.Digest())
	})
	t.Run("should return nil msg if sender is restricted", func(t *testing.T) {
		cfg.ComplianceConfig.RestrictedAddresses = []string{tx.From}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForTokenSentToTSS(tx, ethcommon.HexToAddress(tx.From), receipt.BlockNumber.Uint64())
		require.Nil(t, msg)
	})
	t.Run("should return nil msg if receiver is restricted", func(t *testing.T) {
		txCopy := &ethrpc.Transaction{}
		*txCopy = *tx
		message := hex.EncodeToString(ethcommon.HexToAddress(testutils.OtherAddress).Bytes())
		txCopy.Input = message // use other address as receiver
		cfg.ComplianceConfig.RestrictedAddresses = []string{testutils.OtherAddress}
		config.LoadComplianceConfig(cfg)
		msg := ob.BuildInboundVoteMsgForTokenSentToTSS(txCopy, ethcommon.HexToAddress(txCopy.From), receipt.BlockNumber.Uint64())
		require.Nil(t, msg)
	})
	t.Run("should return nil msg on donation transaction", func(t *testing.T) {
		msg := ob.BuildInboundVoteMsgForTokenSentToTSS(txDonation,
			ethcommon.HexToAddress(txDonation.From), receiptDonation.BlockNumber.Uint64())
		require.Nil(t, msg)
	})
}
