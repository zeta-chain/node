package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/zetacore/e2e/runner"
	"github.com/zeta-chain/zetacore/e2e/utils"
	crosschaintypes "github.com/zeta-chain/zetacore/x/crosschain/types"
)

const payloadMessageERC20 = "this is a test ERC20 deposit and call payload"

func TestV2ERC20DepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount, ok := big.NewInt(0).SetString(args[0], 10)
	require.True(r, ok, "Invalid amount specified for TestV2ERC20DepositAndCall")

	allowance, err := r.ERC20.Allowance(&bind.CallOpts{}, r.Account.EVMAddress(), r.GatewayEVMAddr)
	require.NoError(r, err)

	// approve 1000*1e18 if allowance is zero
	if allowance.Cmp(big.NewInt(0)) == 0 {
		tx, err := r.ERC20.Approve(r.EVMAuth, r.GatewayEVMAddr, big.NewInt(0).Mul(big.NewInt(1e18), big.NewInt(1000)))
		require.NoError(r, err)
		receipt := utils.MustWaitForTxReceipt(r.Ctx, r.EVMClient, tx, r.Logger, r.ReceiptTimeout)
		require.True(r, receipt.Status == 1, "approval failed")
	}

	// perform the deposit
	tx := r.V2ERC20DepositAndCall(r.TestDAppV2ZEVMAddr, amount, []byte(payloadMessageERC20))

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	require.Equal(r, crosschaintypes.CctxStatus_OutboundMined, cctx.CctxStatus.Status)

	// check the payload was received on the contract
	message, err := r.TestDAppV2ZEVM.LastMessage(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, payloadMessageERC20, message)

	// check the amount was received on the contract
	amountReceived, err := r.TestDAppV2ZEVM.LastAmount(&bind.CallOpts{})
	require.NoError(r, err)
	require.Equal(r, amount.String(), amountReceived.String())
}
