package e2etests

import (
	"time"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/runner/ton"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
	crosschainTypes "github.com/zeta-chain/node/x/crosschain/types"
)

// we need to use this send mode due to how wallet V5 works
//
//	https://github.com/tonkeeper/w5/blob/main/contracts/wallet_v5.fc#L82
//	https://docs.ton.org/develop/smart-contracts/guidelines/message-modes-cookbook
const tonDepositSendCode = toncontracts.SendFlagSeparateFees + toncontracts.SendFlagIgnoreErrors

func TestTONDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// Given deployer
	ctx, deployer, chain := r.Ctx, r.TONDeployer, chains.TONLocalnet

	// Given amount
	amount := parseUint(r, args[0])

	// https://github.com/zeta-chain/protocol-contracts-ton/blob/main/contracts/gateway.fc#L28
	// (will be optimized & dynamic in the future)
	depositFee := math.NewUint(10_000_000)

	// Given TON Gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given sample wallet with a balance of 50 TON
	sender, err := deployer.CreateWallet(ctx, ton.TONCoins(50))
	require.NoError(r, err)

	// Given sample EVM address
	recipient := sample.EthAddress()

	// ACT
	r.Logger.Info(
		"Sending deposit of %s TON from %s to zEVM %s",
		amount.String(),
		sender.GetAddress().ToRaw(),
		recipient.Hex(),
	)

	err = gw.SendDeposit(ctx, sender, amount, recipient, tonDepositSendCode)

	// ASSERT
	require.NoError(r, err)

	// Wait for CCTX mining
	filter := func(cctx *crosschainTypes.CrossChainTx) bool {
		return cctx.InboundParams.SenderChainId == chain.ChainId &&
			cctx.InboundParams.Sender == sender.GetAddress().ToRaw()
	}

	cctxs := r.WaitForSpecificCCTX(filter, time.Minute)
	require.Len(r, cctxs, 1)

	cctx := r.WaitForMinedCCTXFromIndex(cctxs[0].Index)

	// Check CCTX
	expectedDeposit := amount.Sub(depositFee)

	require.Equal(r, sender.GetAddress().ToRaw(), cctx.InboundParams.Sender)
	require.Equal(r, expectedDeposit.Uint64(), cctx.InboundParams.Amount.Uint64())

	// Check sender's balance
	balance, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, recipient)
	require.NoError(r, err)

	r.Logger.Info("Recipient's zEVM TON balance after deposit: %d", balance.Uint64())

	require.Equal(r, expectedDeposit.Uint64(), balance.Uint64())
}
