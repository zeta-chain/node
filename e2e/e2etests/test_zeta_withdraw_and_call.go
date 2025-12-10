package e2etests

import (
	"math/big"
	"strings"

	// "github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	// crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

// TestZetaWithdrawAndCall tests that ZETA withdraw and call through gateway
// is not supported in V2 - no CCTX should be created.
func TestZetaWithdrawAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	amount := utils.ParseBigInt(r, args[0])
	gasLimit := utils.ParseBigInt(r, args[1])

	payload := strings.ToLower(r.ZetaEthAddr.String())

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.ZETAWithdrawAndCall(
		r.TestDAppV2EVMAddr,
		amount,
		[]byte(payload),
		evmChainID,
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		gasLimit,
	)

	// ZETA withdraws through gateway are not supported in V2, verify no CCTX is created
	utils.EnsureNoCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient)

	// // wait for the cctx to be mined
	// cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	// r.Logger.CCTX(*cctx, "withdraw")
	// utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	// r.AssertTestDAppEVMCalled(true, payload, amount)
	//
	// // check expected sender was used
	// senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(
	// 	&bind.CallOpts{},
	// 	[]byte(payload),
	// )
	// require.NoError(r, err)
	// require.Equal(r, r.ZEVMAuth.From, senderForMsg)
}
