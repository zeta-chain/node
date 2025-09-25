package legacy

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	testcontract "github.com/zeta-chain/node/e2e/contracts/depositor"
	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
)

func TestMultipleERC20Deposit(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	// parse the deposit amount and count
	depositAmount := utils.ParseBigInt(r, args[0])
	numberOfDeposits := utils.ParseBigInt(r, args[1])
	require.NotZero(r, numberOfDeposits.Int64())

	initialBal, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	txhash := multipleDeposits(r, depositAmount, numberOfDeposits)
	cctxs := utils.WaitCctxsMinedByInboundHash(
		r.Ctx,
		txhash.Hex(),
		r.CctxClient,
		int(numberOfDeposits.Int64()),
		r.Logger,
		r.CctxTimeout,
	)
	require.Len(r, cctxs, 3)

	// check new balance is increased by amount * count
	bal, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)

	diff := big.NewInt(0).Sub(bal, initialBal)
	total := depositAmount.Mul(depositAmount, numberOfDeposits)

	require.Equal(r, 0, diff.Cmp(total), "balance difference is not correct")
}

func multipleDeposits(r *runner.E2ERunner, amount, count *big.Int) ethcommon.Hash {
	// deploy depositor
	depositorAddr, txDeploy, depositor, err := testcontract.DeployDepositor(r.EVMAuth, r.EVMClient, r.ERC20CustodyAddr)
	require.NoError(r, err)
	r.WaitForTxReceiptOnZEVM(txDeploy)

	fullAmount := big.NewInt(0).Mul(amount, count)

	// approve
	tx, err := r.ERC20.Approve(r.EVMAuth, depositorAddr, fullAmount)
	require.NoError(r, err)
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	r.Logger.Info("ERC20 Approve receipt tx hash: %s", tx.Hash().Hex())

	// deposit
	tx, err = depositor.RunDeposits(r.EVMAuth, r.EVMAddress().Bytes(), r.ERC20Addr, amount, []byte{}, count)
	require.NoError(r, err)

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.Info("Deposits receipt tx hash: %s", tx.Hash().Hex())

	for _, log := range receipt.Logs {
		event, err := r.ERC20Custody.ParseDeposited(*log)
		if err != nil {
			continue
		}
		r.Logger.Info("Multiple deposit event: ")
		r.Logger.Info("  Amount: %d, ", event.Amount)
	}
	r.Logger.Info("gas limit %d", r.ZEVMAuth.GasLimit)
	return tx.Hash()
}
