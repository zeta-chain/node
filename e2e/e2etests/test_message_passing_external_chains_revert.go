package e2etests

import (
	"math/big"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/contracts/testdapp"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	cctxtypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// TestMessagePassingRevertSuccessExternalChains tests message passing with successful revert between external EVM chains
// TODO: Use two external EVM chains for these tests
// https://github.com/zeta-chain/node/issues/2185
func TestMessagePassingRevertSuccessExternalChains(r *runner.E2ERunner, args []string) {
	if len(args) != 1 {
		panic("TestMessagePassingRevertSuccess requires exactly one argument for the amount.")
	}

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	if !ok {
		panic("Invalid amount specified for TestMessagePassingRevertSuccess.")
	}

	chainID, err := r.EVMClient.ChainID(r.Ctx)
	if err != nil {
		panic(err)
	}

	auth := r.EVMAuth

	tx, err := r.ZetaEth.Approve(auth, r.EvmTestDAppAddr, amount)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Approve tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("tx failed")
	}
	r.Logger.Info("Approve tx receipt: %d", receipt.Status)

	r.Logger.Info("Calling TestDApp.SendHello on contract address %s", r.EvmTestDAppAddr.Hex())
	testDApp, err := testdapp.NewTestDApp(r.EvmTestDAppAddr, r.EVMClient)
	if err != nil {
		panic(err)
	}

	res2, err := r.BankClient.SupplyOf(r.Ctx, &banktypes.QuerySupplyOfRequest{
		Denom: "azeta",
	})
	if err != nil {
		panic(err)
	}
	r.Logger.Info("$$$ Before: SUPPLY OF AZETA: %d", res2.Amount.Amount)

	tx, err = testDApp.SendHelloWorld(auth, r.EvmTestDAppAddr, chainID, amount, true)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("TestDApp.SendHello tx hash: %s", tx.Hash().Hex())
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("TestDApp.SendHello tx receipt: status %d", receipt.Status)

	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.String(), r.CctxClient, r.Logger, r.CctxTimeout)
	if cctx.CctxStatus.Status != cctxtypes.CctxStatus_Reverted {
		panic("expected cctx to be reverted")
	}
	outTxHash := cctx.GetCurrentOutTxParam().OutboundTxHash
	receipt, err = r.EVMClient.TransactionReceipt(r.Ctx, ethcommon.HexToHash(outTxHash))
	if err != nil {
		panic(err)
	}
	for _, log := range receipt.Logs {
		event, err := r.ConnectorEth.ParseZetaReverted(*log)
		if err == nil {
			r.Logger.Info("ZetaReverted event: ")
			r.Logger.Info("  Dest Addr: %s", ethcommon.BytesToAddress(event.DestinationAddress).Hex())
			r.Logger.Info("  Dest Chain: %d", event.DestinationChainId)
			r.Logger.Info("  RemainingZetaValue: %d", event.RemainingZetaValue)
			r.Logger.Info("  Message: %x", event.Message)
		}
	}
	res3, err := r.BankClient.SupplyOf(r.Ctx, &banktypes.QuerySupplyOfRequest{
		Denom: "azeta",
	})
	if err != nil {
		panic(err)
	}
	r.Logger.Info("$$$ After: SUPPLY OF AZETA: %d", res3.Amount.Amount.BigInt())
	r.Logger.Info("$$$ Diff: SUPPLY OF AZETA: %d", res3.Amount.Amount.Sub(res2.Amount.Amount).BigInt())
}
