package e2etests

import (
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/zeta-chain/zetacore/e2e/contracts/testdapp"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
)

// deployFunc is a function that deploys a contract
type deployFunc func(r *runner.E2ERunner) (ethcommon.Address, error)

// deployMap maps contract names to deploy functions
var deployMap = map[string]deployFunc{
	"testdapp_zevm": deployZEVMTestDApp,
	"testdapp_evm":  deployEVMTestDApp,
}

// TestDeployContract deploys the specified contract
func TestDeployContract(r *runner.E2ERunner, args []string) {
	availableContractNames := make([]string, 0, len(deployMap))
	for contractName := range deployMap {
		availableContractNames = append(availableContractNames, contractName)
	}
	availableContractNamesMessage := fmt.Sprintf("Available contract names: %v", availableContractNames)

	if len(args) != 1 {
		panic("TestDeployContract requires exactly one argument for the contract name. " + availableContractNamesMessage)
	}
	contractName := args[0]

	deployFunc, ok := deployMap[contractName]
	if !ok {
		panic(fmt.Sprintf("Unknown contract name: %s, %s", contractName, availableContractNamesMessage))
	}

	addr, err := deployFunc(r)
	if err != nil {
		panic(err)
	}

	r.Logger.Print("%s deployed at %s", contractName, addr.Hex())
}

// deployZEVMTestDApp deploys the TestDApp contract on ZetaChain
func deployZEVMTestDApp(r *runner.E2ERunner) (ethcommon.Address, error) {
	addr, tx, _, err := testdapp.DeployTestDApp(
		r.ZEVMAuth,
		r.ZEVMClient,
		r.ConnectorZEVMAddr,
		r.WZetaAddr,
	)
	if err != nil {
		return addr, err
	}

	// Wait for the transaction to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		return addr, fmt.Errorf("contract deployment failed")
	}

	return addr, nil
}

// deployEVMTestDApp deploys the TestDApp contract on Ethereum
func deployEVMTestDApp(r *runner.E2ERunner) (ethcommon.Address, error) {
	addr, tx, _, err := testdapp.DeployTestDApp(
		r.EVMAuth,
		r.EVMClient,
		r.ConnectorEthAddr,
		r.ZetaEthAddr,
	)
	if err != nil {
		return addr, err
	}

	// Wait for the transaction to be mined
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		return addr, fmt.Errorf("contract deployment failed")
	}

	return addr, nil
}
