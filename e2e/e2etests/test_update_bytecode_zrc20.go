package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/contracts/testzrc20"
	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	"github.com/zeta-chain/zetacore/pkg/coin"
	"github.com/zeta-chain/zetacore/testutil/sample"
	fungibletypes "github.com/zeta-chain/zetacore/x/fungible/types"
)

// TestUpdateBytecodeZRC20 tests updating the bytecode of a zrc20 and interact with it
func TestUpdateBytecodeZRC20(r *runner.E2ERunner, _ []string) {
	// Random approval
	approved := sample.EthAddress()
	tx, err := r.ETHZRC20.Approve(r.ZEVMAuth, approved, big.NewInt(1e10))
	require.NoError(r, err)

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Deploy the TestZRC20 contract
	r.Logger.Info("Deploying contract with new bytecode")
	newZRC20Address, tx, newZRC20Contract, err := testzrc20.DeployTestZRC20(
		r.ZEVMAuth,
		r.ZEVMClient,
		big.NewInt(5),
		// #nosec G115 test - always in range
		uint8(coin.CoinType_Gas),
	)
	require.NoError(r, err)

	// Wait for the contract to be deployed
	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	// Get the code hash of the new contract
	codeHashRes, err := r.FungibleClient.CodeHash(r.Ctx, &fungibletypes.QueryCodeHashRequest{
		Address: newZRC20Address.String(),
	})
	require.NoError(r, err)

	r.Logger.Info("New contract code hash: %s", codeHashRes.CodeHash)

	// Get current info of the ZRC20
	name, err := r.ETHZRC20.Name(&bind.CallOpts{})
	require.NoError(r, err)

	symbol, err := r.ETHZRC20.Symbol(&bind.CallOpts{})
	require.NoError(r, err)

	decimals, err := r.ETHZRC20.Decimals(&bind.CallOpts{})
	require.NoError(r, err)

	totalSupply, err := r.ETHZRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)

	balance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	approval, err := r.ETHZRC20.Allowance(&bind.CallOpts{}, r.EVMAddress(), approved)
	require.NoError(r, err)

	r.Logger.Info("Updating the bytecode of the ZRC20")
	msg := fungibletypes.NewMsgUpdateContractBytecode(
		r.ZetaTxServer.MustGetAccountAddressFromName(utils.AdminPolicyName),
		r.ETHZRC20Addr.Hex(),
		codeHashRes.CodeHash,
	)
	res, err := r.ZetaTxServer.BroadcastTx(utils.AdminPolicyName, msg)
	require.NoError(r, err)

	r.Logger.Info("Update zrc20 bytecode tx hash: %s", res.TxHash)

	// Get new info of the ZRC20
	r.Logger.Info("Checking the state of the ZRC20 remains the same")
	newName, err := r.ETHZRC20.Name(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, name, newName)

	newSymbol, err := r.ETHZRC20.Symbol(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, symbol, newSymbol)

	newDecimals, err := r.ETHZRC20.Decimals(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, decimals, newDecimals)

	newTotalSupply, err := r.ETHZRC20.TotalSupply(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, 0, totalSupply.Cmp(newTotalSupply))

	newBalance, err := r.ETHZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	require.Equal(r, 0, balance.Cmp(newBalance))

	newApproval, err := r.ETHZRC20.Allowance(&bind.CallOpts{}, r.EVMAddress(), approved)
	require.NoError(r, err)
	require.Equal(r, 0, approval.Cmp(newApproval))

	r.Logger.Info("Can interact with the new code of the contract")

	testZRC20Contract, err := testzrc20.NewTestZRC20(r.ETHZRC20Addr, r.ZEVMClient)
	require.NoError(r, err)

	tx, err = testZRC20Contract.UpdateNewField(r.ZEVMAuth, big.NewInt(1e10))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	newField, err := testZRC20Contract.NewField(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, 0, newField.Cmp(big.NewInt(1e10)))

	r.Logger.Info("Interacting with the bytecode contract doesn't disrupt the zrc20 contract")
	tx, err = newZRC20Contract.UpdateNewField(r.ZEVMAuth, big.NewInt(1e5))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	newField, err = newZRC20Contract.NewField(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, 0, newField.Cmp(big.NewInt(1e5)), "new field value mismatch on bytecode contract")

	newField, err = testZRC20Contract.NewField(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, 0, newField.Cmp(big.NewInt(1e10)), "new field value mismatch on zrc20 contract")

	// can continue to operate the ZRC20
	r.Logger.Info("Checking the ZRC20 can continue to operate after state change")
	tx, err = r.ETHZRC20.Transfer(r.ZEVMAuth, approved, big.NewInt(1e14))
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	newBalance, err = r.ETHZRC20.BalanceOf(&bind.CallOpts{}, approved)
	require.NoError(r, err)
	require.Equal(r, 0, newBalance.Cmp(big.NewInt(1e14)))
}
