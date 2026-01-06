package e2etests

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/testutil/sample"
)

func TestERC20DepositRestricted(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveERC20OnEVM(r.GatewayEVMAddr)

	// perform the deposit on restricted address
	tx := r.ERC20Deposit(
		ethcommon.HexToAddress(sample.RestrictedEVMAddressTest),
		amount,
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for 5 zeta blocks
	r.WaitForBlocks(5)

	// no cctx should be created
	utils.EnsureNoCctxMinedByInboundHash(r.Ctx, tx.Hash().String(), r.CctxClient)
}
