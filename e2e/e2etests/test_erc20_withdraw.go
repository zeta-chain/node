package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

func TestERC20Withdraw(r *runner.E2ERunner, args []string) {
	approvedAmount := big.NewInt(1e18)
	if len(args) != 1 {
		panic("TestERC20Withdraw requires exactly one argument for the withdrawal amount.")
	}

	withdrawalAmount, ok := new(big.Int).SetString(args[0], 10)
	if !ok {
		panic("Invalid withdrawal amount specified for TestERC20Withdraw.")
	}

	if withdrawalAmount.Cmp(approvedAmount) >= 0 {
		panic("Withdrawal amount must be less than the approved amount (1e18).")
	}

	// approve
	tx, err := r.ETHZRC20.Approve(r.ZevmAuth, r.ZRC20Addr, approvedAmount)
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status == 0 {
		panic("approve failed")
	}
	r.Logger.Info("eth zrc20 approve receipt: status %d", receipt.Status)

	// withdraw
	tx, err = r.ZRC20.Withdraw(r.ZevmAuth, r.DeployerAddress.Bytes(), withdrawalAmount)
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := r.ZRC20.ParseWithdrawal(*log)
		if err != nil {
			continue
		}
		r.Logger.Info(
			"  logs: from %s, to %x, value %d, gasfee %d",
			event.From.Hex(),
			event.To,
			event.Value,
			event.Gasfee,
		)
	}

	// verify the withdraw value
	cctx := utils.WaitCctxMinedByInTxHash(r.Ctx, receipt.TxHash.Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	verifyTransferAmountFromCCTX(r, cctx, withdrawalAmount.Int64())
}

// verifyTransferAmountFromCCTX verifies the transfer amount from the CCTX on Goerli
func verifyTransferAmountFromCCTX(r *runner.E2ERunner, cctx *crosschaintypes.CrossChainTx, amount int64) {
	r.Logger.Info("outTx hash %s", cctx.GetCurrentOutTxParam().OutboundTxHash)

	receipt, err := r.GoerliClient.TransactionReceipt(
		r.Ctx,
		ethcommon.HexToHash(cctx.GetCurrentOutTxParam().OutboundTxHash),
	)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Receipt txhash %s status %d", receipt.TxHash, receipt.Status)
	for _, log := range receipt.Logs {
		event, err := r.ERC20.ParseTransfer(*log)
		if err != nil {
			continue
		}
		r.Logger.Info("  logs: from %s, to %s, value %d", event.From.Hex(), event.To.Hex(), event.Value)
		if event.Value.Int64() != amount {
			panic("value is not correct")
		}
	}
}
