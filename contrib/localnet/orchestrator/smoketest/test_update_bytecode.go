package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/zeta-chain/zetacore/common"
	"github.com/zeta-chain/zetacore/contrib/localnet/orchestrator/smoketest/contracts/testzrc20"
	"github.com/zeta-chain/zetacore/testutil/sample"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// TestUpdateBytecode tests updating the bytecode of a zrc20 and interact with it
func (sm *SmokeTest) TestUpdateBytecode() {
	startTime := time.Now()
	defer func() {
		fmt.Printf("test finishes in %s\n", time.Since(startTime))
	}()
	LoudPrintf("Testing updating ZRC20 bytecode swap...\n")

	// Random approval
	approved := sample.EthAddress()
	tx, err := sm.ETHZRC20.Approve(sm.zevmAuth, approved, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	receipt := MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status != 1 {
		panic("approval failed")
	}

	// Deploy the TestZRC20 contract
	fmt.Println("Deploying contract with new bytecode")
	newZRC20Address, _, newZRC20Contract, err := testzrc20.DeployTestZRC20(sm.zevmAuth, sm.zevmClient, big.NewInt(5), uint8(common.CoinType_Gas))
	if err != nil {
		panic(err)
	}

	// Get current info of the ZRC20
	name, err := sm.ETHZRC20.Name(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	symbol, err := sm.ETHZRC20.Symbol(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	decimals, err := sm.ETHZRC20.Decimals(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	totalSupply, err := sm.ETHZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	balance, err := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	approval, err := sm.ETHZRC20.Allowance(&bind.CallOpts{}, DeployerAddress, approved)
	if err != nil {
		panic(err)
	}

	fmt.Println("Updating the bytecode of the ZRC20")
	msg := fungibletypes.NewMsgUpdateContractBytecode(
		FungibleAdminAddress,
		sm.ETHZRC20Addr,
		newZRC20Address,
	)
	res, err := sm.zetaTxServer.BroadcastTx(FungibleAdminName, msg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Update zrc20 bytecode tx hash: %s\n", res.TxHash)

	// Get new info of the ZRC20
	fmt.Println("Checking the state of the ZRC20 remains the same")
	newName, err := sm.ETHZRC20.Name(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if name != newName {
		panic("name shouldn't change upon bytecode update")
	}
	newSymbol, err := sm.ETHZRC20.Symbol(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if symbol != newSymbol {
		panic("symbol shouldn't change upon bytecode update")
	}
	newDecimals, err := sm.ETHZRC20.Decimals(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if decimals != newDecimals {
		panic("decimals shouldn't change upon bytecode update")
	}
	newTotalSupply, err := sm.ETHZRC20.TotalSupply(&bind.CallOpts{})
	if err != nil {
		panic(err)
	}
	if totalSupply.Cmp(newTotalSupply) != 0 {
		panic("total supply shouldn't change upon bytecode update")
	}
	newBalance, err := sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, DeployerAddress)
	if err != nil {
		panic(err)
	}
	if balance.Cmp(newBalance) != 0 {
		panic("balance shouldn't change upon bytecode update")
	}
	newApproval, err := sm.ETHZRC20.Allowance(&bind.CallOpts{}, DeployerAddress, approved)
	if err != nil {
		panic(err)
	}
	if approval.Cmp(newApproval) != 0 {
		panic("approval shouldn't change upon bytecode update")
	}

	fmt.Println("Can interact with the new code of the contract")
	testZRC20Contract, err := testzrc20.NewTestZRC20(sm.ETHZRC20Addr, sm.zevmClient)
	if err != nil {
		panic(err)
	}
	tx, err = testZRC20Contract.UpdateNewField(sm.zevmAuth, big.NewInt(1e10))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
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

	fmt.Println("Interacting with the bytecode contract doesn't disrupt the zrc20 contract")
	tx, err = newZRC20Contract.UpdateNewField(sm.zevmAuth, big.NewInt(1e5))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
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
	fmt.Println("Checking the ZRC20 can continue to operate after state change")
	tx, err = sm.ETHZRC20.Transfer(sm.zevmAuth, approved, big.NewInt(1e14))
	if err != nil {
		panic(err)
	}
	receipt = MustWaitForTxReceipt(sm.zevmClient, tx)
	if receipt.Status != 1 {
		panic("transfer failed")
	}
	newBalance, err = sm.ETHZRC20.BalanceOf(&bind.CallOpts{}, approved)
	if err != nil {
		panic(err)
	}
	if newBalance.Cmp(big.NewInt(1e14)) != 0 {
		panic("balance not updated")
	}
}
