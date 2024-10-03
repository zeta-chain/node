package e2etests

import (
	"time"

	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/runner/ton"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
	toncontracts "github.com/zeta-chain/node/pkg/contracts/ton"
	testcontract "github.com/zeta-chain/node/testutil/contracts"
	crosschainTypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestTONDepositAndCall(r *runner.E2ERunner, args []string) {
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

	// Given sample zEVM contract
	contractAddr, _, contract, err := testcontract.DeployExample(r.ZEVMAuth, r.ZEVMClient)
	require.NoError(r, err)
	r.Logger.Info("Example zevm contract deployed at: %s", contractAddr.String())

	// Given call data
	callData := []byte("hello from TON!")

	// ACT
	r.Logger.Info(
		"Sending deposit of %s TON from %s to zEVM %s and calling contract with %q",
		amount.String(),
		sender.GetAddress().ToRaw(),
		contractAddr.Hex(),
		string(callData),
	)

	err = gw.SendDepositAndCall(ctx, sender, amount, contractAddr, callData, tonDepositSendCode)

	// ASSERT
	require.NoError(r, err)

	// Wait for CCTX mining
	filter := func(cctx *crosschainTypes.CrossChainTx) bool {
		return cctx.InboundParams.SenderChainId == chain.ChainId &&
			cctx.InboundParams.Sender == sender.GetAddress().ToRaw()
	}

	cctxs := r.WaitForSpecificCCTX(filter, time.Minute)
	require.Len(r, cctxs, 1)

	r.WaitForMinedCCTXFromIndex(cctxs[0].Index)

	expectedDeposit := amount.Sub(depositFee)

	// check if example contract has been called, bar value should be set to amount
	utils.MustHaveCalledExampleContract(r, contract, expectedDeposit.BigInt())

	// Check sender's balance
	balance, err := r.TONZRC20.BalanceOf(&bind.CallOpts{}, contractAddr)
	require.NoError(r, err)

	r.Logger.Info("Contract's zEVM TON balance after deposit: %d", balance.Uint64())

	require.Equal(r, expectedDeposit.Uint64(), balance.Uint64())
}
