package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestERC20WithdrawAndArbitraryCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	payload := randomPayload(r)

	r.AssertTestDAppEVMCalled(false, payload, amount)

	r.ApproveERC20ZRC20(r.GatewayZEVMAddr)
	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	revertAddress := r.EVMAddress()

	// perform the withdraw
	tx := r.ERC20WithdrawAndArbitraryCall(
		r.TestDAppV2EVMAddr,
		amount,
		r.EncodeERC20Call(r.ERC20Addr, amount, payload),
		gatewayzevm.RevertOptions{
			RevertAddress:    revertAddress,
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	// wait for the withdraw tx to mine on ZEVM (user has paid amount + fees),
	// then capture the post-payment ERC20-ZRC20 balance.
	receipt := utils.MustWaitForTxReceipt(r.Ctx, r.ZEVMClient, tx, r.Logger, r.ReceiptTimeout)
	utils.RequireTxSuccessful(r, receipt)
	revertBalanceBefore, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")

	// V2 arbitrary calls via ERC20Custody.withdrawAndCall route to
	// gateway.executeWithERC20, which has the same sender==0 arbitrary-call
	// branch as gateway.execute. The signer cancels these CCTXs (TSS
	// self-transfer); the observer reports the outbound failed; the CCTX
	// terminates as Reverted via the standard V2 revert flow; the principal
	// is refunded to the revert address.
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

	// destination dApp must not have been called
	r.AssertTestDAppEVMCalled(false, payload, amount)

	// revert address should receive at least the principal back
	revertBalanceAfter, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, revertAddress)
	require.NoError(r, err)
	require.EqualValues(r, new(big.Int).Add(revertBalanceBefore, amount), revertBalanceAfter)
}
