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

func TestZetaDepositAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	r.ApproveZetaOnEVM(r.GatewayEVMAddr)

	payload := randomPayload(r)

	r.AssertTestDAppEVMCalled(false, payload, amount)

	// perform the deposit
	// Deposit reverts and calls EVM contract
	tx := r.ZetaDepositAndCall(r.TestDAppV2ZEVMAddr, amount, []byte("revert"), gatewayevm.RevertOptions{
		RevertAddress:    r.TestDAppV2EVMAddr,
		CallOnRevert:     true,
		RevertMessage:    []byte(payload),
		OnRevertGasLimit: big.NewInt(200000),
	})

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "zeta_deposit_and_call_revert_with_call")

	if r.IsV2ZETAEnabled() {
		// V2 ZETA flows enabled: deposit and call should revert
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

		// check the payload was received on the contract
		r.AssertTestDAppEVMCalled(true, payload, big.NewInt(0))

		// check expected sender was used
		senderForMsg, err := r.TestDAppV2EVM.SenderWithMessage(
			&bind.CallOpts{},
			[]byte(payload),
		)
		require.NoError(r, err)
		require.Equal(r, r.EVMAuth.From, senderForMsg)
	} else {
		// V2 ZETA flows disabled: deposit should be aborted
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Aborted)
		require.Equal(r, cctx.CctxStatus.StatusMessage, crosschaintypes.ErrZetaThroughGateway.Error())
	}
}
