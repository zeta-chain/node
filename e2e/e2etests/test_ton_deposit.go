package e2etests

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	"github.com/zeta-chain/node/testutil/sample"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestTONDeposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// Given deployer
	ctx, deployer, chain := r.Ctx, r.TONDeployer, chains.TONLocalnet

	// Given amount
	amount := parseUint(r, args[0])

	// Given approx deposit fee
	depositFee, err := r.TONGateway.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDeposit)
	require.NoError(r, err)

	// Given sample wallet with a balance of 50 TON
	sender, err := deployer.CreateWallet(ctx, toncontracts.Coins(50))
	require.NoError(r, err)

	// Given sample EVM address
	recipient := sample.EthAddress()

	// ACT
	err = r.TONDeposit(sender, amount, recipient)

	// ASSERT
	require.NoError(r, err)

	// Wait for CCTX mining
	filter := func(cctx *cctypes.CrossChainTx) bool {
		return cctx.InboundParams.SenderChainId == chain.ChainId &&
			cctx.InboundParams.Sender == sender.GetAddress().ToRaw()
	}

	cctx := r.WaitForSpecificCCTX(filter, time.Minute)

	// Check CCTX
	expectedDeposit := amount.Sub(depositFee)

	require.Equal(r, sender.GetAddress().ToRaw(), cctx.InboundParams.Sender)
	require.Equal(r, expectedDeposit.Uint64(), cctx.InboundParams.Amount.Uint64())

	// Check receiver's balance
	balance, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, recipient)
	require.NoError(r, err)

	r.Logger.Info("Recipient's zEVM TON balance after deposit: %d", balance.Uint64())

	require.Equal(r, expectedDeposit.Uint64(), balance.Uint64())
}
