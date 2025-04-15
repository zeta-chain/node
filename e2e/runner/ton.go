package runner

import (
	"encoding/hex"
	"math/big"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/query"
	eth "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/tonkeeper/tongo/ton"
	"github.com/tonkeeper/tongo/wallet"

	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

// we need to use this send mode due to how wallet V5 works
//
//	https://github.com/tonkeeper/w5/blob/main/contracts/wallet_v5.fc#L82
//	https://docs.ton.org/develop/smart-contracts/guidelines/message-modes-cookbook
const tonDepositSendCode = toncontracts.SendFlagSeparateFees + toncontracts.SendFlagIgnoreErrors

// currently implemented only for DepositAndCall,
// can be adopted for all TON ops
type tonOpts struct {
	expectedStatus cctypes.CctxStatus
}

type TONOpt func(t *tonOpts)

func TONExpectStatus(status cctypes.CctxStatus) TONOpt {
	return func(t *tonOpts) { t.expectedStatus = status }
}

// LogCCTXDetails logs all details of a cross-chain transaction
func (r *E2ERunner) LogCCTXDetails(cctx *cctypes.CrossChainTx) {
	if cctx == nil {
		r.Logger.Info("❌ CCTX is nil, cannot log details")
		return
	}

	r.Logger.Info("📋 CCTX Detailed Information:")
	r.Logger.Info("  - ID: %s", cctx.Index)
	r.Logger.Info("  - Creator: %s", cctx.Creator)

	// Inbound parameters
	if cctx.InboundParams != nil {
		r.Logger.Info("  - Inbound Parameters:")
		r.Logger.Info("    - SenderChainId: %d", cctx.InboundParams.SenderChainId)
		r.Logger.Info("    - Sender: %s", cctx.InboundParams.Sender)
		r.Logger.Info("    - TxOrigin: %s", cctx.InboundParams.TxOrigin)
		r.Logger.Info("    - CoinType: %s", cctx.InboundParams.CoinType.String())
		r.Logger.Info("    - Asset: %s", cctx.InboundParams.Asset)
		r.Logger.Info("    - Amount: %s", cctx.InboundParams.Amount.String())
		r.Logger.Info("    - ObservedHash: %s", cctx.InboundParams.ObservedHash)
		r.Logger.Info("    - ObservedExternalHeight: %d", cctx.InboundParams.ObservedExternalHeight)
		r.Logger.Info("    - BallotIndex: %s", cctx.InboundParams.BallotIndex)
		r.Logger.Info("    - FinalizedZetaHeight: %d", cctx.InboundParams.FinalizedZetaHeight)
		r.Logger.Info("    - TxFinalizationStatus: %s", cctx.InboundParams.TxFinalizationStatus.String())
		r.Logger.Info("    - IsCrossChainCall: %t", cctx.InboundParams.IsCrossChainCall)
		r.Logger.Info("    - InboundStatus: %s", cctx.InboundParams.Status.String())
		r.Logger.Info("    - ConfirmationMode: %s", cctx.InboundParams.ConfirmationMode.String())
	} else {
		r.Logger.Info("  - Inbound Parameters: None")
	}

	// Outbound parameters
	if len(cctx.OutboundParams) > 0 {
		r.Logger.Info("  - Outbound Parameters (%d):", len(cctx.OutboundParams))
		for i, outbound := range cctx.OutboundParams {
			r.Logger.Info("    - Outbound #%d:", i+1)
			r.Logger.Info("      - Receiver: %s", outbound.Receiver)
			r.Logger.Info("      - ReceiverChainId: %d", outbound.ReceiverChainId)
			r.Logger.Info("      - CoinType: %s", outbound.CoinType.String())
			r.Logger.Info("      - Amount: %s", outbound.Amount.String())
			r.Logger.Info("      - TssNonce: %d", outbound.TssNonce)
			r.Logger.Info("      - Hash: %s", outbound.Hash)
			r.Logger.Info("      - BallotIndex: %s", outbound.BallotIndex)
			r.Logger.Info("      - ObservedExternalHeight: %d", outbound.ObservedExternalHeight)
			r.Logger.Info("      - GasUsed: %d", outbound.GasUsed)
			r.Logger.Info("      - EffectiveGasPrice: %s", outbound.EffectiveGasPrice.String())
			r.Logger.Info("      - EffectiveGasLimit: %d", outbound.EffectiveGasLimit)
			r.Logger.Info("      - TssPubkey: %s", outbound.TssPubkey)
			r.Logger.Info("      - TxFinalizationStatus: %s", outbound.TxFinalizationStatus.String())

			// Log call options if available
			if outbound.CallOptions != nil {
				r.Logger.Info("      - CallOptions:")
				r.Logger.Info("        - GasLimit: %d", outbound.CallOptions.GasLimit)
				r.Logger.Info("        - IsArbitraryCall: %t", outbound.CallOptions.IsArbitraryCall)
			}

			r.Logger.Info("      - ConfirmationMode: %s", outbound.ConfirmationMode.String())
		}
	} else {
		r.Logger.Info("  - Outbound Parameters: None")
	}

	// CCTX Status
	if cctx.CctxStatus != nil {
		r.Logger.Info("  - Status:")
		r.Logger.Info("    - Status: %s", cctx.CctxStatus.Status.String())
		r.Logger.Info("    - StatusMessage: %s", cctx.CctxStatus.StatusMessage)
		if cctx.CctxStatus.ErrorMessage != "" {
			r.Logger.Info("    - ErrorMessage: %s", cctx.CctxStatus.ErrorMessage)
		}
		r.Logger.Info("    - LastUpdateTimestamp: %d", cctx.CctxStatus.LastUpdateTimestamp)
		r.Logger.Info("    - CreatedTimestamp: %d", cctx.CctxStatus.CreatedTimestamp)
		r.Logger.Info("    - IsAbortRefunded: %t", cctx.CctxStatus.IsAbortRefunded)

		if cctx.CctxStatus.ErrorMessageRevert != "" {
			r.Logger.Info("    - ErrorMessageRevert: %s", cctx.CctxStatus.ErrorMessageRevert)
		}
		if cctx.CctxStatus.ErrorMessageAbort != "" {
			r.Logger.Info("    - ErrorMessageAbort: %s", cctx.CctxStatus.ErrorMessageAbort)
		}
	} else {
		r.Logger.Info("  - Status: None")
	}

	// Additional metadata
	r.Logger.Info("  - ZetaFees: %s", cctx.ZetaFees.String())
	r.Logger.Info("  - ProtocolContractVersion: %s", cctx.ProtocolContractVersion)

	if cctx.RelayedMessage != "" {
		r.Logger.Info("  - RelayedMessage: %s", cctx.RelayedMessage)
	}

	// Revert Options
	r.Logger.Info("  - Revert Options:")
	r.Logger.Info("    - RevertAddress: %s", cctx.RevertOptions.RevertAddress)
	r.Logger.Info("    - AbortAddress: %s", cctx.RevertOptions.AbortAddress)
	r.Logger.Info("    - CallOnRevert: %t", cctx.RevertOptions.CallOnRevert)
	if len(cctx.RevertOptions.RevertMessage) > 0 {
		r.Logger.Info("    - RevertMessage: %s", hex.EncodeToString(cctx.RevertOptions.RevertMessage))
	}
	r.Logger.Info("    - RevertGasLimit: %s", cctx.RevertOptions.RevertGasLimit.String())
}

// TONDeposit deposit TON to Gateway contract
func (r *E2ERunner) TONDeposit(
	gw *toncontracts.Gateway,
	sender *wallet.Wallet,
	amount math.Uint,
	zevmRecipient eth.Address,
) (*cctypes.CrossChainTx, error) {
	// Log all existing CCTXs before starting
	r.Logger.Info("📋 Logging all existing CCTXs before deposit...")
	if err := r.TONDumpCCTXs(); err != nil {
		r.Logger.Info("⚠️ Failed to dump CCTXs: %v", err)
	}

	chain := chains.TONTestnet

	require.NotNil(r, r.TONGateway, "TON Gateway is not initialized")

	require.NotNil(r, sender, "Sender wallet is nil")
	require.False(r, amount.IsZero())
	require.NotEqual(r, (eth.Address{}).String(), zevmRecipient.String())

	r.Logger.Info(
		"Sending deposit of %s TON from %s to zEVM %s",
		amount.String(),
		sender.GetAddress().ToRaw(),
		zevmRecipient.Hex(),
	)

	// Debug information about addresses
	r.Logger.Info("🔍 Debug Information:")
	r.Logger.Info("  - TON address: %s", sender.GetAddress().ToRaw())
	r.Logger.Info("  - ZETA address: %s", zevmRecipient.Hex())
	r.Logger.Info("  - Gateway address: %s", gw.AccountID().ToRaw())

	// Send TX
	r.Logger.Info("📤 Sending TON transaction to blockchain...")
	err := gw.SendDeposit(r.Ctx, sender, amount, zevmRecipient, tonDepositSendCode)
	if err != nil {
		r.Logger.Error("Failed to send TON deposit: %v", err)
		return nil, errors.Wrap(err, "failed to send TON deposit")
	}
	r.Logger.Info("✅ TON deposit transaction sent successfully")

	// Get sender account details for reference
	senderAddr := sender.GetAddress()
	r.Logger.Info("📋 TON Transaction Details:")
	r.Logger.Info("  - Sender Address: %s", senderAddr.ToRaw())
	r.Logger.Info("  - Sender Address (human): %s", senderAddr.ToHuman(false, true))
	r.Logger.Info("  - Gateway: %s", gw.AccountID().ToRaw())
	r.Logger.Info("  - Amount: %s TON", amount.String())
	r.Logger.Info("  - Chain ID: %d", chain.ChainId)

	// Verify chain params are set
	chainParams, err := r.ObserverClient.GetChainParamsForChain(r.Ctx, &observertypes.QueryGetChainParamsForChainRequest{
		ChainId: chain.ChainId,
	})
	if err != nil {
		r.Logger.Print("⚠️ Unable to get chain params for TON: %v", err)
	} else {
		r.Logger.Print("✅ Chain params for TON are set")
		r.Logger.Print("🔍 Gateway address in chain params: %s", chainParams.ChainParams.GatewayAddress)
		r.Logger.Print("🔍 Expected gateway address: %s", gw.AccountID().ToRaw())
		r.Logger.Print("🔍 Gateway address match: %v", chainParams.ChainParams.GatewayAddress == gw.AccountID().ToRaw())
	}

	// Create a filter function to find matching CCTX
	senderAddress := sender.GetAddress().ToRaw()
	expectedChainId := chain.ChainId

	r.Logger.Info("🔍 Filter criteria for CCTX:")
	r.Logger.Info("  - Expected chain ID: %d", expectedChainId)
	r.Logger.Info("  - Expected sender: %s", senderAddress)

	filter := func(cctx *cctypes.CrossChainTx) bool {
		// Just check if it's from TON
		return cctx != nil &&
			cctx.InboundParams != nil &&
			cctx.InboundParams.SenderChainId == expectedChainId
	}

	// Wait for cctx to be mined
	r.Logger.Info("⏳ Waiting for CCTX to be processed...")
	cctx := r.WaitForSpecificCCTX(filter, cctypes.CctxStatus_OutboundMined, r.CctxTimeout)
	r.Logger.Info("✅ CCTX processed successfully")

	// Log detailed CCTX information
	r.LogCCTXDetails(cctx)

	r.Logger.Info("Transaction status: %s", cctx.CctxStatus.Status.String())

	return cctx, nil
}

// TONDepositAndCall deposit TON to Gateway contract with call data.
func (r *E2ERunner) TONDepositAndCall(
	gw *toncontracts.Gateway,
	sender *wallet.Wallet,
	amount math.Uint,
	zevmRecipient eth.Address,
	callData []byte,
	opts ...TONOpt,
) (*cctypes.CrossChainTx, error) {
	// Log all existing CCTXs before starting
	r.Logger.Info("📋 Logging all existing CCTXs before deposit-and-call...")
	if err := r.TONDumpCCTXs(); err != nil {
		r.Logger.Info("⚠️ Failed to dump CCTXs: %v", err)
	}

	cfg := &tonOpts{expectedStatus: cctypes.CctxStatus_OutboundMined}
	for _, opt := range opts {
		opt(cfg)
	}

	chain := chains.TONTestnet

	require.NotNil(r, r.TONGateway, "TON Gateway is not initialized")

	require.NotNil(r, sender, "Sender wallet is nil")
	require.False(r, amount.IsZero())
	require.NotEqual(r, (eth.Address{}).String(), zevmRecipient.String())
	require.NotEmpty(r, callData)

	r.Logger.Info(
		"Sending deposit of %s TON from %s to zEVM %s and calling contract with %q",
		amount.String(),
		sender.GetAddress().ToRaw(),
		zevmRecipient.Hex(),
		string(callData),
	)

	err := gw.SendDepositAndCall(r.Ctx, sender, amount, zevmRecipient, callData, tonDepositSendCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send TON deposit and call")
	}

	// Create a filter function to find matching CCTX
	senderAddress := sender.GetAddress().ToRaw()
	expectedChainId := chain.ChainId
	expectedRelayedMessage := hex.EncodeToString(callData)

	r.Logger.Info("🔍 Filter criteria for CCTX:")
	r.Logger.Info("  - Expected chain ID: %d", expectedChainId)
	r.Logger.Info("  - Expected sender: %s", senderAddress)
	r.Logger.Info("  - Expected relayed message: %s", expectedRelayedMessage)

	filter := func(cctx *cctypes.CrossChainTx) bool {
		return cctx != nil &&
			cctx.InboundParams != nil &&
			cctx.InboundParams.SenderChainId == expectedChainId &&
			cctx.InboundParams.Sender == senderAddress &&
			cctx.RelayedMessage == expectedRelayedMessage
	}

	// Wait for cctx with a 10-minute timeout
	r.Logger.Info("⏳ Waiting for CCTX to be processed (10 minute timeout)...")
	cctx := r.WaitForSpecificCCTX(filter, cfg.expectedStatus, 10*time.Minute)

	// Log detailed CCTX information
	r.LogCCTXDetails(cctx)

	r.Logger.Info("Transaction status: %s", cctx.CctxStatus.Status.String())

	return cctx, nil
}

// SendWithdrawTONZRC20 sends withdraw tx of TON ZRC20 tokens
func (r *E2ERunner) SendWithdrawTONZRC20(
	to ton.AccountID,
	amount *big.Int,
	approveAmount *big.Int,
) *ethtypes.Transaction {
	// approve
	tx, err := r.TONZRC20.Approve(r.ZEVMAuth, r.TONZRC20Addr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve")

	// withdraw
	tx, err = r.TONZRC20.Withdraw(r.ZEVMAuth, []byte(to.ToRaw()), amount)
	require.NoError(r, err)
	r.Logger.EVMTransaction(*tx, "withdraw")

	// wait for tx receipt
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "withdraw")
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	return tx
}

// WithdrawTONZRC20 withdraws an amount of ZRC20 TON tokens and waits for the cctx to be mined
func (r *E2ERunner) WithdrawTONZRC20(to ton.AccountID, amount *big.Int, approveAmount *big.Int) *cctypes.CrossChainTx {
	// Log all existing CCTXs before starting
	r.Logger.Info("📋 Logging all existing CCTXs before withdrawal...")
	if err := r.TONDumpCCTXs(); err != nil {
		r.Logger.Info("⚠️ Failed to dump CCTXs: %v", err)
	}

	tx := r.SendWithdrawTONZRC20(to, amount, approveAmount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctypes.CctxStatus_OutboundMined)

	// Log detailed CCTX information
	r.LogCCTXDetails(cctx)

	return cctx
}

// TONDumpCCTXs dumps all cross-chain transactions in the system
func (r *E2ERunner) TONDumpCCTXs() error {
	r.Logger.Info("📋 Dumping TON-related cross-chain transactions...")

	// Get all CCTXs with pagination
	var allCctxs []*cctypes.CrossChainTx
	nextKey := []byte{}
	pageSize := uint64(50)

	for {
		resp, err := r.CctxClient.CctxAll(
			r.Ctx,
			&cctypes.QueryAllCctxRequest{
				Pagination: &query.PageRequest{
					Key:        nextKey,
					Limit:      pageSize,
					CountTotal: true,
				},
			},
		)
		if err != nil {
			return errors.Wrap(err, "failed to get CCTXs")
		}

		// Filter to TON-related transactions
		for _, cctx := range resp.CrossChainTx {
			isTonRelated := false

			// Check inbound
			if cctx.InboundParams != nil && cctx.InboundParams.SenderChainId == chains.TONTestnet.ChainId {
				isTonRelated = true
			}

			// Check outbound
			if !isTonRelated && len(cctx.OutboundParams) > 0 {
				for _, outbound := range cctx.OutboundParams {
					if outbound.ReceiverChainId == chains.TONTestnet.ChainId {
						isTonRelated = true
						break
					}
				}
			}

			if isTonRelated {
				allCctxs = append(allCctxs, cctx)
			}
		}

		r.Logger.Info("Processed %d CCTXs (page of %d)", len(resp.CrossChainTx), pageSize)

		if len(resp.Pagination.NextKey) == 0 {
			r.Logger.Info("Total CCTXs found: %d, TON-related: %d", resp.Pagination.Total, len(allCctxs))
			break
		}

		nextKey = resp.Pagination.NextKey
	}

	if len(allCctxs) == 0 {
		r.Logger.Info("No TON-related CCTXs found in the system")
		return nil
	}

	// Log a summary of all CCTXs
	r.Logger.Info("📊 TON CCTX Summary (%d transactions):", len(allCctxs))
	for i, cctx := range allCctxs {
		statusStr := "Unknown"
		if cctx.CctxStatus != nil {
			statusStr = cctx.CctxStatus.Status.String()
		}

		// Basic info
		r.Logger.Info("[%d] ID: %s", i+1, cctx.Index)
		r.Logger.Info("  - Status: %s", statusStr)

		// Inbound info
		if cctx.InboundParams != nil {
			r.Logger.Info("  - From: Chain %d, Sender %s",
				cctx.InboundParams.SenderChainId,
				cctx.InboundParams.Sender)
			r.Logger.Info("  - Amount: %s (CoinType: %s)",
				cctx.InboundParams.Amount.String(),
				cctx.InboundParams.CoinType.String())
			r.Logger.Info("  - Hash: %s", cctx.InboundParams.ObservedHash)
		}

		// Outbound info (all outbounds)
		if len(cctx.OutboundParams) > 0 {
			r.Logger.Info("  - Outbounds (%d):", len(cctx.OutboundParams))
			for j, outbound := range cctx.OutboundParams {
				r.Logger.Info("    [%d] To: Chain %d, Receiver %s",
					j+1,
					outbound.ReceiverChainId,
					outbound.Receiver)
				r.Logger.Info("      - Hash: %s", outbound.Hash)
				r.Logger.Info("      - Status: %s", outbound.TxFinalizationStatus.String())
				r.Logger.Info("      - Amount: %s", outbound.Amount.String())
			}
		}

		r.Logger.Info("  -----------------")
	}

	// Get detailed information for each TON-related CCTX
	for _, cctx := range allCctxs {
		r.LogCCTXDetails(cctx)
	}

	// Print the exact format of sender address in the transaction
	for _, cctx := range allCctxs {
		if cctx.InboundParams != nil &&
			cctx.InboundParams.SenderChainId == chains.TONTestnet.ChainId {
			r.Logger.Info("Found TON tx with sender: %s", cctx.InboundParams.Sender)
		}
	}

	return nil
}
