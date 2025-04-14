package runner

import (
	"encoding/hex"
	"math/big"
	"time"

	"cosmossdk.io/math"
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

// TONDeposit deposit TON to Gateway contract
func (r *E2ERunner) TONDeposit(
	gw *toncontracts.Gateway,
	sender *wallet.Wallet,
	amount math.Uint,
	zevmRecipient eth.Address,
) (*cctypes.CrossChainTx, error) {
	chain := chains.TONLocalnet

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

	// Log TON Gateway and chain information for debugging
	r.Logger.Info("TON Chain Information - Chain ID: %d, Gateway Address: %s",
		chain.ChainId,
		r.TONGateway.String())

	// Send TX
	err := gw.SendDeposit(r.Ctx, sender, amount, zevmRecipient, tonDepositSendCode)
	if err != nil {
		r.Logger.Error("Failed to send TON deposit: %v", err)
		return nil, errors.Wrap(err, "failed to send TON deposit")
	}

	r.Logger.Info("TON deposit transaction sent successfully, waiting for CCTX...")

	filter := func(cctx *cctypes.CrossChainTx) bool {
		match := cctx.InboundParams.SenderChainId == chain.ChainId &&
			cctx.InboundParams.Sender == sender.GetAddress().ToRaw()

		if match {
			r.Logger.Info("Found matching CCTX: ID=%s, Status=%s", cctx.Index, cctx.CctxStatus.Status.String())
			r.Logger.Info("CCTX details - Inbound hash: %s", cctx.InboundParams.ObservedHash)
			if len(cctx.OutboundParams) > 0 {
				r.Logger.Info("CCTX outbound details - Chain: %d, Status: %s",
					cctx.OutboundParams[0].ReceiverChainId,
					cctx.OutboundParams[0].TxFinalizationStatus.String())
			}
		} else {
			r.Logger.Info("CCTX doesn't match filter - CCTX Chain ID: %d, Sender: %s",
				cctx.InboundParams.SenderChainId,
				cctx.InboundParams.Sender)
		}

		return match
	}

	// Wait for cctx
	r.Logger.Info("Waiting for CCTX with expected status: %s", cctypes.CctxStatus_OutboundMined.String())
	cctx := r.WaitForSpecificCCTX(filter, cctypes.CctxStatus_OutboundMined, time.Minute)

	// Log detailed CCTX information
	r.Logger.Info("CCTX Processing Complete: ID=%s, Status=%s", cctx.Index, cctx.CctxStatus.Status.String())
	r.Logger.Info("CCTX Details: InboundTxHash=%s", cctx.InboundParams.ObservedHash)
	if len(cctx.OutboundParams) > 0 && cctx.OutboundParams[0] != nil {
		r.Logger.Info("CCTX Outbound Details: Hash=%s, Receiver=%s, Amount=%s",
			cctx.OutboundParams[0].Hash,
			cctx.OutboundParams[0].Receiver,
			cctx.OutboundParams[0].Amount.String())
	}

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
	cfg := &tonOpts{expectedStatus: cctypes.CctxStatus_OutboundMined}
	for _, opt := range opts {
		opt(cfg)
	}

	chain := chains.TONLocalnet

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

	// Log TON Gateway and chain information for debugging
	r.Logger.Info("TON Chain Information - Chain ID: %d, Gateway Address: %s, Expected Status: %s",
		chain.ChainId,
		r.TONGateway.String(),
		cfg.expectedStatus.String())

	err := gw.SendDepositAndCall(r.Ctx, sender, amount, zevmRecipient, callData, tonDepositSendCode)
	if err != nil {
		r.Logger.Error("Failed to send TON deposit and call: %v", err)
		return nil, errors.Wrap(err, "failed to send TON deposit and call")
	}

	r.Logger.Info("TON deposit and call transaction sent successfully, waiting for CCTX...")

	filter := func(cctx *cctypes.CrossChainTx) bool {
		match := cctx.InboundParams.SenderChainId == chain.ChainId &&
			cctx.InboundParams.Sender == sender.GetAddress().ToRaw() &&
			cctx.RelayedMessage == hex.EncodeToString(callData)

		if match {
			r.Logger.Info("Found matching CCTX: ID=%s, Status=%s", cctx.Index, cctx.CctxStatus.Status.String())
			r.Logger.Info("CCTX details - Inbound hash: %s", cctx.InboundParams.ObservedHash)
			if len(cctx.OutboundParams) > 0 {
				r.Logger.Info("CCTX outbound details - Chain: %d, Status: %s",
					cctx.OutboundParams[0].ReceiverChainId,
					cctx.OutboundParams[0].TxFinalizationStatus.String())
			}
		} else {
			r.Logger.Info("CCTX doesn't match filter - CCTX Chain ID: %d, Sender: %s, Message: %s",
				cctx.InboundParams.SenderChainId,
				cctx.InboundParams.Sender,
				cctx.RelayedMessage)
		}

		return match
	}

	// Wait for cctx
	r.Logger.Info("Waiting for CCTX with expected status: %s", cfg.expectedStatus.String())
	cctx := r.WaitForSpecificCCTX(filter, cfg.expectedStatus, time.Minute)

	// Log detailed CCTX information
	r.Logger.Info("CCTX Processing Complete: ID=%s, Status=%s", cctx.Index, cctx.CctxStatus.Status.String())
	r.Logger.Info("CCTX Details: InboundTxHash=%s", cctx.InboundParams.ObservedHash)
	if len(cctx.OutboundParams) > 0 && cctx.OutboundParams[0] != nil {
		r.Logger.Info("CCTX Outbound Details: Hash=%s, Receiver=%s, Amount=%s, Message=%s",
			cctx.OutboundParams[0].Hash,
			cctx.OutboundParams[0].Receiver,
			cctx.OutboundParams[0].Amount.String(),
			cctx.RelayedMessage)
	}

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
	tx := r.SendWithdrawTONZRC20(to, amount, approveAmount)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctypes.CctxStatus_OutboundMined)

	return cctx
}
