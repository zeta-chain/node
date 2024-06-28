package runner

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	zetabitcoin "github.com/zeta-chain/zetacore/zetaclient/chains/bitcoin"
)

// DepositSolWithAmount deposits Sol on ZetaChain with a specific amount
func (runner *E2ERunner) DepositSolWithAmount(amount float64) (txHash *chainhash.Hash) {
	runner.Logger.Print("⏳ depositing Sol into ZEVM")

	// list deployer utxos
	utxos, err := runner.ListDeployerUTXOs()
	if err != nil {
		panic(err)
	}

	spendableAmount := 0.0
	spendableUTXOs := 0
	for _, utxo := range utxos {
		if utxo.Spendable {
			spendableAmount += utxo.Amount
			spendableUTXOs++
		}
	}

	if spendableAmount < amount {
		panic(fmt.Errorf(
			"not enough spendable BTC to run the test; have %f, require %f",
			spendableAmount,
			amount,
		))
	}

	runner.Logger.Info("ListUnspent:")
	runner.Logger.Info("  spendableAmount: %f", spendableAmount)
	runner.Logger.Info("  spendableUTXOs: %d", spendableUTXOs)
	runner.Logger.Info("Now sending two txs to TSS address...")

	amount = amount + zetabitcoin.DefaultDepositorFee
	txHash, err = runner.SendToTSSFromDeployerToDeposit(amount, utxos)
	if err != nil {
		panic(err)
	}
	runner.Logger.Info("send BTC to TSS txHash: %s", txHash.String())

	return txHash
}
