package runner

import (
	"context"
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
	cctx := r.tonWaitForInboundCCTX(waitFrom, filter)

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

	err := gw.SendDepositAndCall(r.Ctx, sender, amount, zevmRecipient, callData, tonDepositSendCode)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send TON deposit and call")
	}

	filter := func(cctx *cctypes.CrossChainTx) bool {
		return cctx.InboundParams.SenderChainId == chain.ChainId &&
			cctx.InboundParams.Sender == sender.GetAddress().ToRaw() &&
			cctx.RelayedMessage == hex.EncodeToString(callData)
	}

	// Wait for cctx
	cctx := r.WaitForSpecificCCTX(filter, cfg.expectedStatus, time.Minute)

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
) *cctypes.CrossChainTx {
	var (
		timeout  = 2 * time.Minute
		interval = time.Second
	)

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

			utils.RequireCCTXStatus(r, cctx, cctypes.CctxStatus_OutboundMined)

			return cctx
		}

		time.Sleep(interval)
	}
}
