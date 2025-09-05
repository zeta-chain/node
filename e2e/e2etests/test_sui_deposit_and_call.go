package e2etests

import (
	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/e2e/runner"
	"github.com/zeta-chain/node/e2e/utils"
	"github.com/zeta-chain/node/pkg/coin"
	"github.com/zeta-chain/node/pkg/contracts/sui"
	crosschaintypes "github.com/zeta-chain/node/x/crosschain/types"
)

func TestSuiDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	oldBalance, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)

	payload := randomPayload(r)

	// make the deposit transaction
	resp := r.SuiDepositAndCallSUI(r.TestDAppV2ZEVMAddr, math.NewUintFromBigInt(amount), []byte(payload))

	r.Logger.Info("Sui deposit and call tx: %s", resp.Digest)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, resp.Digest, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	require.EqualValues(r, coin.CoinType_Gas, cctx.InboundParams.CoinType)
	require.EqualValues(r, amount.Uint64(), cctx.InboundParams.Amount.Uint64())
	require.True(r, cctx.InboundParams.IsCrossChainCall)

	newBalance, err := r.SUIZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)
	require.EqualValues(r, oldBalance.Add(oldBalance, amount).Uint64(), newBalance.Uint64())

	// check sender passed in the call
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err)

	sender, err := sui.EncodeAddress(signer.Address())
	require.NoError(r, err)

	actualSender, err := r.TestDAppV2ZEVM.GetSenderWithMessage(&bind.CallOpts{}, payload)
	require.NoError(r, err)
	require.EqualValues(r, sender, actualSender)

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, amount)
}
