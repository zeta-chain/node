package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/e2e/contracts/testzrc20"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/testutil/sample"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// TestUpdateBytecode tests updating the bytecode of a zrc20 and interact with it
func TestUpdateBytecode(r *runner.E2ERunner) {
	// Random approval
	approved := sample.EthAddress()
	tx, err := r.ETHZRC20.Approve(r.ZevmAuth, approved, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("approval failed")
	}

	// Deploy the TestZRC20 contract
	r.Logger.Info("Deploying contract with new bytecode")
	newZRC20Address, tx, newZRC20Contract, err := testzrc20.DeployTestZRC20(
		r.ZevmAuth,
		r.ZevmClient,
		big.NewInt(5),
		// #nosec G701 test - always in range
		uint8(common.CoinType_Gas),
	)
	if err != nil {
		panic(err)
	}

	// Wait for the contract to be deployed
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("contract deployment failed")
	}

	// Get the code hash of the new contract
	codeHashRes, err := r.FungibleClient.CodeHash(r.Ctx, &fungibletypes.QueryCodeHashRequest{
		Address: newZRC20Address.String(),
	})
	if err != nil {
		panic(err)
	}
	r.Logger.Info("New contract code hash: %s", codeHashRes.CodeHash)

	// Get current info of the ZRC20
	name, err := r.ETHZRC20.Name(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	symbol, err := r.ETHZRC20.Symbol(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	decimals, err := r.ETHZRC20.Decimals(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	totalSupply, err := r.ETHZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	balance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	approval, err := r.ETHZRC20.Allowance(&bind.CallOpts{}, r.DeployerAddress, approved)
	if err != nil {
		panic(err)
	}

	r.Logger.Info("Updating the bytecode of the ZRC20")
	msg := fungibletypes.NewMsgUpdateContractBytecode(
		r.ZetaTxServer.GetAccountAddress(0),
		r.ETHZRC20Addr.Hex(),
		codeHashRes.CodeHash,
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	r.Logger.Info("Update zrc20 bytecode tx hash: %s", res.TxHash)

	// Get new info of the ZRC20
	r.Logger.Info("Checking the state of the ZRC20 remains the same")
	newName, err := r.ETHZRC20.Name(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if name != newName {
		panic("name shouldn't change upon bytecode update")
	}
	newSymbol, err := r.ETHZRC20.Symbol(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if symbol != newSymbol {
		panic("symbol shouldn't change upon bytecode update")
	}
	newDecimals, err := r.ETHZRC20.Decimals(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if decimals != newDecimals {
		panic("decimals shouldn't change upon bytecode update")
	}
	newTotalSupply, err := r.ETHZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if totalSupply.Cmp(newTotalSupply) != 0 {
		panic("total supply shouldn't change upon bytecode update")
	}
	newBalance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.DeployerAddress)
	if err != nil {
		panic(err)
	}
	if balance.Cmp(newBalance) != 0 {
		panic("balance shouldn't change upon bytecode update")
	}
	newApproval, err := r.ETHZRC20.Allowance(&bind.CallOpts{}, r.DeployerAddress, approved)
	if err != nil {
		panic(err)
	}
	if approval.Cmp(newApproval) != 0 {
		panic("approval shouldn't change upon bytecode update")
	}

	r.Logger.Info("Can interact with the new code of the contract")
	testZRC20Contract, err := testzrc20.NewTestZRC20(r.ETHZRC20Addr, r.ZevmClient)
	if err != nil {
		panic(err)
	}
	tx, err = testZRC20Contract.UpdateNewField(r.ZevmAuth, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("update new field failed")
	}
	newField, err := testZRC20Contract.NewField(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if newField.Cmp(big.NewInt(1e10)) != 0 {
		panic("new field value mismatch")
	}

	r.Logger.Info("Interacting with the bytecode contract doesn't disrupt the zrc20 contract")
	tx, err = newZRC20Contract.UpdateNewField(r.ZevmAuth, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("update new field failed")
	}
	newField, err = newZRC20Contract.NewField(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if newField.Cmp(big.NewInt(1e5)) != 0 {
		panic("new field value mismatch on bytecode contract")
	}
	newField, err = testZRC20Contract.NewField(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if newField.Cmp(big.NewInt(1e10)) != 0 {
		panic("new field value mismatch on zrc20 contract")
	}

	// can continue to operate the ZRC20
	r.Logger.Info("Checking the ZRC20 can continue to operate after state change")
	tx, err = r.ETHZRC20.Transfer(r.ZevmAuth, approved, big.NewInt(1e14))
	if err != nil {
		panic(err)
	}
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZevmClient, tx, r.Logger, r.ReceiptTimeout)
	if receipt.Status != 1 {
		panic("transfer failed")
	}
	newBalance, err = r.ETHZRC20.BalanceOf(&bind.CallOpts{}, approved)
	if err != nil {
		panic(err)
	}
	if newBalance.Cmp(big.NewInt(1e14)) != 0 {
		panic("balance not updated")
	}
}
