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

func TestSuiTokenDepositAndCall(r *runner.E2ERunner, args []string) {
	require.Len(r, args, 1)

	amount := utils.ParseBigInt(r, args[0])

	oldBalance, err := r.SuiTokenZRC20.BalanceOf(&bind.CallOpts{}, r.TestDAppV2ZEVMAddr)
	require.NoError(r, err)

	payload := randomPayload(r)

	// given sender
	signer, err := r.Account.SuiSigner()
	require.NoError(r, err)
	sender, err := sui.EncodeAddress(signer.Address())
	require.NoError(r, err)

	// make the deposit transaction
	resp := r.SuiFungibleTokenDepositAndCall(r.TestDAppV2ZEVMAddr, math.NewUintFromBigInt(amount), []byte(payload))

	r.Logger.Info("Sui deposit and call tx: %s", resp.Digest)

	// wait for the cctx to be mined
	cctx := utils.WaitCctxMinedByInboundHash(r.Ctx, resp.Digest, r.CctxClient, r.Logger, r.CctxTimeout)
	r.Logger.CCTX(*cctx, "deposit")
	utils.RequireCCTXStatus(r, cctx, crosschaintypes.CctxStatus_OutboundMined)
	require.EqualValues(r, coin.CoinType_ERC20, cctx.InboundParams.CoinType)
	require.EqualValues(r, amount.Uint64(), cctx.InboundParams.Amount.Uint64())

	// wait for the zrc20 balance to be updated
	change := utils.NewExactChange(amount)
	utils.WaitAndVerifyZRC20BalanceChange(r, r.SuiTokenZRC20, r.TestDAppV2ZEVMAddr, oldBalance, change, r.Logger)

	// check the payload was received on the contract
	r.AssertTestDAppZEVMCalled(true, payload, sender, amount)
}
