package runner

import (
	"context"
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

	gwState, err := r.Clients.TON.GetAccountState(r.Ctx, gw.AccountID())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get TON Gateway account state")
	}

	var (
		lastTxHash = gwState.LastTransHash
		lastLt     = gwState.LastTransLt
	)

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

	// Debug info to help understand test failure
	r.Logger.Info("CCTX Debug Info:")
	r.Logger.Info("  Index: %s", cctx.Index)
	r.Logger.Info("  InboundTxParams.Sender: %s", cctx.InboundParams.Sender)
	r.Logger.Info("  InboundTxParams.SenderChainId: %d", cctx.InboundParams.SenderChainId)
	r.Logger.Info("  InboundTxParams.ObservedHash: %s", cctx.InboundParams.ObservedHash)
	r.Logger.Info("  RelayedMessage: %s", cctx.RelayedMessage)
	r.Logger.Info("  Status: %s", cctx.CctxStatus.Status.String())
	r.Logger.Info("  Sender Bytes Length: %d", len([]byte(cctx.InboundParams.Sender)))

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

type tonWaitFrom struct {
	accountID  ton.AccountID
	lastTxHash ton.Bits256
	lastLt     uint64
}

// waits for specific inbound message for a given account and
func (r *E2ERunner) tonWaitForInboundCCTX(
	from tonWaitFrom,
	filter func(tx *ton.Transaction) bool,
	expectedStatus ...cctypes.CctxStatus,
) *cctypes.CrossChainTx {
	var (
		timeout  = 2 * time.Minute
		interval = time.Second
		status   = cctypes.CctxStatus_OutboundMined // Default status
	)

	// Override default status if provided
	if len(expectedStatus) > 0 {
		status = expectedStatus[0]
	}

	ctx, cancel := context.WithTimeout(r.Ctx, timeout)
	defer cancel()

	client := r.Clients.TON

	for {
		txs, err := client.GetTransactionsSince(ctx, from.accountID, from.lastLt, from.lastTxHash)
		require.NoError(r, err, "failed to getTransactionsSince")

		for i := range txs {
			tx := txs[i]

			if !filter(&tx) {
				continue
			}

			cctx := utils.WaitCctxMinedByInboundHash(
				ctx,
				liteapi.TransactionToHashString(tx),
				r.CctxClient,
				r.Logger,
				r.CctxTimeout,
			)

			utils.RequireCCTXStatus(r, cctx, status)

			return cctx
		}

		time.Sleep(interval)
	}
}
