package e2etests

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/zetaclient/chains/ton/liteapi"
)

func TestTONWithdraw(r *runner.E2ERunner, args []string) {
	// ARRANGE
	require.Len(r, args, 1)

	// Given a deployer
	_, deployer := r.Ctx, r.TONDeployer

	// That donates 100 TON to some zEVM sender
	zevmSender := r.ZEVMAuth.From

	_, err := r.TONDeposit(&deployer.Wallet, toncontracts.Coins(100), zevmSender)
	require.NoError(r, err)

	// Given his ZRC-20 balance
	senderZRC20BalanceBefore, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, zevmSender)
	require.NoError(r, err)
	r.Logger.Print("zEVM sender's ZRC20 TON balance before withdraw: %d", senderZRC20BalanceBefore)

	// Given another TON wallet
	tonRecipient, err := deployer.CreateWallet(r.Ctx, toncontracts.Coins(1))
	require.NoError(r, err)

	tonRecipientBalanceBefore, err := deployer.GetBalanceOf(r.Ctx, tonRecipient.GetAddress())
	require.NoError(r, err)

	r.Logger.Print("Recipient's TON balance before withdrawal: %s", toncontracts.FormatCoins(tonRecipientBalanceBefore))

	// Given amount to withdraw (and approved amount in TON ZRC20 to cover the gas fee)
	amount := parseUint(r, args[0])
	approvedAmount := amount.Add(toncontracts.Coins(1))

	// ACT
	cctx := r.WithdrawTONZRC20(tonRecipient.GetAddress(), amount.BigInt(), approvedAmount.BigInt())

	// ASSERT
	r.Logger.Print(
		"Withdraw TON ZRC20 transaction (with %s) sent: %+v",
		toncontracts.FormatCoins(amount),
		map[string]any{
			"zevm_sender":   zevmSender.Hex(),
			"ton_recipient": tonRecipient.GetAddress().ToRaw(),
			"ton_amount":    toncontracts.FormatCoins(amount),
			"cctx_index":    cctx.Index,
			"ton_hash":      cctx.GetCurrentOutboundParam().Hash,
			"zevm_hash":     cctx.InboundParams.ObservedHash,
		},
	)

	// Make sure that recipient's TON balance has increased
	tonRecipientBalanceAfter, err := deployer.GetBalanceOf(r.Ctx, tonRecipient.GetAddress())
	require.NoError(r, err)

	r.Logger.Print("Recipient's balance after withdrawal: %s", toncontracts.FormatCoins(tonRecipientBalanceAfter))

	// Make sure that sender's ZRC20 balance has decreased
	senderZRC20BalanceAfter, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, zevmSender)
	require.NoError(r, err)
	r.Logger.Print("zEVM sender's ZRC20 TON balance after withdraw: %d", senderZRC20BalanceAfter)
	r.Logger.Print(
		"zEVM sender's ZRC20 TON balance diff: %d",
		big.NewInt(0).Sub(senderZRC20BalanceBefore, senderZRC20BalanceAfter),
	)

	// Make sure that TON withdrawal CCTX contain outgoing message with exact withdrawal amount
	lt, hash, err := liteapi.TransactionHashFromString(cctx.GetCurrentOutboundParam().Hash)
	require.NoError(r, err)

	txs, err := r.Clients.TON.GetTransactions(r.Ctx, 1, r.TONGateway.AccountID(), lt, hash)
	require.NoError(r, err)
	require.Len(r, txs, 1)

	// TON coins that were withdrawn from GW to the recipient
	inMsgAmount := math.NewUint(
		uint64(txs[0].Msgs.OutMsgs.Values()[0].Value.Info.IntMsgInfo.Value.Grams),
	)

	// #nosec G115 always in range
	require.Equal(r, int(amount.Uint64()), int(inMsgAmount.Uint64()))
}

// TODO: Add "withdraw_many_concurrent" test
// https://github.com/zeta-chain/node/issues/3044
