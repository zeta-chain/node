package e2etests

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	testcontract "github.com/zeta-chain/node/testutil/contracts"
	cctypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestTONDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// Given deployer
	ctx, deployer, chain := r.Ctx, r.TONDeployer, chains.TONLocalnet

	// Given amount
	amount := parseUint(r, args[0])

	// Given approx depositAndCall fee
	depositFee, err := r.TONGateway.GetTxFee(ctx, r.Clients.TON, toncontracts.OpDepositAndCall)
	require.NoError(r, err)

	// Given sample wallet with a balance of 50 TON
	sender, err := deployer.CreateWallet(ctx, toncontracts.Coins(50))
	require.NoError(r, err)

	// Given sample zEVM contract
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Example zevm contract deployed at: %s", contractAddr.String())

	// Given call data
	callData := []byte("hello from TON!")

	// ACT
	err = r.TONDepositAndCall(sender, amount, contractAddr, callData)

	// ASSERT
	require.NoError(r, err)

	// Wait for CCTX mining
	filter := func(cctx *cctypes.CrossChainTx) bool {
		return cctx.InboundParams.SenderChainId == chain.ChainId &&
			cctx.InboundParams.Sender == sender.GetAddress().ToRaw()
	}

	r.WaitForSpecificCCTX(filter, time.Minute)

	expectedDeposit := amount.Sub(depositFee)

	// check if example contract has been called, bar value should be set to amount
	utils.MustHaveCalledExampleContract(r, contract, expectedDeposit.BigInt())

	// Check receiver's balance
	balance, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, contractAddr)
	require.NoError(r, err)

	r.Logger.Info("Contract's zEVM TON balance after deposit: %d", balance.Uint64())

	require.Equal(r, expectedDeposit.Uint64(), balance.Uint64())
}
