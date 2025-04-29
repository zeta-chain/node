package runner

import (
	"context"
	"fmt"
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
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
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
	revertGasLimit math.Uint
}

type TONOpt func(t *tonOpts)

func TONExpectStatus(status cctypes.CctxStatus) TONOpt {
	return func(t *tonOpts) { t.expectedStatus = status }
}

// TONSetRevertGasLimit sets a higher gas limit for revert operations
func TONSetRevertGasLimit(gasLimit math.Uint) TONOpt {
	return func(t *tonOpts) { t.revertGasLimit = gasLimit }
}

// TONDeposit deposit TON to Gateway contract
func (r *E2ERunner) TONDeposit(
	gw *toncontracts.Gateway,
	sender *wallet.Wallet,
	amount math.Uint,
	zevmRecipient eth.Address,
) (*cctypes.CrossChainTx, error) {
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

	gwState, err := r.Clients.TON.GetAccountState(r.Ctx, gw.AccountID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get TON Gateway account state")
	}

	var (
		lastTxHash = gwState.LastTransHash
		lastLt     = gwState.LastTransLt
	)

	// Send TX
	err = gw.SendDeposit(r.Ctx, sender, amount, zevmRecipient, tonDepositSendCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send TON deposit")
	}

	filter := func(tx *ton.Transaction) bool {
		msgInfo := tx.Msgs.InMsg.Value.Value.Info.IntMsgInfo
		if msgInfo == nil {
			return false
		}

		from, err := ton.AccountIDFromTlb(msgInfo.Src)
		if err != nil {
			return false
		}

		return from.ToRaw() == sender.GetAddress().ToRaw()
	}

	waitFrom := tonWaitFrom{
		accountID:  gw.AccountID(),
		lastTxHash: ton.Bits256(lastTxHash),
		lastLt:     lastLt,
	}

	// Wait for cctx
	cctx := r.tonWaitForInboundCCTX(waitFrom, filter, cctypes.CctxStatus_OutboundMined)

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

	// If we're expecting a Reverted status, ensure we have enough gas
	if cfg.expectedStatus == cctypes.CctxStatus_Reverted {
		// Log that we're expecting a reverted status, but don't do anything special yet
		r.Logger.Info("Expecting Reverted status for this transaction")
	}

	gwState, err := r.Clients.TON.GetAccountState(r.Ctx, gw.AccountID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get TON Gateway account state")
	}

	var (
		lastTxHash = gwState.LastTransHash
		lastLt     = gwState.LastTransLt
	)

	// Log pre-transaction info
	r.Logger.Info("TON Pre-Transaction: lastTxHash=%v, lastLt=%v", lastTxHash, lastLt)

	// Send TX
	err = gw.SendDepositAndCall(r.Ctx, sender, amount, zevmRecipient, callData, tonDepositSendCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send TON deposit and call")
	}

	filter := func(tx *ton.Transaction) bool {
		msgInfo := tx.Msgs.InMsg.Value.Value.Info.IntMsgInfo
		if msgInfo == nil {
			return false
		}

		from, err := ton.AccountIDFromTlb(msgInfo.Src)
		if err != nil {
			return false
		}

		return from.ToRaw() == sender.GetAddress().ToRaw()
	}

	waitFrom := tonWaitFrom{
		accountID:  gw.AccountID(),
		lastTxHash: ton.Bits256(lastTxHash),
		lastLt:     lastLt,
	}

	// Wait for cctx
	cctx := r.tonWaitForInboundCCTX(waitFrom, filter, cfg.expectedStatus)

	// The relayed message might be stored as a hex string, so we need to check both formats
	if cctx.RelayedMessage != string(callData) {
		// Check if the relayed message is a hex encoding of the call data
		hexEncoded := fmt.Sprintf("%x", callData)
		if cctx.RelayedMessage != hexEncoded {
			require.Equal(r, string(callData), cctx.RelayedMessage,
				"CCTX relayed message doesn't match the callData (also checked hex format)")
		} else {
			r.Logger.Info("CCTX relayed message matched the hex-encoded callData: %s", hexEncoded)
		}
	}

	return cctx, nil
}

// SendWithdrawTONZRC20 sends withdraw tx of TON ZRC20 tokens
func (r *E2ERunner) SendWithdrawTONZRC20(
	to ton.AccountID,
	amount *big.Int,
	approveAmount *big.Int,
) *ethtypes.Transaction {

	tx, err := r.TONZRC20.Approve(r.ZEVMAuth, r.GatewayZEVMAddr, approveAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "approve")

	// Now perform the withdrawal
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

type tonWaitFrom struct {
	accountID  ton.AccountID
	lastTxHash ton.Bits256
	lastLt     uint64
}

// waits for specific inbound message for a given account and
func (r *E2ERunner) tonWaitForInboundCCTX(
	from tonWaitFrom,
	filter func(tx *ton.Transaction) bool,
	expectedStatus cctypes.CctxStatus,
) *cctypes.CrossChainTx {
	var (
		timeout  = 2 * time.Minute
		interval = time.Second
		status   = expectedStatus // Use the passed status directly
	)

	r.Logger.Info("tonWaitForInboundCCTX: Waiting for CCTX with expected status: %s", status.String())

	ctx, cancel := context.WithTimeout(r.Ctx, timeout)
	defer cancel()

	client := r.Clients.TON

	for {
		txs, err := client.GetTransactionsSince(ctx, from.accountID, from.lastLt, from.lastTxHash)
		require.NoError(r, err, "failed to getTransactionsSince")

		r.Logger.Info("tonWaitForInboundCCTX: Found %d transactions since last hash", len(txs))

		for i := range txs {
			tx := txs[i]

			// Apply the filter
			if !filter(&tx) {
				r.Logger.Info("tonWaitForInboundCCTX: Transaction %d filtered out", i)
				continue
			}

			r.Logger.Info("tonWaitForInboundCCTX: Found matching transaction, hash: %s", liteapi.TransactionToHashString(tx))

			// Get the CCTX by inbound hash
			cctx := utils.WaitCctxMinedByInboundHash(
				ctx,
				liteapi.TransactionToHashString(tx),
				r.CctxClient,
				r.Logger,
				r.CctxTimeout,
			)

			r.Logger.Info("tonWaitForInboundCCTX: Got CCTX with status: %s, requiring: %s", cctx.CctxStatus.Status.String(), status.String())

			// Verify the status matches what we expect
			utils.RequireCCTXStatus(r, cctx, status)

			r.Logger.Info("tonWaitForInboundCCTX: CCTX status verified successfully")

			return cctx
		}

		time.Sleep(interval)
	}
}
