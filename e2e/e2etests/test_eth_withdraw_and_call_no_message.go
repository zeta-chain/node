package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts/pkg/gatewayzevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestETHWithdrawAndCallNoMessage(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 2)

	previousGasLimit := r.ZEVMAuth.GasLimit
	r.ZEVMAuth.GasLimit = 10000000
	defer func() {
		r.ZEVMAuth.GasLimit = previousGasLimit
	}()

	amount := utils.ParseBigInt(r, args[0])
	gasLimit := utils.ParseBigInt(r, args[1])

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	// perform the withdraw
	tx := r.ETHWithdrawAndCall(
		r.TestDAppV2EVMAddr,
		amount,
		[]byte{},
		gatewayzevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
		gasLimit,
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "withdraw_and_call")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check called
	messageIndex, err := r.TestDAppV2EVM.GetNoMessageIndex(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.AssertTestDAppEVMCalled(true, messageIndex, amount)

	// check expected sender was used
	senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(
		&bind.CallOpts{},
		[]byte(messageIndex),
	)
	require.NoError(r, err)
	require.Equal(r, r.ZEVMAuth.From, senderForMsg)
}
