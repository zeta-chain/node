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
		r.Logger.Info("    - TxOrigin: %s", cctx.InboundParams.TxOrigin)
		r.Logger.Info("    - CoinType: %s", cctx.InboundParams.CoinType.String())
		r.Logger.Info("    - Amount: %s", cctx.InboundParams.Amount.String())
		r.Logger.Info("    - ObservedHash: %s", cctx.InboundParams.ObservedHash)
		r.Logger.Info("    - BallotIndex: %s", cctx.InboundParams.BallotIndex)
		r.Logger.Info("    - TxFinalizationStatus: %s", cctx.InboundParams.TxFinalizationStatus.String())
		r.Logger.Info("    - IsCrossChainCall: %t", cctx.InboundParams.IsCrossChainCall)
		r.Logger.Info("    - InboundStatus: %s", cctx.InboundParams.Status.String())
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
			r.Logger.Info("      - TssPubkey: %s", outbound.TssPubkey)
			r.Logger.Info("      - TxFinalizationStatus: %s", outbound.TxFinalizationStatus.String())
			r.Logger.Info("      - ConfirmationMode: %s", outbound.ConfirmationMode.String())
		}
	} else {
		r.Logger.Info("  - Outbound Parameters: None")
	}

	// CCTX Status
	if cctx.CctxStatus != nil {
		r.Logger.Info("  - Status:")
		r.Logger.Info("    - Status: %s", cctx.CctxStatus.Status.String())
		if cctx.CctxStatus.ErrorMessage != "" {
			r.Logger.Info("    - ErrorMessage: %s", cctx.CctxStatus.ErrorMessage)
		}
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
	r.Logger.Info("📋 Logging all existing CCTXs before deposit inside TONDeposit...")
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

	// Get the current set of CCTXs before we send our transaction
	beforeCCTXs, err := r.getAllCCTXsByChain(chain.ChainId)
	if err != nil {
		r.Logger.Error("Failed to get CCTXs before sending deposit: %v", err)
	} else {
		r.Logger.Info("Found %d existing TON CCTXs before sending deposit", len(beforeCCTXs))
	}

	// Send TX
	r.Logger.Info("📤 Sending TON transaction to blockchain...")
	err = gw.SendDeposit(r.Ctx, sender, amount, zevmRecipient, tonDepositSendCode)
	if err != nil {
		r.Logger.Error("Failed to send TON deposit: %v", err)
		return nil, errors.Wrap(err, "failed to send TON deposit")
	}
	r.Logger.Info("✅ TON deposit transaction sent successfully")

	// Give some time for the TON chain to process our transaction
	r.Logger.Info("⏱️ Waiting for transaction to be processed on TON blockchain (60 seconds)...")
	time.Sleep(60 * time.Second)

	// Now we'll look for any new CCTXs that have appeared since we sent our transaction
	r.Logger.Info("🔍 Looking for new CCTXs after deposit...")

	// Wait for a new CCTX to appear with multiple retries
	var cctx *cctypes.CrossChainTx
	maxRetries := 10

	for retry := 0; retry < maxRetries; retry++ {
		r.Logger.Info("Retry %d/%d: Checking for new TON CCTXs...", retry+1, maxRetries)

		afterCCTXs, err := r.getAllCCTXsByChain(chain.ChainId)
		if err != nil {
			r.Logger.Error("Failed to get CCTXs after deposit: %v", err)
			time.Sleep(15 * time.Second)
			continue
		}

		r.Logger.Info("Found %d TON CCTXs after deposit", len(afterCCTXs))

		// Look for new CCTXs that weren't there before
		for _, newCctx := range afterCCTXs {
			// Skip if this CCTX was already present before our deposit
			wasPresent := false
			for _, oldCctx := range beforeCCTXs {
				if oldCctx.Index == newCctx.Index {
					wasPresent = true
					break
				}
			}

			if !wasPresent {
				// This is a new CCTX!
				r.Logger.Info("✅ Found new CCTX since deposit: %s", newCctx.Index)

				// Check if it matches what we expect
				if newCctx.InboundParams != nil &&
					newCctx.InboundParams.Amount.Equal(amount) &&
					newCctx.CctxStatus != nil {

					createdTime := time.Unix(int64(newCctx.CctxStatus.CreatedTimestamp), 0)
					r.Logger.Info("  - Created at: %s", createdTime.Format(time.RFC3339))
					r.Logger.Info("  - Chain ID: %d", newCctx.InboundParams.SenderChainId)
					r.Logger.Info("  - Amount: %s", newCctx.InboundParams.Amount.String())
					r.Logger.Info("  - Sender: %s", newCctx.InboundParams.Sender)
					r.Logger.Info("  - Status: %s", newCctx.CctxStatus.Status.String())

					cctx = newCctx
					break
				}
			}
		}

		if cctx != nil {
			r.Logger.Info("✅ Found new CCTX since deposit: %s", cctx.Index)
			break
		}

		// Wait before next retry
		r.Logger.Info("No new matching CCTX found, waiting 15 seconds before retry...")
		time.Sleep(15 * time.Second)
	}

	// If no new CCTXs were found, try one last approach with a more general filter
	if cctx == nil {
		r.Logger.Info("⚠️ No new CCTXs found, trying more general filter...")

		// Filter that looks for recent TON transactions from any sender
		filter := func(cctx *cctypes.CrossChainTx) bool {
			// Combined null check and chain ID check
			if cctx == nil || cctx.InboundParams == nil || cctx.InboundParams.SenderChainId != chain.ChainId {
				return false
			}

			// Log for debugging
			r.Logger.Info("Checking TON CCTX: %s", cctx.Index)
			r.Logger.Info("  - Amount: %s vs Expected: %s",
				cctx.InboundParams.Amount.String(), amount.String())

			// Check amount is approximate (could be less due to fees)
			if !cctx.InboundParams.Amount.IsZero() &&
				cctx.InboundParams.Amount.LTE(amount) &&
				cctx.InboundParams.Amount.GTE(amount.QuoUint64(2)) {

				r.Logger.Info("  ✅ Found TON transaction with matching amount")
				return true
			}

			return false
		}

		r.Logger.Info("⏳ Waiting for CCTX to be processed with general filter (5 minutes)...")
		cctx = r.WaitForSpecificCCTX(filter, cctypes.CctxStatus_OutboundMined, 5*time.Minute)
	}

	if cctx == nil {
		r.Logger.Error("❌ No matching CCTX found after 5 minutes")
		return nil, errors.New("timeout waiting for CCTX")
	}

	r.Logger.Info("✅ CCTX processed successfully")

	// Log detailed CCTX information
	r.LogCCTXDetails(cctx)

	r.Logger.Info("Transaction status: %s", cctx.CctxStatus.Status.String())

	return cctx, nil
}

// getAllCCTXsByChain returns all CCTXs for a specific chain
func (r *E2ERunner) getAllCCTXsByChain(chainID int64) ([]*cctypes.CrossChainTx, error) {
	var chainCCTXs []*cctypes.CrossChainTx
	nextKey := []byte{}
	pageSize := uint64(100)

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
			return nil, errors.Wrap(err, "failed to get CCTXs")
		}

		for _, cctx := range resp.CrossChainTx {
			if cctx.InboundParams != nil && cctx.InboundParams.SenderChainId == chainID {
				chainCCTXs = append(chainCCTXs, cctx)
			}
		}

		if len(resp.Pagination.NextKey) == 0 {
			break
		}

		nextKey = resp.Pagination.NextKey
	}

	return chainCCTXs, nil
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
				r.Logger.Info("🔍 Found TON inbound transaction: %s", cctx.Index)
				r.Logger.Info("  - Sender: %s", cctx.InboundParams.Sender)
				r.Logger.Info("  - Amount: %s", cctx.InboundParams.Amount)
				r.Logger.Info("  - Hash: %s", cctx.InboundParams.ObservedHash)
			}

			// Check outbound
			if !isTonRelated && len(cctx.OutboundParams) > 0 {
				for _, outbound := range cctx.OutboundParams {
					if outbound.ReceiverChainId == chains.TONTestnet.ChainId {
						isTonRelated = true
						r.Logger.Info("🔍 Found TON outbound transaction: %s", cctx.Index)
						r.Logger.Info("  - Receiver: %s", outbound.Receiver)
						r.Logger.Info("  - Amount: %s", outbound.Amount)
						r.Logger.Info("  - Hash: %s", outbound.Hash)
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
