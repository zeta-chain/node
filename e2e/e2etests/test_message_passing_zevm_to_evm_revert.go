package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/contracts/testdapp"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	cctxtypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestMessagePassingZEVMtoEVMRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse the amount
	amount := parseBigInt(r, args[0])

	// Set destination details
	EVMChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	destinationAddress := r.EvmTestDAppAddr

	// Contract call originates from ZEVM chain
	r.ZEVMAuth.Value = amount
	tx, err := r.WZeta.Deposit(r.ZEVMAuth)
	require.NoError(r, err)

	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("wzeta deposit tx hash: %s", tx.Hash().Hex())
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta deposit")
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.WZeta.Approve(r.ZEVMAuth, r.ZevmTestDAppAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("wzeta approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "wzeta approve")
	utils.RequireTxSuccessful(r, receipt)

	testDAppZEVM, err := testdapp.NewTestDApp(r.ZevmTestDAppAddr, r.ZEVMClient)
	require.NoError(r, err)

	// Get ZETA balance before test
	previousBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	require.NoError(r, err)

	// Call the SendHelloWorld function on the ZEVM dapp Contract which would in turn create a new send, to be picked up by the zetanode evm hooks
	// set Do revert to true which adds a message to signal the EVM zetaReceiver to revert the transaction
	tx, err = testDAppZEVM.SendHelloWorld(r.ZEVMAuth, destinationAddress, EVMChainID, amount, true)
	require.NoError(r, err)

	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// New inbound message picked up by zetanode evm hooks and processed directly to initiate a contract call on EVM which would revert the transaction
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Reverted)

	// On finalization the Fungible module calls the onRevert function which in turn calls the onZetaRevert function on the sender contract
	receipt, err = r.ZEVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	require.NoError(r, err)
	utils.RequireTxSuccessful(r, receipt)

	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppZEVM.ParseRevertedHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received HelloWorld event:")
			receivedHelloWorldEvent = true
		}
	}
	require.True(r, receivedHelloWorldEvent, "expected Reverted HelloWorld event")

	// Check ZETA balance on ZEVM TestDApp and check new balance is between previous balance and previous balance + amount
	// New balance is increased because ZETA are sent from the sender but sent back to the contract
	// Contract receive less than the amount because of the gas fee to pay
	newBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceAndAmountZEVM := big.NewInt(0).Add(previousBalanceZEVM, amount)

	// check higher than previous balance and lower than previous balance + amount
	invariant := newBalanceZEVM.Cmp(previousBalanceZEVM) <= 0 || newBalanceZEVM.Cmp(previousBalanceAndAmountZEVM) > 0
	require.False(r,
		invariant,
		"expected new balance to be between %s and %s, got %s",
		previousBalanceZEVM.String(),
		previousBalanceAndAmountZEVM.String(),
		newBalanceZEVM.String(),
	)

	// Check ZETA balance on EVM TestDApp and check new balance is previous balance
	newBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	require.NoError(r, err)
	require.Equal(r, 0, newBalanceEVM.Cmp(previousBalanceEVM), "expected new balance to be equal to previous balance")
}
