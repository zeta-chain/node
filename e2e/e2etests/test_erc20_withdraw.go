package e2etests

import (
	"math/big"

	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
	observertypes "github.com/zeta-chain/node/x/observer/types"
)

func TestERC20Withdraw(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	r.ApproveERC20ZRC20(r.GatewayZEVMAddr)
	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.ERC20Withdraw(r.EVMAddress(), amount, gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// Get chain params for stability pool percentage
	chainParams, err := r.ObserverClient.GetChainParamsForChain(
		r.Ctx,
		&observertypes.QueryGetChainParamsForChainRequest{ChainId: evmChainID.Int64()},
	)
	require.NoError(r, err)

	// Verify gas accounting and log refund amounts
	utils.VerifyOutboundGasAccounting(r, cctx, chainParams.ChainParams.StabilityPoolPercentage, r.Logger)
}
