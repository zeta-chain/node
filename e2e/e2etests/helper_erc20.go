package e2etests

import (
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/runner"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

// verifyTransferAmountFromCCTX verifies the transfer amount from the CCTX on EVM
func verifyTransferAmountFromCCTX(r *runner.E2ERunner, cctx *crosschaintypes.CrossChainTx, amount int64) {
	r.Logger.Info("outTx hash %s", cctx.GetCurrentOutTxParam().OutboundTxHash)

	receipt, err := r.EVMClient.TransactionReceipt(
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
