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

// fungibleModuleAddress is a constant representing the EVM address of the Fungible module account
const fungibleModuleAddress = "0x735b14BB79463307AAcBED86DAf3322B1e6226aB"

func TestMessagePassingEVMtoZEVMRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	fungibleEthAddress := ethcommon.HexToAddress(fungibleModuleAddress)
	require.True(r, fungibleEthAddress != ethcommon.Address{}, "invalid fungible module address")

	// parse the amount
	amount := parseBigInt(r, args[0])

	// Set destination details
	zEVMChainID, err := r.ZEVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	destinationAddress := r.ZevmTestDAppAddr

	// Contract call originates from EVM chain
	tx, err := r.ZetaEth.Approve(r.EVMAuth, r.EvmTestDAppAddr, amount)
	require.NoError(r, err)

	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("Approve tx receipt: %d", receipt.Status)

	testDAppEVM, err := testdapp.NewTestDApp(r.EvmTestDAppAddr, r.EVMClient)
	require.NoError(r, err)

	// Get ZETA balance before test
	previousBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	require.NoError(r, err)

	previousFungibleBalance, err := r.WZeta.BalanceOf(&bind.CallOpts{}, fungibleEthAddress)
	require.NoError(r, err)

	// Call the SendHelloWorld function on the EVM dapp Contract which would in turn create a new send, to be picked up by the zeta-clients
	// set Do revert to true which adds a message to signal the ZEVM zetaReceiver to revert the transaction
	tx, err = testDAppEVM.SendHelloWorld(r.EVMAuth, destinationAddress, zEVMChainID, amount, true)
	require.NoError(r, err)

	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)

	// New inbound message picked up by zeta-clients and voted on by observers to initiate a contract call on zEVM which would revert the transaction
	// A revert transaction is created and gets fialized on the original sender chain.
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	utils.RequireCCTXStatus(r, cctx, cctxtypes.CctxStatus_Reverted)

	// On finalization the Tss address calls the onRevert function which in turn calls the onZetaRevert function on the sender contract
	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(cctx.GetCurrentOutboundParam().Hash))
	require.NoError(r, err)
	utils.RequireTxSuccessful(r, receipt)

	receivedHelloWorldEvent := false
	for _, log := range receipt.Logs {
		_, err := testDAppEVM.ParseRevertedHelloWorldEvent(*log)
		if err == nil {
			r.Logger.Info("Received RevertHelloWorld event:")
			receivedHelloWorldEvent = true
		}
	}
	require.True(r, receivedHelloWorldEvent, "expected Reverted HelloWorld event")

	// Check ZETA balance on ZEVM TestDApp and check new balance is previous balance
	newBalanceZEVM, err := r.WZeta.BalanceOf(&bind.CallOpts{}, r.ZevmTestDAppAddr)
	require.NoError(r, err)
	require.Equal(
		r,
		0,
		newBalanceZEVM.Cmp(previousBalanceZEVM),
		"expected new balance to be %s, got %s",
		previousBalanceZEVM.String(),
		newBalanceZEVM.String(),
	)

	// Check ZETA balance on EVM TestDApp and check new balance is between previous balance and previous balance + amount
	// New balance is increased because ZETA are sent from the sender but sent back to the contract
	// New balance is less than previous balance + amount because of the gas fee to pay
	newBalanceEVM, err := r.ZetaEth.BalanceOf(&bind.CallOpts{}, r.EvmTestDAppAddr)
	require.NoError(r, err)

	previousBalanceAndAmountEVM := big.NewInt(0).Add(previousBalanceEVM, amount)

	// check higher than previous balance and lower than previous balance + amount
	invariant := newBalanceEVM.Cmp(previousBalanceEVM) <= 0 || newBalanceEVM.Cmp(previousBalanceAndAmountEVM) > 0
	require.False(
		r,
		invariant,
		"expected new balance to be between %s and %s, got %s",
		previousBalanceEVM.String(),
		previousBalanceAndAmountEVM.String(),
		newBalanceEVM.String(),
	)

	// Check ZETA balance on Fungible Module and check new balance is previous balance
	newFungibleBalance, err := r.WZeta.BalanceOf(&bind.CallOpts{}, fungibleEthAddress)
	require.NoError(r, err)

	require.Equal(
		r,
		0,
		newFungibleBalance.Cmp(previousFungibleBalance),
		"expected new balance to be %s, got %s",
		previousFungibleBalance.String(),
		newFungibleBalance.String(),
	)
}
