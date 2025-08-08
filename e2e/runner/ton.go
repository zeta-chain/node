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
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/utils"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
	"github.com/zeta-chain/node/zetaclient/chains/ton/rpc"
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

// TONDepositRaw deposits TON to Gateway contract and returns the raw tx. Doesn't wait for cctx to be mined.
func (r *E2ERunner) TONDepositRaw(
	gw *toncontracts.Gateway,
	sender *wallet.Wallet,
	amount math.Uint,
	zevmRecipient eth.Address,
) (ton.Transaction, error) {
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
		return ton.Transaction{}, errors.Wrap(err, "failed to get TON Gateway account state")
	}

	var (
		lastTxHash = gwState.LastTxHash
		lastLt     = gwState.LastTxLT
	)

	// Send TX
	err = gw.SendDeposit(r.Ctx, sender, amount, zevmRecipient, tonDepositSendCode)
	if err != nil {
		return ton.Transaction{}, errors.Wrap(err, "failed to send TON deposit")
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

	// Wait for tx
	tx := r.tonWaitForTx(waitFrom, filter)

	return tx, nil
}

// TONDeposit deposit TON to Gateway contract and wait for cctx to be mined
func (r *E2ERunner) TONDeposit(
	gw *toncontracts.Gateway,
	sender *wallet.Wallet,
	amount math.Uint,
	zevmRecipient eth.Address,
) (*cctypes.CrossChainTx, error) {
	tx, err := r.TONDepositRaw(gw, sender, amount, zevmRecipient)
	if err != nil {
		return nil, err
	}

	txHash := rpc.TransactionToHashString(tx)

	// Wait for cctx
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash, r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctypes.CctxStatus_OutboundMined)

	return cctx, nil
}

// TONDepositAndCall deposit TON to Gateway contract with call data and wait for cctx to be mined
func (r *E2ERunner) TONDepositAndCall(
	gw *toncontracts.Gateway,
	sender *wallet.Wallet,
	amount math.Uint,
	zevmRecipient eth.Address,
	callData []byte,
	opts ...TONOpt,
) (*cctypes.CrossChainTx, error) {
	cfg := r.tonEnsureCallOpts(sender, amount, zevmRecipient, callData, opts...)

	r.Logger.Info(
		"Sending deposit of %s TON from %s to zEVM %s and calling contract with 0x%x",
		amount.String(),
		sender.GetAddress().ToRaw(),
		zevmRecipient.Hex(),
		callData,
	)

	gwState, err := r.Clients.TON.GetAccountState(r.Ctx, gw.AccountID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get TON Gateway account state")
	}

	var (
		lastTxHash = gwState.LastTxHash
		lastLt     = gwState.LastTxLT
	)

	// Log pre-transaction info
	r.Logger.Info("TON depositAndCall: gateway tx before (last): hash=%s, lt=%v", lastTxHash.Hex(), lastLt)

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

	// Wait for tx
	tx := r.tonWaitForTx(waitFrom, filter)

	txHash := rpc.TransactionToHashString(tx)

	// Wait for cctx
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash, r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cfg.expectedStatus)

	// The relayed message might be stored as a hex string, so check both formats
	require.Contains(r,
		[]string{string(callData), fmt.Sprintf("%x", callData)},
		cctx.RelayedMessage,
		"CCTX relayed message mismatch",
	)

	return cctx, nil
}

func (r *E2ERunner) TONCall(
	gw *toncontracts.Gateway,
	sender *wallet.Wallet,
	amount math.Uint,
	zevmRecipient eth.Address,
	callData []byte,
	opts ...TONOpt,
) (*cctypes.CrossChainTx, error) {
	cfg := r.tonEnsureCallOpts(sender, amount, zevmRecipient, callData, opts...)

	r.Logger.Info(
		"Sending call with %s TON from %s to zEVM %s with data 0x%x",
		amount.String(),
		sender.GetAddress().ToRaw(),
		zevmRecipient.Hex(),
		callData,
	)

	gwState, err := r.Clients.TON.GetAccountState(r.Ctx, gw.AccountID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get TON Gateway account state")
	}

	var (
		lastTxHash = gwState.LastTxHash
		lastLt     = gwState.LastTxLT
	)

	// Log pre-transaction info
	r.Logger.Info("TON call: gateway tx before (last): hash=%s, lt=%v", lastTxHash.Hex(), lastLt)

	// Send TX
	err = gw.SendCall(r.Ctx, sender, amount, zevmRecipient, callData, tonDepositSendCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send TON call")
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

	// Wait for tx
	tx := r.tonWaitForTx(waitFrom, filter)

	txHash := rpc.TransactionToHashString(tx)

	// Wait for cctx
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, txHash, r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cfg.expectedStatus)

	// The relayed message might be stored as a hex string, so check both formats
	require.Contains(r,
		[]string{string(callData), fmt.Sprintf("%x", callData)},
		cctx.RelayedMessage,
		"CCTX relayed message mismatch",
	)

	return cctx, nil
}

func (r *E2ERunner) tonEnsureCallOpts(
	sender *wallet.Wallet,
	amount math.Uint,
	zevmRecipient eth.Address,
	callData []byte,
	opts ...TONOpt,
) *tonOpts {
	cfg := &tonOpts{
		expectedStatus: cctypes.CctxStatus_OutboundMined,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	require.NotNil(r, r.TONGateway, "TON Gateway is not initialized")
	require.NotNil(r, sender, "Sender wallet is nil")
	require.False(r, amount.IsZero(), "amount is zero")
	require.NotEqual(r, (eth.Address{}).String(), zevmRecipient.String(), "empty zevm recipient")
	require.NotEmpty(r, callData, "empty call data")

	return cfg
}

// SendWithdrawTONZRC20 sends withdraw tx of TON ZRC20 tokens
func (r *E2ERunner) SendWithdrawTONZRC20(
	to ton.AccountID,
	amount *big.Int,
	revertOptions gatewayzevm.RevertOptions,
) *ethtypes.Transaction {
	if revertOptions.OnRevertGasLimit == nil {
		revertOptions.OnRevertGasLimit = big.NewInt(0)
	}

	// Approve ALL TON-ZRC20 from runner's wallet to the gateway
	r.ApproveTONZRC20(r.GatewayZEVMAddr)

	// Perform the withdrawal
	tx, err := r.GatewayZEVM.Withdraw0(
		r.ZEVMAuth,
		[]byte(to.ToRaw()),
		amount,
		r.TONZRC20Addr,
		revertOptions,
	)
	require.NoError(r, err)
	r.Logger.EVMTransaction(tx, "zevm ton withdraw")

	// wait for tx receipt
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt, "withdraw")
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)

	return tx
}

// WithdrawTONZRC20 withdraws an amount of ZRC20 TON tokens and waits for the cctx to be mined
func (r *E2ERunner) WithdrawTONZRC20(
	recipient ton.AccountID,
	amount *big.Int,
	revertOptions gatewayzevm.RevertOptions,
) *cctypes.CrossChainTx {
	tx := r.SendWithdrawTONZRC20(recipient, amount, revertOptions)

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
func (r *E2ERunner) tonWaitForTx(from tonWaitFrom, filter func(tx *ton.Transaction) bool) ton.Transaction {
	var (
		timeout  = 2 * time.Minute
		interval = time.Second
		attempts = 0
	)

	ctx, cancel := context.WithTimeout(r.Ctx, timeout)
	defer cancel()

	client := r.Clients.TON

	for {
		txs, err := client.GetTransactionsSince(ctx, from.accountID, from.lastLt, from.lastTxHash)
		require.NoError(r, err, "failed to getTransactionsSince")

		r.Logger.Info("tonWaitForInboundCCTX[%d]: Found %d transactions since last hash", attempts, len(txs))

		for i := range txs {
			tx := txs[i]

			// Apply the filter
			if !filter(&tx) {
				r.Logger.Info("tonWaitForInboundCCTX[%d]: Transaction %d filtered out", attempts, i)
				continue
			}

			r.Logger.Info(
				"tonWaitForInboundCCTX[%d]: Found matching transaction: %s",
				attempts,
				rpc.TransactionToHashString(tx),
			)

			return tx
		}

		inc := 500 * time.Millisecond * time.Duration(attempts)
		time.Sleep(interval + inc)

		attempts++
	}
}
