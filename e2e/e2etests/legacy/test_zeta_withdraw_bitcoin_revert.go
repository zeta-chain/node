package legacy

import (
	"math/big"

	"github.com/stretchr/testify/require"
	connectorzevm "github.com/zeta-chain/protocol-contracts-evm/pkg/zetaconnectorzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/chains"
)

func TestZetaWithdrawBTCRevert(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	// parse withdraw amount
	amount := utils.ParseBigInt(r, args[0])

	r.ZEVMAuth.Value = amount
	tx, err := r.WZeta.Deposit(r.ZEVMAuth)
	require.NoError(r, err)

	r.ZEVMAuth.Value = big.NewInt(0)
	r.Logger.Info("Deposit tx hash: %s", tx.Hash().Hex())

	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	r.Logger.EVMReceipt(*receipt, "Deposit")
	utils.RequireTxSuccessful(r, receipt)

	tx, err = r.WZeta.Approve(r.ZEVMAuth, r.ConnectorZEVMAddr, big.NewInt(1e18))
	require.NoError(r, err)

	r.Logger.Info("wzeta.approve tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)

	r.Logger.EVMReceipt(*receipt, "Approve")

	lessThanAmount := amount.Div(amount, big.NewInt(10)) // 1/10 of amount
	tx, err = r.ConnectorZEVM.Send(r.ZEVMAuth, connectorzevm.ZetaInterfacesSendInput{
		DestinationChainId:  big.NewInt(chains.BitcoinRegtest.ChainId),
		DestinationAddress:  r.EVMAddress().Bytes(),
		DestinationGasLimit: big.NewInt(400_000),
		Message:             nil,
		ZetaValueAndGas:     lessThanAmount,
		ZetaParams:          nil,
	})
	require.NoError(r, err)

	r.Logger.Info("send tx hash: %s", tx.Hash().Hex())

	receipt = utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequiredTxFailed(r, receipt)

	r.Logger.EVMReceipt(*receipt, "send")
}
