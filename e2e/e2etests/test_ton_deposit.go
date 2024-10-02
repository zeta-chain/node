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

// TestTONDeposit (!) This boilerplate is a demonstration of E2E capabilities for TON integration
// Actual Deposit test is not implemented yet.
func TestTONDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// Given TON Localnet chain
	chain := chains.TONLocalnet

	// Given deployer
	ctx, deployer := r.Ctx, r.TONDeployer

	// Given amount
	amount := math.NewUintFromBigInt(parseBigInt(r, args[0]))

	// Given TON Gateway
	gw := toncontracts.NewGateway(r.TONGateway)

	// Given sample wallet with a balance of 50 TON
	sender, err := deployer.CreateWallet(ctx, ton.TONCoins(50))
	require.NoError(r, err)

	// Given sample EVM address
	recipient := sample.EthAddress()

	// ACT
	r.Logger.Print(
		"Sending deposit of %s TON from %s to zEVM %s",
		amount.String(),
		sender.GetAddress().ToRaw(),
		recipient.Hex(),
	)

	// we need to include this send mode due to how wallet V5 works
	//  https://github.com/tonkeeper/w5/blob/main/contracts/wallet_v5.fc#L82
	err = gw.SendDeposit(ctx, sender, amount, recipient, toncontracts.SendFlagIgnoreErrors)

	// ASSERT
	require.NoError(r, err)

	// Wait for CCTX mining
	cctxs := catchPendingCCTX(r, chain.ChainId, time.Minute)
	require.Len(r, cctxs, 1)

	cctx := cctxs[0]

	// Check cctx props
	require.NotNil(r, cctx.InboundParams)
	require.NotNil(r, cctx.InboundParams.Sender)
	//require.NoError(r, cctx.InboundParams.)
	//cctx

	// todo validate CCTX

	r.WaitForMinedCCTXFromIndex(cctx.Index)

	// Check sender's balance
	balance, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, recipient)
	require.NoError(r, err)

	r.Logger.Print("recipient's zEVM TON balance after deposit: %d", balance)

	// todo check balance equals to cctx.amount
}

// use another method - this returns pending only pending OUTBOUNDS
func catchPendingCCTX(r *runner.E2ERunner, chainID int64, timeout time.Duration) []*crosschainTypes.CrossChainTx {
	in := &crosschainTypes.QueryListPendingCctxRequest{ChainId: chainID}

	start := time.Now()

	for time.Since(start) < timeout {
		res, err := r.CctxClient.ListPendingCctx(r.Ctx, in)
		if err == nil && len(res.CrossChainTx) > 0 {
			return res.CrossChainTx
		}

		time.Sleep(time.Second)
	}

	r.Logger.Error("Timeout waiting for pending CCTX for chain %d", chainID)
	r.FailNow()

	return nil
}
