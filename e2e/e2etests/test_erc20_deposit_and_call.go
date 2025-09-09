package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestERC20DepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	payload := randomPayload(r)
	sender := r.EVMAddress().Bytes()

	r.AssertTestDAppZEVMCalled(false, payload, sender, amount)

	oldBalance, err := r.ERC20ZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)

	// perform the deposit
	tx := r.ERC20DepositAndCall(
		r.TestDAppV2ZEVMAddr,
		amount,
		[]byte(payload),
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(amount)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.ERC20ZRC20, r.TestDAppV2ZEVMAddr, oldBalance, change, r.Logger)

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, sender, amount)
}
