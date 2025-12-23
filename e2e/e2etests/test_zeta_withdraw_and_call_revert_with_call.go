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

// TestZetaWithdrawAndCallRevertWithCall tests ZETA withdraw and call revert with call through gateway
func TestZetaWithdrawAndCallRevertWithCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	payload := randomPayload(r)
	evmChainID, err := r.EVMClient.ChainID(r.Ctx)
	require.NoError(r, err)

	r.ApproveETHZRC20(r.GatewayZEVMAddr)

	r.AssertTestDAppZEVMCalled(false, payload, nil, amount)

	// perform the withdraw
	tx := r.ZETAWithdrawAndArbitraryCall(
		r.TestDAppV2EVMAddr,
		amount,
		evmChainID,
		r.EncodeERC20CallRevert(r.ZetaEthAddr, amount),
		gatewayzevm.RevertOptions{
			RevertAddress:    r.TestDAppV2ZEVMAddr,
			CallOnRevert:     true,
			RevertMessage:    []byte(payload),
			OnRevertGasLimit: big.NewInt(0),
		},
	)

	if r.IsV2ZETAEnabled() {
		// V2 ZETA flows enabled: withdraw and call should revert with call
		cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient, r.Logger, r.CctxTimeout)
		r.Logger.CCTX(*cctx, "zeta_withdraw_and_call_revert_with_call")
		utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_Reverted)

		r.AssertTestDAppZEVMCalled(true, payload, nil, big.NewInt(0))

		// check expected sender was used
		senderForMsg, err := r.TestDAppV2ZEVM.SenderWithMessage(
			&bind.CallOpts{},
			[]byte(payload),
		)
		require.NoError(r, err)
		require.Equal(r, r.ZEVMAuth.From, senderForMsg)

		newBalance, err := r.ZEVMClient.BalanceAt(r.Ctx, r.TestDAppV2ZEVMAddr, nil)
		require.NoError(r, err)
		require.True(r, newBalance.Cmp(big.NewInt(0)) > 0)
	} else {
		// V2 ZETA flows disabled: tx should revert on GatewayZEVM, no CCTX created
		utils.EnsureNoCctxMinedByInboundHash(r.Ctx, tx.Hash().Hex(), r.CctxClient)
	}
}
