package e2etests

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/zeta-chain/protocol-contracts-evm/pkg/gatewayevm.sol"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestZetaDepositAndCallNoMessage(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveZetaOnEVM(r.GatewayEVMAddr)
	receiverAddress := r.TestDAppV2ZEVMAddr

	oldBalance, err := r.ZEVMClient.BalanceAt(r.Ctx, receiverAddress, nil)
	require.NoError(r, err)

	// perform the deposit
	tx := r.ZetaDepositAndCall(
		receiverAddress,
		amount,
		[]byte{},
		gatewayevm.RevertOptions{OnRevertGasLimit: big.NewInt(0)},
	)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta_deposit_and_call_no_message")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)

	// check the payload was received on the contract
	messageIndex, err := r.TestDAppV2ZEVM.GetNoMessageIndex(&bind.CallOpts{}, r.EVMAddress())
	require.NoError(r, err)
	r.AssertTestDAppZEVMCalled(true, messageIndex, r.EVMAddress().Bytes(), amount)

	// check the balance was updated
	newBalance, err := r.ZEVMClient.BalanceAt(r.Ctx, receiverAddress, nil)
	require.NoError(r, err)
	require.Equal(r, new(big.Int).Add(oldBalance, amount), newBalance)
}
